# Untested on a real Windows machine.
# Installs carrel.exe from the latest GitHub release into %LOCALAPPDATA%\carrel\bin.
$ErrorActionPreference = 'Stop'

function Install-Carrel {
  $repo = if ($env:CARREL_REPO) { $env:CARREL_REPO } else { '{{repo_slug}}' }
  $version = if ($env:CARREL_VERSION) { $env:CARREL_VERSION } else { 'latest' }
  if ($repo -match '[\[{]') { throw 'No repository set. Set CARREL_REPO=owner/carrel.' }

  switch ($env:PROCESSOR_ARCHITECTURE) {
    'AMD64' { $arch = 'amd64' }
    'ARM64' { $arch = 'arm64' }
    default { throw "Unsupported processor: $($env:PROCESSOR_ARCHITECTURE). Download a release by hand from https://github.com/$repo/releases." }
  }

  if ($version -eq 'latest') { $base = "https://github.com/$repo/releases/latest/download" }
  else { $base = "https://github.com/$repo/releases/download/$version" }
  if ($env:CARREL_BASE_URL) { $base = $env:CARREL_BASE_URL }

  $file = "carrel_windows_$arch.zip"
  $tmp = Join-Path ([System.IO.Path]::GetTempPath()) ("carrel-" + [System.Guid]::NewGuid())
  New-Item -ItemType Directory -Path $tmp | Out-Null
  try {
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    Write-Host "Downloading $file"
    Invoke-WebRequest -UseBasicParsing -Uri "$base/$file" -OutFile (Join-Path $tmp $file)
    Invoke-WebRequest -UseBasicParsing -Uri "$base/checksums.txt" -OutFile (Join-Path $tmp 'checksums.txt')

    $want = $null
    foreach ($line in Get-Content (Join-Path $tmp 'checksums.txt')) {
      $parts = $line -split '\s+', 2
      if ($parts.Count -eq 2 -and $parts[1].TrimStart('*') -eq $file) { $want = $parts[0].ToLower() }
    }
    if (-not $want) { throw "$file is not listed in checksums.txt. Nothing was installed." }
    $got = (Get-FileHash -Algorithm SHA256 (Join-Path $tmp $file)).Hash.ToLower()
    if ($want -ne $got) { throw "Checksum does not match for $file. Nothing was installed." }

    $dir = Join-Path $env:LOCALAPPDATA 'carrel\bin'
    New-Item -ItemType Directory -Force -Path $dir | Out-Null
    Expand-Archive -Force -Path (Join-Path $tmp $file) -DestinationPath $tmp
    Copy-Item -Force (Join-Path $tmp 'carrel.exe') (Join-Path $dir 'carrel.exe')
    Write-Host "Installed $dir\carrel.exe"

    $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
    if (-not $userPath) { $userPath = '' }
    if (($userPath -split ';') -notcontains $dir) {
      [Environment]::SetEnvironmentVariable('Path', ($userPath.TrimEnd(';') + ';' + $dir).TrimStart(';'), 'User')
      Write-Host "Added $dir to your PATH. Open a new terminal to use it."
    }
    $env:Path = "$dir;$env:Path"

    Write-Host ''
    & (Join-Path $dir 'carrel.exe') doctor
    Write-Host ''
    Write-Host 'Run: carrel'
  } finally {
    Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
  }
}

Install-Carrel
