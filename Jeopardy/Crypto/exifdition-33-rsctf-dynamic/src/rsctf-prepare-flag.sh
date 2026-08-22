#!/bin/sh
set -eu
temporary=${TMPDIR:-/tmp}/rsctf-exif.$$
exiftool -o "$temporary" -Comment="$RSCTF_FLAG" flag.png >/dev/null
cat "$temporary" > flag.png
rm -f "$temporary"
