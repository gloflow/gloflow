


regular build for server execution:
```console
foo@bar:~$ go build -o build/gf_lang
```

---

web-assembly JS-environment execution:
```console
foo@bar:~$ GOOS=js GOARCH=wasm go build -o build/gf_lang.wasm
```

JavaScript glue code needed to execute the Golang WASM code:
```console
foo@bar:~$ cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" build/
```
