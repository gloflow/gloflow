
This is a compact, data-first, S-expression-style DSL specialized for procedural 3D scene/animation generation with:

  - recursive rule-based generation
    - Rule declaration form: `["rule", "<name>", [<modifiers>], [<body...>]]`.
    - Recursive/iterative invocation by name: `["R1"]` or `"R1"` inside a rule body to continue/grow structures.
    - Iteration limits via modifiers: `["iters_max", N]` to bound recursion.
    - Rule state isolation: rules mutate local state that does not automatically leak to callers (examples use this to control independent growth).
    - Examples exist as named ASTs (e.g., `rules_test__program_ast_lst`).

  - transform and material/state setters
    - Transform/state tokens: `["x", n]`, `["y", n]`, `["z", n]`, `["sx", n]`, `["ry", n]`, etc.
    - State setters: `["set", "color", ["rgb", r,g,b]]`, `["set", "material", ["wireframe", true]]`.
    - Hex and rgb color forms both accepted (e.g., `"#334ea2"` or `["rgb", 0.7,0,0]`).
    - Material properties and shader uniform setting: `["set", "material_prop", ["<shader>", "shader_uniform", ["i", <value>]]]`.
    - State propagation semantics: setters affect subsequent expressions until overridden.

  - conditionals and simple control flow
    - Conditional form: `["if", [<op>, <lhs>, <rhs>], [<then-body...>]]`.
    - Supports comparison operators: `==`, `>`, `<`, `!=`, etc.
    - System variables accessible in conditions: `$i` is the rule iteration counter.
    - Debug/aux functions used inside flow: `["print", [...]]`.
    - Example usage: branching to other rules when `["==", "$i", 4]` or `[" >", "$i", 8]`.

  - arithmetic and random functions
    - Arithmetic expressions inline: `["*", a, b]`, `["+", a, b]`, etc.; operands may be numbers, variables, or nested expressions.
    - Multiplication used both for runtime arithmetic and historically as a compile-time replication macro (comments show expanded semantics).
    - Random number generation: `["rand", [min, max]]` for stochastic variation.
    - Complex/nested numeric expressions allowed as state setter operands (e.g., `["sx", ["*", -0.0010, "$i"]]`).

  - shader embedding with uniform binding to language variables
    - Shader block form: `["shader", "<name>", ["uniforms", [[name, type, default], ...]], ["vertex", `...`], ["fragment", `...`]]`.
    - GLSL source embedded directly as vertex/fragment strings in the AST.
    - Uniforms declared for validation and defaults; language can set those uniforms via `["set","material_prop", ...]`.
    - Examples bind language variables like `$cr`, `$cg`, `$cb` to shader uniforms to pass color/state into GLSL.

  - simple animation primitives and pivot stack operations
    - Animation primitive: `["animate", <properties>, <duration>, <mode>]` (e.g., `["animate", [["x",10], ["y",20]], 2, "repeat"]`).
    - Line drawing API: `["set", "line", ["start"]]` and `"line"` used in rule bodies to create line sequences.
    - Pivot stack operations: `["push", "rotation_pivot", "current_pos"]` and `["pop", "rotation_pivot", "current_pos"]` to localize transforms.
    - Setters+animate combine to animate transforms or camera/global properties; animations can be repeatable.

Additional structural notes (general context and host integration)
  - The DSL is represented as native JS/TS arrays (AST-like) inside `gloflow/gf_lang/test/gf_examples.ts` and returned by `get()` as named example ASTs — the interpreter consumes these arrays directly.
  - Programs begin with `["lang_v", "<version>"]` indicating language versioning.
  - Geometry tokens are simple strings: `"cube"`, `"sphere"`, `"line"`, etc.
  - The design favors concise expression trees (S-expression style) for compact procedural scene specification.

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
