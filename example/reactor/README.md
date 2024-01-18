## Reactor module example
By including this package, you'll turn your plugin into a [Reactor](https://dylibso.com/blog/wasi-command-reactor/) module. This makes sure that you can use WASI (e.g. File Access) in your exported functions.

To test this example, run:

```bash
tinygo build -target wasi -o reactor.wasm .\tiny_main.go
extism call ./reactor.wasm read_file --input "./test.txt" --allow-path . --wasi --log-level info
# => Hello World!
```

If you don't include the pacakge, you'll see this output:
```bash
extism call .\c.wasm read_file --input "./test.txt" --allow-path . --wasi --log-level info
# => 2024/01/18 20:48:48 open ./test.txt: errno 76
```