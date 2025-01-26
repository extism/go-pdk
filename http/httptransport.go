package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	pdk "github.com/extism/go-pdk"
	extismhttp "github.com/extism/go-pdk/internal/http"
	"github.com/extism/go-pdk/internal/memory"
)

// HTTPTransport implement go's http.RoundTripper interface, enabling usage of standard go
// http.Client within a plugin
type HTTPTransport struct {
}

func (t *HTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	convertRequestHeaders := func() map[string]string {
		if len(req.Header) == 0 {
			return nil
		}

		result := map[string]string{}

		for name, values := range req.Header {
			result[name] = strings.Join(values, ",")

		}

		return result
	}

	meta := pdk.HTTPRequestMeta{
		URL:     req.URL.String(),
		Headers: convertRequestHeaders(),
		Method:  req.Method,
	}

	metaData, err := json.Marshal(meta)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request headers: %q", err)
	}

	metaMemory := pdk.AllocateBytes(metaData)
	defer metaMemory.Free()

	var bodyMemoryOffset memory.ExtismPointer
	if req.Body != nil {
		bodyData, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body bytes: %q", err)
		}

		bodyMemory := pdk.AllocateBytes(bodyData)
		defer bodyMemory.Free()

		bodyMemoryOffset = memory.ExtismPointer(bodyMemory.Offset())
	}

	respPointer := extismhttp.ExtismHTTPRequest(
		memory.ExtismPointer(metaMemory.Offset()),
		bodyMemoryOffset,
	)
	respLength := memory.ExtismLengthUnsafe(respPointer)
	respStatus := extismhttp.ExtismHTTPStatusCode()

	headersPointer := extismhttp.ExtismHTTPHeaders()
	respHeaders := map[string]string{}

	if headersPointer != 0 {
		headersLength := memory.ExtismLengthUnsafe(headersPointer)
		headersMemory := memory.NewMemory(headersPointer, headersLength)
		defer headersMemory.Free()
		json.Unmarshal(headersMemory.ReadBytes(), &respHeaders)
	}

	convertResponseHeaders := func() http.Header {
		result := http.Header{}
		for key, value := range respHeaders {
			result.Add(key, value)
		}

		return result
	}

	resp := &http.Response{
		Status:        http.StatusText(int(respStatus)),
		StatusCode:    int(respStatus),
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        convertResponseHeaders(),
		Body:          nil,
		ContentLength: -1,
		Request:       req,
	}

	hasBody := req.Method != "HEAD" && respLength > 0
	if hasBody {
		respMemory := memory.NewMemory(respPointer, respLength)
		respBuf := make([]byte, respMemory.Length())
		respMemory.Load(respBuf)

		resp.Body = io.NopCloser(bytes.NewReader(respBuf))
		resp.ContentLength = int64(respLength)
	}

	return resp, nil
}
