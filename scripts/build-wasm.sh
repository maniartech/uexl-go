#!/usr/bin/env bash
# build-wasm.sh — Build UExL browser WASM binary with TinyGo
# Usage: bash scripts/build-wasm.sh
# Run from the uexl-go directory.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
WASM_SRC="$REPO_ROOT/cmd/uexl-wasm"
PUBLIC_DIR="$REPO_ROOT/../uexl-playground/public"
OUT="$PUBLIC_DIR/uexl.wasm"

if ! command -v tinygo >/dev/null 2>&1; then
  echo "✗ tinygo was not found on PATH. Install TinyGo before running this script."
  exit 1
fi

TINYGOROOT_PATH="$(tinygo env TINYGOROOT)"
WASM_EXEC_SRC="$TINYGOROOT_PATH/targets/wasm_exec.js"

echo "▶ Building UExL WASM..."
echo "  Source : $WASM_SRC"
echo "  Output : $OUT"
echo "  Target : Browser"
echo "  Tool   : tinygo -target wasm -opt s -no-debug"

mkdir -p "$PUBLIC_DIR"

# Build the WASM binary
(cd "$WASM_SRC" && tinygo build -o "$OUT" -target wasm -opt s -no-debug .)

echo "  ✓ uexl.wasm built ($(du -sh "$OUT" | cut -f1))"

# Optional: run wasm-opt (Binaryen) for further ~15% size reduction.
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

# Copy wasm_exec.js from TinyGo
if [ -f "$WASM_EXEC_SRC" ]; then
  cp "$WASM_EXEC_SRC" "$PUBLIC_DIR/wasm_exec.js"
  echo "  ✓ wasm_exec.js copied from $TINYGOROOT_PATH"
else
  echo "  ✗ wasm_exec.js not found at $WASM_EXEC_SRC"
  exit 1
fi

echo "✅ WASM build complete."
