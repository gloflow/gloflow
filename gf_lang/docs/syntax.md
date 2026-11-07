# gf_lang indentation syntax (`.gfl`)

An indentation-based surface syntax that parses into the same AST as the JSON-like `.gf` programs (`gf_lang.ParseSyntax`, `go/gf_lang/gf_syntax_parse.go`). `.gfl` files are loaded by `ParseProgramASTfromFile` (by extension); `.gf` keeps working.

## Rules
- One statement per line. 4-space indentation nests (no tabs). A trailing `:` opens a block.
- `key=value` tokens are transforms, collected into the leading modifier list: `x=9 z=2 cube` -> `[["x",9,"z",2],"cube"]`. A line without them is a flat list: `R1` -> `["R1"]`.
- Keywords: `repeat N:` -> `["*",N,scope]`, `rule Name [max=N]:` -> `["rule",Name,[["iters_max",N]],[body]]`, `if a op b:` -> `["if",[op,a,b],[body]]`, `shader Name:`.
- `repeat` body is a single scope: transform-only lines and geometry are flattened into one list. `rule`/`if` bodies are lists of statements.
- Setters: `name = value` -> `["set",name,value]` (e.g. `color = rgb(0.7 0 0)`, `material = wireframe(true)`). Other setters (`material_prop`, `push`, `pop`) use the plain form: `set material_prop (...)`.
- `(a op b)` is a binary infix expression (ops: `+ - * / % == != < > <= >= && ||`), no precedence, nest with parentheses. `name(args)` is a call -> `[name, args...]`.
- Atoms: numbers, `true`/`false`, `$vars`, `#hex`, `"strings"`, bare words. Commas are whitespace. Comments: `//`. Multi-line text: `"""` ... `"""` (dedented).
- `a | b | c` chains segments into one statement (`x=9 z=2 | y=3 cube` -> `[["x",9,"z",2],["y",3],"cube"]`). Each segment is transforms + geometry, a setter, or a list. A line ending with `|` continues on the next line. A chain with a single non-transform segment starts with `| `.
- `{...}` is a nested statement (a list that starts with a list), same content as a chain: `{z=-9 R1}` -> `[["z",-9],"R1"]`.
- A line with just a quoted string/number/bool is that atom (`"R1"` -> `"R1"`), while the bare word `R1` is the list `["R1"]`.
- `shader` block: `uniform <name> <type> <default>` lines become `["uniforms",[...]]`; other lines (`vertex """..."""`) become `[word, string]` elements.

## Example
```
lang_v 0.1

rule R1 max=20:
    repeat 2:
        x=1 y=1.5 rz=0.05 cube
    if $i == 4:
        R2
    sx=(-0.001 * $i)
    R1
```
See `examples/*.gfl`. Tests: `go test -run Syntax ./go/gf_lang`.

## Tooling
- **WASM:** the parser is pure Go (no OS dependencies) and is part of the `gf_lang_web` WASM build, exposed as `gf_lang_parse_gfl(src)` -> `{ast_json_str, error_str}`. The tests also pass with `GOOS=js GOARCH=wasm`.
- **Translating `.gf` -> `.gfl`:** `go run ./go/gf_lang_convert examples/*.gf examples/primary/*.gf` (uses `gf_lang.FormatSyntax`, which verifies that the output parses back to exactly the same AST; comments of the `.gf` are not carried over). `TestSyntaxFormatAllExamples` fails if a `.gfl` is out of date.
- **Dashboard:** an example is a `<name>.gf` and/or `<name>.gfl`; the toolbar `.gf` / `.gfl` buttons switch between them. `/v1/lang/examples` returns `code_str` (.gf) and `gfl_code_str` (.gfl) for each.
