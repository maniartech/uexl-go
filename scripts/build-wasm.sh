#!/usr/bin/env bash
# build-wasm.sh — Build UExL WASM binary and copy wasm_exec.js
# Usage: bash scripts/build-wasm.sh
# Run from the uexl-go directory.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
WASM_SRC="$REPO_ROOT/../uexl-playground/wasm"
PUBLIC_DIR="$REPO_ROOT/../uexl-playground/public"
OUT="$PUBLIC_DIR/uexl.wasm"
GOROOT_PATH="$(go env GOROOT)"
# Locate wasm_exec.js — path changed in Go 1.21+: lib/wasm/ (was misc/wasm/)
WASM_EXEC_SRC="$GOROOT_PATH/lib/wasm/wasm_exec.js"
if [ ! -f "$WASM_EXEC_SRC" ]; then
  WASM_EXEC_SRC="$GOROOT_PATH/misc/wasm/wasm_exec.js"
fi

echo "▶ Building UExL WASM..."
echo "  Source : $WASM_SRC"
echo "  Output : $OUT"

mkdir -p "$PUBLIC_DIR"

# Build the WASM binary
# -s -w  : strip debug symbols + DWARF (~30% size reduction)
# -trimpath: remove local file paths from binary
(cd "$WASM_SRC" && GOOS=js GOARCH=wasm go build -trimpath -ldflags "-s -w" -o "$OUT" .)

echo "  ✓ uexl.wasm built ($(du -sh "$OUT" | cut -f1))"

# Optional: run wasm-opt (Binaryen) for further ~15% size reduction.
# Feature flags required for Go 1.21+ generated WASM.
# Install: https://github.com/WebAssembly/binaryen/releases
if command -v wasm-opt &>/dev/null; then
  wasm-opt -Oz \
    --enable-bulk-memory \
    --enable-nontrapping-float-to-int \
    --enable-sign-ext \
    --enable-mutable-globals \
    "$OUT" -o "$OUT.tmp" && mv "$OUT.tmp" "$OUT"
  echo "  ✓ wasm-opt applied ($(du -sh "$OUT" | cut -f1) after optimisation)"
else
  echo "  INFO: wasm-opt not found, skipping (install Binaryen for extra ~15% reduction)"
fi

# Copy wasm_exec.js from Go SDK
if [ -f "$WASM_EXEC_SRC" ]; then
  cp "$WASM_EXEC_SRC" "$PUBLIC_DIR/wasm_exec.js"
  echo "  ✓ wasm_exec.js copied from $GOROOT_PATH"
else
  echo "  ✗ wasm_exec.js not found at $WASM_EXEC_SRC"
  echo "    Try: cp \"\$(go env GOROOT)/misc/wasm/wasm_exec.js\" ../uexl-playground/public/"
  exit 1
fi

echo "✅ WASM build complete."
