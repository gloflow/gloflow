# Instructions for coding agents / LLMs

## Building Rust code

Build the Rust crates in `rust/` with `./run_build_rust.sh` (repo root). Do not run `cargo build` by hand
and copy files around; the script does the build and puts the outputs where the rest of the system expects them.

- Builds `gf_images_jobs` and `gf_images_jobs_py` (`gf_core` is built as their dependency). Crates are standalone, there is no Cargo workspace.
- Output goes to `rust/build/`: `libgf_images_jobs.a` / `.so` (used by Go) and `gf_images_jobs_py.so` (Python extension, no `lib` prefix).
- Options: `--static` (MUSL static build), `--debug`, `--clean`, `--verbose`.
- Needs `cargo` (nightly toolchain, see `Dockerfile__gf_builder_rust__ubuntu`). `--static` also needs the `x86_64-unknown-linux-musl` target and `musl-tools`.
- If you add a crate or change an output lib, update `run_build_rust.sh`.
- `rust/deprecated/` is not built.
