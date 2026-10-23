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

const maxTeks = 2000 // batas karakter Capaian Pembelajaran / KKTP per sel

var nilaiExcelHeaders = []string{
	"NIS", "Nama Peserta Didik", "Nilai Akhir", "Capaian Pembelajaran", "KKTP",
}

func isSemesterValid(s string) bool {
	return s == "Ganjil" || s == "Genap"
}

func errorJSON(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}

// konteks = hasil validasi 4 filter utama + nama-namanya (utk template Excel).
type konteks struct {
	IDTahunPelajaran int
	Semester         string
	IDRombel         int
	IDMapel          int

	IDKelas    int
	NamaTahun  string
	NamaKelas  string
	NamaRombel string
	NamaMapel  string
}

func resolveKonteks(idTP int, semester string, idRombel, idMapel, idKelasCatat int) (*konteks, *fiber.Error) {
	if idTP <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Tahun pelajaran wajib dipilih.")
	}

	if !isSemesterValid(semester) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Semester wajib dipilih: Ganjil atau Genap.")
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
			return nil, fiber.NewError(
				fiber.StatusBadRequest,
				"Mata pelajaran tidak ditemukan, nonaktif, atau merupakan mapel induk (nilai diisi di sub-mapelnya).",
			)
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return &konteks{
		IDTahunPelajaran: idTP,
		Semester:         semester,
		IDRombel:         idRombel,
		IDMapel:          idMapel,
		IDKelas:          idKelas,
		NamaTahun:        namaTahun,
		NamaKelas:        namaKelas,
		NamaRombel:       rombel.NamaRombel,
		NamaMapel:        namaMapel,
	}, nil
}

// validNilai: 0..100 (juga menolak NaN).
func validNilai(v float64) bool {
	return v >= 0 && v <= 100
}

// GetDaftarHandler = daftar siswa 1 rombel + nilainya.
func GetDaftarHandler(c *fiber.Ctx) error {
	k, ferr := resolveKonteks(
		c.QueryInt("id_tahunpelajaran", 0),
		c.Query("semester"),
		c.QueryInt("id_rombel", 0),
		c.QueryInt("id_mapel", 0),
		c.QueryInt("id_kelas_catat", 0),
	)
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	list, err := GetDaftar(k.IDTahunPelajaran, k.Semester, k.IDRombel, k.IDMapel)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    list,
	})
}

// BulkHandler = simpan nilai 1 rombel + 1 semester + 1 mapel sekaligus.
func BulkHandler(c *fiber.Ctx) error {
	var req BulkRequest

	if err := c.BodyParser(&req); err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "Payload tidak valid.")
	}

	k, ferr := resolveKonteks(req.IDTahunPelajaran, req.Semester, req.IDRombel, req.IDMapel, req.IDKelasCatat)
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

		if item.NilaiAkhir != nil && !validNilai(*item.NilaiAkhir) {
			return errorJSON(c, fiber.StatusBadRequest, "Nilai akhir NIS "+item.Nis+" harus berupa angka 0 sampai 100.")
		}

		if utf8.RuneCountInString(item.Capaian) > maxTeks || utf8.RuneCountInString(item.Kktp) > maxTeks {
			return errorJSON(c, fiber.StatusBadRequest,
				fmt.Sprintf("Capaian Pembelajaran / KKTP NIS %s melebihi %d karakter.", item.Nis, maxTeks))
		}
	}

	saved, cleared, err := BulkSave(k.IDTahunPelajaran, k.Semester, k.IDKelas, k.IDMapel, req.Items)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal menyimpan nilai: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Nilai berhasil disimpan.",
		"saved":   saved,
		"cleared": cleared,
	})
}

// TemplateHandler = template Excel yang SUDAH berisi NIS + nama siswa rombel
// (dan nilai yang sudah tersimpan, kalau ada), tinggal diisi lalu diimpor.
func TemplateHandler(c *fiber.Ctx) error {
	k, ferr := resolveKonteks(
		c.QueryInt("id_tahunpelajaran", 0),
		c.Query("semester"),
		c.QueryInt("id_rombel", 0),
		c.QueryInt("id_mapel", 0),
		c.QueryInt("id_kelas_catat", 0),
	)
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	list, err := GetDaftar(k.IDTahunPelajaran, k.Semester, k.IDRombel, k.IDMapel)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Nilai"
	f.SetSheetName("Sheet1", sheet)

	// kolom NIS dipaksa Text supaya nol di depan tidak hilang saat dibuka di Excel
	if textStyle, errStyle := f.NewStyle(&excelize.Style{NumFmt: 49}); errStyle == nil {
		f.SetColStyle(sheet, "A:A", textStyle)
	}

	for i, h := range nilaiExcelHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, s := range list {
		row := i + 2

		values := []interface{}{s.Nis, s.NamaSiswa, "", s.Capaian, s.Kktp}
		if s.NilaiAkhir != nil {
			values[2] = *s.NilaiAkhir
		}

		for j, v := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheet, cell, v)
		}
	}

	if headerStyle, errStyle := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"2563EB"}, Pattern: 1},
	}); errStyle == nil {
		lastCell, _ := excelize.CoordinatesToCellName(len(nilaiExcelHeaders), 1)
		f.SetCellStyle(sheet, "A1", lastCell, headerStyle)
	}

	f.SetColWidth(sheet, "A", "A", 14)
	f.SetColWidth(sheet, "B", "B", 34)
	f.SetColWidth(sheet, "C", "C", 12)
	f.SetColWidth(sheet, "D", "E", 60)

	f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	// Sheet Petunjuk: konteks + aturan. Import TIDAK membaca sheet ini.
	f.NewSheet("Petunjuk")
	f.SetCellValue("Petunjuk", "A1", "Konteks file ini (pilihan di layar Input Nilai saat mengunduh, jangan diubah):")
	f.SetCellValue("Petunjuk", "A2", "Tahun Pelajaran : "+k.NamaTahun)
	f.SetCellValue("Petunjuk", "A3", "Semester : "+k.Semester)
	f.SetCellValue("Petunjuk", "A4", "Rombel : "+k.NamaRombel+" (Kelas "+k.NamaKelas+")")
	f.SetCellValue("Petunjuk", "A5", "Mata Pelajaran : "+k.NamaMapel)
	f.SetCellValue("Petunjuk", "A7", "Aturan pengisian:")
	f.SetCellValue("Petunjuk", "A8", "1. Isi hanya kolom Nilai Akhir, Capaian Pembelajaran dan KKTP. Kolom NIS dan Nama jangan diubah.")
	f.SetCellValue("Petunjuk", "A9", "2. Nilai Akhir berupa angka 0 sampai 100 (contoh: 85 atau 85.5).")
	f.SetCellValue("Petunjuk", "A10", fmt.Sprintf("3. Capaian Pembelajaran dan KKTP maksimal %d karakter per sel.", maxTeks))
	f.SetCellValue("Petunjuk", "A11", "4. Baris yang ketiga kolom isiannya kosong dilewati (nilai lama siswa itu tidak diubah).")
	f.SetCellValue("Petunjuk", "A12", "5. NIS harus anggota rombel yang dipilih dan tidak boleh muncul dua kali.")
	f.SetCellValue("Petunjuk", "A13", "6. Impor bersifat semua-atau-tidak-sama-sekali: kalau ada satu baris salah, tidak ada yang disimpan.")

	if idx, errIdx := f.GetSheetIndex(sheet); errIdx == nil {
		f.SetActiveSheet(idx)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal membuat template.")
	}

	filename := fmt.Sprintf("template-nilai-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return c.Send(buf.Bytes())
}

// ImportHandler = impor nilai dari Excel (form-data: file + 4 filter). All-or-nothing.
func ImportHandler(c *fiber.Ctx) error {
	idTP, _ := strconv.Atoi(c.FormValue("id_tahunpelajaran"))
	idRombel, _ := strconv.Atoi(c.FormValue("id_rombel"))
	idMapel, _ := strconv.Atoi(c.FormValue("id_mapel"))

	idKelasCatat, _ := strconv.Atoi(c.FormValue("id_kelas_catat"))

	k, ferr := resolveKonteks(idTP, c.FormValue("semester"), idRombel, idMapel, idKelasCatat)
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
	var items []BulkItem
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
		nilaiStr := get(2)
		capaian := get(3)
		kktp := get(4)

		if nis == "" && nilaiStr == "" && capaian == "" && kktp == "" {
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

		// tidak ada isian sama sekali -> dilewati (nilai lama tidak disentuh)
		if nilaiStr == "" && capaian == "" && kktp == "" {
			skipped++
			continue
		}

		var nilaiPtr *float64
		if nilaiStr != "" {
			// terima "85,5" (koma desimal Indonesia) maupun "85.5"
			v, errParse := strconv.ParseFloat(strings.ReplaceAll(nilaiStr, ",", "."), 64)
			if errParse != nil || !validNilai(v) {
				errs = append(errs, errItem{rowNo, "Nilai Akhir (" + nilaiStr + ") harus berupa angka 0 sampai 100."})
				continue
			}

			nilaiPtr = &v
		}

		if utf8.RuneCountInString(capaian) > maxTeks || utf8.RuneCountInString(kktp) > maxTeks {
			errs = append(errs, errItem{rowNo, fmt.Sprintf("Capaian Pembelajaran / KKTP melebihi %d karakter.", maxTeks)})
			continue
		}

		items = append(items, BulkItem{Nis: nis, NilaiAkhir: nilaiPtr, Capaian: capaian, Kktp: kktp})
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

	saved, _, err := BulkSave(k.IDTahunPelajaran, k.Semester, k.IDKelas, k.IDMapel, items)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal menyimpan nilai: "+err.Error())
	}

	message := fmt.Sprintf("Impor %d data nilai berhasil.", saved)
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
