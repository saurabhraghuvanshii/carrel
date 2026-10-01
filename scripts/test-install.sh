#!/bin/sh
# Tests install.sh against a fake release served on localhost. Needs python3.
set -eu

here="$(cd "$(dirname "$0")" && pwd)"
work="$(mktemp -d)"
pid=""
cleanup() { [ -n "$pid" ] && kill "$pid" 2>/dev/null; rm -rf "$work"; }
trap cleanup EXIT INT TERM

case "$(uname -s)" in Linux) os=linux ;; Darwin) os=darwin ;; *) echo "skip: unsupported OS"; exit 0 ;; esac
case "$(uname -m)" in x86_64|amd64) arch=amd64 ;; aarch64|arm64) arch=arm64 ;; *) echo "skip: unsupported arch"; exit 0 ;; esac
file="carrel_${os}_${arch}.tar.gz"

mkdir -p "$work/rel" "$work/pkg"
cat > "$work/pkg/carrel" <<'FAKE'
#!/bin/sh
echo "carrel doctor: checking the tools needed to run your code"
echo "  ok       java   fake"
echo "  missing  g++    install it and make sure it is on your PATH"
FAKE
chmod +x "$work/pkg/carrel"
tar -czf "$work/rel/$file" -C "$work/pkg" carrel
(cd "$work/rel" && if command -v sha256sum >/dev/null 2>&1; then sha256sum "$file"; else shasum -a 256 "$file"; fi) > "$work/rel/checksums.txt"

port="$(python3 -c 'import socket; s=socket.socket(); s.bind(("127.0.0.1",0)); print(s.getsockname()[1])')"
python3 -m http.server "$port" --bind 127.0.0.1 --directory "$work/rel" >/dev/null 2>&1 &
pid=$!
i=0
until python3 -c "import urllib.request; urllib.request.urlopen('http://127.0.0.1:$port/checksums.txt')" 2>/dev/null; do
  i=$((i + 1)); [ "$i" -lt 50 ] || { echo "FAIL: fake server did not start"; exit 1; }
  sleep 0.1
done

pass=0
ok() { pass=$((pass + 1)); echo "ok   $1"; }
bad() { echo "FAIL $1"; exit 1; }

run() {
  env HOME="$work/home" CARREL_REPO=test/carrel CARREL_BASE_URL="http://127.0.0.1:$port" \
    CARREL_INSTALL_DIR="$work/bin" "$@" sh "$here/install.sh"
}

out="$(run 2>&1)" || bad "success path: $out"
[ -x "$work/bin/carrel" ] || bad "binary not installed"
echo "$out" | grep -q "C++ is missing" || bad "missing compiler hint not shown"
echo "$out" | tail -n 1 | grep -q "Run: carrel" || bad "last line is not Run: carrel"
ok "installs, runs doctor, hints at missing g++"

rm -rf "$work/bin"
cp "$work/rel/checksums.txt" "$work/good.txt"
sed 's/^./0/' "$work/good.txt" > "$work/rel/checksums.txt"
if out="$(run 2>&1)"; then bad "wrong checksum should fail"; fi
echo "$out" | grep -q "checksum does not match" || bad "wrong checksum message: $out"
[ ! -e "$work/bin/carrel" ] || bad "binary installed despite bad checksum"
cp "$work/good.txt" "$work/rel/checksums.txt"
ok "wrong checksum stops the install"

mkdir -p "$work/nonet"
for t in sh uname mkdir mktemp rm id awk cut; do ln -sf "$(command -v "$t")" "$work/nonet/$t"; done
if out="$(run PATH="$work/nonet" 2>&1)"; then bad "no curl or wget should fail"; fi
echo "$out" | grep -q "needs curl or wget" || bad "no downloader message: $out"
ok "no curl or wget fails clearly"

mkdir -p "$work/fakeos"
printf '#!/bin/sh\necho Plan9\n' > "$work/fakeos/uname"
chmod +x "$work/fakeos/uname"
if out="$(run PATH="$work/fakeos:$PATH" 2>&1)"; then bad "unsupported OS should fail"; fi
echo "$out" | grep -q "install.ps1" || bad "unsupported OS message: $out"
ok "unsupported OS points to install.ps1"

if command -v shellcheck >/dev/null 2>&1; then
  shellcheck "$here/install.sh" "$here/test-install.sh" && ok "shellcheck clean"
fi
echo "$pass checks passed"
