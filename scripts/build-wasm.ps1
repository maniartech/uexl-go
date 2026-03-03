# build-wasm.ps1 — Build UExL WASM binary and copy wasm_exec.js
# Usage: powershell -ExecutionPolicy Bypass -File scripts\build-wasm.ps1
# Run from the uexl-go directory.

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$ScriptDir  = Split-Path -Parent $MyInvocation.MyCommand.Path
$RepoRoot   = Split-Path -Parent $ScriptDir
$WasmSrc    = Join-Path $RepoRoot "..\uexl-playground\wasm"
$PublicDir  = Join-Path $RepoRoot "..\uexl-playground\public"
$OutFile    = Join-Path $PublicDir "uexl.wasm"

# Resolve GOROOT from go env
$GoRoot = (go env GOROOT).Trim()
# Locate wasm_exec.js — path changed in Go 1.21+: lib/wasm/ (was misc/wasm/)
$WasmExecSrc = Join-Path $GoRoot "lib\wasm\wasm_exec.js"
if (-not (Test-Path $WasmExecSrc)) {
    $WasmExecSrc = Join-Path $GoRoot "misc\wasm\wasm_exec.js"
}

Write-Host "▶ Building UExL WASM..."
Write-Host "  Source : $WasmSrc"
Write-Host "  Output : $OutFile"

# Ensure output directory exists
New-Item -ItemType Directory -Force -Path $PublicDir | Out-Null

# Build the WASM binary
# -s -w  : strip debug symbols + DWARF (~30% size reduction)
# -trimpath: remove local file paths from binary
Push-Location $WasmSrc
try {
    $env:GOOS    = "js"
    $env:GOARCH  = "wasm"
    go build -trimpath -ldflags "-s -w" -o $OutFile .
    if ($LASTEXITCODE -ne 0) { throw "go build failed with exit code $LASTEXITCODE" }
} finally {
    # Always restore env vars
    Remove-Item Env:\GOOS   -ErrorAction SilentlyContinue
    Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue
    Pop-Location
}

$size = [math]::Round((Get-Item $OutFile).Length / 1MB, 2)
Write-Host "  OK uexl.wasm built ($size MB)"

# Optional: run wasm-opt (Binaryen) for further ~15% size reduction.
# Install: https://github.com/WebAssembly/binaryen/releases
# Feature flags required for Go 1.21+ generated WASM.
if (Get-Command wasm-opt -ErrorAction SilentlyContinue) {
    $tmpFile = $OutFile + ".tmp"
    wasm-opt -Oz `
        --enable-bulk-memory `
        --enable-nontrapping-float-to-int `
        --enable-sign-ext `
        --enable-mutable-globals `
        $OutFile -o $tmpFile
    if ($LASTEXITCODE -eq 0) {
        Move-Item $tmpFile $OutFile -Force
        $optSize = [math]::Round((Get-Item $OutFile).Length / 1MB, 2)
        Write-Host "  OK wasm-opt applied ($optSize MB after optimisation)"
    } else {
        Remove-Item $tmpFile -ErrorAction SilentlyContinue
        Write-Host "  WARN wasm-opt failed, keeping original"
    }
} else {
    Write-Host "  INFO wasm-opt not found, skipping (install Binaryen for extra ~15% reduction)"
}

# Copy wasm_exec.js from Go SDK
if (Test-Path $WasmExecSrc) {
    Copy-Item $WasmExecSrc (Join-Path $PublicDir "wasm_exec.js") -Force
    Write-Host "  OK wasm_exec.js copied from $GoRoot"
} else {
    Write-Host "  ERROR: wasm_exec.js not found at $WasmExecSrc"
    exit 1
}

Write-Host "WASM build complete."
