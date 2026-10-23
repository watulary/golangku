package kenaikan

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"

	"sekolah/internal/auth"
	"sekolah/internal/config"
)

// fotoDefaultFor mengembalikan nama file foto default sesuai jenis kelamin.
func fotoDefaultFor(jenisKelamin string) string {
	if jenisKelamin == "P" {
		return "P.jpg"
	}
	return "L.jpg"
}

// buildFotoURL menghasilkan URL foto siswa dari kolom foto pesertadidik,
// jatuh ke foto default sesuai jenis kelamin kalau siswa belum pernah upload foto.
func buildFotoURL(foto, jenisKelamin string) string {
	if foto == "" {
		foto = fotoDefaultFor(jenisKelamin)
	}
	return config.App.AppURL + "/storage/fotosiswa/" + foto
}

// GetAdminListHandler mengembalikan daftar kenaikan untuk halaman admin
// (pagination + search + sort), admin-only tanpa filter ownership.
func GetAdminListHandler(c *fiber.Ctx) error {
	search := c.Query("search", "")
	sortBy := c.Query("sort_by", "")
	sortDir := c.Query("sort_dir", "asc")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	items, err := GetAllAdmin(search, sortBy, sortDir, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil data kenaikan",
		})
	}
	total, err := CountAdmin(search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menghitung data kenaikan",
		})
	}
	for i := range items {
		items[i].FotoURL = buildFotoURL(items[i].Foto, items[i].JenisKelamin)
	}
	totalPage := (total + limit - 1) / limit
	if totalPage <= 0 {
		totalPage = 1
	}

	return c.JSON(fiber.Map{
		"success": true, "message": "Berhasil mengambil data kenaikan",
		"data": items,
		"pagination": fiber.Map{
			"page": page, "limit": limit, "total": total, "total_page": totalPage,
		},
	})
}

// GetByIDHandler mengembalikan satu data kenaikan berdasarkan nis (dipakai form edit).
func GetByIDHandler(c *fiber.Ctx) error {
	nis := c.Params("nis")
	item, err := GetByID(nis)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false, "message": "Data kenaikan tidak ditemukan",
		})
	}
	item.FotoURL = buildFotoURL(item.Foto, item.JenisKelamin)
	return c.JSON(fiber.Map{
		"success": true, "message": "Berhasil mengambil data kenaikan", "data": item,
	})
}

// CreateHandler menyimpan satu baris data kenaikan baru dari form tambah manual.
func CreateHandler(c *fiber.Ctx) error {
	var data Kenaikan
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Data tidak valid",
		})
	}
	data.NIS = strings.TrimSpace(data.NIS)
	if data.NIS == "" || data.IDTahunPelajaran == 0 || data.IDKelas == 0 || data.IDRombel == 0 || data.Ket == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Semua field wajib diisi",
		})
	}

	adaSiswa, err := ExistsPesertaDidik(data.NIS)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal memeriksa data peserta didik",
		})
	}
	if !adaSiswa {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "NIS tidak terdaftar di data peserta didik",
		})
	}

	sudahAda, err := Exists(data.NIS)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal memeriksa data kenaikan",
		})
	}
	if sudahAda {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "NIS sudah terdaftar di data kenaikan",
		})
	}

	if err := Create(data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menyimpan data kenaikan",
		})
	}
	return c.JSON(fiber.Map{
		"success": true, "message": "Data kenaikan berhasil ditambahkan",
	})
}

// UpdateHandler memperbarui data kenaikan. nis read-only (tidak diubah lewat sini).
// Cek changed dulu sebelum update DB, response menyertakan field changed.
func UpdateHandler(c *fiber.Ctx) error {
	nis := c.Params("nis")

	old, err := GetByID(nis)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false, "message": "Data kenaikan tidak ditemukan",
		})
	}

	var data Kenaikan
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Data tidak valid",
		})
	}
	if data.IDTahunPelajaran == 0 || data.IDKelas == 0 || data.IDRombel == 0 || data.Ket == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Semua field wajib diisi",
		})
	}

	var changed []string
	if old.IDTahunPelajaran != data.IDTahunPelajaran {
		changed = append(changed, "id_tahunpelajaran")
	}
	if old.IDKelas != data.IDKelas {
		changed = append(changed, "id_kelas")
	}
	if old.IDRombel != data.IDRombel {
		changed = append(changed, "id_rombel")
	}
	if old.Ket != data.Ket {
		changed = append(changed, "ket")
	}

	if len(changed) == 0 {
		return c.JSON(fiber.Map{
			"success": true, "message": "Tidak ada perubahan data.", "changed": changed,
		})
	}

	if err := Update(nis, data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal memperbarui data kenaikan",
		})
	}
	return c.JSON(fiber.Map{
		"success": true, "message": "Data kenaikan berhasil diperbarui", "changed": changed,
	})
}

// DeleteHandler menghapus satu baris data kenaikan.
func DeleteHandler(c *fiber.Ctx) error {
	nis := c.Params("nis")
	if err := Delete(nis); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menghapus data kenaikan",
		})
	}
	return c.JSON(fiber.Map{
		"success": true, "message": "Data kenaikan berhasil dihapus",
	})
}

// BulkDeleteHandler menghapus banyak baris sekaligus (checklist + DashboardBulkAction).
func BulkDeleteHandler(c *fiber.Ctx) error {
	var body struct {
		NIS []string `json:"nis"`
	}
	if err := c.BodyParser(&body); err != nil || len(body.NIS) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Tidak ada data yang dipilih",
		})
	}
	if err := BulkDelete(body.NIS); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menghapus data kenaikan",
		})
	}
	return c.JSON(fiber.Map{
		"success": true, "message": strconv.Itoa(len(body.NIS)) + " data kenaikan berhasil dihapus",
	})
}

// TemplateHandler menghasilkan file Excel kosong (header saja) untuk diisi lalu diimpor.
func TemplateHandler(c *fiber.Ctx) error {
	f := excelize.NewFile()
	defer f.Close()
	sheet := "Kenaikan"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"ID Tahun Pelajaran", "ID Kelas", "ID Rombel", "NIS", "Keterangan"}
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"2563EB"}, Pattern: 1},
	})
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	lastCell, _ := excelize.CoordinatesToCellName(len(headers), 1)
	f.SetCellStyle(sheet, "A1", lastCell, headerStyle)
	f.SetPanes(sheet, &excelize.Panes{Freeze: true, Split: false, XSplit: 0, YSplit: 1, TopLeftCell: "A2", ActivePane: "bottomLeft"})

	// Sheet Petunjuk: daftar referensi id_tahunpelajaran/id_kelas/id_rombel yang valid saat ini.
	f.NewSheet("Petunjuk")
	f.SetCellValue("Petunjuk", "A1", "Aturan pengisian:")
	f.SetCellValue("Petunjuk", "A2", "1. ID Tahun Pelajaran, ID Kelas, ID Rombel harus sesuai data referensi yang berlaku (lihat daftar di bawah).")
	f.SetCellValue("Petunjuk", "A3", "2. NIS harus sudah terdaftar di data Peserta Didik.")
	f.SetCellValue("Petunjuk", "A4", "3. NIS tidak boleh duplikat, dan tidak boleh sudah ada di data Kenaikan.")
	f.SetCellValue("Petunjuk", "A5", "4. Keterangan wajib diisi (contoh: Naik ke Kelas VIII, Tinggal Kelas VII, Lulus).")

	f.SetCellValue("Petunjuk", "A7", "Referensi Tahun Pelajaran (id - nama):")
	f.SetCellValue("Petunjuk", "A8", "Referensi Kelas (id - nama):")
	f.SetCellValue("Petunjuk", "A9", "Referensi Rombel (id - nama):")

	activeIdx, _ := f.GetSheetIndex(sheet)
	f.SetActiveSheet(activeIdx)

	buf, err := f.WriteToBuffer()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal membuat template",
		})
	}
	filename := fmt.Sprintf("template-kenaikan-kelas-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(buf.Bytes())
}

// ImportHandler mengimpor data kenaikan dari file Excel, all-or-nothing.
func ImportHandler(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "File excel wajib diunggah",
		})
	}
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Gagal membaca file excel",
		})
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Format file excel tidak valid",
		})
	}
	defer f.Close()

	sheet := f.GetSheetList()[0]
	rows, err := f.GetRows(sheet)
	if err != nil || len(rows) <= 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "File excel kosong",
		})
	}

	type errItem struct {
		Row     int    `json:"row"`
		Message string `json:"message"`
	}
	var errs []errItem
	var dataList []Kenaikan
	nisSeen := map[string]bool{}

	for i, row := range rows {
		if i == 0 {
			continue // lewati header
		}
		get := func(idx int) string {
			if idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}
		idTP := get(0)
		idKelas := get(1)
		idRombel := get(2)
		nis := get(3)
		ket := get(4)

		if idTP == "" && idKelas == "" && idRombel == "" && nis == "" && ket == "" {
			continue // lewati baris kosong
		}

		rowNo := i + 1
		if idTP == "" {
			errs = append(errs, errItem{rowNo, "ID Tahun Pelajaran kosong."})
		}
		if idKelas == "" {
			errs = append(errs, errItem{rowNo, "ID Kelas kosong."})
		}
		if idRombel == "" {
			errs = append(errs, errItem{rowNo, "ID Rombel kosong."})
		}
		if nis == "" {
			errs = append(errs, errItem{rowNo, "NIS kosong."})
		}
		if ket == "" {
			errs = append(errs, errItem{rowNo, "Keterangan kosong."})
		}

		if nis != "" {
			if nisSeen[nis] {
				errs = append(errs, errItem{rowNo, "NIS (" + nis + ") terdapat data ganda dalam file."})
			}
			nisSeen[nis] = true

			adaSiswa, errChk := ExistsPesertaDidik(nis)
			if errChk == nil && !adaSiswa {
				errs = append(errs, errItem{rowNo, "NIS (" + nis + ") tidak terdaftar di data peserta didik."})
			}
			sudahAda, errChk2 := Exists(nis)
			if errChk2 == nil && sudahAda {
				errs = append(errs, errItem{rowNo, "NIS (" + nis + ") sudah terdaftar di data kenaikan."})
			}
		}

		idTPInt, e1 := strconv.Atoi(idTP)
		idKelasInt, e2 := strconv.Atoi(idKelas)
		idRombelInt, e3 := strconv.Atoi(idRombel)
		if idTP != "" && e1 != nil {
			errs = append(errs, errItem{rowNo, "ID Tahun Pelajaran harus berupa angka."})
		}
		if idKelas != "" && e2 != nil {
			errs = append(errs, errItem{rowNo, "ID Kelas harus berupa angka."})
		}
		if idRombel != "" && e3 != nil {
			errs = append(errs, errItem{rowNo, "ID Rombel harus berupa angka."})
		}

		dataList = append(dataList, Kenaikan{
			IDTahunPelajaran: idTPInt, IDKelas: idKelasInt, IDRombel: idRombelInt, NIS: nis, Ket: ket,
		})
	}

	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Impor gagal, periksa kembali data.", "errors": errs,
		})
	}
	if len(dataList) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "File excel kosong",
		})
	}

	for _, d := range dataList {
		if err := Create(d); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false, "message": "Gagal menyimpan data kenaikan",
			})
		}
	}

	return c.JSON(fiber.Map{
		"success": true, "message": "Unggah " + strconv.Itoa(len(dataList)) + " data berhasil.",
	})
}

// ExportHandler mengunduh seluruh data kenaikan (dengan join lengkap) sebagai Excel.
func ExportHandler(c *fiber.Ctx) error {
	search := c.Query("search", "")
	items, err := GetAllForExport(search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil data kenaikan",
		})
	}

	f := excelize.NewFile()
	defer f.Close()
	sheet := "Kenaikan"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{
		"No", "ID Tahun Pelajaran", "Tahun Pelajaran", "ID Kelas", "Kelas", "ID Rombel", "Rombel",
		"NIS", "NISN", "No Peserta Ujian", "Nama Siswa", "L/P", "Tempat Lahir", "Tanggal Lahir", "Keterangan",
	}
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"2563EB"}, Pattern: 1},
	})
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	lastCell, _ := excelize.CoordinatesToCellName(len(headers), 1)
	f.SetCellStyle(sheet, "A1", lastCell, headerStyle)

	for i, item := range items {
		row := i + 2
		vals := []interface{}{
			i + 1, item.IDTahunPelajaran, item.NamaTahunPelajaran, item.IDKelas, item.NamaKelas,
			item.IDRombel, item.NamaRombel, item.NIS, item.NISN, item.NoUjianSD, item.NamaSiswa,
			item.JenisKelamin, item.TempatLahirSiswa, item.TanggalLahirSiswa, item.Ket,
		}
		for j, v := range vals {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheet, cell, v)
		}
	}
	f.SetPanes(sheet, &excelize.Panes{Freeze: true, Split: false, XSplit: 0, YSplit: 1, TopLeftCell: "A2", ActivePane: "bottomLeft"})

	buf, err := f.WriteToBuffer()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal membuat file excel",
		})
	}

	filename := fmt.Sprintf("data-kenaikan-kelas-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(buf.Bytes())
}

// --- kenaikan_waktu (pengaturan pengumuman) ---

// GetWaktuHandler mengembalikan pengaturan pengumuman kenaikan kelas saat ini.
func GetWaktuHandler(c *fiber.Ctx) error {
	w, err := GetWaktu()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil pengaturan pengumuman kenaikan",
		})
	}
	return c.JSON(fiber.Map{
		"success": true, "message": "Berhasil mengambil pengaturan pengumuman kenaikan", "data": w,
	})
}

// UpdateWaktuHandler memperbarui pengaturan pengumuman kenaikan kelas (single-row).
func UpdateWaktuHandler(c *fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false, "message": "Tidak terautentikasi",
		})
	}

	old, err := GetWaktu()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal memuat data pengumuman kenaikan",
		})
	}

	var data KenaikanWaktu
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Data tidak valid",
		})
	}

	var changed []string
	if old.Aturwaktu != data.Aturwaktu {
		changed = append(changed, "aturwaktu")
	}
	if old.Judul != data.Judul {
		changed = append(changed, "judul")
	}
	if old.Nomor != data.Nomor {
		changed = append(changed, "nomor")
	}
	if old.Pengantar != data.Pengantar {
		changed = append(changed, "pengantar")
	}
	if old.Penutup != data.Penutup {
		changed = append(changed, "penutup")
	}
	if old.Akhir != data.Akhir {
		changed = append(changed, "akhir")
	}
	if old.Tanggal != data.Tanggal {
		changed = append(changed, "tanggal")
	}
	if old.Mengetahui != data.Mengetahui {
		changed = append(changed, "mengetahui")
	}
	if old.Kepsek != data.Kepsek {
		changed = append(changed, "kepsek")
	}
	if old.NIPKepsek != data.NIPKepsek {
		changed = append(changed, "nipkepsek")
	}
	if old.TTDSkl != data.TTDSkl {
		changed = append(changed, "ttd_skl")
	}
	if old.TahunPelajaran != data.TahunPelajaran {
		changed = append(changed, "tahun_pelajaran")
	}

	if len(changed) == 0 {
		return c.JSON(fiber.Map{
			"success": true, "message": "Tidak ada perubahan data.", "changed": changed,
		})
	}

	if err := UpdateWaktu(claims.IDPengguna, data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal memperbarui pengaturan pengumuman kenaikan",
		})
	}
	return c.JSON(fiber.Map{
		"success": true, "message": "Ubah data kenaikan kelas berhasil.", "changed": changed,
	})
}

// GetSiswaHandler mengembalikan data kenaikan milik siswa yang login (role Pesertadidik).
// Kalau siswa belum punya record kenaikan, "kenaikan" dikembalikan null (bukan error).
func GetSiswaHandler(c *fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil || claims.NIS == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false, "message": "Tidak terautentikasi",
		})
	}

	w, err := GetWaktu()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil pengaturan pengumuman kenaikan",
		})
	}

	item, err := GetByID(claims.NIS)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil data kenaikan",
		})
	}
	if item != nil {
		item.FotoURL = buildFotoURL(item.Foto, item.JenisKelamin)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Berhasil mengambil data kenaikan",
		"data": fiber.Map{
			"waktu":    w,
			"kenaikan": item, // null kalau belum ada record
		},
	})
}
