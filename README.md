# Extism Go PDK

This library can be used to write [Extism Plug-ins](https://extism.org/docs/concepts/plug-in) in Go.

## Install

Include the library with Go get:

```bash
go get github.com/extism/go-pdk
```

## Reference Documentation

You can find the reference documentation for this library on [pkg.go.dev](https://pkg.go.dev/github.com/extism/go-pdk).

## Getting Started

The goal of writing an [Extism plug-in](https://extism.org/docs/concepts/plug-in) is to compile your Go code to a Wasm module with exported functions that the host application can invoke. The first thing you should understand is creating an export. Let's write a simple program that exports a `greet` function which will take a name as a string and return a greeting string. Paste this into your `main.go`:

```go
package main

import (
	"github.com/extism/go-pdk"
)

//export greet
func greet() int32 {
	input := pdk.Input()
	greeting := `Hello, ` + string(input) + `!`
	pdk.OutputString(greeting)
	return 0
}

func main() {}
```

Some things to note about this code:

1. The `//export greet` comment is required. This marks the greet function as an export with the name `greet` that can be called by the host.
2. We need a `main` but it is unused.
3. Exports in the Go PDK are coded to the raw ABI. You get parameters from the host by calling [pdk.Input* functions](https://pkg.go.dev/github.com/extism/go-pdk#Input) and you send returns back with the [pdk.Output* functions](https://pkg.go.dev/github.com/extism/go-pdk#Output).
4. An Extism export expects an i32 return code. `0` is success and `1` is a failure.

Install the `tinygo` compiler: 

See https://tinygo.org/getting-started/install/ for instructions for your platform. 

> Note: while the core Go toolchain has support to target WebAssembly, we find `tinygo` to work well for plug-in code. Please open issues on this repository if you try building with `go build` instead & have problems!

Compile this with the command:

```bash
tinygo build -o plugin.wasm -target wasi main.go
```

We can now test `plugin.wasm` using the [Extism CLI](https://github.com/extism/cli)'s `run`
command:

```bash
extism call plugin.wasm greet --input "Benjamin" --wasi
# => Hello, Benjamin!
```

> **Note**: Currently `wasi` must be provided for all Go plug-ins even if they don't need system access, however this will eventually be optional.

> **Note**: We also have a web-based, plug-in tester called the [Extism Playground](https://playground.extism.org/)

### More Exports: Error Handling

Suppose we want to re-write our greeting module to never greet Benjamins. We can use [pdk.SetError](https://pkg.go.dev/github.com/extism/go-pdk#SetError) or [pdk.SetErrorString](https://pkg.go.dev/github.com/extism/go-pdk#SetErrorString):

```go
//export greet
func greet() int32 {
	name := string(pdk.Input())
	if name == "Benjamin" {
		pdk.SetError(errors.New("Sorry, we don't greet Benjamins!"))
		return 1
	}
	greeting := `Hello, ` + name + `!`
	pdk.OutputString(greeting)
	return 0
}
```

Now when we try again:

```bash
extism call plugin.wasm greet --input="Benjamin" --wasi
# => Error: Sorry, we don't greet Benjamins!
echo $? # print last status code
# => 1
extism call plugin.wasm greet --input="Zach" --wasi
# => Hello, Zach!
echo $?
# => 0
```

### Json

Extism export functions simply take bytes in and bytes out. Those can be whatever you want them to be. A common and simple way to get more complex types to and from the host is with json:


```go
type Add struct {
	A int `json:"a"`
	B int `json:"b"`
}

type Sum struct {
	Sum int `json:"sum"`
}

//export add
func add() int32 {
	params := Add{}
	// use json input helper, which automatically unmarshals the plugin input into your struct
	err := pdk.InputJSON(&params)
	if err != nil {
		pdk.SetError(err)
		return 1
	}
	sum := Sum{Sum: params.A + params.B}
	// use json output helper, which automatically marshals your struct to the plugin output
	output, err := pdk.OutputJSON(sum)
	if err != nil {
		pdk.SetError(err)
		return 1
	}
	return 0
}
```

```bash
extism call plugin.wasm add --input='{"a": 20, "b": 21}' --wasi
# => {"sum":41}
```

## Configs

Configs are key-value pairs that can be passed in by the host when creating a
plug-in. These can be useful to statically configure the plug-in with some data that exists across every function call. Here is a trivial example using [pdk.GetConfig](https://pkg.go.dev/github.com/extism/go-pdk#GetConfig):

```go
//export greet
func greet() int32 {
	user, ok := pdk.GetConfig("user")
	if !ok {
		pdk.SetErrorString("This plug-in requires a 'user' key in the config")
		return 1
	}
	greeting := `Hello, ` + user + `!`
	pdk.OutputString(greeting)
	return 0
}
```

To test it, the [Extism CLI](https://github.com/extism/cli) has a `--config` option that lets you pass in `key=value` pairs:

```bash
extism call plugin.wasm greet --config user=Benjamin
# => Hello, Benjamin!
```

## Variables

Variables are another key-value mechanism but it's a mutable data store that
will persist across function calls. These variables will persist as long as the
host has loaded and not freed the plug-in. 

```go
//export count
func count() int32 {
	count := pdk.GetVarInt("count")
	count = count + 1
	pdk.SetVarInt("count", count)
	pdk.OutputString(strconv.Itoa(count))
	return 0
}
```

> **Note**: Use the untyped variants [pdk.SetVar(string, []byte)](https://pkg.go.dev/github.com/extism/go-pdk#SetVar) and [pdk.GetVar(string) []byte](https://pkg.go.dev/github.com/extism/go-pdk#GetVar) to handle your own types.

## Logging

Because Wasm modules by default do not have access to the system, printing to stdout won't work (unless you use WASI).
Extism provides a simple [logging function](https://pkg.go.dev/github.com/extism/go-pdk#Log) that allows you to use the host application to log without having to give the plug-in permission to make syscalls.

```go
//export log_stuff
func logStuff() int32 {
	pdk.Log(pdk.LogInfo, "An info log!")
	pdk.Log(pdk.LogDebug, "A debug log!")
	pdk.Log(pdk.LogWarn, "A trace log!")
	pdk.Log(pdk.LogError, "An error log!")
	pdk.Log(pdk.LogTrace, "A trace log!")
	return 0
}
```

From [Extism CLI](https://github.com/extism/cli):

```bash
extism call plugin.wasm log_stuff --wasi --log-level=trace
2023/10/12 12:11:23 Calling function : log_stuff
2023/10/12 12:11:23 An info log!
2023/10/12 12:11:23 A debug log!
2023/10/12 12:11:23 A trace log!
2023/10/12 12:11:23 An error log!
```

> *Note*: From the CLI you need to pass a level with `--log-level`. If you are running the plug-in in your own host using one of our SDKs, you need to make sure that you call `set_log_file` to `"stdout"` or some file location.

## HTTP

Sometimes it is useful to let a plug-in [make HTTP calls](https://pkg.go.dev/github.com/extism/go-pdk#HTTPRequest.Send). [See this example](example/http.go)

```go
//export http_get
func httpGet() int32 {
	// create an HTTP Request (withuot relying on WASI), set headers as needed
	req := pdk.NewHTTPRequest("GET", "https://jsonplaceholder.typicode.com/todos/1")
	req.SetHeader("some-name", "some-value")
	req.SetHeader("another", "again")
	// send the request, get response back (can check status on response via res.Status())
	res := req.Send()

	pdk.OutputMemory(res.Memory())

	return 0
}
```

By default, Extism modules cannot make HTTP requests unless you specify which hosts it can connect to. You can use `--alow-host` in the Extism CLI to set this:

```
extism call plugin.wasm http_get --wasi --allow-host='*.typicode.com'
# => { "userId": 1, "id": 1, "title": "delectus aut autem", "completed": false }
```

## Imports (Host Functions)

Like any other code module, Wasm not only let's you export functions to the outside world, you can
import them too. Host Functions allow a plug-in to import functions defined in the host. For example,
if you host application is written in Python, it can pass a Python function down to your Go plug-in
where you can invoke it.

This topic can get fairly complicated and we have not yet fully abstracted the Wasm knowledge you need
to do this correctly. So we recommend reading our [concept doc on Host Functions](https://extism.org/docs/concepts/host-functions) before you get started.

### A Simple Example

Host functions have a similar interface as exports. You just need to declare them as extern on the top of your main.go. You only declare the interface as it is the host's responsibility to provide the implementation:

```go
//go:wasmimport extism:host/user a_python_func
func aPythonFunc(uint64) uint64
```

We should be able to call this function as a normal Go function. Note that we need to manually handle the pointer casting:

```go
//export hello_from_python
func helloFromPython() int32 {
    msg := "An argument to send to Python"
    mem := pdk.AllocateString(msg)
    defer mem.Free()
    ptr := aPythonFunc(mem.Offset())
    rmem := pdk.FindMemory(ptr)
    response := string(rmem.ReadBytes())
    pdk.OutputString(response)
    return 0
}

```

### Testing it out

We can't really test this from the Extism CLI as something must provide the implementation. So let's
write out the Python side here. Check out the [docs for Host SDKs](https://extism.org/docs/concepts/host-sdk) to implement a host function in a language of your choice.

```python
from extism import host_fn, Plugin

@host_fn()
def a_python_func(input: str) -> str:
    # just printing this out to prove we're in Python land
    print("Hello from Python!")

    # let's just add "!" to the input string
    # but you could imagine here we could add some
    # applicaiton code like query or manipulate the database
    # or our application APIs
    return input + "!"
```

Now when we load the plug-in we pass the host function:
 
```python
manifest = {"wasm": [{"path": "/path/to/plugin.wasm"}]}
plugin = Plugin(manifest, functions=[a_python_func], wasi=True)
result = plugin.call('hello_from_python', b'').decode('utf-8')
print(result)
```

```bash
python3 app.py
# => Hello from Python!
# => An argument to send to Python!
```

## Reactor modules

Since TinyGo doesn't support [Reactor modules](https://dylibso.com/blog/wasi-command-reactor/) yet, If you want to use WASI inside your Reactor module functions (exported functions other than `main`), you'll need to import `wasi-reactor` module which makes sure libc and go runtime are properly initialized:

```go
package main

import (
	"os"

	"github.com/extism/go-pdk"
	_ "github.com/extism/go-pdk/wasi-reactor"
)

//export read_file
func read_file() {
	name := pdk.InputString()

	content, err := os.ReadFile(name)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return
	}

	pdk.Output(content)
}

func main() {}
```

```bash
tinygo build -target wasi -o reactor.wasm .\tiny_main.go
extism call ./reactor.wasm read_file --input "./test.txt" --allow-path . --wasi --log-level info
# => Hello World!
```

Note: this is not required if you only have the `main` function.

### Reach Out!

Have a question or just want to drop in and say hi? [Hop on the Discord](https://extism.org/discord)!
