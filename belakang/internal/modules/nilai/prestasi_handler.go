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

var prestasiExcelHeaders = []string{"NIS", "Nama Peserta Didik", "Prestasi", "Nilai", "Keterangan"}

type prestasiKonteks struct {
	IDTahunPelajaran int
	Semester         string
	IDRombel         int

	IDKelas    int
	NamaTahun  string
	NamaKelas  string
	NamaRombel string
}

func resolvePrestasiKonteks(idTP int, semester string, idRombel, idKelasCatat int) (*prestasiKonteks, *fiber.Error) {
	if idTP <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Tahun pelajaran wajib dipilih.")
	}

	if !isSemesterValid(semester) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Semester wajib dipilih: Ganjil atau Genap.")
	}

	if idRombel <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Rombel wajib dipilih.")
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

	return &prestasiKonteks{
		IDTahunPelajaran: idTP,
		Semester:         semester,
		IDRombel:         idRombel,
		IDKelas:          idKelas,
		NamaTahun:        namaTahun,
		NamaKelas:        namaKelas,
		NamaRombel:       rombel.NamaRombel,
	}, nil
}

func validasiItemPrestasi(it *PrestasiBulkItem) string {
	it.Prestasi = strings.TrimSpace(it.Prestasi)
	it.Nilai = strings.TrimSpace(it.Nilai)
	it.Keterangan = strings.TrimSpace(it.Keterangan)

	switch {
	case it.Prestasi == "":
		return "Prestasi NIS " + it.Nis + " belum diisi."
	case utf8.RuneCountInString(it.Prestasi) > maxPrestasi:
		return fmt.Sprintf("Prestasi NIS %s melebihi %d karakter.", it.Nis, maxPrestasi)
	case utf8.RuneCountInString(it.Nilai) > maxNilaiPrestasi:
		return fmt.Sprintf("Nilai NIS %s melebihi %d karakter.", it.Nis, maxNilaiPrestasi)
	case utf8.RuneCountInString(it.Keterangan) > maxKeteranganPrestasi:
		return fmt.Sprintf("Keterangan NIS %s melebihi %d karakter.", it.Nis, maxKeteranganPrestasi)
	}

	return ""
}

func nisList(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for nis := range m {
		out = append(out, nis)
	}

	return out
}

// GetDaftarPrestasiHandler = siswa rombel (untuk dropdown) + prestasi yang sudah tercatat.
func GetDaftarPrestasiHandler(c *fiber.Ctx) error {
	k, ferr := resolvePrestasiKonteks(c.QueryInt("id_tahunpelajaran", 0), c.Query("semester"), c.QueryInt("id_rombel", 0), c.QueryInt("id_kelas_catat", 0))
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	siswa, err := GetSiswaRombel(k.IDRombel)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	data, err := GetDaftarPrestasi(k.IDTahunPelajaran, k.Semester, k.IDRombel)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"success": true, "data": data, "siswa": siswa})
}

// BulkPrestasiHandler = simpan (ganti seluruh) prestasi 1 rombel + 1 semester.
func BulkPrestasiHandler(c *fiber.Ctx) error {
	var req PrestasiBulkRequest

	if err := c.BodyParser(&req); err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "Payload tidak valid.")
	}

	k, ferr := resolvePrestasiKonteks(req.IDTahunPelajaran, req.Semester, req.IDRombel, req.IDKelasCatat)
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	nisRombel, err := GetNISRombel(k.IDRombel)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	seen := map[string]bool{}

	for i := range req.Items {
		it := &req.Items[i]
		it.Nis = strings.TrimSpace(it.Nis)

		if it.Nis == "" {
			return errorJSON(c, fiber.StatusBadRequest, "Ada baris prestasi yang belum dipilih siswanya.")
		}

		if !nisRombel[it.Nis] {
			return errorJSON(c, fiber.StatusBadRequest, "NIS "+it.Nis+" bukan anggota rombel ini.")
		}

		if msg := validasiItemPrestasi(it); msg != "" {
			return errorJSON(c, fiber.StatusBadRequest, msg)
		}

		key := it.Nis + "|" + strings.ToLower(it.Prestasi)
		if seen[key] {
			return errorJSON(c, fiber.StatusBadRequest, "Prestasi \""+it.Prestasi+"\" untuk NIS "+it.Nis+" terisi dua kali.")
		}
		seen[key] = true
	}

	saved, err := SavePrestasi(k.IDTahunPelajaran, k.Semester, k.IDKelas, nisList(nisRombel), req.Items)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal menyimpan prestasi: "+err.Error())
	}

	return c.JSON(fiber.Map{"success": true, "message": "Catatan prestasi berhasil disimpan.", "saved": saved})
}

// TemplatePrestasiHandler = template Excel berisi prestasi yang sudah ada (atau daftar siswa kalau belum ada).
func TemplatePrestasiHandler(c *fiber.Ctx) error {
	k, ferr := resolvePrestasiKonteks(c.QueryInt("id_tahunpelajaran", 0), c.Query("semester"), c.QueryInt("id_rombel", 0), c.QueryInt("id_kelas_catat", 0))
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	data, err := GetDaftarPrestasi(k.IDTahunPelajaran, k.Semester, k.IDRombel)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	if len(data) == 0 {
		siswa, err := GetSiswaRombel(k.IDRombel)
		if err != nil {
			return errorJSON(c, fiber.StatusInternalServerError, err.Error())
		}

		for _, s := range siswa {
			data = append(data, PrestasiRow{Nis: s.Nis, NamaSiswa: s.NamaSiswa})
		}
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Prestasi"
	f.SetSheetName("Sheet1", sheet)

	if textStyle, errStyle := f.NewStyle(&excelize.Style{NumFmt: 49}); errStyle == nil {
		f.SetColStyle(sheet, "A:A", textStyle)
	}

	for i, h := range prestasiExcelHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, r := range data {
		values := []interface{}{r.Nis, r.NamaSiswa, r.Prestasi, r.Nilai, r.Keterangan}

		for j, v := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue(sheet, cell, v)
		}
	}

	if headerStyle, errStyle := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"2563EB"}, Pattern: 1},
	}); errStyle == nil {
		lastCell, _ := excelize.CoordinatesToCellName(len(prestasiExcelHeaders), 1)
		f.SetCellStyle(sheet, "A1", lastCell, headerStyle)
	}

	f.SetColWidth(sheet, "A", "A", 14)
	f.SetColWidth(sheet, "B", "B", 34)
	f.SetColWidth(sheet, "C", "C", 50)
	f.SetColWidth(sheet, "D", "D", 14)
	f.SetColWidth(sheet, "E", "E", 40)

	f.SetPanes(sheet, &excelize.Panes{
		Freeze: true, Split: false, XSplit: 0, YSplit: 1, TopLeftCell: "A2", ActivePane: "bottomLeft",
	})

	f.NewSheet("Petunjuk")
	f.SetCellValue("Petunjuk", "A1", "Konteks file ini (pilihan di layar Catatan Prestasi saat mengunduh, jangan diubah):")
	f.SetCellValue("Petunjuk", "A2", "Tahun Pelajaran : "+k.NamaTahun)
	f.SetCellValue("Petunjuk", "A3", "Semester : "+k.Semester)
	f.SetCellValue("Petunjuk", "A4", "Rombel : "+k.NamaRombel+" (Kelas "+k.NamaKelas+")")
	f.SetCellValue("Petunjuk", "A6", "Aturan pengisian:")
	f.SetCellValue("Petunjuk", "A7", "1. Satu baris = satu prestasi. Siswa dengan banyak prestasi ditulis di beberapa baris dengan NIS yang sama.")
	f.SetCellValue("Petunjuk", "A8", "2. Kolom Prestasi wajib diisi kalau Nilai/Keterangan diisi. Baris tanpa Prestasi, Nilai, dan Keterangan dilewati.")
	f.SetCellValue("Petunjuk", "A9", fmt.Sprintf("3. Prestasi maks. %d, Nilai maks. %d, Keterangan maks. %d karakter.", maxPrestasi, maxNilaiPrestasi, maxKeteranganPrestasi))
	f.SetCellValue("Petunjuk", "A10", "4. NIS harus anggota rombel yang dipilih.")
	f.SetCellValue("Petunjuk", "A11", "5. PENTING: impor MENGGANTI seluruh prestasi rombel ini pada tahun pelajaran dan semester yang dipilih dengan isi file. Unduh template dulu supaya prestasi lama ikut terbawa.")
	f.SetCellValue("Petunjuk", "A12", "6. Impor semua-atau-tidak-sama-sekali: kalau ada satu baris salah, tidak ada yang disimpan.")

	if idx, errIdx := f.GetSheetIndex(sheet); errIdx == nil {
		f.SetActiveSheet(idx)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal membuat template.")
	}

	filename := fmt.Sprintf("template-prestasi-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return c.Send(buf.Bytes())
}

// ImportPrestasiHandler = impor Excel (form-data: file + 3 filter). Mengganti seluruh prestasi periode itu.
func ImportPrestasiHandler(c *fiber.Ctx) error {
	idTP, _ := strconv.Atoi(c.FormValue("id_tahunpelajaran"))
	idRombel, _ := strconv.Atoi(c.FormValue("id_rombel"))

	idKelasCatat, _ := strconv.Atoi(c.FormValue("id_kelas_catat"))

	k, ferr := resolvePrestasiKonteks(idTP, c.FormValue("semester"), idRombel, idKelasCatat)
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
	var items []PrestasiBulkItem
	seen := map[string]bool{}

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

		nis, prestasi, nilai, ket := get(0), get(2), get(3), get(4)

		if prestasi == "" && nilai == "" && ket == "" {
			continue // baris kosong / siswa tanpa prestasi
		}

		rowNo := i + 1

		if nis == "" {
			errs = append(errs, errItem{rowNo, "NIS kosong."})
			continue
		}

		if !nisRombel[nis] {
			errs = append(errs, errItem{rowNo, "NIS (" + nis + ") bukan anggota rombel " + k.NamaRombel + "."})
			continue
		}

		it := PrestasiBulkItem{Nis: nis, Prestasi: prestasi, Nilai: nilai, Keterangan: ket}
		if msg := validasiItemPrestasi(&it); msg != "" {
			errs = append(errs, errItem{rowNo, msg})
			continue
		}

		key := nis + "|" + strings.ToLower(it.Prestasi)
		if seen[key] {
			errs = append(errs, errItem{rowNo, "Prestasi \"" + it.Prestasi + "\" untuk NIS " + nis + " terdapat data ganda dalam file."})
			continue
		}
		seen[key] = true

		items = append(items, it)
	}

	if len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Impor gagal, periksa kembali data.",
			"errors":  errs,
		})
	}

	if len(items) == 0 {
		return errorJSON(c, fiber.StatusUnprocessableEntity, "Tidak ada prestasi yang terisi di file excel.")
	}

	saved, err := SavePrestasi(k.IDTahunPelajaran, k.Semester, k.IDKelas, nisList(nisRombel), items)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal menyimpan prestasi: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Impor %d catatan prestasi berhasil.", saved),
		"saved":   saved,
	})
}
