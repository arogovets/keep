#!/bin/sh
set -eu

if ! command -v go >/dev/null 2>&1; then
  echo "Go 1.22 or later is required. Install Go, then rerun ./install-wsl.sh." >&2
  exit 1
fi

version=$(go env GOVERSION)
version=${version#go}
major=${version%%.*}
minor=${version#*.}
minor=${minor%%.*}
if [ "$major" -lt 1 ] || { [ "$major" -eq 1 ] && [ "$minor" -lt 22 ]; }; then
  echo "Go 1.22 or later is required; found $(go env GOVERSION)." >&2
  exit 1
fi

case "$(uname -m)" in
  x86_64) goarch=amd64 ;;
  aarch64|arm64) goarch=arm64 ;;
  *) echo "Unsupported WSL architecture: $(uname -m)" >&2; exit 1 ;;
esac

bindir=${BINDIR:-"$HOME/.local/bin"}
mkdir -p "$bindir"
GOOS=windows GOARCH="$goarch" go build -ldflags "-s -w -X main.version=dev" -o "$bindir/keep.exe" ./cmd/keep

cat > "$bindir/keep" <<'WRAPPER'
#!/bin/sh
script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
exec "$script_dir/keep.exe" "$@"
WRAPPER
chmod +x "$bindir/keep"

echo "Installed keep.exe and the WSL keep wrapper in $bindir"
case ":${PATH:-}:" in
  *:"$bindir":*) ;;
  *) echo "Add it to PATH with: export PATH=\"$bindir:\$PATH\"" ;;
esac
