# gf_lang scene graph and headless rendering

The interpreter does not draw. It calls out through `GFexternAPI`. For server-side rendering the
`GFsceneRecorder` (`go/gf_lang/gf_scene.go`) implements that API and records everything into a
`GFscene`, a plain JSON-serializable description. The headless renderer (`rust/gf_lang_render`, wgpu)
turns the scene into a PNG.

```
program.gf --(Go interpreter + GFsceneRecorder)--> scene.json --(Rust, wgpu)--> out.png
```

## Usage

Build everything with `./run_build_gf_lang.sh` in the gloflow root (see its `--help`), then:

```console
cd gf_lang/go/build
./gf_lang_scene -in ../../examples/shader_scene.gf -out scene.json          # scene graph only
./gf_lang_scene -in ../../examples/shader_scene.gf -render out.png -fit     # run + render
../../rust/build/gf_lang_render --scene scene.json --out out.png --width 1024 --height 768 --fit
```

`gf_lang_scene` flags: `-in`, `-out`, `-render`, `-renderer` (path to `gf_lang_render`), `-width`, `-height`,
`-fit`, `-max-entities` (default 1,000,000; the scene gets `"truncated": true` if hit).

`gf_lang_render` reads the scene from `--scene FILE` (or stdin), flags: `--out`, `--width`, `--height`,
`--samples 1|2|4` (MSAA, default 4), `--fit`, and the explicit camera `--camera-pos x,y,z --camera-target x,y,z
[--camera-up x,y,z] [--fov DEG]` (up defaults to +Y, fov to 75). `--camera-pos` and `--camera-target` go together,
and the explicit camera takes precedence over `--fit`.

`--fit` frames all entities, looking from the scene camera's direction. The camera flags render from an exact
view instead (the dashboard passes the user's trackball view this way, through the `camera` field of `POST /v1/lang/render`). The browser engine has interactive
trackball controls, a static render does not, so without `--fit` you get the fixed default camera
(position 10,10,10 looking at the origin, fov 75).

In Go: `gf_lang.SceneRun(programAST, maxEntities)` returns the `*GFscene`.

The renderer needs a Vulkan (or Metal/DX12/GL) driver. On servers without a GPU install a software one,
on Ubuntu `apt install mesa-vulkan-drivers` (lavapipe). Tests: `cd rust/gf_lang_render && cargo test`
(the render test is skipped if there is no adapter; `GF_RENDER_REQUIRE_GPU=1` makes it fail instead).

## Scene format (version 1)

```json
{
  "version": 1,
  "background": [0.5, 0.5, 0.5],
  "camera": {"position": [10,10,10], "target": [0,0,0], "fov_deg": 75, "near": 0.1, "far": 10000},
  "shaders": {
    "name": {"uniforms": [{"name": "i", "type": "float", "default": 0.1}], "vertex": "...", "fragment": "..."}
  },
  "entities": [
    {"kind": "cube", "position": [0,0,0], "rotation": [0,0,0], "scale": [1,1,1], "color": [0,0,1],
     "wireframe": true, "material": {"type": "basic"}},
    {"kind": "sphere", "...": "...",
     "material": {"type": "shader", "shader": "name", "uniforms": {"i": 0.3}}},
    {"kind": "line", "points": [[0,0,0],[1,2,3]], "color": [0,0,0], "...": "..."}
  ],
  "truncated": false
}
```

- `background`: omitted means black (the browser engine default).
- Positions are final world-space values: the recorder applies the coord origin stack and rotation pivots, mirroring `ts/engine/gf_engine.ts`.
- `rotation`: euler XYZ in radians (three.js order; model matrix = T * Rx*Ry*Rz * S).
- `cube`: unit cube centered at the origin. `sphere`: radius **10** (like the browser engine), scaled by `scale`.
- `basic` material: flat unlit `color` (like `MeshBasicMaterial`), pushed back in depth so wireframe edges drawn over it are visible.
- `wireframe`: black edges overlay. `line`: a segment (black, 1px), made from two consecutive `line` calls.
- `shader` material: `uniforms` is a snapshot of the uniform values at the time the entity was created.
- An entity created from a cube may also carry an `animation` (recorded from `animate`); the renderer ignores it (stills only).
- Unknown versions are rejected by the renderer; bump `GFsceneVersionInt` (Go) and `SUPPORTED_VERSION` (Rust) together on breaking changes.

## Custom shaders

Shaders are written for three.js `ShaderMaterial` (GLSL ES 1.00) in the gf_lang program. The renderer
translates them to Vulkan GLSL 450 (`rust/gf_lang_render/src/shader.rs`) and compiles them with naga:

- built-ins like three.js: uniforms `projectionMatrix`, `modelViewMatrix`, `viewMatrix`, `modelMatrix`, `normalMatrix`, `cameraPosition`; vertex attributes `position`, `normal`, `uv`
- `varying`, `attribute`, `uniform`, `precision`, `gl_FragColor` are rewritten; fragment varyings are matched to vertex varyings by name
- user uniform types: `float`, `int`, `vec2`, `vec3`, `vec4`. Other types (matrices, arrays, samplers/textures) are rejected with an error
- compile errors are reported with the translated source

In `.gf` files (pure JSON) shader code has to be a single-line string with `\n` escapes (the TS examples use template literals).

## Differences from the browser engine

- Colors are written as-is (no color-space conversion); 4x MSAA like `antialias: true`.
- Sphere wireframe shows latitude/longitude lines only (no triangulation diagonals).
- Animations and `rpc_call`/`rpc_serve` are not supported while recording (RPC calls are no-ops).
- Camera is static (see `--fit`).
- The recorder mirrors some browser engine quirks on purpose: spheres rotate around the world origin (ignoring pivots), lines are black, once a shader material is set it stays active, shader uniform changes apply to the active shader.
