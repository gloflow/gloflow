#!/usr/bin/env bash
# Builds and runs gf_solo (full gf_solo server) locally.
#
# Usage: ./run_solo.sh [--no-build] [--rust] [-- <extra args for "gf_solo start service">]
#   --no-build  skip "go build", run the existing binary
#   --rust      build with the "withrust" tag (needs ./run_build_rust.sh run first)
#
# Binary: build/gf_apps/gf_solo/gf_solo (gitignored)
# Config: build/gf_apps/gf_solo/config/gf_solo.yaml
#         (the binary is run from build/gf_apps/gf_solo/, so ./config/... paths resolve)
#
# The web UI is served from static/ and templates/ in that dir
# (built with: python3 ops/cli__build.py -run=build_web -app=gf_solo).
# Needs MongoDB and Postgres reachable as configured in the yaml (or GF_* ENV vars,
# ex. GF_PORT, GF_ENV, GF_MONGODB_HOST, GF_SQL_HOST).

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
APP_DIR="$ROOT_DIR/build/gf_apps/gf_solo"
BIN_PATH="$APP_DIR/gf_solo"

BUILD=1; RUST=0
EXTRA_ARGS=()
while [ $# -gt 0 ]; do
    case "$1" in
        --no-build) BUILD=0 ;;
        --rust)     RUST=1 ;;
        -h|--help)  sed -n '2,15p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
        --)         shift; EXTRA_ARGS=("$@"); break ;;
        *) echo "unknown option: $1" >&2; exit 1 ;;
    esac
    shift
done

if [ "$BUILD" = 1 ]; then
    echo ">>> BUILD gf_solo"
    BUILD_ARGS=(-buildvcs=false)
    export LD_LIBRARY_PATH="$ROOT_DIR/rust/build${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}"
    [ "$RUST" = 1 ] && BUILD_ARGS+=(-tags withrust)
    (cd "$ROOT_DIR" && go build "${BUILD_ARGS[@]}" -o "$BIN_PATH" ./go/gf_apps/gf_solo)
fi

if [ ! -x "$BIN_PATH" ]; then
    echo "ERROR!! - $BIN_PATH not found. run without --no-build." >&2
    exit 1
fi

echo ">>> RUN gf_solo"
cd "$APP_DIR"
export LD_LIBRARY_PATH="$ROOT_DIR/rust/build${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}"
exec "$BIN_PATH" start service "${EXTRA_ARGS[@]}"
