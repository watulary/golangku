package pesertadidik

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"sekolah/internal/auth"
	"sekolah/internal/config"
)

const (
	fotoDir      = "storage/fotosiswa"
	fotoDirArsip = "storage/fotosiswaarsip"
	fotoDefaultL = "L.jpg"
	fotoDefaultP = "P.jpg"
)

// fotoDefaultFor menentukan nama file default berdasarkan jenis_kelamin.
// Selain "P", dianggap "L" (termasuk kalau jenis_kelamin kosong/belum diisi).
func fotoDefaultFor(jenisKelamin string) string {
	if jenisKelamin == "P" {
		return fotoDefaultP
	}
	return fotoDefaultL
}

func buildFotoURL(foto, jenisKelamin string) string {
	if foto == "" {
		foto = fotoDefaultFor(jenisKelamin)
	}
	return config.App.AppURL + "/storage/fotosiswa/" + foto
}

// buildFotoURLArsip sama seperti buildFotoURL tapi mengarah ke folder khusus arsip
// (storage/fotosiswaarsip). Foto siswa yang sudah diarsipkan disalin fisik ke folder ini
// (lihat copyFotoToArsip, dipanggil dari ArsipkanSiswa) supaya tidak lagi bergantung pada
// file di storage/fotosiswa milik data aktif — yang bisa berubah/terhapus kalau foto siswa
// diganti belakangan. Foto default tetap diambil dari folder aktif (asset bersama, tidak
// digandakan). Kalau file belum ada di folder arsip (mis. data arsip lama sebelum fitur
// copy-otomatis ini ada), fallback ke folder aktif supaya foto lama yang masih ada di sana
// tetap tampil, bukan langsung 404.
func buildFotoURLArsip(foto, jenisKelamin string) string {
	if foto == "" {
		return config.App.AppURL + "/storage/fotosiswa/" + fotoDefaultFor(jenisKelamin)
	}
	if _, err := os.Stat(filepath.Join(fotoDirArsip, foto)); err == nil {
		return config.App.AppURL + "/storage/fotosiswaarsip/" + foto
	}
	return config.App.AppURL + "/storage/fotosiswa/" + foto
}

// ==================== PESERTA DIDIK AKTIF ====================

func GetAdminListHandler(c *fiber.Ctx) error {
	search := c.Query("search", "")
	idTahunPelajaran := c.QueryInt("id_tahunpelajaran", 0)
	idKelas := c.QueryInt("id_kelas", 0)
	idRombel := c.QueryInt("id_rombel", 0)
	sortBy := c.Query("sort_by", "nama_siswa")
	sortDir := c.Query("sort_dir", "asc")
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)
	offset := (page - 1) * limit

	data, err := GetAllAdmin(search, idTahunPelajaran, idKelas, idRombel, sortBy, sortDir, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data peserta didik"})
	}
	total, err := CountAdmin(search, idTahunPelajaran, idKelas, idRombel)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menghitung data peserta didik"})
	}

	for i := range data {
		data[i].FotoURL = buildFotoURL(data[i].Foto, data[i].JenisKelamin)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Berhasil mengambil data peserta didik",
		"data":    data,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func GetByIDHandler(c *fiber.Ctx) error {
	nis := c.Params("nis")

	data, err := GetByID(nis)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data peserta didik tidak ditemukan"})
	}
	data.FotoURL = buildFotoURL(data.Foto, data.JenisKelamin)

	return c.JSON(fiber.Map{"success": true, "message": "Berhasil mengambil data peserta didik", "data": data})
}

func CreateHandler(c *fiber.Ctx) error {
	var data PesertaDidik
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Data tidak valid"})
	}
	if data.NIS == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "NIS wajib diisi"})
	}

	if _, err := GetByID(data.NIS); err == nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "NIS sudah terdaftar"})
	}

	if err := Create(&data); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menyimpan data peserta didik"})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "message": "Peserta didik berhasil ditambahkan", "data": data})
}

// UpdateHandler cek dulu field mana yang berubah sebelum update ke DB
func UpdateHandler(c *fiber.Ctx) error {
	nis := c.Params("nis")

	existing, err := GetByID(nis)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data peserta didik tidak ditemukan"})
	}

	var data PesertaDidik
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Data tidak valid"})
	}
	data.NIS = nis
	// pertahankan foto lama, upload foto ditangani endpoint terpisah
	data.Foto = existing.Foto

	changed := diffChangedFields(existing, &data)
	if len(changed) == 0 {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Tidak ada perubahan data",
			"changed": changed,
			"data":    existing,
		})
	}

	if err := Update(nis, &data); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal memperbarui data peserta didik"})
	}

	data.FotoURL = buildFotoURL(data.Foto, data.JenisKelamin)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data peserta didik berhasil diperbarui",
		"changed": changed,
		"data":    data,
	})
}

// diffChangedFields membandingkan field-field db-tagged antara data lama & baru, mengembalikan nama kolom yang berubah
func diffChangedFields(oldData, newData *PesertaDidik) []string {
	changed := []string{}
	vOld := reflect.ValueOf(oldData).Elem()
	vNew := reflect.ValueOf(newData).Elem()
	t := vOld.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			continue
		}
		oldVal := vOld.Field(i).Interface()
		newVal := vNew.Field(i).Interface()
		if oldVal != newVal {
			changed = append(changed, dbTag)
		}
	}
	return changed
}

func DeleteHandler(c *fiber.Ctx) error {
	nis := c.Params("nis")

	if _, err := GetByID(nis); err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data peserta didik tidak ditemukan"})
	}

	if err := Delete(nis); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menghapus data peserta didik"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Data peserta didik berhasil dihapus"})
}

// UploadFotoHandler mengganti foto siswa (uuid baru tiap upload, hapus file lama setelah update DB sukses)
func UploadFotoHandler(c *fiber.Ctx) error {
	nis := c.Params("nis")

	existing, err := GetByID(nis)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data peserta didik tidak ditemukan"})
	}

	file, err := c.FormFile("foto")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "File foto wajib diunggah"})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
	default:
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Foto harus berformat JPG, JPEG, PNG atau WEBP"})
	}

	if file.Size > 100*1024 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Ukuran foto maksimal 100 KB"})
	}

	newFilename := uuid.Must(uuid.NewV7()).String() + ext
	savePath := filepath.Join(fotoDir, newFilename)

	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menyimpan foto"})
	}

	if err := UpdateFoto(nis, newFilename); err != nil {
		os.Remove(savePath)
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal memperbarui foto"})
	}

	oldFoto := existing.Foto
	if oldFoto != "" {
		if masihDipakai, err := FotoMasihDipakai(oldFoto); err != nil || masihDipakai {
			// gagal cek ATAU masih direferensikan snapshot mutasi/arsip -> jangan hapus file fisik
		} else {
			go deleteFotoFile(oldFoto)
		}
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"message":  "Foto peserta didik berhasil diperbarui",
		"foto":     newFilename,
		"foto_url": buildFotoURL(newFilename, existing.JenisKelamin),
	})
}

func deleteFotoFile(filename string) {
	path := filepath.Join(fotoDir, filename)

	// Beri waktu proses lain (antivirus, thumbnail cache, cloud sync, dsb)
	// melepas lock file sebelum mulai mencoba menghapus.
	time.Sleep(3 * time.Second)

	const maxAttempt = 30
	const retryDelay = 1 * time.Second

	for attempt := 1; attempt <= maxAttempt; attempt++ {
		err := os.Remove(path)

		if err == nil || os.IsNotExist(err) {
			return
		}

		time.Sleep(retryDelay)
	}
}

// copyFotoToArsip menyalin fisik file foto dari folder aktif (storage/fotosiswa) ke folder
// khusus arsip (storage/fotosiswaarsip). Dipanggil setelah ArsipkanSiswa berhasil, supaya foto
// arsip punya salinan sendiri yang independen dari data aktif. Best-effort: kalau source belum
// ada / gagal disalin, proses arsip tetap dianggap berhasil (foto arsip tinggal fallback default).
func copyFotoToArsip(filename string) {
	dstPath := filepath.Join(fotoDirArsip, filename)
	if _, err := os.Stat(dstPath); err == nil {
		return // sudah ada (mis. arsip ulang di tahun pelajaran berbeda), tidak perlu disalin lagi
	}

	srcPath := filepath.Join(fotoDir, filename)
	src, err := os.Open(srcPath)
	if err != nil {
		return
	}
	defer src.Close()

	if err := os.MkdirAll(fotoDirArsip, 0755); err != nil {
		return
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return
	}
	defer dst.Close()

	io.Copy(dst, src)
}

// ==================== MUTASI ====================

type mutasiRequest struct {
	StatusKeluar  string `json:"status_keluar"`
	AlasanKeluar  string `json:"alasan_keluar"`
	TanggalKeluar string `json:"tanggal_keluar"`
	Tujuan        string `json:"tujuan"`
	NomorSurat    string `json:"nomor_surat"`
}

func MutasiHandler(c *fiber.Ctx) error {
	nis := c.Params("nis")

	var req mutasiRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Data tidak valid"})
	}
	if req.StatusKeluar == "" || req.TanggalKeluar == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Status dan tanggal keluar wajib diisi"})
	}

	if err := MutasiSiswa(nis, req.StatusKeluar, req.AlasanKeluar, req.TanggalKeluar, req.Tujuan, req.NomorSurat); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal memproses mutasi: " + err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": fmt.Sprintf("Siswa NIS %s berhasil dimutasi", nis)})
}

func GetMutasiListHandler(c *fiber.Ctx) error {
	search := c.Query("search", "")
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)
	offset := (page - 1) * limit

	data, err := GetAllMutasi(search, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data mutasi"})
	}
	total, err := CountMutasi(search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menghitung data mutasi"})
	}

	for i := range data {
		data[i].FotoURL = buildFotoURL(data[i].Foto, data[i].JenisKelamin)
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"message":    "Berhasil mengambil data mutasi",
		"data":       data,
		"pagination": fiber.Map{"page": page, "limit": limit, "total": total},
	})
}

// ==================== ARSIP ====================

func ArsipHandler(c *fiber.Ctx) error {
	nis := c.Params("nis")

	if err := ArsipkanSiswa(nis); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal memproses arsip: " + err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": fmt.Sprintf("Siswa NIS %s berhasil diarsipkan (lulus)", nis)})
}

func GetArsipListHandler(c *fiber.Ctx) error {
	search := c.Query("search", "")
	idTahunPelajaran := c.QueryInt("id_tahunpelajaran", 0)
	idKelas := c.QueryInt("id_kelas", 0)
	idRombel := c.QueryInt("id_rombel", 0)
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)
	offset := (page - 1) * limit

	data, err := GetAllArsip(search, idTahunPelajaran, idKelas, idRombel, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data arsip"})
	}
	total, err := CountArsip(search, idTahunPelajaran, idKelas, idRombel)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menghitung data arsip"})
	}

	for i := range data {
		data[i].FotoURL = buildFotoURLArsip(data[i].Foto, data[i].JenisKelamin)
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"message":    "Berhasil mengambil data arsip",
		"data":       data,
		"pagination": fiber.Map{"page": page, "limit": limit, "total": total},
	})
}

func GetArsipByIDHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	data, err := GetArsipByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data arsip tidak ditemukan"})
	}
	data.FotoURL = buildFotoURLArsip(data.Foto, data.JenisKelamin)

	return c.JSON(fiber.Map{"success": true, "message": "Berhasil mengambil detail arsip", "data": data})
}

// DeleteFotoHandler menghapus foto siswa (kembali ke foto default) dan menghapus file lama dari disk
func DeleteFotoHandler(c *fiber.Ctx) error {
	nis := c.Params("nis")

	existing, err := GetByID(nis)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data peserta didik tidak ditemukan"})
	}
	if existing.Foto == "" {
		return c.JSON(fiber.Map{
			"success":  true,
			"message":  "Foto sudah menggunakan foto default",
			"foto_url": buildFotoURL("", existing.JenisKelamin),
		})
	}

	if err := DeleteFoto(nis); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menghapus foto"})
	}

	oldFoto := existing.Foto
	if masihDipakai, err := FotoMasihDipakai(oldFoto); err != nil || masihDipakai {
		// gagal cek ATAU masih direferensikan snapshot mutasi/arsip -> jangan hapus file fisik
	} else {
		go deleteFotoFile(oldFoto)
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"message":  "Foto berhasil dihapus, kembali ke foto default",
		"foto_url": buildFotoURL("", existing.JenisKelamin),
	})
}

func CheckDuplicateHandler(c *fiber.Ctx) error {
	field := c.Query("field")
	value := c.Query("value")
	exclude := c.Query("exclude", "")

	exists, nis, nama, err := CheckDuplicate(field, value, exclude)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Gagal memeriksa data ganda"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"exists":     exists,
			"nis":        nis,
			"nama_siswa": nama,
		},
	})
}

func GetMutasiByIDHandler(c *fiber.Ctx) error {
	nis := c.Params("nis")

	data, err := GetMutasiByID(nis)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data mutasi tidak ditemukan"})
	}
	data.FotoURL = buildFotoURL(data.Foto, data.JenisKelamin)

	return c.JSON(fiber.Map{"success": true, "message": "Berhasil mengambil detail mutasi", "data": data})
}

type bulkNisRequest struct {
	NIS []string `json:"nis"`
}

func BulkDeleteHandler(c *fiber.Ctx) error {
	var req bulkNisRequest
	if err := c.BodyParser(&req); err != nil || len(req.NIS) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Daftar NIS wajib diisi"})
	}

	if err := BulkDelete(req.NIS); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menghapus data peserta didik"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("%d data peserta didik berhasil dihapus", len(req.NIS)),
	})
}

type bulkMutasiRequest struct {
	NIS           []string `json:"nis"`
	StatusKeluar  string   `json:"status_keluar"`
	AlasanKeluar  string   `json:"alasan_keluar"`
	TanggalKeluar string   `json:"tanggal_keluar"`
	Tujuan        string   `json:"tujuan"`
	NomorSurat    string   `json:"nomor_surat"`
}

func BulkMutasiHandler(c *fiber.Ctx) error {
	var req bulkMutasiRequest
	if err := c.BodyParser(&req); err != nil || len(req.NIS) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Daftar NIS wajib diisi"})
	}
	if req.StatusKeluar == "" || req.TanggalKeluar == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Status dan tanggal keluar wajib diisi"})
	}

	result := BulkMutasi(req.NIS, req.StatusKeluar, req.AlasanKeluar, req.TanggalKeluar, req.Tujuan, req.NomorSurat)

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("%d berhasil dimutasi, %d gagal", len(result.Success), len(result.Failed)),
		"data":    result,
	})
}

func BulkArsipHandler(c *fiber.Ctx) error {
	var req bulkNisRequest
	if err := c.BodyParser(&req); err != nil || len(req.NIS) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Daftar NIS wajib diisi"})
	}

	result := BulkArsipkan(req.NIS)

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("%d berhasil diarsipkan, %d gagal", len(result.Success), len(result.Failed)),
		"data":    result,
	})
}

func GetMeHandler(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*auth.Claims)

	data, err := GetByID(claims.NIS)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(fiber.Map{"success": true, "data": nil})
		}
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	data.FotoURL = buildFotoURL(data.Foto, data.JenisKelamin)

	return c.JSON(fiber.Map{"success": true, "data": data})
}
