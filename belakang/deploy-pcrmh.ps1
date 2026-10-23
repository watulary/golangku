<#
.SYNOPSIS
    Alias lama untuk deploy.ps1 (PC rumah).

.DESKRIPSI
    Dulu file ini berisi salinan script yang sama dengan path hardcoded
    (D:\sekolah\...). Sekarang deploy.ps1 sudah otomatis mendeteksi lokasi
    project, jadi file ini cukup meneruskan parameter ke deploy.ps1.
    Dipertahankan supaya perintah lama tetap jalan.

.CARA PAKAI
    .\deploy-pcrmh.ps1
    .\deploy-pcrmh.ps1 -Upload
#>

param(
    [switch]$Upload
)

& (Join-Path $PSScriptRoot "deploy.ps1") -Upload:$Upload
exit $LASTEXITCODE
