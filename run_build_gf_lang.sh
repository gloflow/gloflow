#!/usr/bin/env bash
# Builds gf_lang: the Go parts (interpreter server, scene recorder CLI, WASM), the Rust headless
# renderer (wgpu), and the TS code of the experiments/examples (browser test page).
#
# Usage: ./run_build_gf_lang.sh [--skip-go] [--skip-rust] [--skip-ts] [--debug] [--clean] [--verbose]
#   --skip-go    don't build the Go binaries/WASM
#   --skip-rust  don't build the Rust renderer
#   --skip-ts    don't build the TS experiments/examples
#   --debug      debug build of the Rust renderer instead of release
#   --clean      "cargo clean" the renderer first
#   --verbose    verbose cargo/go/tsc output
#
# Output:
#   gf_lang/go/build/gf_lang             interpreter server (gf_lang_server)
#   gf_lang/go/build/gf_lang_scene       CLI: runs a .gf program and writes the scene graph JSON (or renders it)
#   gf_lang/go/build/gf_lang_web.wasm    interpreter for the browser (+ wasm_exec.js in gf_lang/build/)
#   gf_lang/rust/build/gf_lang_render    headless wgpu renderer: scene JSON -> PNG
#   gf_lang/build/                       gf_lang.js (TS experiments/examples), gf_lang_test.html, gf_ide.css, gf_lang_web.wasm
#
# Needs: go, cargo (Rust), tsc. The Rust renderer needs a Vulkan driver at runtime; on machines without
# a GPU install a software one (Ubuntu: apt install mesa-vulkan-drivers).

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LANG_DIR="$ROOT_DIR/gf_lang"
GO_BUILD_DIR="$LANG_DIR/go/build"
RUST_CRATE_DIR="$LANG_DIR/rust/gf_lang_render"
RUST_BUILD_DIR="$LANG_DIR/rust/build"
WEB_BUILD_DIR="$LANG_DIR/build"

DO_GO=1; DO_RUST=1; DO_TS=1; DEBUG=0; CLEAN=0; VERBOSE=0
for arg in "$@"; do
    case "$arg" in
        --skip-go)   DO_GO=0 ;;
        --skip-rust) DO_RUST=0 ;;
        --skip-ts)   DO_TS=0 ;;
        --debug)     DEBUG=1 ;;
        --clean)     CLEAN=1 ;;
        --verbose)   VERBOSE=1 ;;
        -h|--help)   sed -n '2,22p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
        *) echo "unknown option: $arg" >&2; exit 1 ;;
    esac
done

# rustup installs cargo in ~/.cargo/bin, which non-login shells often don't have in PATH
if ! command -v cargo >/dev/null 2>&1 && [ -x "$HOME/.cargo/bin/cargo" ]; then
    export PATH="$HOME/.cargo/bin:$PATH"
fi

need() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "ERROR!! - '$1' not found in PATH. $2" >&2
        exit 1
    fi
}

# check all needed tools first, so a missing one doesn't fail the build half way
[ "$DO_GO" = 1 ]   && need go    "Install Go (https://go.dev/dl/)."
[ "$DO_RUST" = 1 ] && need cargo "Install Rust (https://rustup.rs)."
[ "$DO_TS" = 1 ]   && need tsc   "Install TypeScript: npm install -g typescript."

mkdir -p "$GO_BUILD_DIR" "$RUST_BUILD_DIR" "$WEB_BUILD_DIR"

#------------------------------------------------------------
# GO
if [ "$DO_GO" = 1 ]; then
    GO_ARGS=()
    [ "$VERBOSE" = 1 ] && GO_ARGS+=(-v)

    cd "$ROOT_DIR"
    echo ">>> BUILD GO gf_lang_server"
    go build "${GO_ARGS[@]}" -o "$GO_BUILD_DIR/gf_lang" ./gf_lang/go/gf_lang_server

    echo ">>> BUILD GO gf_lang_scene"
    go build "${GO_ARGS[@]}" -o "$GO_BUILD_DIR/gf_lang_scene" ./gf_lang/go/gf_lang_scene

    echo ">>> BUILD GO gf_lang_web (WASM)"
    GOOS=js GOARCH=wasm go build "${GO_ARGS[@]}" -o "$GO_BUILD_DIR/gf_lang_web.wasm" ./gf_lang/go/gf_lang_web
    cp "$GO_BUILD_DIR/gf_lang_web.wasm" "$WEB_BUILD_DIR/gf_lang_web.wasm"

    # JS glue code for the WASM (location depends on the Go version)
    GOROOT_DIR="$(go env GOROOT)"
    for f in "$GOROOT_DIR/lib/wasm/wasm_exec.js" "$GOROOT_DIR/misc/wasm/wasm_exec.js"; do
        if [ -f "$f" ]; then cp "$f" "$WEB_BUILD_DIR/wasm_exec.js"; break; fi
    done
    [ -f "$WEB_BUILD_DIR/wasm_exec.js" ] || echo "WARNING!! - wasm_exec.js not found in GOROOT" >&2
fi

#------------------------------------------------------------
# RUST
if [ "$DO_RUST" = 1 ]; then
    echo ">>> BUILD RUST gf_lang_render"
    CARGO_ARGS=()
    SUBDIR="debug"
    [ "$DEBUG" = 1 ] || { CARGO_ARGS+=(--release); SUBDIR="release"; }
    [ "$VERBOSE" = 1 ] && CARGO_ARGS+=(-vv)

    pushd "$RUST_CRATE_DIR" >/dev/null
    [ "$CLEAN" = 1 ] && cargo clean
    cargo build "${CARGO_ARGS[@]}"
    popd >/dev/null

    cp -v "$RUST_CRATE_DIR/target/$SUBDIR/gf_lang_render" "$RUST_BUILD_DIR/"
fi

#------------------------------------------------------------
# TS - experiments and examples (browser test page)
if [ "$DO_TS" = 1 ]; then
    echo ">>> BUILD TS gf_lang (experiments/examples)"
    cd "$LANG_DIR"
    tsc --module system --target es2016 --outFile "$WEB_BUILD_DIR/gf_lang.js" \
        ts/gf_lang.ts \
        test/gf_lang_test.ts

    cp test/gf_lang_test.html "$WEB_BUILD_DIR/gf_lang_test.html"
    cp css/gf_ide.css         "$WEB_BUILD_DIR/gf_ide.css"
fi

echo ">>> DONE"
