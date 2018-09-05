#!/bin/bash
set -euo pipefail

SCRIPT=$(readlink -f "$0")
DIR=$(dirname "$SCRIPT")

GO="go"
GOFMT="gofmt"

cd "${DIR}/_c"
"$GO" tool cgo -godefs "ctype.go" \
    | awk -f "ctype.awk" \
    | "$GOFMT" > "../fuse/ctype.go"
