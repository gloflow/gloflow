# gf_lang headless rendering: design notes

Why and how gf_lang scenes are rendered on the server. For usage and the scene format see [scene.md](scene.md).

## Background

- The Go interpreter never renders. It calls out through `GFexternAPI` (`go/gf_lang/gf_extern_api.go`): `InitEngineFun`, `SetStateFun`, `CreateCube/Sphere/LineFun`, `AnimateFun`, and the RPC hooks. Each geometry call is 12 floats: position, rotation, scale, RGB.
- In the browser, `ts/engine/gf_engine.ts` is the only place that uses three.js: flat unlit `MeshBasicMaterial`s, `LineSegments` wireframes, GLSL `ShaderMaterial`s, tween animations, and camera/background state.
- The interpreter already runs on the server, so only the rendering step was missing.

## Decision

The `GFexternAPI` is already a clean seam, so rendering is split in two steps:

1. Run the interpreter with a recording extern API (`GFsceneRecorder`). It collects a scene description (JSON) instead of drawing.
2. Render that scene description with a renderer written in Rust on `wgpu`.

```
program.gf --(Go: interpreter + GFsceneRecorder)--> scene.json --(Rust: gf_lang_render, wgpu)--> PNG
```

Scope of the first version: stills only (no video), custom shaders supported, Rust code in `gf_lang/rust`.

## Options considered

| Option | Notes |
|---|---|
| Headless Chromium (Puppeteer/Playwright) + the existing three.js engine | Least work and pixel-identical, but heavy: slow startup, big image, needs a browser pool. Software GL via SwiftShader on GPU-less servers. |
| Node + `headless-gl` | Lighter, but WebGL1 only and aging. Not recommended. |
| Pure Go software rasterizer | No dependencies and fast, but custom GLSL shaders can't run. |
| **Rust + `wgpu` (chosen)** | Real GPU rendering (or Mesa lavapipe software Vulkan), real shaders, small footprint, no browser. Cost: a second renderer to maintain. |
| CPU rasterizer in Rust | Same shader limitation as the Go one. Could be added later as a fast path for scenes without shaders. |

## How Go hands the scene to Rust

Separate process, scene as JSON on stdin or a file (`gf_lang_scene -render` runs `gf_lang_render` and pipes the scene in).

- Decoupled deploy and scaling; a renderer crash can't take down the interpreter; timeouts are easy.
- Batch, not streaming: the whole program runs first, so the entity cap can abort before rendering.
- Alternatives not taken: cgo into a Rust `staticlib` (lowest latency, but complicates cross-compiling and Docker builds, and FFI panics are dangerous), or Rust-as-WASM inside Go via `wazero` (sandboxed, but CPU only).
- The scene format is the contract and carries a `version`. Both sides reject/bump together (`GFsceneVersionInt` in Go, `SUPPORTED_VERSION` in Rust).

## Design points

- **Flat scene, not a hierarchy.** The recorder applies the coord-origin stack and rotation pivots, so entities carry final world-space transforms. This mirrors `gf_engine.ts` and keeps the renderer simple. A real hierarchy would only matter for animating subtrees.
- **Faithful to the browser engine**, including quirks (sphere rotation around the world origin, black lines, a shader staying active once set). Differences are listed in [scene.md](scene.md).
- **Shaders:** three.js prepends built-ins to `ShaderMaterial` code, so user code uses `modelViewMatrix`, `position`, etc. without declaring them. The renderer does the same and rewrites GLSL ES 1.00 to GLSL 450 for naga. Supported uniform types are `float`, `int`, `vec2`-`vec4`.
- **Rendering:** instanced draws for basic meshes (one draw per mesh kind), one draw per entity for shader materials. 4x MSAA, depth with a polygon offset on filled geometry so wireframes show.
- **Camera:** the browser camera is moved by trackball controls, which a static render doesn't have. Default is (10,10,10) looking at the origin; `--fit` frames all entities.
- **Reproducibility / safety:** rules with multiple definitions pick a random one, and there is no seed exposed yet. The server executes untrusted programs, so keep `iters_max`, the entity cap, and render timeouts in place.

## Testing

- Go: `go test ./gf_lang/go/gf_lang -run "Scene|ParseHex"`.
- Rust: `cd rust/gf_lang_render && cargo test`. Shader translation and scene parsing need no GPU; the render test skips without an adapter (`GF_RENDER_REQUIRE_GPU=1` makes it fail).
- Not done yet: golden-image comparison against the three.js screenshots in `renders/`. Output is not pixel-identical to the browser, so this needs a tolerance.

## Known issues and gaps

- `GFstateChange` has duplicate JSON tags (`RotationXf` and `RotationZf` are both `rotation_x_f`), so `encoding/json` drops both fields. The recorder reads the struct directly and is unaffected; the WASM/JS path may be.
- Animations are recorded (`animation` on cubes) but ignored by the renderer; `rpc_call`/`rpc_serve` are no-ops while recording.
- `.gf` files are JSON, so shaders must be single-line strings with `\n` escapes.
- No HTTP endpoint yet (`/v1/lang/run` is still a stub).

## Next steps

1. HTTP endpoint `/v1/lang/render`: run, render via a pool of renderer processes with timeouts, return PNG; cache by hash of (code, seed, params).
2. Seed for the interpreter RNG, for reproducible renders.
3. Golden-image tests against the browser output.
4. Animation sampling by time, then video (ffmpeg) frames.
5. Optional CPU fast path for scenes without custom shaders.
6. Longer term: make the Rust renderer the single canonical one (wgpu also targets WebGPU in the browser), to avoid maintaining two renderers.
