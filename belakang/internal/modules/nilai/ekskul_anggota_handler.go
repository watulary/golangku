package nilai

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

var keanggotaanExcelHeaders = []string{"NIS", "Nama Peserta Didik", "Anggota"}

// normalisasiAnggota: "Ya"/"Tidak" (case-insensitive, spasi diabaikan) -> true/false.
// "" dianggap sama seperti "Tidak" (bukan anggota) supaya template kosong tetap bisa diimpor.
func normalisasiAnggota(nilai string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(nilai)) {
	case "", "tidak":
		return false, true
	case "ya":
		return true, true
	default:
		return false, false
	}
}

// resolveKeanggotaanKonteks memakai validasi yang sama dengan resolveEkskulKonteks
// (tahun pelajaran, semester, rombel, ekskul wajib & harus ada), tanpa opsi kelas-catat
// karena keanggotaan per periode tidak mengenal "dicatat sebagai kelas".
func resolveKeanggotaanKonteks(idTP int, semester string, idRombel, idEkskul int) (*ekskulKonteks, *fiber.Error) {
	return resolveEkskulKonteks(idTP, semester, idRombel, idEkskul, 0)
}

// GetDaftarKeanggotaanHandler = siswa 1 rombel + status keanggotaan ekskul pada periode ini.
func GetDaftarKeanggotaanHandler(c *fiber.Ctx) error {
	k, ferr := resolveKeanggotaanKonteks(
		c.QueryInt("id_tahunpelajaran", 0),
		c.Query("semester"),
		c.QueryInt("id_rombel", 0),
		c.QueryInt("id_ekskul", 0),
	)
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	list, err := GetDaftarKeanggotaan(k.IDTahunPelajaran, k.Semester, k.IDRombel, k.IDEkskul)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"success": true, "data": list})
}

// SimpanKeanggotaanHandler = simpan keanggotaan ekskul 1 rombel + 1 periode + 1 ekskul sekaligus.
func SimpanKeanggotaanHandler(c *fiber.Ctx) error {
	var req KeanggotaanSimpanRequest

	if err := c.BodyParser(&req); err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "Payload tidak valid.")
	}

	k, ferr := resolveKeanggotaanKonteks(req.IDTahunPelajaran, req.Semester, req.IDRombel, req.IDEkskul)
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	nisRombelMap, err := GetNISRombel(k.IDRombel)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	nisRombel := make([]string, 0, len(nisRombelMap))
	for nis := range nisRombelMap {
		nisRombel = append(nisRombel, nis)
	}

	seen := make(map[string]bool, len(req.NisAnggota))
	nisAnggota := make([]string, 0, len(req.NisAnggota))

	for _, nis := range req.NisAnggota {
		nis = strings.TrimSpace(nis)
		if nis == "" {
			continue
		}

		if !nisRombelMap[nis] {
			return errorJSON(c, fiber.StatusBadRequest, "NIS "+nis+" bukan anggota rombel ini.")
		}

		if seen[nis] {
			continue
		}
		seen[nis] = true

		nisAnggota = append(nisAnggota, nis)
	}

	if err := SaveKeanggotaanEkskul(k.IDTahunPelajaran, k.Semester, k.IDEkskul, nisRombel, nisAnggota); err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal menyimpan keanggotaan ekskul: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Keanggotaan ekstrakurikuler berhasil disimpan.",
	})
}

// TemplateKeanggotaanHandler = template Excel berisi NIS + nama siswa 1 rombel,
// kolom Anggota sudah terisi Ya/Tidak sesuai status keanggotaan periode ini saat ini.
func TemplateKeanggotaanHandler(c *fiber.Ctx) error {
	k, ferr := resolveKeanggotaanKonteks(
		c.QueryInt("id_tahunpelajaran", 0),
		c.Query("semester"),
		c.QueryInt("id_rombel", 0),
		c.QueryInt("id_ekskul", 0),
	)
	if ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	list, err := GetDaftarKeanggotaan(k.IDTahunPelajaran, k.Semester, k.IDRombel, k.IDEkskul)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Keanggotaan Ekskul"
	f.SetSheetName("Sheet1", sheet)

	// kolom NIS dipaksa Text supaya nol di depan tidak hilang saat dibuka di Excel
	if textStyle, errStyle := f.NewStyle(&excelize.Style{NumFmt: 49}); errStyle == nil {
		f.SetColStyle(sheet, "A:A", textStyle)
	}

	for i, h := range keanggotaanExcelHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, s := range list {
		anggota := "Tidak"
		if s.Anggota {
			anggota = "Ya"
		}

		values := []interface{}{s.Nis, s.NamaSiswa, anggota}

		for j, v := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue(sheet, cell, v)
		}
	}

	if headerStyle, errStyle := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"2563EB"}, Pattern: 1},
	}); errStyle == nil {
		lastCell, _ := excelize.CoordinatesToCellName(len(keanggotaanExcelHeaders), 1)
		f.SetCellStyle(sheet, "A1", lastCell, headerStyle)
	}

	f.SetColWidth(sheet, "A", "A", 14)
	f.SetColWidth(sheet, "B", "B", 34)
	f.SetColWidth(sheet, "C", "C", 14)

	if dv := excelize.NewDataValidation(true); dv != nil {
		dv.Sqref = fmt.Sprintf("C2:C%d", len(list)+1)
		if err := dv.SetDropList([]string{"Ya", "Tidak"}); err == nil {
			f.AddDataValidation(sheet, dv)
		}
	}

	f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	f.NewSheet("Petunjuk")
	f.SetCellValue("Petunjuk", "A1", "Konteks file ini (pilihan di layar Keanggotaan Ekstrakurikuler saat mengunduh, jangan diubah):")
	f.SetCellValue("Petunjuk", "A2", "Tahun Pelajaran : "+k.NamaTahun)
	f.SetCellValue("Petunjuk", "A3", "Semester : "+k.Semester)
	f.SetCellValue("Petunjuk", "A4", "Rombel : "+k.NamaRombel+" (Kelas "+k.NamaKelas+")")
	f.SetCellValue("Petunjuk", "A5", "Ekstrakurikuler : "+k.NamaEkskul)
	f.SetCellValue("Petunjuk", "A7", "Aturan pengisian:")
	f.SetCellValue("Petunjuk", "A8", "1. Isi hanya kolom Anggota. Kolom NIS dan Nama jangan diubah.")
	f.SetCellValue("Petunjuk", "A9", "2. Anggota diisi Ya (anggota periode ini) atau Tidak (bukan anggota periode ini).")
	f.SetCellValue("Petunjuk", "A10", "3. Kolom Anggota yang dikosongkan dianggap Tidak.")
	f.SetCellValue("Petunjuk", "A11", "4. NIS harus anggota rombel yang dipilih dan tidak boleh muncul dua kali.")
	f.SetCellValue("Petunjuk", "A12", "5. Impor MENGGANTI seluruh keanggotaan ekskul ini pada rombel & periode ini sesuai isian file (bukan menambah).")
	f.SetCellValue("Petunjuk", "A13", "6. Impor bersifat semua-atau-tidak-sama-sekali: kalau ada satu baris salah, tidak ada yang disimpan.")

	if idx, errIdx := f.GetSheetIndex(sheet); errIdx == nil {
		f.SetActiveSheet(idx)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal membuat template.")
	}

	filename := fmt.Sprintf("template-keanggotaan-ekskul-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return c.Send(buf.Bytes())
}

// ImportKeanggotaanHandler = impor keanggotaan ekskul dari Excel (form-data: file + filter).
// MENGGANTI seluruh keanggotaan rombel & periode & ekskul ini sesuai isian file. All-or-nothing.
func ImportKeanggotaanHandler(c *fiber.Ctx) error {
	idTP, _ := strconv.Atoi(c.FormValue("id_tahunpelajaran"))
	idRombel, _ := strconv.Atoi(c.FormValue("id_rombel"))
	idEkskul, _ := strconv.Atoi(c.FormValue("id_ekskul"))

	k, ferr := resolveKeanggotaanKonteks(idTP, c.FormValue("semester"), idRombel, idEkskul)
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

	nisRombelMap, err := GetNISRombel(k.IDRombel)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	nisRombel := make([]string, 0, len(nisRombelMap))
	for nis := range nisRombelMap {
		nisRombel = append(nisRombel, nis)
	}

	type errItem struct {
		Row     int    `json:"row"`
		Message string `json:"message"`
	}

	var errs []errItem
	var nisAnggota []string
	seen := map[string]bool{}

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
		anggotaTeks := get(2)

		if nis == "" && anggotaTeks == "" {
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

		if !nisRombelMap[nis] {
			errs = append(errs, errItem{rowNo, "NIS (" + nis + ") bukan anggota rombel " + k.NamaRombel + "."})
			continue
		}

		anggota, valid := normalisasiAnggota(anggotaTeks)
		if !valid {
			errs = append(errs, errItem{rowNo, "Anggota NIS (" + nis + ") harus diisi Ya atau Tidak."})
			continue
		}

		if anggota {
			nisAnggota = append(nisAnggota, nis)
		}
	}

	if len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Impor gagal, periksa kembali data.",
			"errors":  errs,
		})
	}

	if err := SaveKeanggotaanEkskul(k.IDTahunPelajaran, k.Semester, k.IDEkskul, nisRombel, nisAnggota); err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal menyimpan keanggotaan ekskul: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Impor keanggotaan berhasil, %d siswa tercatat sebagai anggota.", len(nisAnggota)),
		"saved":   len(nisAnggota),
	})
}
