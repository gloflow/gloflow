#!/usr/bin/env bash
# Builds the Go apps (gf_solo, gf_solo_images) via ops/cli__build.py.
#
# Usage: ./run_build.sh [--in-cont] [--fetch-deps] [--app <name>]...
#   --in-cont     build inside the gf_builder_go container (default: build on the host)
#   --fetch-deps  also fetch Go dependencies
#   --app <name>  build only this app (gf_solo | gf_solo_images). repeatable. default: both
#
# Output:
#   build/gf_apps/gf_solo/gf_solo
#   build/gf_apps/gf_solo_images/gf_solo_images
#
# Host builds use the "withrust" tag, so the Rust libs must be built first (./run_build_rust.sh).
# Needs python3 with: colored delegator.py requests boto3 bs4 (ex. in a venv put first in PATH).

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLI="$ROOT_DIR/ops/cli__build.py"

IN_CONT=0; FETCH_DEPS=0
APPS=()
while [ $# -gt 0 ]; do
    case "$1" in
        --in-cont)    IN_CONT=1 ;;
        --fetch-deps) FETCH_DEPS=1 ;;
        --app)        shift; APPS+=("${1:?--app needs a name}") ;;
        -h|--help)    sed -n '2,14p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
        *) echo "unknown option: $1" >&2; exit 1 ;;
    esac
    shift
done
[ ${#APPS[@]} -gt 0 ] || APPS=(gf_solo gf_solo_images)

ARGS=(-run=build_go)
[ "$IN_CONT" = 1 ] || ARGS+=(-build_outof_cont)
[ "$FETCH_DEPS" = 1 ] && ARGS+=(-fetch_deps)

if [ "$IN_CONT" = 0 ] && [ ! -f "$ROOT_DIR/rust/build/libgf_images_jobs.so" ]; then
    echo "ERROR!! - rust libs not found in rust/build/. run ./run_build_rust.sh first." >&2
    exit 1
fi

for app in "${APPS[@]}"; do
    echo ">>> BUILD $app"
    python3 "$CLI" "${ARGS[@]}" "-app=$app"
done

echo ">>> DONE"
