#!/usr/bin/env bash
# Builds the Rust crates in rust/ and copies the output libs into rust/build/.
#
# Usage: ./run_build_rust.sh [--static] [--debug] [--clean] [--verbose] [--only CRATE]
#   --static   MUSL static build (target x86_64-unknown-linux-musl)
#   --debug    debug build instead of release
#   --clean    run "cargo clean" in each crate first
#   --verbose  cargo -vv
#   --only     build and copy just this crate: gf_images_jobs, gf_images_jobs_py or gf_ml_shapes
#              (e.g. to build the dataset generator gf_ml_shapes alone)
#
# Output (rust/build/):
#   libgf_images_jobs.a, libgf_images_jobs.so  - for Go (static/dynamic linking)
#   gf_images_jobs_py.so                       - Python extension (no "lib" prefix)
#   gf_ml_shapes                               - CLI binary, synthetic shapes dataset generator
#
# gf_core is built as a dependency of the other crates.

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RUST_DIR="$ROOT_DIR/rust"
BUILD_DIR="$RUST_DIR/build"

STATIC=0; DEBUG=0; CLEAN=0; VERBOSE=0; ONLY=""
EXPECT_ONLY=0
for arg in "$@"; do
    if [ "$EXPECT_ONLY" = 1 ]; then ONLY="$arg"; EXPECT_ONLY=0; continue; fi
    case "$arg" in
        --static)  STATIC=1 ;;
        --debug)   DEBUG=1 ;;
        --clean)   CLEAN=1 ;;
        --verbose) VERBOSE=1 ;;
        --only)    EXPECT_ONLY=1 ;;
        --only=*)  ONLY="${arg#*=}" ;;
        -h|--help) sed -n '2,16p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
        *) echo "unknown option: $arg" >&2; exit 1 ;;
    esac
done

if [ "$EXPECT_ONLY" = 1 ]; then echo "--only needs a crate name" >&2; exit 1; fi
case "$ONLY" in
    ""|gf_images_jobs|gf_images_jobs_py|gf_ml_shapes) ;;
    *) echo "unknown crate for --only: $ONLY" >&2; exit 1 ;;
esac
want() { [ -z "$ONLY" ] || [ "$ONLY" = "$1" ]; }

if ! command -v cargo >/dev/null 2>&1; then
    echo "ERROR!! - cargo not found. install Rust (https://rustup.rs)." >&2
    echo "          GF crates are built with the nightly toolchain: rustup default nightly" >&2
    exit 1
fi

CARGO_ARGS=()
TARGET_SUBDIR="release"
if [ "$STATIC" = 1 ]; then
    # needs: rustup target add x86_64-unknown-linux-musl ; apt install musl-tools
    CARGO_ARGS+=(--target x86_64-unknown-linux-musl)
    TARGET_SUBDIR="x86_64-unknown-linux-musl/debug"
    [ "$DEBUG" = 1 ] || { CARGO_ARGS+=(--release); TARGET_SUBDIR="x86_64-unknown-linux-musl/release"; }
elif [ "$DEBUG" = 1 ]; then
    TARGET_SUBDIR="debug"
else
    CARGO_ARGS+=(--release)
fi
[ "$VERBOSE" = 1 ] && CARGO_ARGS+=(-vv)

build_crate() {
    local name="$1"
    echo ">>> BUILD $name"
    pushd "$RUST_DIR/$name" >/dev/null
    [ "$CLEAN" = 1 ] && cargo clean
    if [ "$STATIC" = 1 ] && [ "$name" = "gf_images_jobs_py" ]; then
        # the py extension is a dynamic lib, so it can't be crt-static
        RUSTFLAGS='-C target-feature=-crt-static' cargo build "${CARGO_ARGS[@]}"
    else
        cargo build "${CARGO_ARGS[@]}"
    fi
    popd >/dev/null
}

want gf_images_jobs    && build_crate gf_images_jobs
want gf_images_jobs_py && build_crate gf_images_jobs_py
want gf_ml_shapes      && build_crate gf_ml_shapes

echo ">>> COPY libs into $BUILD_DIR"
mkdir -p "$BUILD_DIR"

JOBS_OUT="$RUST_DIR/gf_images_jobs/target/$TARGET_SUBDIR"
PY_OUT="$RUST_DIR/gf_images_jobs_py/target/$TARGET_SUBDIR"
SHAPES_OUT="$RUST_DIR/gf_ml_shapes/target/$TARGET_SUBDIR"

if want gf_images_jobs; then
    cp -v "$JOBS_OUT/libgf_images_jobs.a" "$BUILD_DIR/"
    if [ -f "$JOBS_OUT/libgf_images_jobs.so" ]; then cp -v "$JOBS_OUT/libgf_images_jobs.so" "$BUILD_DIR/"; fi
fi
# Python requires extension libs without the "lib" prefix
if want gf_images_jobs_py; then cp -v "$PY_OUT/libgf_images_jobs_py.so" "$BUILD_DIR/gf_images_jobs_py.so"; fi
if want gf_ml_shapes; then cp -v "$SHAPES_OUT/gf_ml_shapes" "$BUILD_DIR/"; fi

echo ">>> DONE"
