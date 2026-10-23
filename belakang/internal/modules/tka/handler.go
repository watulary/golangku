package tka

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
	"sekolah/internal/database"
)

type errItem struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}

// GetAllAdminHandler mengembalikan daftar TKA untuk halaman admin (pagination + search + sort).
func GetAllAdminHandler(c *fiber.Ctx) error {
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
			"success": false, "message": "Gagal mengambil data TKA",
		})
	}
	total, err := CountAdmin(search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menghitung data TKA",
		})
	}
	totalPage := (total + limit - 1) / limit
	if totalPage <= 0 {
		totalPage = 1
	}

	return c.JSON(fiber.Map{
		"success": true, "message": "Berhasil mengambil data TKA",
		"data": items,
		"pagination": fiber.Map{
			"page": page, "limit": limit, "total": total, "total_page": totalPage,
		},
	})
}

// BulkDeleteHandler menghapus banyak baris TKA sekaligus (checklist + DashboardBulkAction).
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
			"success": false, "message": "Gagal menghapus data TKA",
		})
	}
	return c.JSON(fiber.Map{
		"success": true, "message": strconv.Itoa(len(body.NIS)) + " data TKA berhasil dihapus",
	})
}

// TemplateHandler menghasilkan file Excel kosong (header + sheet Petunjuk) untuk diisi lalu diimpor.
func TemplateHandler(c *fiber.Ctx) error {
	f := excelize.NewFile()
	defer f.Close()
	sheet := "TKA"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"ID Tahun Pelajaran", "ID Kelas", "ID Rombel", "NIS", "No Peserta",
		"Nilai Matematika", "Keterangan MTK", "Nilai Bhs Indonesia", "Keterangan Bind"}
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

	// Semua kolom (bukan cuma NIS) diberi format Text, supaya Excel tidak
	// "mencerdaskan" nilai apa pun yang diketik user saat mengisi template.
	textStyle, _ := f.NewStyle(&excelize.Style{NumFmt: 49})
	firstCol, _ := excelize.CoordinatesToCellName(1, 2)
	lastColBottom, _ := excelize.CoordinatesToCellName(len(headers), 1000)
	f.SetCellStyle(sheet, firstCol, lastColBottom, textStyle)

	f.NewSheet("Petunjuk")
	f.SetCellValue("Petunjuk", "A1", "Aturan pengisian:")
	f.SetCellValue("Petunjuk", "A2", "1. Jangan mengubah urutan kolom.")
	f.SetCellValue("Petunjuk", "A3", "2. Kolom NIS harus sudah terdaftar di data Peserta Didik.")
	f.SetCellValue("Petunjuk", "A4", "3. Semua kolom wajib diisi.")

	row := 6
	f.SetCellValue("Petunjuk", fmt.Sprintf("A%d", row), "Referensi ID Tahun Pelajaran:")
	row++
	tahunList, _ := getTahunPelajaranAktifList()
	for _, t := range tahunList {
		f.SetCellValue("Petunjuk", fmt.Sprintf("A%d", row), fmt.Sprintf("%d - %s", t.ID, t.Nama))
		row++
	}

	row++
	f.SetCellValue("Petunjuk", fmt.Sprintf("A%d", row), "Referensi ID Kelas:")
	row++
	kelasList, _ := getKelasAktifList()
	for _, k := range kelasList {
		f.SetCellValue("Petunjuk", fmt.Sprintf("A%d", row), fmt.Sprintf("%d - %s", k.ID, k.Nama))
		row++
	}

	row++
	f.SetCellValue("Petunjuk", fmt.Sprintf("A%d", row), "Referensi ID Rombel:")
	row++
	rombelList, _ := getRombelAktifList()
	for _, r := range rombelList {
		f.SetCellValue("Petunjuk", fmt.Sprintf("A%d", row), fmt.Sprintf("%d - %s", r.ID, r.Nama))
		row++
	}

	filename := fmt.Sprintf("template-tka-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	buf, err := f.WriteToBuffer()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal membuat template",
		})
	}
	return c.Send(buf.Bytes())
}

// ImportHandler memvalidasi seluruh baris Excel dulu (all-or-nothing), baru simpan dalam 1 transaksi.
func ImportHandler(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file_excel")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "File tidak ditemukan",
		})
	}
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Gagal membuka file",
		})
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "File excel tidak valid",
		})
	}
	defer f.Close()
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil || len(rows) <= 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "File excel kosong",
		})
	}

	var errs []errItem
	var dataList []Tka
	nisSeen := map[string]bool{}
	nopesSeen := map[string]bool{}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		get := func(idx int) string {
			if idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}
		idTP, idKelas, idRombel := get(0), get(1), get(2)
		nis, nopes, mtk, ketMtk, bind, ketBind := get(3), get(4), get(5), get(6), get(7), get(8)

		if idTP == "" && idKelas == "" && idRombel == "" && nis == "" && nopes == "" &&
			mtk == "" && ketMtk == "" && bind == "" && ketBind == "" {
			continue
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
		if nopes == "" {
			errs = append(errs, errItem{rowNo, "Nopes kosong."})
		}
		if mtk == "" {
			errs = append(errs, errItem{rowNo, "Nilai MTK kosong."})
		}
		if ketMtk == "" {
			errs = append(errs, errItem{rowNo, "Ket MTK kosong."})
		}
		if bind == "" {
			errs = append(errs, errItem{rowNo, "Nilai Bind kosong."})
		}
		if ketBind == "" {
			errs = append(errs, errItem{rowNo, "Ket Bind kosong."})
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
			sudahAda, errChk2 := ExistsByNIS(nis)
			if errChk2 == nil && sudahAda {
				errs = append(errs, errItem{rowNo, "NIS (" + nis + ") sudah terdaftar di data TKA."})
			}
		}
		if nopes != "" {
			if nopesSeen[nopes] {
				errs = append(errs, errItem{rowNo, "Nopes (" + nopes + ") terdapat data ganda dalam file."})
			}
			nopesSeen[nopes] = true

			sudahAda, errChk := ExistsByNopes(nopes)
			if errChk == nil && sudahAda {
				errs = append(errs, errItem{rowNo, "Nopes (" + nopes + ") sudah terdaftar di data TKA."})
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

		dataList = append(dataList, Tka{
			IDTahunPelajaran: idTPInt, IDKelas: idKelasInt, IDRombel: idRombelInt,
			NIS: nis, Nopes: nopes, MTK: mtk, KetMTK: ketMtk, Bind: bind, KetBind: ketBind,
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

	tx, err := database.DB.Beginx()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal memulai transaksi",
		})
	}
	for _, d := range dataList {
		if err := Insert(tx, d); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false, "message": "Gagal menyimpan data TKA",
			})
		}
	}
	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal menyimpan data",
		})
	}

	return c.JSON(fiber.Map{
		"success": true, "message": "Unggah " + strconv.Itoa(len(dataList)) + " data berhasil.",
	})
}

// ExportHandler mengunduh seluruh data TKA (dengan join lengkap) sebagai Excel.
func ExportHandler(c *fiber.Ctx) error {
	search := c.Query("search", "")
	items, err := GetAllForExport(search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil data TKA",
		})
	}

	f := excelize.NewFile()
	defer f.Close()
	sheet := "TKA"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{
		"No", "ID Tahun Pelajaran", "Tahun Pelajaran", "ID Kelas", "Kelas", "ID Rombel", "Rombel",
		"NIS", "NISN", "No Peserta", "Nama Siswa", "L/P", "Tempat Lahir", "Tanggal Lahir",
		"Matematika", "Keterangan", "Bhs Indonesia", "Keterangan",
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

	// Semua kolom data (bukan cuma NIS/NISN/No Peserta) diberi format Text,
	// supaya Excel tidak strip angka nol depan atau "mencerdaskan" nilai apa pun.
	textStyle, _ := f.NewStyle(&excelize.Style{NumFmt: 49})
	for i, item := range items {
		row := i + 2
		vals := []interface{}{
			i + 1, item.IDTahunPelajaran, item.NamaTahunPelajaran, item.IDKelas, item.NamaKelas,
			item.IDRombel, item.NamaRombel, item.NIS, item.NISN, item.Nopes, item.NamaSiswa,
			item.JenisKelamin, item.TempatLahirSiswa, item.TanggalLahirSiswa,
			item.MTK, item.KetMTK, item.Bind, item.KetBind,
		}
		for j, v := range vals {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheet, cell, v)
			f.SetCellStyle(sheet, cell, cell, textStyle)
		}
	}
	f.SetPanes(sheet, &excelize.Panes{Freeze: true, Split: false, XSplit: 0, YSplit: 1, TopLeftCell: "A2", ActivePane: "bottomLeft"})

	buf, err := f.WriteToBuffer()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal membuat file excel",
		})
	}

	filename := fmt.Sprintf("data-tka-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(buf.Bytes())
}

// --- tka_waktu (pengaturan pengumuman) ---

func GetWaktuHandler(c *fiber.Ctx) error {
	w, err := GetWaktu()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil pengaturan pengumuman TKA",
		})
	}
	return c.JSON(fiber.Map{
		"success": true, "message": "Berhasil mengambil pengaturan pengumuman TKA", "data": w,
	})
}

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
			"success": false, "message": "Gagal memuat data pengumuman TKA",
		})
	}

	var data TkaWaktu
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Data tidak valid",
		})
	}

	var changed []string
	if old.AturWaktu != data.AturWaktu {
		changed = append(changed, "aturwaktu")
	}
	if old.AturForm != data.AturForm {
		changed = append(changed, "aturform")
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
	if old.TtdSKL != data.TtdSKL {
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
			"success": false, "message": "Gagal memperbarui pengaturan pengumuman TKA",
		})
	}
	return c.JSON(fiber.Map{
		"success": true, "message": "Ubah data TKA berhasil.", "changed": changed,
	})
}

// GetMeHandler mengembalikan data TKA milik siswa yang login (role Pesertadidik).
// Kalau siswa belum punya record TKA, "tka" dikembalikan null (bukan error).
func GetMeHandler(c *fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil || claims.NIS == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false, "message": "Tidak terautentikasi",
		})
	}

	w, err := GetWaktu()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil pengaturan pengumuman TKA",
		})
	}

	item, err := GetByNIS(claims.NIS)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Gagal mengambil data TKA",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Berhasil mengambil data TKA",
		"data": fiber.Map{
			"waktu": w,
			"tka":   item, // null kalau belum ada record
		},
	})
}
