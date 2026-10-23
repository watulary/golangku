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

var ujianExcelHeaders = []string{"NIS", "Nama Peserta Didik", "Nilai Akhir", "Capaian Pembelajaran"}

type ujianKonteks struct {
	IDTahunPelajaran int
	Jenis            string
	IDRombel         int
	IDMapel          int

	IDKelas    int
	NamaTahun  string
	NamaKelas  string
	NamaRombel string
	NamaMapel  string
}

func namaJenisUjian(j string) string {
	if j == UjianDaerah {
		return "Ujian Sekolah Standar Daerah"
	}

	return "Ujian Sekolah"
}

func resolveUjianKonteks(idTP int, jenis string, idRombel, idMapel, idKelasCatat int) (*ujianKonteks, *fiber.Error) {
	if idTP <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Tahun pelajaran wajib dipilih.")
	}

	if !isJenisUjianValid(jenis) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Jenis ujian wajib dipilih: Sekolah atau Standar Daerah.")
	}

	if idRombel <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Rombel wajib dipilih.")
	}

	if idMapel <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Mata pelajaran wajib dipilih.")
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

	namaMapel, err := GetNamaMapelDinilai(idMapel)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fiber.NewError(fiber.StatusBadRequest,
				"Mata pelajaran tidak ditemukan, nonaktif, atau merupakan mapel induk (nilai diisi di sub-mapelnya).")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return &ujianKonteks{
		IDTahunPelajaran: idTP, Jenis: jenis, IDRombel: idRombel, IDMapel: idMapel,
		IDKelas: idKelas, NamaTahun: namaTahun, NamaKelas: namaKelas,
		NamaRombel: rombel.NamaRombel, NamaMapel: namaMapel,
	}, nil
}

func GetDaftarUjianHandler(c *fiber.Ctx) error {
	k, ferr := resolveUjianKonteks(
		c.QueryInt("id_tahunpelajaran", 0), c.Query("jenis"),
		c.QueryInt("id_rombel", 0), c.QueryInt("id_mapel", 0), c.QueryInt("id_kelas_catat", 0),
	)
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	list, err := GetDaftarUjian(k.IDTahunPelajaran, k.Jenis, k.IDRombel, k.IDMapel)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"success": true, "data": list})
}

func BulkUjianHandler(c *fiber.Ctx) error {
	var req UjianBulkRequest

	if err := c.BodyParser(&req); err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "Payload tidak valid.")
	}

	k, ferr := resolveUjianKonteks(req.IDTahunPelajaran, req.Jenis, req.IDRombel, req.IDMapel, req.IDKelasCatat)
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
		it := &req.Items[i]
		it.Nis = strings.TrimSpace(it.Nis)

		if it.Nis == "" {
			return errorJSON(c, fiber.StatusBadRequest, "Ada item tanpa NIS.")
		}

		if !nisRombel[it.Nis] {
			return errorJSON(c, fiber.StatusBadRequest, "NIS "+it.Nis+" bukan anggota rombel ini.")
		}

		if seen[it.Nis] {
			return errorJSON(c, fiber.StatusBadRequest, "NIS "+it.Nis+" terkirim lebih dari sekali.")
		}

		seen[it.Nis] = true

		if it.NilaiAkhir != nil && !validNilai(*it.NilaiAkhir) {
			return errorJSON(c, fiber.StatusBadRequest, "Nilai akhir NIS "+it.Nis+" harus berupa angka 0 sampai 100.")
		}

		if utf8.RuneCountInString(it.Capaian) > maxTeks {
			return errorJSON(c, fiber.StatusBadRequest,
				fmt.Sprintf("Capaian Pembelajaran NIS %s melebihi %d karakter.", it.Nis, maxTeks))
		}
	}

	saved, cleared, err := BulkSaveUjian(k.IDTahunPelajaran, k.Jenis, k.IDKelas, k.IDMapel, req.Items)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal menyimpan nilai ujian: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true, "message": "Nilai ujian berhasil disimpan.", "saved": saved, "cleared": cleared,
	})
}

func TemplateUjianHandler(c *fiber.Ctx) error {
	k, ferr := resolveUjianKonteks(
		c.QueryInt("id_tahunpelajaran", 0), c.Query("jenis"),
		c.QueryInt("id_rombel", 0), c.QueryInt("id_mapel", 0), c.QueryInt("id_kelas_catat", 0),
	)
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	list, err := GetDaftarUjian(k.IDTahunPelajaran, k.Jenis, k.IDRombel, k.IDMapel)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Nilai Ujian"
	f.SetSheetName("Sheet1", sheet)

	if textStyle, errStyle := f.NewStyle(&excelize.Style{NumFmt: 49}); errStyle == nil {
		f.SetColStyle(sheet, "A:A", textStyle)
	}

	for i, h := range ujianExcelHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, s := range list {
		values := []interface{}{s.Nis, s.NamaSiswa, "", s.Capaian}
		if s.NilaiAkhir != nil {
			values[2] = *s.NilaiAkhir
		}

		for j, v := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue(sheet, cell, v)
		}
	}

	if headerStyle, errStyle := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"2563EB"}, Pattern: 1},
	}); errStyle == nil {
		lastCell, _ := excelize.CoordinatesToCellName(len(ujianExcelHeaders), 1)
		f.SetCellStyle(sheet, "A1", lastCell, headerStyle)
	}

	f.SetColWidth(sheet, "A", "A", 14)
	f.SetColWidth(sheet, "B", "B", 34)
	f.SetColWidth(sheet, "C", "C", 12)
	f.SetColWidth(sheet, "D", "D", 70)

	f.NewSheet("Petunjuk")
	f.SetCellValue("Petunjuk", "A1", "Konteks file ini (pilihan di layar Input Nilai Ujian saat mengunduh, jangan diubah):")
	f.SetCellValue("Petunjuk", "A2", "Tahun Pelajaran : "+k.NamaTahun)
	f.SetCellValue("Petunjuk", "A3", "Jenis Ujian : "+namaJenisUjian(k.Jenis))
	f.SetCellValue("Petunjuk", "A4", "Rombel : "+k.NamaRombel+" (Kelas "+k.NamaKelas+")")
	f.SetCellValue("Petunjuk", "A5", "Mata Pelajaran : "+k.NamaMapel)
	f.SetCellValue("Petunjuk", "A7", "Aturan pengisian:")
	f.SetCellValue("Petunjuk", "A8", "1. Isi hanya kolom Nilai Akhir dan Capaian Pembelajaran. Kolom NIS dan Nama jangan diubah.")
	f.SetCellValue("Petunjuk", "A9", "2. Nilai Akhir berupa angka 0 sampai 100 (contoh: 85 atau 85.5).")
	f.SetCellValue("Petunjuk", "A10", fmt.Sprintf("3. Capaian Pembelajaran maksimal %d karakter per sel.", maxTeks))
	f.SetCellValue("Petunjuk", "A11", "4. Baris yang kedua kolom isiannya kosong dilewati (nilai lama siswa itu tidak diubah).")
	f.SetCellValue("Petunjuk", "A12", "5. Impor bersifat semua-atau-tidak-sama-sekali: kalau ada satu baris salah, tidak ada yang disimpan.")

	if idx, errIdx := f.GetSheetIndex(sheet); errIdx == nil {
		f.SetActiveSheet(idx)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal membuat template.")
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="template-nilai-ujian-%s.xlsx"`, time.Now().Format("2006-01-02-15-04-05")))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return c.Send(buf.Bytes())
}

func ImportUjianHandler(c *fiber.Ctx) error {
	idTP, _ := strconv.Atoi(c.FormValue("id_tahunpelajaran"))
	idRombel, _ := strconv.Atoi(c.FormValue("id_rombel"))
	idMapel, _ := strconv.Atoi(c.FormValue("id_mapel"))
	idKelasCatat, _ := strconv.Atoi(c.FormValue("id_kelas_catat"))

	k, ferr := resolveUjianKonteks(idTP, c.FormValue("jenis"), idRombel, idMapel, idKelasCatat)
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

	var items []UjianBulkItem

	seen := map[string]bool{}
	skipped := 0

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

		nis, nilaiStr, capaian := get(0), get(2), get(3)

		if nis == "" && nilaiStr == "" && capaian == "" {
			continue
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

		if nilaiStr == "" && capaian == "" {
			skipped++
			continue
		}

		var nilaiPtr *float64

		if nilaiStr != "" {
			v, errParse := strconv.ParseFloat(strings.ReplaceAll(nilaiStr, ",", "."), 64)
			if errParse != nil || !validNilai(v) {
				errs = append(errs, errItem{rowNo, "Nilai Akhir (" + nilaiStr + ") harus berupa angka 0 sampai 100."})
				continue
			}

			nilaiPtr = &v
		}

		if utf8.RuneCountInString(capaian) > maxTeks {
			errs = append(errs, errItem{rowNo, fmt.Sprintf("Capaian Pembelajaran melebihi %d karakter.", maxTeks)})
			continue
		}

		items = append(items, UjianBulkItem{Nis: nis, NilaiAkhir: nilaiPtr, Capaian: capaian})
	}

	if len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false, "message": "Impor gagal, periksa kembali data.", "errors": errs,
		})
	}

	if len(items) == 0 {
		return errorJSON(c, fiber.StatusUnprocessableEntity, "Tidak ada data nilai yang terisi di file excel.")
	}

	saved, _, err := BulkSaveUjian(k.IDTahunPelajaran, k.Jenis, k.IDKelas, k.IDMapel, items)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal menyimpan nilai ujian: "+err.Error())
	}

	message := fmt.Sprintf("Impor %d data nilai ujian berhasil.", saved)
	if skipped > 0 {
		message += fmt.Sprintf(" %d baris tanpa isian dilewati.", skipped)
	}

	return c.JSON(fiber.Map{"success": true, "message": message, "saved": saved, "skipped": skipped})
}
