package nilai

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

const (
	maxNilaiEkskul      = 20  // karakter, mis. "Sangat Baik" / "A"
	maxKeteranganEkskul = 255 // karakter
)

var ekskulExcelHeaders = []string{"NIS", "Nama Peserta Didik", "Nilai", "Keterangan"}

type ekskulKonteks struct {
	IDTahunPelajaran int
	Semester         string
	IDRombel         int
	IDEkskul         int

	IDKelas    int
	NamaTahun  string
	NamaKelas  string
	NamaRombel string
	NamaEkskul string
}

func resolveEkskulKonteks(idTP int, semester string, idRombel, idEkskul, idKelasCatat int) (*ekskulKonteks, *fiber.Error) {
	if idTP <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Tahun pelajaran wajib dipilih.")
	}

	if !isSemesterValid(semester) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Semester wajib dipilih: Ganjil atau Genap.")
	}

	if idRombel <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Rombel wajib dipilih.")
	}

	if idEkskul <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Ekstrakurikuler wajib dipilih.")
	}

	namaTahun, err := GetNamaTahunPelajaran(idTP)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Tahun pelajaran tidak ditemukan.")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	rombel, err := GetRombel(idRombel)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Rombel tidak ditemukan.")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	idKelas, namaKelas, ferr := pilihKelasCatat(idKelasCatat, rombel.IDKelas, rombel.NamaKelas)
	if ferr != nil {
		return nil, ferr
	}

	namaEkskul, err := GetNamaEkskul(idEkskul)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Ekstrakurikuler tidak ditemukan atau nonaktif.")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return &ekskulKonteks{
		IDTahunPelajaran: idTP,
		Semester:         semester,
		IDRombel:         idRombel,
		IDEkskul:         idEkskul,
		IDKelas:          idKelas,
		NamaTahun:        namaTahun,
		NamaKelas:        namaKelas,
		NamaRombel:       rombel.NamaRombel,
		NamaEkskul:       namaEkskul,
	}, nil
}

// GetEkskulOpsiHandler = ekstrakurikuler aktif (dropdown).
func GetEkskulOpsiHandler(c *fiber.Ctx) error {
	data, err := GetEkskulAktif()
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"success": true, "data": data})
}

// GetDaftarEkskulHandler = anggota ekskul di 1 rombel + nilainya.
func GetDaftarEkskulHandler(c *fiber.Ctx) error {
	k, ferr := resolveEkskulKonteks(
		c.QueryInt("id_tahunpelajaran", 0),
		c.Query("semester"),
		c.QueryInt("id_rombel", 0),
		c.QueryInt("id_ekskul", 0),
		c.QueryInt("id_kelas_catat", 0),
	)
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	list, err := GetDaftarEkskul(k.IDTahunPelajaran, k.Semester, k.IDRombel, k.IDEkskul)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"success": true, "data": list})
}

func validasiItemEkskul(item *EkskulBulkItem) string {
	item.Nilai = strings.TrimSpace(item.Nilai)
	item.Keterangan = strings.TrimSpace(item.Keterangan)

	if utf8.RuneCountInString(item.Nilai) > maxNilaiEkskul {
		return fmt.Sprintf("Nilai NIS %s melebihi %d karakter.", item.Nis, maxNilaiEkskul)
	}

	if utf8.RuneCountInString(item.Keterangan) > maxKeteranganEkskul {
		return fmt.Sprintf("Keterangan NIS %s melebihi %d karakter.", item.Nis, maxKeteranganEkskul)
	}

	return ""
}

// BulkEkskulHandler = simpan nilai ekskul 1 rombel + 1 semester + 1 ekskul sekaligus.
func BulkEkskulHandler(c *fiber.Ctx) error {
	var req EkskulBulkRequest

	if err := c.BodyParser(&req); err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "Payload tidak valid.")
	}

	k, ferr := resolveEkskulKonteks(req.IDTahunPelajaran, req.Semester, req.IDRombel, req.IDEkskul, req.IDKelasCatat)
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	if len(req.Items) == 0 {
		return errorJSON(c, fiber.StatusBadRequest, "Tidak ada data nilai yang dikirim.")
	}

	nisRombel, err := GetNISRombel(k.IDRombel)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	seen := make(map[string]bool, len(req.Items))

	for i := range req.Items {
		item := &req.Items[i]
		item.Nis = strings.TrimSpace(item.Nis)

		if item.Nis == "" {
			return errorJSON(c, fiber.StatusBadRequest, "Ada item tanpa NIS.")
		}

		if !nisRombel[item.Nis] {
			return errorJSON(c, fiber.StatusBadRequest, "NIS "+item.Nis+" bukan anggota rombel ini.")
		}

		if seen[item.Nis] {
			return errorJSON(c, fiber.StatusBadRequest, "NIS "+item.Nis+" terkirim lebih dari sekali.")
		}
		seen[item.Nis] = true

		if msg := validasiItemEkskul(item); msg != "" {
			return errorJSON(c, fiber.StatusBadRequest, msg)
		}
	}

	saved, cleared, err := BulkSaveEkskul(k.IDTahunPelajaran, k.Semester, k.IDKelas, k.IDEkskul, req.Items)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal menyimpan nilai ekskul: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Nilai ekstrakurikuler berhasil disimpan.",
		"saved":   saved,
		"cleared": cleared,
	})
}

// TemplateEkskulHandler = template Excel berisi NIS + nama anggota ekskul rombel (dan nilai yang sudah ada).
func TemplateEkskulHandler(c *fiber.Ctx) error {
	k, ferr := resolveEkskulKonteks(
		c.QueryInt("id_tahunpelajaran", 0),
		c.Query("semester"),
		c.QueryInt("id_rombel", 0),
		c.QueryInt("id_ekskul", 0),
		c.QueryInt("id_kelas_catat", 0),
	)
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	list, err := GetDaftarEkskul(k.IDTahunPelajaran, k.Semester, k.IDRombel, k.IDEkskul)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Nilai Ekskul"
	f.SetSheetName("Sheet1", sheet)

	// kolom NIS dipaksa Text supaya nol di depan tidak hilang saat dibuka di Excel
	if textStyle, errStyle := f.NewStyle(&excelize.Style{NumFmt: 49}); errStyle == nil {
		f.SetColStyle(sheet, "A:A", textStyle)
	}

	for i, h := range ekskulExcelHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, s := range list {
		values := []interface{}{s.Nis, s.NamaSiswa, s.Nilai, s.Keterangan}

		for j, v := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue(sheet, cell, v)
		}
	}

	if headerStyle, errStyle := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"2563EB"}, Pattern: 1},
	}); errStyle == nil {
		lastCell, _ := excelize.CoordinatesToCellName(len(ekskulExcelHeaders), 1)
		f.SetCellStyle(sheet, "A1", lastCell, headerStyle)
	}

	f.SetColWidth(sheet, "A", "A", 14)
	f.SetColWidth(sheet, "B", "B", 34)
	f.SetColWidth(sheet, "C", "C", 16)
	f.SetColWidth(sheet, "D", "D", 50)

	f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	f.NewSheet("Petunjuk")
	f.SetCellValue("Petunjuk", "A1", "Konteks file ini (pilihan di layar Nilai Ekstrakurikuler saat mengunduh, jangan diubah):")
	f.SetCellValue("Petunjuk", "A2", "Tahun Pelajaran : "+k.NamaTahun)
	f.SetCellValue("Petunjuk", "A3", "Semester : "+k.Semester)
	f.SetCellValue("Petunjuk", "A4", "Rombel : "+k.NamaRombel+" (Kelas "+k.NamaKelas+")")
	f.SetCellValue("Petunjuk", "A5", "Ekstrakurikuler : "+k.NamaEkskul)
	f.SetCellValue("Petunjuk", "A7", "Aturan pengisian:")
	f.SetCellValue("Petunjuk", "A8", "1. Isi hanya kolom Nilai dan Keterangan. Kolom NIS dan Nama jangan diubah.")
	f.SetCellValue("Petunjuk", "A9", fmt.Sprintf("2. Nilai berupa teks bebas (contoh: A, B, Baik, Sangat Baik), maksimal %d karakter.", maxNilaiEkskul))
	f.SetCellValue("Petunjuk", "A10", fmt.Sprintf("3. Keterangan maksimal %d karakter.", maxKeteranganEkskul))
	f.SetCellValue("Petunjuk", "A11", "4. Baris yang Nilai dan Keterangannya kosong dilewati (nilai lama siswa itu tidak diubah).")
	f.SetCellValue("Petunjuk", "A12", "5. NIS harus anggota rombel yang dipilih dan tidak boleh muncul dua kali.")
	f.SetCellValue("Petunjuk", "A13", "6. Impor bersifat semua-atau-tidak-sama-sekali: kalau ada satu baris salah, tidak ada yang disimpan.")

	if idx, errIdx := f.GetSheetIndex(sheet); errIdx == nil {
		f.SetActiveSheet(idx)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal membuat template.")
	}

	filename := fmt.Sprintf("template-nilai-ekskul-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return c.Send(buf.Bytes())
}

// ImportEkskulHandler = impor nilai ekskul dari Excel (form-data: file + 4 filter). All-or-nothing.
func ImportEkskulHandler(c *fiber.Ctx) error {
	idTP, _ := strconv.Atoi(c.FormValue("id_tahunpelajaran"))
	idRombel, _ := strconv.Atoi(c.FormValue("id_rombel"))
	idEkskul, _ := strconv.Atoi(c.FormValue("id_ekskul"))

	idKelasCatat, _ := strconv.Atoi(c.FormValue("id_kelas_catat"))

	k, ferr := resolveEkskulKonteks(idTP, c.FormValue("semester"), idRombel, idEkskul, idKelasCatat)
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "File excel wajib diunggah.")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "Gagal membaca file excel.")
	}
	defer file.Close()

	xf, err := excelize.OpenReader(file)
	if err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "Format file excel tidak valid.")
	}
	defer xf.Close()

	sheets := xf.GetSheetList()
	if len(sheets) == 0 {
		return errorJSON(c, fiber.StatusBadRequest, "File excel kosong.")
	}

	rows, err := xf.GetRows(sheets[0])
	if err != nil || len(rows) <= 1 {
		return errorJSON(c, fiber.StatusBadRequest, "File excel kosong.")
	}

	nisRombel, err := GetNISRombel(k.IDRombel)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	type errItem struct {
		Row     int    `json:"row"`
		Message string `json:"message"`
	}

	var errs []errItem
	var items []EkskulBulkItem
	seen := map[string]bool{}
	skipped := 0

	for i, row := range rows {
		if i == 0 {
			continue // header
		}

		get := func(idx int) string {
			if idx < len(row) {
				return strings.TrimSpace(row[idx])
			}

			return ""
		}

		nis := get(0)
		nilai := get(2)
		keterangan := get(3)

		if nis == "" && nilai == "" && keterangan == "" {
			continue // baris kosong
		}

		rowNo := i + 1

		if nis == "" {
			errs = append(errs, errItem{rowNo, "NIS kosong."})
			continue
		}

		if seen[nis] {
			errs = append(errs, errItem{rowNo, "NIS (" + nis + ") terdapat data ganda dalam file."})
			continue
		}
		seen[nis] = true

		if !nisRombel[nis] {
			errs = append(errs, errItem{rowNo, "NIS (" + nis + ") bukan anggota rombel " + k.NamaRombel + "."})
			continue
		}

		if nilai == "" && keterangan == "" {
			skipped++
			continue
		}

		item := EkskulBulkItem{Nis: nis, Nilai: nilai, Keterangan: keterangan}
		if msg := validasiItemEkskul(&item); msg != "" {
			errs = append(errs, errItem{rowNo, msg})
			continue
		}

		items = append(items, item)
	}

	if len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Impor gagal, periksa kembali data.",
			"errors":  errs,
		})
	}

	if len(items) == 0 {
		return errorJSON(c, fiber.StatusUnprocessableEntity, "Tidak ada data nilai yang terisi di file excel.")
	}

	saved, _, err := BulkSaveEkskul(k.IDTahunPelajaran, k.Semester, k.IDKelas, k.IDEkskul, items)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal menyimpan nilai ekskul: "+err.Error())
	}

	message := fmt.Sprintf("Impor %d nilai ekstrakurikuler berhasil.", saved)
	if skipped > 0 {
		message += fmt.Sprintf(" %d baris tanpa isian dilewati.", skipped)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": message,
		"saved":   saved,
		"skipped": skipped,
	})
}
