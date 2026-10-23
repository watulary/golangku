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

var p5ExcelHeaders = []string{"NIS", "Nama Peserta Didik", "Nilai (BB / MB / BSH / SB)", "Elemen", "Sub Elemen", "Target Pencapaian"}

type p5Konteks struct {
	IDTahunPelajaran int
	Semester         string
	IDRombel         int
	Dimensi          int

	IDKelas    int
	NamaTahun  string
	NamaKelas  string
	NamaRombel string
}

func resolveP5Konteks(idTP int, semester string, idRombel, dimensi, idKelasCatat int) (*p5Konteks, *fiber.Error) {
	if idTP <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Tahun pelajaran wajib dipilih.")
	}

	if !isSemesterValid(semester) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Semester wajib dipilih: Ganjil atau Genap.")
	}

	if idRombel <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Rombel wajib dipilih.")
	}

	if !isDimensiP5Valid(dimensi) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Dimensi wajib dipilih.")
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

	return &p5Konteks{
		IDTahunPelajaran: idTP, Semester: semester, IDRombel: idRombel, Dimensi: dimensi,
		IDKelas: idKelas, NamaTahun: namaTahun, NamaKelas: namaKelas, NamaRombel: rombel.NamaRombel,
	}, nil
}

// validasiItemP5 merapikan & memeriksa 1 item; kembalian "" kalau valid.
func validasiItemP5(it *P5BulkItem) string {
	it.Nis = strings.TrimSpace(it.Nis)
	it.Nilai = strings.ToUpper(strings.TrimSpace(it.Nilai))
	it.Elemen = strings.TrimSpace(it.Elemen)
	it.SubElemen = strings.TrimSpace(it.SubElemen)
	it.Target = strings.TrimSpace(it.Target)

	switch {
	case !isNilaiP5Valid(it.Nilai):
		return "Nilai harus BB, MB, BSH, atau SB."
	case utf8.RuneCountInString(it.Elemen) > maxElemenP5 || utf8.RuneCountInString(it.SubElemen) > maxElemenP5:
		return fmt.Sprintf("Elemen / Sub Elemen melebihi %d karakter.", maxElemenP5)
	case utf8.RuneCountInString(it.Target) > maxTargetP5:
		return fmt.Sprintf("Target Pencapaian melebihi %d karakter.", maxTargetP5)
	}

	return ""
}

func GetDaftarP5Handler(c *fiber.Ctx) error {
	k, ferr := resolveP5Konteks(
		c.QueryInt("id_tahunpelajaran", 0), c.Query("semester"),
		c.QueryInt("id_rombel", 0), c.QueryInt("dimensi", 0), c.QueryInt("id_kelas_catat", 0),
	)
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	list, err := GetDaftarP5(k.IDTahunPelajaran, k.Semester, k.IDRombel, k.Dimensi)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"success": true, "data": list})
}

func BulkP5Handler(c *fiber.Ctx) error {
	var req P5BulkRequest

	if err := c.BodyParser(&req); err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "Payload tidak valid.")
	}

	k, ferr := resolveP5Konteks(req.IDTahunPelajaran, req.Semester, req.IDRombel, req.Dimensi, req.IDKelasCatat)
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	if len(req.Items) == 0 {
		return errorJSON(c, fiber.StatusBadRequest, "Tidak ada data yang dikirim.")
	}

	nisRombel, err := GetNISRombel(k.IDRombel)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	seen := make(map[string]bool, len(req.Items))

	for i := range req.Items {
		it := &req.Items[i]

		if msg := validasiItemP5(it); msg != "" {
			return errorJSON(c, fiber.StatusBadRequest, "NIS "+it.Nis+": "+msg)
		}

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
	}

	saved, cleared, err := BulkSaveP5(k.IDTahunPelajaran, k.Semester, k.IDKelas, k.Dimensi, req.Items)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal menyimpan data P5: "+err.Error())
	}

	return c.JSON(fiber.Map{"success": true, "message": "Data P5 berhasil disimpan.", "saved": saved, "cleared": cleared})
}

func TemplateP5Handler(c *fiber.Ctx) error {
	k, ferr := resolveP5Konteks(
		c.QueryInt("id_tahunpelajaran", 0), c.Query("semester"),
		c.QueryInt("id_rombel", 0), c.QueryInt("dimensi", 0), c.QueryInt("id_kelas_catat", 0),
	)
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	list, err := GetDaftarP5(k.IDTahunPelajaran, k.Semester, k.IDRombel, k.Dimensi)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "P5"
	f.SetSheetName("Sheet1", sheet)

	if textStyle, errStyle := f.NewStyle(&excelize.Style{NumFmt: 49}); errStyle == nil {
		f.SetColStyle(sheet, "A:A", textStyle)
	}

	for i, h := range p5ExcelHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, s := range list {
		for j, v := range []interface{}{s.Nis, s.NamaSiswa, s.Nilai, s.Elemen, s.SubElemen, s.Target} {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue(sheet, cell, v)
		}
	}

	if headerStyle, errStyle := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"2563EB"}, Pattern: 1},
	}); errStyle == nil {
		lastCell, _ := excelize.CoordinatesToCellName(len(p5ExcelHeaders), 1)
		f.SetCellStyle(sheet, "A1", lastCell, headerStyle)
	}

	f.SetColWidth(sheet, "A", "A", 14)
	f.SetColWidth(sheet, "B", "B", 34)
	f.SetColWidth(sheet, "C", "C", 24)
	f.SetColWidth(sheet, "D", "F", 50)

	f.NewSheet("Petunjuk")
	f.SetCellValue("Petunjuk", "A1", "Konteks file ini (pilihan di layar Projek P5 saat mengunduh, jangan diubah):")
	f.SetCellValue("Petunjuk", "A2", "Tahun Pelajaran : "+k.NamaTahun)
	f.SetCellValue("Petunjuk", "A3", "Semester : "+k.Semester)
	f.SetCellValue("Petunjuk", "A4", "Rombel : "+k.NamaRombel+" (dicatat sebagai kelas "+k.NamaKelas+")")
	f.SetCellValue("Petunjuk", "A5", fmt.Sprintf("Dimensi : %d. %s", k.Dimensi, DimensiP5[k.Dimensi-1]))
	f.SetCellValue("Petunjuk", "A7", "Aturan pengisian:")
	f.SetCellValue("Petunjuk", "A8", "1. Kolom NIS dan Nama jangan diubah.")
	f.SetCellValue("Petunjuk", "A9", "2. Nilai diisi BB, MB, BSH, atau SB (huruf besar/kecil bebas).")
	f.SetCellValue("Petunjuk", "A10", fmt.Sprintf("3. Elemen dan Sub Elemen maksimal %d karakter, Target Pencapaian maksimal %d karakter.", maxElemenP5, maxTargetP5))
	f.SetCellValue("Petunjuk", "A11", "4. Baris yang semua kolom isiannya kosong dilewati (data lama siswa itu tidak diubah).")
	f.SetCellValue("Petunjuk", "A12", "5. Impor bersifat semua-atau-tidak-sama-sekali: kalau ada satu baris salah, tidak ada yang disimpan.")

	if idx, errIdx := f.GetSheetIndex(sheet); errIdx == nil {
		f.SetActiveSheet(idx)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal membuat template.")
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="template-p5-%s.xlsx"`, time.Now().Format("2006-01-02-15-04-05")))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return c.Send(buf.Bytes())
}

func ImportP5Handler(c *fiber.Ctx) error {
	idTP, _ := strconv.Atoi(c.FormValue("id_tahunpelajaran"))
	idRombel, _ := strconv.Atoi(c.FormValue("id_rombel"))
	dimensi, _ := strconv.Atoi(c.FormValue("dimensi"))
	idKelasCatat, _ := strconv.Atoi(c.FormValue("id_kelas_catat"))

	k, ferr := resolveP5Konteks(idTP, c.FormValue("semester"), idRombel, dimensi, idKelasCatat)
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

	var items []P5BulkItem

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

		it := P5BulkItem{Nis: get(0), Nilai: get(2), Elemen: get(3), SubElemen: get(4), Target: get(5)}

		if it.Nis == "" && it.Nilai == "" && it.Elemen == "" && it.SubElemen == "" && it.Target == "" {
			continue
		}

		rowNo := i + 1

		if it.Nis == "" {
			errs = append(errs, errItem{rowNo, "NIS kosong."})
			continue
		}

		if seen[it.Nis] {
			errs = append(errs, errItem{rowNo, "NIS (" + it.Nis + ") terdapat data ganda dalam file."})
			continue
		}

		seen[it.Nis] = true

		if !nisRombel[it.Nis] {
			errs = append(errs, errItem{rowNo, "NIS (" + it.Nis + ") bukan anggota rombel " + k.NamaRombel + "."})
			continue
		}

		if it.Nilai == "" && it.Elemen == "" && it.SubElemen == "" && it.Target == "" {
			skipped++
			continue
		}

		if msg := validasiItemP5(&it); msg != "" {
			errs = append(errs, errItem{rowNo, msg})
			continue
		}

		items = append(items, it)
	}

	if len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false, "message": "Impor gagal, periksa kembali data.", "errors": errs,
		})
	}

	if len(items) == 0 {
		return errorJSON(c, fiber.StatusUnprocessableEntity, "Tidak ada data yang terisi di file excel.")
	}

	saved, _, err := BulkSaveP5(k.IDTahunPelajaran, k.Semester, k.IDKelas, k.Dimensi, items)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal menyimpan data P5: "+err.Error())
	}

	message := fmt.Sprintf("Impor %d data P5 berhasil.", saved)
	if skipped > 0 {
		message += fmt.Sprintf(" %d baris tanpa isian dilewati.", skipped)
	}

	return c.JSON(fiber.Map{"success": true, "message": message, "saved": saved, "skipped": skipped})
}
