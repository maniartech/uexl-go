# build-wasm.ps1 — Build UExL browser WASM binary with TinyGo
# Usage: powershell -ExecutionPolicy Bypass -File scripts\build-wasm.ps1
# Run from the uexl-go directory.

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$ScriptDir  = Split-Path -Parent $MyInvocation.MyCommand.Path
$RepoRoot   = Split-Path -Parent $ScriptDir
$WasmSrc    = Join-Path $RepoRoot "cmd\uexl-wasm"
$PublicDir  = Join-Path $RepoRoot "..\uexl-playground\public"
$OutFile    = Join-Path $PublicDir "uexl.wasm"

if (-not (Get-Command tinygo -ErrorAction SilentlyContinue)) {
    throw "tinygo was not found on PATH. Install TinyGo before running this script."
}

$TinyGoVersion = (tinygo version 2>$null)
if (-not $TinyGoVersion) {
    throw "Unable to determine TinyGo version."
}

$TinyGoTarget = "wasm"
$TinyGoOptLevel = "s"

# TinyGo ships the correct browser runtime helper for its wasm target.
$TinyGoRoot = (tinygo env TINYGOROOT).Trim()
$WasmExecSrc = Join-Path $TinyGoRoot "targets\wasm_exec.js"

Write-Host "▶ Building UExL WASM..."
Write-Host "  Source : $WasmSrc"
Write-Host "  Output : $OutFile"
Write-Host "  Target : Browser"
Write-Host "  Tool   : tinygo -target $TinyGoTarget -opt $TinyGoOptLevel -no-debug"

# Ensure output directory exists
New-Item -ItemType Directory -Force -Path $PublicDir | Out-Null

# Build the WASM binary with TinyGo.
Push-Location $WasmSrc
try {
    $TinyGoArgs = @(
        "build"
        "-o", $OutFile
        "-target", $TinyGoTarget
        "-opt", $TinyGoOptLevel
        "-no-debug"
        "."
    )
    & tinygo @TinyGoArgs
    if ($LASTEXITCODE -ne 0) { throw "tinygo build failed with exit code $LASTEXITCODE" }
} finally {
    Pop-Location
}

$size = [math]::Round((Get-Item $OutFile).Length / 1MB, 2)
Write-Host "  OK uexl.wasm built ($size MB)"

# Optional: run wasm-opt (Binaryen) for further size reduction.
# Install: https://github.com/WebAssembly/binaryen/releases
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

if (Test-Path $WasmExecSrc) {
    Copy-Item $WasmExecSrc (Join-Path $PublicDir "wasm_exec.js") -Force
    Write-Host "  OK wasm_exec.js copied from $TinyGoRoot"
} else {
    Write-Host "  ERROR: wasm_exec.js not found at $WasmExecSrc"
    exit 1
}

Write-Host "WASM build complete."
