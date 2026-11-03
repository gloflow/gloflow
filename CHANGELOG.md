# Changelog

All notable changes to gloflow are documented here. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and the project uses [Semantic Versioning](https://semver.org/) (while below 1.0.0: a minor bump may break things).

Newest entry first; its version must equal the contents of `VERSION`. See [AGENTS.md](AGENTS.md).
The Go code reports `VERSION` (embedded by the root `gloflow` package, `gloflow.go`) in `gf_core.GetSysReleseInfo`.

## [0.15.0] - 2026-11-03

### Changed
- The `POST /v1/a/ue` (user event) handler moved from `go/gf_apps/gf_analytics_lib/` into `gf_solo` itself (`go/gf_apps/gf_solo/gf_solo_service/gf_solo_handlers.go`). The endpoint behaves the same; its metrics are now in group `gf_solo__user_events` (service `gf_solo`) instead of `main`/`gf_analytics`.

### Removed
- `go/gf_apps/gf_analytics_lib/` (including `gf_analytics_main`). Its remaining setup (the `gf_domains_lib` mongo indexes, domains aggregation and handlers, and the `/a` static files route) now runs directly from `gf_solo` (`gf_solo_service.go`), unchanged. The `analytics__py_stats_dirs` and `analytics__run_indexer` config keys are now unused.

## [0.14.0] - 2026-11-03

### Removed
- `py/src/gf_ml/deprecated/` (the `__gf_ml_worker` TensorFlow/Keras CIFAR-10 experiment, its tests and generated data). The palette test's input images moved to `py/src/gf_apps/gf_images/gf_images_palette/test/data/input/`.
- The `gf_ml_worker` app from the build (`meta/gf_meta.py`) and from `ops/cli__build.py`, so `-app gf_ml_worker` no longer works. Its tracked output `build/gf_apps/gf_ml_worker/` is deleted.
- Packing of `gf_images_jobs_py.so` (Python extension) into the `gf_solo` and `gf_solo_images` builds, removed for now. The Rust code is expected to be used directly by the Go code. The `rust/build/` output is unchanged.
- Commented-out TensorFlow and ML worker lines in `meta/gf_meta.py` and in the `gf_solo` and `gf_solo_images` `Dockerfile_ubuntu` files.
- `build/gf_apps/gf_solo/deprecated/__Dockerfile_alpine`.

## [0.13.1] - 2026-10-31

### Changed
- `gf_images_jobs_py` (Python extension) upgraded from pyo3/numpy 0.15 to 0.23, so it builds against Python 3.12 and 3.13 (pyo3 0.15 failed with `PyUnicode_READY` not found). The `view_numpy_arr_*` functions now take their arrays read-only; the Python API is otherwise unchanged.

## [0.13.0] - 2026-10-31

### Changed
- `gf_ml_shapes` now draws with `tiny-skia` (pure Rust 2D rasterizer) instead of its own scanline fill. Shapes have anti-aliased edges (`ANTI_ALIAS` in `canvas.rs`; `false` gives hard edges) and are centred on a continuous point, so images differ slightly from before, even with the same seed. Regenerate the train and eval sets together.

## [0.12.0] - 2026-10-31

### Added
- `run_build_rust.sh --only <crate>`: builds and copies just one crate (`gf_images_jobs`, `gf_images_jobs_py` or `gf_ml_shapes`).

## [0.11.0] - 2026-11-03

First entry of this changelog. It summarises all the work since 0.10.0 ("nation_genesis", 2022-11), which was only
tracked in git history and in the release description in `go/gf_core/gf_sys_release_info.go` (entries for 0.10.0 and earlier stay there).

### Added
- Versioning: `VERSION` and `CHANGELOG.md` at the repo root (the Go code reports `VERSION` through `gf_core.GetSysReleseInfo`, replacing the hardcoded `0.10.0.0`).
- `gf_images`: JPEG quality and output format (png as well as jpeg) options for image jobs; default dirs for storing processed images; `gf_images_solo` / `gf_solo_images` standalone images server (pure API); new flows DB functions; images added to flows automatically for uploads and external images.
- `gf_identity`: generalised session functions so external auth systems other than Auth0 can be used (SSO callback), user agent stored with sessions, `gf_session` merged into `gf_identity_core`.
- Policy hooks: policy logic is delegated to external hooks (`gf_core` hooks, `GFappHooks`) instead of built-in policies.
- CORS handling in `gf_rpc_lib` behind a flag, with Prometheus counters; per-handler flag to disable metrics.
- `gf_bookmarks`: Mongo to SQL migration with verification; `gf_tagger`: endpoint to fetch all tags.
- `gf_errors` now have an `env` column with the environment the error came from.
- Web: `gf_code_editor`, `gf_lang` (wasm compiler, GitHub Actions), `gf_maps`, new image viewer/flows browser controls, tag and share UI, user events.
- Rust: `run_build_rust.sh`, and `gf_ml_shapes` (synthetic shapes dataset generator for ML training); `run_build.sh`, `run_solo.sh`, `run_solo_images.sh`.
- `AGENTS.md` with build instructions for agents.

### Changed
- **Breaking:** storage moved from MongoDB to SQL: images are only written to SQL (reads fall back to Mongo for old data), bookmarks are read from SQL only, `gf_errors` are stored in a SQL table.
- GitHub Actions: Go, web and Rust workflows trigger the container build; `gf_solo` config loading refactored with an external config plugin.
- `gf_core` error type accepts any error; improved self-IP detection (cached); `gf_rpc_lib` no longer stores HTTP requests in the DB.
- The Rust crates no longer depend on TensorFlow (`gf_ml` temporarily deprecated).
- `gf_identity`: logged-in GET and handler error responses are not cached.

### Removed
- Login attempts in `gf_identity`; Go to Python plugin calls; the `bimg`/libvips dependency (ffmpeg is now optional, checked at runtime).
- Deprecated or moved to `deprecated/`: `gf_crawl`, Elasticsearch and crawl code in `gf_analytics`, `gf_publisher` UI, `gf_tagger` leftovers, policies, `p2p`, `ml_worker`, `web3_monitor` web, plot ops, image stats functions.
