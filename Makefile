.PHONY: example
example:
	tinygo build -o example/tiny_countvowels.wasm -target wasi ./example/countvowels
	tinygo build -o example/tiny_http.wasm        -target wasi ./example/http

	GOOS=wasip1 GOARCH=wasm go build -o example/std_countvowels.wasm ./example/countvowels
	GOOS=wasip1 GOARCH=wasm go build -o example/std_http.wasm        ./example/http

test:
	extism call example/tiny_countvowels.wasm count_vowels --wasi --input "this is a test" --set-config '{"thing": "1234"}'
	extism call example/tiny_http.wasm        http_get     --wasi --log-level info --allow-host "jsonplaceholder.typicode.com"

	extism call example/std_countvowels.wasm count_vowels --wasi --input "this is a test" --set-config '{"thing": "1234"}'
	extism call example/std_http.wasm        http_get     --wasi --log-level info --allow-host "jsonplaceholder.typicode.com"