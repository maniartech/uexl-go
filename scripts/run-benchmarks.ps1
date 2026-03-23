# Run UExL benchmarks and save all artifacts to benchmarks/results/
# Usage: .\scripts\run-benchmarks.ps1 [-BenchFilter <pattern>] [-BenchTime <duration>] [-CpuProfile] [-MemProfile]
# Examples:
#   .\scripts\run-benchmarks.ps1
#   .\scripts\run-benchmarks.ps1 -BenchFilter "BenchmarkPipe" -BenchTime 10s
#   .\scripts\run-benchmarks.ps1 -CpuProfile -MemProfile

param(
    [string]$BenchFilter = ".",
    [string]$BenchTime   = "5s",
    [switch]$CpuProfile,
    [switch]$MemProfile
)

$ErrorActionPreference = "Stop"

$repoRoot   = Split-Path -Parent $PSScriptRoot
$benchDir   = Join-Path $repoRoot "benchmarks"
$resultsDir = Join-Path $benchDir "results"

# Ensure results directory exists
New-Item -ItemType Directory -Force -Path $resultsDir | Out-Null

$timestamp  = Get-Date -Format "yyyyMMdd_HHmmss"
$outputFile = Join-Path $resultsDir "bench_${timestamp}.txt"

Push-Location $benchDir
try {
    $args = @(
        "test",
        "-bench=$BenchFilter",
        "-benchmem",
        "-benchtime=$BenchTime",
        "-count=1"
    )

    if ($CpuProfile) {
        $args += "-cpuprofile=$(Join-Path $resultsDir "cpu_${timestamp}.prof")"
    }

    if ($MemProfile) {
        $args += "-memprofile=$(Join-Path $resultsDir "mem_${timestamp}.prof")"
    }

    Write-Host "Running: go $($args -join ' ')" -ForegroundColor Cyan
    & go @args | Tee-Object -FilePath $outputFile

    Write-Host ""
    Write-Host "Results saved to: $outputFile" -ForegroundColor Green
} finally {
    Pop-Location
}
