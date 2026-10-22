# thin launcher: prefer built rabbitmq.exe beside this script
$ErrorActionPreference = "Stop"
$here = Split-Path -Parent $MyInvocation.MyCommand.Path
$exe = Join-Path $here "rabbitmq.exe"
if (-not (Test-Path $exe)) {
  $env:PATH = "D:\tools\go\bin;" + $env:PATH
  $env:GOPROXY = "goproxy.cn,direct"
  Push-Location $here
  try { go build -o rabbitmq.exe . } finally { Pop-Location }
}
& $exe @args
exit $LASTEXITCODE
