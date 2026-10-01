#!/bin/sh
# Fills the formula and manifest from a release's checksums.txt.
# Usage: sh packaging/fill.sh 0.1.0 dist/checksums.txt out/
set -eu
[ $# -eq 3 ] || { echo "usage: fill.sh <version without v> <checksums.txt> <out dir>" >&2; exit 1; }
version=$1 sums=$2 out=$3
here="$(cd "$(dirname "$0")" && pwd)"
mkdir -p "$out"

sha() {
  v="$(awk -v f="carrel_$1" '$2 == f { print $1 }' "$sums")"
  [ -n "$v" ] || { echo "carrel_$1 is not in $sums" >&2; exit 1; }
  printf '%s' "$v"
}

script="s/{{version}}/$version/g"
for p in darwin_arm64 darwin_amd64 linux_arm64 linux_amd64; do
  script="$script;s/{{sha_$p}}/$(sha "$p.tar.gz")/g"
done
script="$script;s/{{sha_windows_amd64}}/$(sha windows_amd64.zip)/g"

sed "$script" "$here/carrel.rb" > "$out/carrel.rb"
sed "$script" "$here/carrel.json" > "$out/carrel.json"
if grep -q '{{' "$out/carrel.rb" "$out/carrel.json"; then echo "placeholder left" >&2; exit 1; fi
echo "wrote $out/carrel.rb and $out/carrel.json"
