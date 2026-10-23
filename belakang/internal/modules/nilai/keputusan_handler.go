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

var keputusanExcelHeaders = []string{
	"NIS", "Nama Peserta Didik", "Status (Lulus / Tidak Lulus)", "No. Ijazah", "No. SKHUS", "Tanggal (dd-mm-yyyy)",
}

type keputusanKonteks struct {
	IDTahunPelajaran int
	IDRombel         int
	NamaTahun        string
	NamaKelas        string
	NamaRombel       string
}

func resolveKeputusanKonteks(idTP, idRombel int) (*keputusanKonteks, *fiber.Error) {
	if idTP <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Tahun pelajaran wajib dipilih.")
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

	return &keputusanKonteks{
		IDTahunPelajaran: idTP, IDRombel: idRombel, NamaTahun: namaTahun,
		NamaKelas: rombel.NamaKelas, NamaRombel: rombel.NamaRombel,
	}, nil
}

// normalisasiStatus: terima "lulus" / "LULUS" / "tidak lulus" (huruf bebas) -> nilai baku; ok=false kalau tak dikenal.
func normalisasiStatus(s string) (string, bool) {
	switch strings.ToLower(strings.Join(strings.Fields(s), " ")) {
	case "":
		return "", true
	case "lulus":
		return StatusLulus, true
	case "tidak lulus":
		return StatusTidakLulus, true
	}

	return "", false
}

func GetDaftarKeputusanHandler(c *fiber.Ctx) error {
	k, ferr := resolveKeputusanKonteks(c.QueryInt("id_tahunpelajaran", 0), c.QueryInt("id_rombel", 0))
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	list, err := GetDaftarKeputusan(k.IDRombel)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"success": true, "data": list})
}

func BulkKeputusanHandler(c *fiber.Ctx) error {
	var req KeputusanBulkRequest

	if err := c.BodyParser(&req); err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "Payload tidak valid.")
	}

	k, ferr := resolveKeputusanKonteks(req.IDTahunPelajaran, req.IDRombel)
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
		it.Nis = strings.TrimSpace(it.Nis)
		it.NoIjazah = strings.TrimSpace(it.NoIjazah)
		it.NoSKHUS = strings.TrimSpace(it.NoSKHUS)

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

		status, ok := normalisasiStatus(it.Status)
		if !ok {
			return errorJSON(c, fiber.StatusBadRequest, "Status NIS "+it.Nis+" harus Lulus atau Tidak Lulus.")
		}

		it.Status = status

		tgl, ok := parseTanggalKeputusan(it.Tanggal)
		if !ok {
			return errorJSON(c, fiber.StatusBadRequest, "Tanggal NIS "+it.Nis+" tidak valid.")
		}

		it.Tanggal = tgl

		if utf8.RuneCountInString(it.NoIjazah) > maxNoDokumen || utf8.RuneCountInString(it.NoSKHUS) > maxNoDokumen {
			return errorJSON(c, fiber.StatusBadRequest,
				fmt.Sprintf("No. Ijazah / SKHUS NIS %s melebihi %d karakter.", it.Nis, maxNoDokumen))
		}
	}

	saved, cleared, err := BulkSaveKeputusan(k.IDTahunPelajaran, req.Items)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal menyimpan keputusan: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true, "message": "Keputusan akhir berhasil disimpan.", "saved": saved, "cleared": cleared,
	})
}

func TemplateKeputusanHandler(c *fiber.Ctx) error {
	k, ferr := resolveKeputusanKonteks(c.QueryInt("id_tahunpelajaran", 0), c.QueryInt("id_rombel", 0))
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	list, err := GetDaftarKeputusan(k.IDRombel)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Keputusan"
	f.SetSheetName("Sheet1", sheet)

	// semua kolom Text: NIS tidak kehilangan nol di depan, tanggal tidak diubah Excel jadi format lain
	if textStyle, errStyle := f.NewStyle(&excelize.Style{NumFmt: 49}); errStyle == nil {
		f.SetColStyle(sheet, "A:F", textStyle)
	}

	for i, h := range keputusanExcelHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, s := range list {
		tgl := ""
		if t, e := time.Parse("2006-01-02", s.Tanggal); e == nil {
			tgl = t.Format("02-01-2006")
		}

		values := []interface{}{s.Nis, s.NamaSiswa, s.Status, s.NoIjazah, s.NoSKHUS, tgl}
		for j, v := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue(sheet, cell, v)
		}
	}

	if headerStyle, errStyle := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"2563EB"}, Pattern: 1},
	}); errStyle == nil {
		lastCell, _ := excelize.CoordinatesToCellName(len(keputusanExcelHeaders), 1)
		f.SetCellStyle(sheet, "A1", lastCell, headerStyle)
	}

	f.SetColWidth(sheet, "A", "A", 14)
	f.SetColWidth(sheet, "B", "B", 34)
	f.SetColWidth(sheet, "C", "C", 28)
	f.SetColWidth(sheet, "D", "E", 26)
	f.SetColWidth(sheet, "F", "F", 20)

	f.NewSheet("Petunjuk")
	f.SetCellValue("Petunjuk", "A1", "Konteks file ini (pilihan di layar Keputusan Akhir saat mengunduh, jangan diubah):")
	f.SetCellValue("Petunjuk", "A2", "Tahun Pelajaran : "+k.NamaTahun)
	f.SetCellValue("Petunjuk", "A3", "Rombel : "+k.NamaRombel+" (Kelas "+k.NamaKelas+")")
	f.SetCellValue("Petunjuk", "A5", "Aturan pengisian:")
	f.SetCellValue("Petunjuk", "A6", "1. Kolom NIS dan Nama jangan diubah.")
	f.SetCellValue("Petunjuk", "A7", "2. Status diisi Lulus atau Tidak Lulus (huruf besar/kecil bebas), boleh dikosongkan.")
	f.SetCellValue("Petunjuk", "A8", fmt.Sprintf("3. No. Ijazah dan No. SKHUS maksimal %d karakter.", maxNoDokumen))
	f.SetCellValue("Petunjuk", "A9", "4. Tanggal berformat dd-mm-yyyy (contoh: 20-06-2027).")
	f.SetCellValue("Petunjuk", "A10", "5. Baris yang semua kolom isiannya kosong dilewati (data lama siswa itu tidak diubah).")
	f.SetCellValue("Petunjuk", "A11", "6. Impor bersifat semua-atau-tidak-sama-sekali: kalau ada satu baris salah, tidak ada yang disimpan.")

	if idx, errIdx := f.GetSheetIndex(sheet); errIdx == nil {
		f.SetActiveSheet(idx)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal membuat template.")
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="template-keputusan-akhir-%s.xlsx"`, time.Now().Format("2006-01-02-15-04-05")))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return c.Send(buf.Bytes())
}

func ImportKeputusanHandler(c *fiber.Ctx) error {
	idTP, _ := strconv.Atoi(c.FormValue("id_tahunpelajaran"))
	idRombel, _ := strconv.Atoi(c.FormValue("id_rombel"))

	k, ferr := resolveKeputusanKonteks(idTP, idRombel)
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

	var items []KeputusanBulkItem

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

		nis, statusStr, ijazah, skhus, tglStr := get(0), get(2), get(3), get(4), get(5)

		if nis == "" && statusStr == "" && ijazah == "" && skhus == "" && tglStr == "" {
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

		if statusStr == "" && ijazah == "" && skhus == "" && tglStr == "" {
			skipped++
			continue
		}

		status, ok := normalisasiStatus(statusStr)
		if !ok {
			errs = append(errs, errItem{rowNo, "Status (" + statusStr + ") harus Lulus atau Tidak Lulus."})
			continue
		}

		tgl, ok := parseTanggalKeputusan(tglStr)
		if !ok {
			errs = append(errs, errItem{rowNo, "Tanggal (" + tglStr + ") harus berformat dd-mm-yyyy."})
			continue
		}

		if utf8.RuneCountInString(ijazah) > maxNoDokumen || utf8.RuneCountInString(skhus) > maxNoDokumen {
			errs = append(errs, errItem{rowNo, fmt.Sprintf("No. Ijazah / SKHUS melebihi %d karakter.", maxNoDokumen)})
			continue
		}

		items = append(items, KeputusanBulkItem{Nis: nis, Status: status, NoIjazah: ijazah, NoSKHUS: skhus, Tanggal: tgl})
	}

	if len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false, "message": "Impor gagal, periksa kembali data.", "errors": errs,
		})
	}

	if len(items) == 0 {
		return errorJSON(c, fiber.StatusUnprocessableEntity, "Tidak ada data yang terisi di file excel.")
	}

	saved, _, err := BulkSaveKeputusan(k.IDTahunPelajaran, items)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal menyimpan keputusan: "+err.Error())
	}

	message := fmt.Sprintf("Impor %d data keputusan berhasil.", saved)
	if skipped > 0 {
		message += fmt.Sprintf(" %d baris tanpa isian dilewati.", skipped)
	}

	return c.JSON(fiber.Map{"success": true, "message": message, "saved": saved, "skipped": skipped})
}
