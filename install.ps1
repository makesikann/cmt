$ErrorActionPreference = "Stop"

$repo = "makesikann/cmt"
$apiUrl = "https://api.github.com/repos/$repo/releases/latest"

Write-Host "Son sürüm denetleniyor..."
try {
    $release = Invoke-RestMethod -Uri $apiUrl
} catch {
    Write-Error "GitHub API erişimi başarısız. Sürüm bulunamıyor."
    exit 1
}

$asset = $release.assets | Where-Object { $_.name -match "windows" -and $_.name -match "x86_64" }

if (!$asset) {
    Write-Error "Uygun Windows sürümü bulunamadı."
    exit 1
}

$downloadUrl = $asset.browser_download_url
$destPath = "$env:TEMP\cmt.zip"
$extractPath = "$env:TEMP\cmt_extract"

Write-Host "İndiriliyor: $downloadUrl..."
Invoke-WebRequest -Uri $downloadUrl -OutFile $destPath

Write-Host "Arşiv çıkartılıyor..."
if (Test-Path $extractPath) { Remove-Item -Path $extractPath -Recurse -Force }
New-Item -ItemType Directory -Path $extractPath | Out-Null
Expand-Archive -Path $destPath -DestinationPath $extractPath -Force

$binDir = "$env:USERPROFILE\bin"
if (!(Test-Path $binDir)) {
    New-Item -ItemType Directory -Force -Path $binDir | Out-Null
    Write-Host "PATH ayarlarınızda bulunabilmesi için $binDir dizinini Ortam Değişkenlerine eklemeyi unutmayın."
}

$exePath = Get-ChildItem -Path $extractPath -Filter "cmt.exe" -Recurse | Select-Object -First 1 -ExpandProperty FullName
if ($exePath) {
    Copy-Item -Path $exePath -Destination "$binDir\cmt.exe" -Force
    Write-Host "✅ Kurulum tamamlandı! cmt komutu başarıyla şu konuma yüklendi: $binDir\cmt.exe"
} else {
    Write-Error "Arşiv içinde cmt.exe bulunamadı."
}

# Temizlik
Remove-Item -Path $destPath -Force
Remove-Item -Path $extractPath -Recurse -Force
