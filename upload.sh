#!/usr/bin/env bash
set -euo pipefail

# Upload all EPUB files in ./output to /Books on crosspoint.local.

UPLOAD_URL="http://192.168.1.62/upload?path=/"

shopt -s nullglob
files=(./output/*.epub)

if (( ${#files[@]} == 0 )); then
  echo "No EPUB files found in ./output" >&2
  exit 0
fi

for file in "${files[@]}"; do
  echo "Uploading ${file} -> ${UPLOAD_URL}"
  curl --fail --show-error --location \
    -X POST \
    -F "file=@${file}" \
    "${UPLOAD_URL}"
  echo
  echo "Uploaded $(basename "$file")"
done
