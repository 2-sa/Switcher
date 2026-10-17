param([Parameter(Mandatory = $true)][string]$Version)

$ErrorActionPreference = 'Stop'
if ($env:GITHUB_ACTIONS -ne 'true') { throw 'Run this installation smoke test only on a disposable GitHub runner.' }
$root = Split-Path $PSScriptRoot -Parent
$installer = Join-Path $root 'bin/BearSwitcher-Setup.exe'
$registryKey = 'HKCU:/Software/Microsoft/Windows/CurrentVersion/Uninstall/BearSwitcher'
if (Test-Path -LiteralPath $registryKey) { throw 'Refusing to replace an existing installation in the test environment.' }
$testDir = Join-Path $env:RUNNER_TEMP ('BearSwitcher install smoke ' + [guid]::NewGuid().ToString('N'))
$dataDir = Join-Path $env:APPDATA 'TcNo Account Switcher'
$marker = Join-Path $dataDir ('installer-preserve-' + [guid]::NewGuid().ToString('N') + '.txt')
New-Item -ItemType Directory -Force -Path $dataDir | Out-Null
[IO.File]::WriteAllText($marker, 'existing account data must survive')

function Invoke-Setup([string]$File, [string]$Arguments, [int]$ExpectedExit = 0) {
  $process = Start-Process -FilePath $File -ArgumentList $Arguments -WindowStyle Hidden -Wait -PassThru
  if ($process.ExitCode -ne $ExpectedExit) { throw "Installer process exited $($process.ExitCode), expected ${ExpectedExit}: $File" }
}

function Assert-Installed {
  $exe = Join-Path $testDir 'TcNo-Acc-Switcher.exe'
  if ((Get-FileHash -LiteralPath $exe).Hash -ne (Get-FileHash -LiteralPath (Join-Path $root 'bin/TcNo-Acc-Switcher.exe')).Hash) {
    throw 'Installed executable differs from the verified update payload.'
  }
  $entry = Get-ItemProperty -LiteralPath $registryKey
  if ($entry.DisplayName -ne 'BearSwitcher' -or $entry.DisplayVersion -ne $Version -or $entry.InstallLocation -ne $testDir) {
    throw 'Windows uninstall registration is incorrect.'
  }
  $shortcut = Join-Path ([Environment]::GetFolderPath('Programs')) 'BearSwitcher/BearSwitcher.lnk'
  $shell = New-Object -ComObject WScript.Shell
  if ($shell.CreateShortcut($shortcut).TargetPath -ne $exe) { throw 'Start menu shortcut has the wrong target.' }
  if ([IO.File]::ReadAllText($marker) -ne 'existing account data must survive') { throw 'Existing user data changed.' }
}

try {
  # NSIS requires /D last, without quotes, even when the directory contains spaces.
  Invoke-Setup $installer "/S /D=$testDir"
  Assert-Installed
  $created = $false
  $runningApp = [Threading.Mutex]::new($false, 'TcNo-Acc-Switcher-Singleton', [ref]$created)
  try {
    if (-not $created) { throw 'An unexpected application instance is already running on the test runner.' }
    Invoke-Setup $installer "/S /D=$testDir" 2
  } finally {
    $runningApp.Dispose()
  }
  $customFile = Join-Path $testDir 'keep-custom-file.txt'
  [IO.File]::WriteAllText($customFile, 'not owned by the installer')
  Invoke-Setup $installer "/S /D=$testDir"
  Assert-Installed
  if ([IO.File]::ReadAllText($customFile) -ne 'not owned by the installer') { throw 'Upgrade removed an unowned file.' }
  Invoke-Setup (Join-Path $testDir 'Uninstall.exe') '/S'
  if (Test-Path -LiteralPath (Join-Path $testDir 'TcNo-Acc-Switcher.exe')) { throw 'Uninstall left the application installed.' }
  if (Test-Path -LiteralPath $registryKey) { throw 'Uninstall left its Windows registration.' }
  if ([IO.File]::ReadAllText($customFile) -ne 'not owned by the installer') { throw 'Uninstall deleted an unowned file.' }
  if ([IO.File]::ReadAllText($marker) -ne 'existing account data must survive') { throw 'Uninstall deleted existing user data.' }
  Remove-Item -LiteralPath $customFile
  Remove-Item -LiteralPath $testDir
  Write-Host 'PASS: install, running-app refusal, upgrade, uninstall, executable hash, shortcuts, registry, and data preservation.'
} finally {
  Remove-Item -LiteralPath $marker -ErrorAction SilentlyContinue
}
