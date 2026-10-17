param(
  [Parameter(Mandatory = $true)][string]$Version,
  [Parameter(Mandatory = $true)][string]$FileVersion
)

$ErrorActionPreference = 'Stop'
$root = Split-Path $PSScriptRoot -Parent
$nsisDir = Join-Path $root 'build/windows/nsis'
$compiler = Join-Path ${env:ProgramFiles(x86)} 'NSIS/makensis.exe'
if (-not (Test-Path -LiteralPath $compiler)) { throw 'NSIS compiler is not installed on the build runner.' }

wails3 generate webview2bootstrapper -dir $nsisDir
if ($LASTEXITCODE -ne 0) { throw 'WebView2 bootstrapper download failed.' }
$bootstrapper = Join-Path $nsisDir 'MicrosoftEdgeWebview2Setup.exe'
$signature = Get-AuthenticodeSignature -LiteralPath $bootstrapper
if ($signature.Status -ne 'Valid' -or $signature.SignerCertificate.Subject -notmatch 'O=Microsoft Corporation') {
  throw 'WebView2 bootstrapper is not signed by Microsoft with a valid certificate.'
}

Push-Location $nsisDir
try {
  & $compiler /INPUTCHARSET UTF8 "/DAPP_VERSION=$Version" "/DFILE_VERSION=$FileVersion" fork.nsi
  if ($LASTEXITCODE -ne 0) { throw 'BearSwitcher installer compilation failed.' }
} finally {
  Pop-Location
}
$installer = Join-Path $root 'bin/BearSwitcher-Setup.exe'
if (-not (Test-Path -LiteralPath $installer)) { throw 'Installer output is missing.' }
Write-Host "Installer ready: $installer"
