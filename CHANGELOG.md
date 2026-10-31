# Changelog

All notable changes to gloflow are documented here. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and the project uses [Semantic Versioning](https://semver.org/) (while below 1.0.0: a minor bump may break things).

Newest entry first; its version must equal the contents of `VERSION`. See [AGENTS.md](AGENTS.md).
The Go code reports `VERSION` (embedded by the root `gloflow` package, `gloflow.go`) in `gf_core.GetSysReleseInfo`.

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
