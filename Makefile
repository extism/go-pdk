.PHONY: example
example:
	tinygo build -o example/tiny_countvowels.wasm -target wasi ./example/countvowels
	tinygo build -o example/tiny_http.wasm        -target wasi ./example/http
	tinygo build -o example/tiny_reactor.wasm     -target wasi ./example/reactor

	GOOS=wasip1 GOARCH=wasm go build -tags std -o example/std_countvowels.wasm ./example/countvowels
	GOOS=wasip1 GOARCH=wasm go build -tags std -o example/std_http.wasm        ./example/http

test:
	extism call example/tiny_countvowels.wasm count_vowels --wasi --input "this is a test" --set-config '{"thing": "1234"}'
	extism call example/tiny_http.wasm        http_get     --wasi --log-level info --allow-host "jsonplaceholder.typicode.com"
	extism call example/tiny_reactor.wasm read_file --input "example/reactor/test.txt" --allow-path ./example/reactor --wasi --log-level info

	extism call example/std_countvowels.wasm _start     --wasi --input "this is a test" --set-config '{"thing": "1234"}'
	extism call example/std_http.wasm        _start     --wasi --log-level info --allow-host "jsonplaceholder.typicode.com"