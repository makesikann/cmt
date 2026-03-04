$ErrorActionPreference = "Stop"

$repo = "makesikann/cmt"
$apiUrl = "https://api.github.com/repos/$repo/releases/latest"

Write-Host "Checking for the latest version..."
try {
    $release = Invoke-RestMethod -Uri $apiUrl
} catch {
    Write-Error "GitHub API access failed. Version not found."
    exit 1
}

$asset = $release.assets | Where-Object { $_.name -match "windows" -and $_.name -match "x86_64" }

if (!$asset) {
    Write-Error "No suitable Windows version found."
    exit 1
}

$downloadUrl = $asset.browser_download_url
$destPath = "$env:TEMP\cmt.zip"
$extractPath = "$env:TEMP\cmt_extract"

Write-Host "Downloading: $downloadUrl..."
Invoke-WebRequest -Uri $downloadUrl -OutFile $destPath

Write-Host "Extracting archive..."
if (Test-Path $extractPath) { Remove-Item -Path $extractPath -Recurse -Force }
New-Item -ItemType Directory -Path $extractPath | Out-Null
Expand-Archive -Path $destPath -DestinationPath $extractPath -Force

$binDir = "$env:USERPROFILE\bin"
if (!(Test-Path $binDir)) {
    New-Item -ItemType Directory -Force -Path $binDir | Out-Null
    Write-Host "Please add $binDir to your PATH environment variable to use the cmt command."
}

$exePath = Get-ChildItem -Path $extractPath -Filter "cmt.exe" -Recurse | Select-Object -First 1 -ExpandProperty FullName
if ($exePath) {
    Copy-Item -Path $exePath -Destination "$binDir\cmt.exe" -Force
    Write-Host "Installation completed! cmt has been successfully installed to: $binDir\cmt.exe"
} else {
    Write-Error "cmt.exe not found in the archive."
}

# Cleanup
Remove-Item -Path $destPath -Force
Remove-Item -Path $extractPath -Recurse -Force
