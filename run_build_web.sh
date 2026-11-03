#!/usr/bin/env bash
# Builds the web code (ts/js/css/html) of the user UI (all pages of gf_solo) via ops/cli__build.py.
#
# Usage: ./run_build_web.sh [--in-cont] [--page <name>] [--app <name>]
#   --in-cont      build inside the gf_builder_web container (default: build on the host)
#   --page <name>  build only this page (ex. gf_home_main), to speed up the build
#   --app <name>   web app to build (default: gf_solo, which contains the pages of all apps)
#
# Output:
#   web/build/gf_apps/gf_solo/   (and the per-app dirs under web/build/)
#
# Host builds need these tools in PATH: tsc, sass, uglifyjs.
# Needs python3 with: colored delegator.py requests boto3 bs4 (ex. in a venv put first in PATH).
# --in-cont needs docker (run with sudo) and the glofloworg/gf_builder_web:latest image.

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLI="$ROOT_DIR/ops/cli__build.py"

IN_CONT=0; PAGE=""; APP="gf_solo"
while [ $# -gt 0 ]; do
    case "$1" in
        --in-cont)  IN_CONT=1 ;;
        --page)     shift; PAGE="${1:?--page needs a name}" ;;
        --app)      shift; APP="${1:?--app needs a name}" ;;
        -h|--help)  sed -n '2,15p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
        *) echo "unknown option: $1" >&2; exit 1 ;;
    esac
    shift
done

ARGS=(-run=build_web "-app=$APP")
[ "$IN_CONT" = 1 ] || ARGS+=(-build_outof_cont)
[ -z "$PAGE" ] || ARGS+=("-page_name=$PAGE")

if [ "$IN_CONT" = 0 ]; then
    for tool in tsc sass uglifyjs; do
        command -v "$tool" >/dev/null 2>&1 || { echo "ERROR!! - '$tool' not found in PATH (or use --in-cont)." >&2; exit 1; }
    done
fi

echo ">>> BUILD WEB $APP${PAGE:+ (page $PAGE)}"
python3 "$CLI" "${ARGS[@]}"
echo ">>> DONE"
