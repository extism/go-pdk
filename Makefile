.PHONY: example
example:
	tinygo build -o example/example.wasm -target wasi example/main.go

test:
	extism call example/example.wasm count_vowels --wasi --input "this is a test" --set-config '{"thing": "1234"}'