# Website Sekolah (Go + Fiber + Vue + Vite + Tailwind + MySQL)

Aplikasi website sekolah: halaman publik (berita, pengumuman, galeri, profil,
PPDB/SPMB) dan dashboard admin (data peserta didik, pegawai, nilai, absensi,
buku induk, kelulusan, dll).

## Stack

| Bagian | Teknologi |
|---|---|
| Backend | Go 1.26, Fiber v2, sqlx, JWT, bcrypt |
| Frontend | Vue 3, Vite 8, Tailwind CSS 4, Pinia, Vue Router |
| Database | MySQL 8 (Laragon) |
| Deploy | Frontend di-embed ke binary Go (`go:embed`), binary Linux untuk VPS |

## Struktur

```
Default Project/
├── belakang/          # backend Go (module "sekolah")
│   ├── cmd/server/    # main.go + frontend/dist (hasil embed)
│   ├── internal/      # auth, middleware, modules, routes
│   ├── storage/       # file upload (TIDAK di-commit)
│   ├── .env           # konfigurasi lokal (TIDAK di-commit)
│   └── deploy.ps1     # build + copy dist + binary Linux
├── depan/             # frontend Vue + Vite
│   └── src/           # pages, components, layouts, stores
├── sql tambahan/      # migrasi SQL per modul
└── Changelog.txt      # catatan perubahan
```

## Prasyarat

- Go 1.21+ (go.mod memakai 1.26; toolchain diunduh otomatis)
- Node.js 20+ dan npm
- MySQL 8 / Laragon (port 3306)

## Setup pertama

1. Salin konfigurasi backend:

   ```powershell
   Copy-Item belakang\.env.example belakang\.env
   ```

   Lalu isi `DB_PASS` dan `JWT_SECRET` di `belakang\.env`.

2. Siapkan database (pilih salah satu):

   - Import dump: `mysql -u root -p smpn254opencode < sql_smpn254_sch_id.sql`
   - Atau buat database kosong lalu jalankan file di `sql tambahan/` sesuai kebutuhan.

3. Install dependency frontend:

   ```powershell
   cd depan
   npm install
   ```

## Menjalankan saat pengembangan

Terminal 1 — backend (http://localhost:3000):

```powershell
cd belakang
go run ./cmd/server
```

Terminal 2 — frontend (http://localhost:5173):

```powershell
cd depan
npm run dev
```

Frontend dev memanggil API ke `http://localhost:3000/api` (lihat `depan/.env`).
Backend juga menyajikan hasil build terakhir di port 3000, tetapi saat
pengembangan gunakan port 5173 agar hot reload aktif.

## Build & deploy

```powershell
cd belakang
.\deploy.ps1              # build frontend + copy dist + binary Linux (web254)
.\deploy.ps1 -Upload      # sekaligus upload ke VPS (butuh SSH key)
```

Script otomatis mendeteksi lokasi project, jadi tidak perlu edit path saat
folder dipindah. Host VPS default `smpn254.sch.id` dan bisa dioverride:

```powershell
$env:SEKOLAH_VPS_HOST = "1.2.3.4"
.\deploy.ps1 -Upload
```

Yang perlu ada di server: binary `web254`, file `.env` produksi, dan folder
`storage/` (upload pengguna). Jalankan lewat systemd/service manager.

## Catatan

- `belakang/.env` berisi password DB dan JWT secret: jangan di-commit.
- Folder `belakang/storage/` (~385 MB) dan dump SQL besar tidak masuk git.
  Backup keduanya secara terpisah (mis. zip harian + `mysqldump`).
- Build produksi memakai `cmd/server/frontend/dist`; folder itu dihasilkan
  oleh `deploy.ps1`, jadi `depan/dist` tidak perlu di-commit.
