
---

# Build
regular build for server execution:
```console
foo@bar:~$ cd gf_lang_server
foo@bar:~$ go build -o ../build/gf_lang
```

---

web-assembly JS-environment execution:
```console
foo@bar:~$ cd gf_lang_web
foo@bar:~$ GOOS=js GOARCH=wasm go build -o ../build/gf_lang_web.wasm
```

JavaScript glue code needed to execute the Golang WASM code:
```console
foo@bar:~$ cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" build/
```

---

# start local web server to test
start it in project root - gloflow
serves the necessary compiled files and other web code from the GF project

```console
foo@bar:~$ python3 -m http.server
```

the browser URL is:
> http://localhost:8000/gf_lang/test/gf_lang_test.html

