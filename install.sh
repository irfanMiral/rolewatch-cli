#!/bin/sh
# Install the rolewatch CLI: download the prebuilt binary for this OS/arch from the
# RoleWatch download server and place it on PATH.
#   curl -fsSL https://rolewatch.irfanmiral.com/install.sh | sh
# (rolewatch.irfanmiral.com serves this same script and the binaries under /dl/.)
set -eu

REPO="rolewatch-cli"
BIN="rolewatch"
INSTALL_DIR="${ROLEWATCH_INSTALL_DIR:-/usr/local/bin}"

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
  linux | darwin) ;;
  *) echo "rolewatch: unsupported OS '$os' (use 'go install $REPO/cmd/$BIN@latest')" >&2; exit 1 ;;
esac

arch=$(uname -m)
case "$arch" in
  x86_64 | amd64) arch=amd64 ;;
  arm64 | aarch64) arch=arm64 ;;
  *) echo "rolewatch: unsupported arch '$arch' (use 'go install $REPO/cmd/$BIN@latest')" >&2; exit 1 ;;
esac

asset="${BIN}_${os}_${arch}"
url="https://rolewatch.irfanmiral.com/dl/${asset}"

tmp=$(mktemp)
trap 'rm -f "$tmp"' EXIT
echo "Downloading ${asset} …"
curl -fsSL "$url" -o "$tmp"
chmod +x "$tmp"

target="${INSTALL_DIR}/${BIN}"
if [ -w "$INSTALL_DIR" ]; then
  mv "$tmp" "$target"
else
  echo "Installing to ${target} (requires sudo) …"
  sudo mv "$tmp" "$target"
fi
trap - EXIT

echo "Installed ${BIN} → ${target}"
echo "Next: ${BIN} auth login    (opens your browser to sign in)"
