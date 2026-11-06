#!/usr/bin/env bash
# Builds and runs the gf_lang server (gf_lang/go/gf_lang_server).
#
# Usage: [GF_LANG_PORT=5010] ./run_gf_lang_server.sh [--no-build] [-h|--help]
#   --no-build  skip "go build", run the existing binary
#
# Binary: gf_lang/go/build/gf_lang (same output as ./run_build_gf_lang.sh, gitignored)
#
# The server runs a test program (tests/first_scene.gf), then runs the Python debug analyzer
# (gf_lang/py/gf_debug_analyzer.py, needs python3), then serves the gf_lang service on
# HTTP port 5010, or GF_LANG_PORT (blocking). It writes serialized_output.json and
# state_history.json into its working dir.
#
# The service also serves the dev dashboard (editor + render) at /, from gf_lang/build/.
# Build it with ./run_build_gf_lang.sh.
#
# The binary is run from gf_lang/go/gf_lang_server/ because the paths in main.go
# (./tests/first_scene.gf, ./../../py/) are relative to the working dir.

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVER_DIR="$ROOT_DIR/gf_lang/go/gf_lang_server"
BIN_DIR="$ROOT_DIR/gf_lang/go/build"
BIN_PATH="$BIN_DIR/gf_lang"

BUILD=1
while [ $# -gt 0 ]; do
    case "$1" in
        --no-build) BUILD=0 ;;
        -h|--help)  sed -n '2,19p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
        *) echo "unknown option: $1" >&2; exit 1 ;;
    esac
    shift
done

if [ "$BUILD" = 1 ]; then
    echo ">>> BUILD gf_lang_server"
    mkdir -p "$BIN_DIR"
    (cd "$ROOT_DIR" && go build -buildvcs=false -o "$BIN_PATH" ./gf_lang/go/gf_lang_server)
fi

if [ ! -x "$BIN_PATH" ]; then
    echo "ERROR!! - $BIN_PATH not found. run without --no-build." >&2
    exit 1
fi

PORT="${GF_LANG_PORT:-5010}"
if [ ! -f "$ROOT_DIR/gf_lang/build/index.html" ]; then
    echo "WARNING!! - the dashboard is not built. run ./run_build_gf_lang.sh" >&2
fi

echo ">>> RUN gf_lang_server (port $PORT)"
echo ""
echo "    open in the browser:  http://localhost:$PORT/   (dev dashboard)"
echo "    JSON API:             POST http://localhost:$PORT/v1/lang/run"
echo ""
cd "$SERVER_DIR"
exec "$BIN_PATH"
