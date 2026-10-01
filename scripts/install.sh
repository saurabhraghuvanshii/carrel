#!/bin/sh
# Installs the carrel binary from the latest GitHub release into ~/.local/bin.
set -eu

REPO_DEFAULT="{{repo_slug}}"

say() { printf '%s\n' "$*"; }
fail() { printf 'carrel install: %s\n' "$*" >&2; exit 1; }

fetch() {
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL -o "$2" "$1"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "$2" "$1"
  else
    fail "needs curl or wget to download the release."
  fi
}

sha256() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | cut -d ' ' -f 1
  elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$1" | cut -d ' ' -f 1
  else
    fail "needs sha256sum or shasum to check the download. Nothing was installed."
  fi
}

main() {
  REPO="${CARREL_REPO:-$REPO_DEFAULT}"
  VERSION="${CARREL_VERSION:-latest}"
  DIR="${CARREL_INSTALL_DIR:-$HOME/.local/bin}"

  case "$REPO" in
    *"["*|*"{"*) fail "no repository set. Set CARREL_REPO=owner/carrel." ;;
  esac

  case "$(uname -s)" in
    Linux) os=linux ;;
    Darwin) os=darwin ;;
    *) fail "this script is for Mac and Linux. On Windows use install.ps1, or download a release by hand from https://github.com/$REPO/releases." ;;
  esac
  case "$(uname -m)" in
    x86_64|amd64) arch=amd64 ;;
    aarch64|arm64) arch=arm64 ;;
    *) fail "unsupported processor: $(uname -m). Download a release by hand from https://github.com/$REPO/releases." ;;
  esac

  if [ "$VERSION" = latest ]; then
    base="https://github.com/$REPO/releases/latest/download"
  else
    base="https://github.com/$REPO/releases/download/$VERSION"
  fi
  base="${CARREL_BASE_URL:-$base}"
  file="carrel_${os}_${arch}.tar.gz"

  if [ "$(id -u)" = 0 ]; then
    say "Note: running as root. Carrel does not need admin rights."
  fi

  tmp="$(mktemp -d)"
  trap 'rm -rf "$tmp"' EXIT INT TERM

  say "Downloading $file"
  fetch "$base/$file" "$tmp/$file" || fail "could not download $base/$file"
  fetch "$base/checksums.txt" "$tmp/checksums.txt" || fail "could not download $base/checksums.txt"

  want="$(awk -v f="$file" '$2 == f || $2 == "*" f { print $1 }' "$tmp/checksums.txt")"
  [ -n "$want" ] || fail "$file is not listed in checksums.txt. Nothing was installed."
  got="$(sha256 "$tmp/$file")"
  [ "$want" = "$got" ] || fail "checksum does not match for $file. Nothing was installed."

  tar -xzf "$tmp/$file" -C "$tmp" carrel || fail "could not unpack $file"
  mkdir -p "$DIR"
  install -m 755 "$tmp/carrel" "$DIR/carrel"
  say "Installed $DIR/carrel"

  case ":$PATH:" in
    *":$DIR:"*) ;;
    *) say "" ; say "$DIR is not on your PATH. Add this line to your shell profile:" ; say "  export PATH=\"$DIR:\$PATH\"" ;;
  esac

  say ""
  report="$("$DIR/carrel" doctor 2>&1 || true)"
  say "$report"
  if printf '%s\n' "$report" | grep -i 'java' | grep -qi 'not found\|missing'; then
    say "Java is missing: install a JDK (version 17 or newer) to practise in Java."
  fi
  if printf '%s\n' "$report" | grep -i 'g++' | grep -qi 'not found\|missing'; then
    say "C++ is missing: install g++ (on a Mac, run xcode-select --install) to practise in C++."
  fi
  say ""
  say "Run: carrel"
}

main "$@"
