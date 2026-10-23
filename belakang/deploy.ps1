<#
.SYNOPSIS
    Build project Website Sekolah (Vue + Go embed) untuk deploy ke VPS.

.DESKRIPSI
    Menjalankan urutan:
    1. Build frontend Vue (npm run build)
    2. Copy hasil dist/ ke folder Go (cmd/server/frontend/dist)
    3. Cross-compile binary Go untuk Linux (amd64)
    4. (Opsional) Upload otomatis ke VPS via scp, kalau parameter -Upload dipakai

.CARA PAKAI
    Build saja (tanpa upload), dari folder mana saja:
        .\deploy.ps1

    Build + upload otomatis ke VPS (butuh SSH key sudah diset):
        .\deploy.ps1 -Upload

.CATATAN
    - Path TIDAK lagi hardcoded. Script menghitung sendiri lokasi project
      dari letak file ini, jadi aman walau folder dipindah / di-rename.
    - Host VPS default "smpn254.sch.id". Bisa dioverride tanpa mengedit file:
        $env:SEKOLAH_VPS_HOST = "1.2.3.4"
        .\deploy.ps1 -Upload
    - Fitur -Upload butuh OpenSSH client (scp/ssh) dan SSH key yang sudah
      terpasang, supaya tidak diminta password berkali-kali.
      Kalau belum, jangan pakai -Upload: upload manual lewat File Manager aaPanel.
#>

param(
    [switch]$Upload
)

# ========== KONFIGURASI ==========
# Path dihitung dari lokasi script (folder belakang), bukan hardcoded.
$BackendDir   = $PSScriptRoot
$ProjectRoot  = Split-Path -Parent $BackendDir
$FrontendDir  = Join-Path $ProjectRoot "depan"
$BackendMain  = Join-Path $BackendDir "cmd\server"
$EmbedDistDir = Join-Path $BackendMain "frontend\dist"
$BinaryName   = "web254"

# Konfigurasi VPS (hanya dipakai kalau -Upload dipakai)
$VpsUser      = "root"
$VpsHost      = if ($env:SEKOLAH_VPS_HOST) { $env:SEKOLAH_VPS_HOST } else { "smpn254.sch.id" }
$VpsPath      = "/www/wwwroot/smpn254.sch.id"
$SystemdName  = "web254"

# ========== JANGAN DIUBAH DI BAWAH INI ==========
$ErrorActionPreference = "Stop"

function Write-Step($msg) {
    Write-Host ""
    Write-Host "==> $msg" -ForegroundColor Cyan
}

function Fail($msg) {
    Write-Host ""
    Write-Host "GAGAL: $msg" -ForegroundColor Red
    exit 1
}

function Assert-Command($name) {
    if (-not (Get-Command $name -ErrorAction SilentlyContinue)) {
        Fail "$name tidak ditemukan di PATH. Pastikan sudah terpasang dan terminal di-restart."
    }
}

try {
    # 0. Validasi awal
    if (-not (Test-Path $FrontendDir)) { Fail "Folder frontend tidak ditemukan: $FrontendDir" }
    if (-not (Test-Path $BackendMain)) { Fail "Folder backend tidak ditemukan: $BackendMain" }
    Assert-Command "npm"
    Assert-Command "go"

    # 1. Build frontend
    Write-Step "Build frontend Vue ($FrontendDir)..."
    Push-Location $FrontendDir
    try {
        npm run build
        if ($LASTEXITCODE -ne 0) { Fail "npm run build gagal (exit code $LASTEXITCODE)" }
    }
    finally {
        Pop-Location
    }

    if (-not (Test-Path (Join-Path $FrontendDir "dist\index.html"))) {
        Fail "Folder dist/ tidak ditemukan setelah build. Cek error di atas."
    }

    # 2. Copy dist ke folder Go
    Write-Step "Menyalin dist/ ke folder Go untuk go:embed..."
    if (Test-Path $EmbedDistDir) {
        Remove-Item $EmbedDistDir -Recurse -Force
    }
    Copy-Item (Join-Path $FrontendDir "dist") $EmbedDistDir -Recurse

    # 3. Build binary Linux
    Write-Step "Compile binary Go untuk Linux (amd64)..."
    $oldGOOS = $env:GOOS
    $oldGOARCH = $env:GOARCH
    Push-Location $BackendDir
    try {
        $env:GOOS = "linux"
        $env:GOARCH = "amd64"
        go build -ldflags="-s -w" -o $BinaryName ./cmd/server
        $buildExitCode = $LASTEXITCODE
    }
    finally {
        $env:GOOS = $oldGOOS
        $env:GOARCH = $oldGOARCH
        Pop-Location
    }

    if ($buildExitCode -ne 0) { Fail "go build gagal (exit code $buildExitCode)" }

    $binaryPath = Join-Path $BackendDir $BinaryName
    if (-not (Test-Path $binaryPath)) { Fail "File binary '$BinaryName' tidak ditemukan setelah build." }

    $size = (Get-Item $binaryPath).Length / 1MB
    Write-Host "Binary berhasil dibuat: $binaryPath ($([math]::Round($size,1)) MB)" -ForegroundColor Green

    # 4. Upload (opsional)
    if ($Upload) {
        if ([string]::IsNullOrWhiteSpace($VpsHost) -or $VpsHost -eq "ISI_IP_ATAU_DOMAIN_VPS") {
            Fail "Host VPS belum diisi. Set `$env:SEKOLAH_VPS_HOST atau edit `$VpsHost di bagian KONFIGURASI."
        }
        Assert-Command "scp"
        Assert-Command "ssh"

        Write-Step "Upload binary ke VPS ($VpsUser@$VpsHost)..."
        scp $binaryPath "${VpsUser}@${VpsHost}:${VpsPath}/${BinaryName}"
        if ($LASTEXITCODE -ne 0) { Fail "scp gagal (exit code $LASTEXITCODE). Cek koneksi SSH / SSH key." }

        Write-Step "Set permission & restart service di VPS..."
        ssh "${VpsUser}@${VpsHost}" "chmod +x ${VpsPath}/${BinaryName} && systemctl restart ${SystemdName} && systemctl status ${SystemdName} --no-pager -l"
        if ($LASTEXITCODE -ne 0) { Fail "Perintah remote gagal (exit code $LASTEXITCODE)." }

        Write-Host ""
        Write-Host "Deploy selesai. Cek website untuk memastikan semuanya jalan normal." -ForegroundColor Green
    }
    else {
        Write-Host ""
        Write-Host "Build selesai. Binary siap di: $binaryPath" -ForegroundColor Green
        Write-Host "Upload manual via File Manager aaPanel ke: $VpsPath" -ForegroundColor Yellow
        Write-Host "Jangan lupa: set permission execute (755) + 'systemctl restart $SystemdName' setelah upload." -ForegroundColor Yellow
    }
}
catch {
    Write-Host ""
    Write-Host "ERROR: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}
