package nilai

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

var naikKelasExcelHeaders = []string{
	"NIS", "Nama Peserta Didik", "Status (Naik / Tinggal)",
}

type naikKelasKonteks struct {
	IDTahunPelajaran int
	IDRombel         int
	Tingkat          int
	NamaTahun        string
	NamaKelas        string
	NamaRombel       string
}

func resolveNaikKelasKonteks(idTP, idRombel int) (*naikKelasKonteks, *fiber.Error) {
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

	tingkat := TingkatKelas(rombel.NamaKelas)
	if tingkat != 7 && tingkat != 8 {
		return nil, fiber.NewError(fiber.StatusBadRequest,
			"Keputusan Naik Kelas hanya utk Kelas 7 dan 8. Kelas 9 memakai menu Keputusan Akhir (Lulus/Tidak Lulus).")
	}

	return &naikKelasKonteks{
		IDTahunPelajaran: idTP, IDRombel: idRombel, Tingkat: tingkat, NamaTahun: namaTahun,
		NamaKelas: rombel.NamaKelas, NamaRombel: rombel.NamaRombel,
	}, nil
}

// normalisasiStatusNaikKelas: terima "naik"/"NAIK"/"tinggal" (huruf bebas) -> nilai baku; ok=false kalau tak dikenal.
func normalisasiStatusNaikKelas(s string) (string, bool) {
	switch strings.ToLower(strings.Join(strings.Fields(s), " ")) {
	case "":
		return "", true
	case "naik":
		return StatusNaik, true
	case "tinggal":
		return StatusTinggal, true
	}

	return "", false
}

func GetDaftarNaikKelasHandler(c *fiber.Ctx) error {
	k, ferr := resolveNaikKelasKonteks(c.QueryInt("id_tahunpelajaran", 0), c.QueryInt("id_rombel", 0))
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	list, err := GetDaftarNaikKelas(k.IDRombel, k.Tingkat)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"success": true, "data": list, "tingkat": k.Tingkat})
}

func BulkNaikKelasHandler(c *fiber.Ctx) error {
	var req NaikKelasBulkRequest

	if err := c.BodyParser(&req); err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "Payload tidak valid.")
	}

	k, ferr := resolveNaikKelasKonteks(req.IDTahunPelajaran, req.IDRombel)
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

		status, ok := normalisasiStatusNaikKelas(it.Status)
		if !ok {
			return errorJSON(c, fiber.StatusBadRequest, "Status NIS "+it.Nis+" harus Naik atau Tinggal.")
		}

		it.Status = status
	}

	saved, cleared, err := BulkSaveNaikKelas(k.IDTahunPelajaran, k.Tingkat, req.Items)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal menyimpan keputusan: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true, "message": "Keputusan naik kelas berhasil disimpan.", "saved": saved, "cleared": cleared,
	})
}

func TemplateNaikKelasHandler(c *fiber.Ctx) error {
	k, ferr := resolveNaikKelasKonteks(c.QueryInt("id_tahunpelajaran", 0), c.QueryInt("id_rombel", 0))
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	list, err := GetDaftarNaikKelas(k.IDRombel, k.Tingkat)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "NaikKelas"
	f.SetSheetName("Sheet1", sheet)

	// semua kolom Text: NIS tidak kehilangan nol di depan
	if textStyle, errStyle := f.NewStyle(&excelize.Style{NumFmt: 49}); errStyle == nil {
		f.SetColStyle(sheet, "A:C", textStyle)
	}

	for i, h := range naikKelasExcelHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, s := range list {
		values := []interface{}{s.Nis, s.NamaSiswa, s.Status}
		for j, v := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue(sheet, cell, v)
		}
	}

	// dropdown Naik/Tinggal, boleh kosong
	if len(list) > 0 {
		if dv := excelize.NewDataValidation(true); dv != nil {
			dv.Sqref = fmt.Sprintf("C2:C%d", len(list)+1)
			if err := dv.SetDropList([]string{"Naik", "Tinggal"}); err == nil {
				f.AddDataValidation(sheet, dv)
			}
		}
	}

	if headerStyle, errStyle := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"2563EB"}, Pattern: 1},
	}); errStyle == nil {
		lastCell, _ := excelize.CoordinatesToCellName(len(naikKelasExcelHeaders), 1)
		f.SetCellStyle(sheet, "A1", lastCell, headerStyle)
	}

	f.SetColWidth(sheet, "A", "A", 14)
	f.SetColWidth(sheet, "B", "B", 34)
	f.SetColWidth(sheet, "C", "C", 26)

	f.NewSheet("Petunjuk")
	f.SetCellValue("Petunjuk", "A1", "Konteks file ini (pilihan di layar Keputusan Naik Kelas saat mengunduh, jangan diubah):")
	f.SetCellValue("Petunjuk", "A2", "Tahun Pelajaran : "+k.NamaTahun)
	f.SetCellValue("Petunjuk", "A3", "Rombel : "+k.NamaRombel+" (Kelas "+k.NamaKelas+")")
	f.SetCellValue("Petunjuk", "A5", "Aturan pengisian:")
	f.SetCellValue("Petunjuk", "A6", "1. Kolom NIS dan Nama jangan diubah.")
	f.SetCellValue("Petunjuk", "A7", "2. Status diisi Naik atau Tinggal (huruf besar/kecil bebas), boleh dikosongkan.")
	f.SetCellValue("Petunjuk", "A8", "3. Baris yang kolom Status-nya kosong dilewati (data lama siswa itu tidak diubah).")
	f.SetCellValue("Petunjuk", "A9", "4. Impor bersifat semua-atau-tidak-sama-sekali: kalau ada satu baris salah, tidak ada yang disimpan.")

	if idx, errIdx := f.GetSheetIndex(sheet); errIdx == nil {
		f.SetActiveSheet(idx)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal membuat template.")
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="template-naikkelas-%s.xlsx"`, time.Now().Format("2006-01-02-15-04-05")))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return c.Send(buf.Bytes())
}

func ImportNaikKelasHandler(c *fiber.Ctx) error {
	idTP, _ := strconv.Atoi(c.FormValue("id_tahunpelajaran"))
	idRombel, _ := strconv.Atoi(c.FormValue("id_rombel"))

	k, ferr := resolveNaikKelasKonteks(idTP, idRombel)
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

	var items []NaikKelasBulkItem

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

		nis, statusStr := get(0), get(2)

		if nis == "" && statusStr == "" {
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

		if statusStr == "" {
			skipped++
			continue
		}

		status, ok := normalisasiStatusNaikKelas(statusStr)
		if !ok {
			errs = append(errs, errItem{rowNo, "Status (" + statusStr + ") harus Naik atau Tinggal."})
			continue
		}

		items = append(items, NaikKelasBulkItem{Nis: nis, Status: status})
	}

	if len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false, "message": "Impor gagal, periksa kembali data.", "errors": errs,
		})
	}

	if len(items) == 0 {
		return errorJSON(c, fiber.StatusUnprocessableEntity, "Tidak ada data yang terisi di file excel.")
	}

	saved, _, err := BulkSaveNaikKelas(k.IDTahunPelajaran, k.Tingkat, items)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal menyimpan keputusan: "+err.Error())
	}

	message := fmt.Sprintf("Impor %d data keputusan naik kelas berhasil.", saved)
	if skipped > 0 {
		message += fmt.Sprintf(" %d baris tanpa isian dilewati.", skipped)
	}

	return c.JSON(fiber.Map{"success": true, "message": message, "saved": saved, "skipped": skipped})
}
