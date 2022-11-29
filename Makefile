.PHONY: example
example:
	tinygo build -o example/example.wasm -target wasi example/main.go
	tinygo build -o example/http.wasm -target wasi example/http.go

test:
	extism call example/example.wasm count_vowels --wasi --input "this is a test" --set-config '{"thing": "1234"}'	
	@echo ""
	extism call example/http.wasm http_get --wasi