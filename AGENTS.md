# Instructions for coding agents / LLMs

## Versioning and changelog (required for every change)

The version lives in two files at the repo root:

- `VERSION`: the current semantic version, one line, no `v` prefix (e.g. `0.12.0`).
- `CHANGELOG.md`: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) format, newest entry first.

With every change that is committed (code, config defaults, API, schema, build; not pure typo fixes in comments):

1. Pick the new version (semver; while `0.x`): **minor** for new features, behaviour changes or anything breaking
   (API, DB schema, config, CLI flags); **patch** for bug fixes, refactors, docs and tests. Use `1.0.0` only when the user says so.
2. Write the new `VERSION`.
3. Add a new section at the top of `CHANGELOG.md`, directly under the intro: `## [x.y.z] - YYYY-MM-DD` with
   `### Added` / `### Changed` / `### Fixed` / `### Removed` subsections as needed. Describe what changed for a
   user/operator in plain words (not the commit message). Never edit or reorder older entries.
4. Check that `VERSION` is semver and equals the newest changelog entry.

If several changes land in one working session before a commit, use one version bump for the whole commit.
Do not add an `[Unreleased]` section: the top entry must always match `VERSION`.

## Building Rust code

Build the Rust crates in `rust/` with `./run_build_rust.sh` (repo root). Do not run `cargo build` by hand
and copy files around; the script does the build and puts the outputs where the rest of the system expects them.

- Builds `gf_images_jobs`, `gf_images_jobs_py` and `gf_ml_shapes` (`gf_core` is built as a dependency of the first two). Crates are standalone, there is no Cargo workspace.
- Output goes to `rust/build/`: `libgf_images_jobs.a` / `.so` (used by Go) and `gf_images_jobs_py.so` (Python extension, no `lib` prefix), `gf_ml_shapes` (CLI binary).
- Options: `--static` (MUSL static build), `--debug`, `--clean`, `--verbose`, `--only <crate>` (build and copy just that crate, e.g. `--only gf_ml_shapes`).
- Needs `cargo` (nightly toolchain, see `Dockerfile__gf_builder_rust__ubuntu`). `--static` also needs the `x86_64-unknown-linux-musl` target and `musl-tools`.
- `gf_ml_shapes`: synthetic shapes image dataset generator for ML training. Pure Rust, depends on `tiny-skia` (drawing) and `png` (encoding).
  Run `rust/build/gf_ml_shapes <train|eval> [--out DIR] [--images-per-class N] [--seed N]` (default `--out data_conv_std`).
  It can be built on its own with `./run_build_rust.sh --only gf_ml_shapes` (e.g. in a container that only has this crate's sources).
  Tests: `cd rust/gf_ml_shapes && cargo test`.
- If you add a crate or change an output lib, update `run_build_rust.sh`.
- `rust/deprecated/` is not built.
