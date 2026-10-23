package spmb

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

var spmbExcelHeaders = []string{
	"no_peserta", "nama_peserta", "kelamin", "asal_sekolah", "jalur", "lapor_diri",
}

func GetHandler(c *fiber.Ctx) error {
	data, err := Get()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func GetAdminListHandler(c *fiber.Ctx) error {
	search := c.Query("search", "")
	jalur, _ := strconv.Atoi(c.Query("jalur", "0"))
	laporDiri, err := strconv.Atoi(c.Query("lapor_diri", "-1"))
	if err != nil {
		laporDiri = -1
	}
	sortBy := c.Query("sort_by", "nama_peserta")
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

	data, err := GetAllAdmin(search, jalur, laporDiri, sortBy, sortDir, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	total, err := CountAdmin(search, jalur, laporDiri)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func GetByIDHandler(c *fiber.Ctx) error {
	noPeserta := c.Params("id")
	data, err := GetByID(noPeserta)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak ditemukan",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func CreateHandler(c *fiber.Ctx) error {
	var data SPMB
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak valid",
		})
	}

	if err := Create(data); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data peserta berhasil ditambahkan",
	})
}

// UpdateHandler membandingkan data lama vs data baru dulu (cek changed)
// sebelum melakukan UPDATE ke DB. Kalau tidak ada perubahan, DB tidak
// disentuh dan response mengembalikan changed: [] supaya frontend bisa
// menampilkan toast "Tidak ada perubahan data."
func UpdateHandler(c *fiber.Ctx) error {
	noPeserta := c.Params("id")

	existing, err := GetByID(noPeserta)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak ditemukan",
		})
	}

	var data SPMB
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak valid",
		})
	}

	var changed []string
	if existing.NamaPeserta != data.NamaPeserta {
		changed = append(changed, "nama_peserta")
	}
	if existing.Kelamin != data.Kelamin {
		changed = append(changed, "kelamin")
	}
	if existing.AsalSekolah != data.AsalSekolah {
		changed = append(changed, "asal_sekolah")
	}
	if existing.Jalur != data.Jalur {
		changed = append(changed, "jalur")
	}
	if existing.LaporDiri != data.LaporDiri {
		changed = append(changed, "lapor_diri")
	}

	if len(changed) == 0 {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Tidak ada perubahan data",
			"changed": []string{},
		})
	}

	if err := Update(noPeserta, data); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data peserta berhasil diperbarui",
		"changed": changed,
	})
}

func DeleteHandler(c *fiber.Ctx) error {
	noPeserta := c.Params("id")

	if err := Delete(noPeserta); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data peserta berhasil dihapus",
	})
}

func TemplateHandler(c *fiber.Ctx) error {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "SPMB"
	f.SetSheetName("Sheet1", sheet)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"2563EB"}, Pattern: 1},
	})

	for i, h := range spmbExcelHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	petunjuk := "Petunjuk"
	f.NewSheet(petunjuk)
	rows := [][]interface{}{
		{"PETUNJUK PENGISIAN TEMPLATE SPMB"},
		{},
		{"Kolom", "Keterangan"},
		{"no_peserta", "Nomor peserta (unik, wajib diisi)"},
		{"nama_peserta", "Nama lengkap peserta"},
		{"kelamin", "L = Laki-laki, P = Perempuan"},
		{"asal_sekolah", "Nama sekolah asal"},
		{"jalur", "1 = Prestasi Akademik"},
		{"", "2 = Prestasi Non Akademik"},
		{"", "3 = Penyandang Disabilitas"},
		{"", "4 = KJP Plus"},
		{"", "5 = Domisili"},
		{"", "6 = Mutasi"},
		{"", "7 = Tahap Kedua"},
		{"lapor_diri", "0 = Belum Lapor"},
		{"", "1 = Sudah Lapor"},
		{"", "2 = Undur Diri"},
	}
	for i, row := range rows {
		cell, _ := excelize.CoordinatesToCellName(1, i+1)
		f.SetSheetRow(petunjuk, cell, &row)
	}

	f.SetActiveSheet(0)

	filename := fmt.Sprintf("template-spmb-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return f.Write(c.Response().BodyWriter())
}

func ImportHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "File tidak ditemukan",
		})
	}

	fileContent, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuka file",
		})
	}
	defer fileContent.Close()

	f, err := excelize.OpenReader(fileContent)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "File tidak valid",
		})
	}
	defer f.Close()

	rows, err := f.GetRows("SPMB")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Sheet 'SPMB' tidak ditemukan",
		})
	}

	type importError struct {
		Row     int    `json:"row"`
		Message string `json:"message"`
	}

	var errorsList []importError
	var validData []SPMB
	seenInFile := make(map[string]bool)
	var noPesertaList []string

	validKelamin := map[string]bool{"L": true, "P": true}
	validJalur := map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true}
	validLaporDiri := map[int]bool{0: true, 1: true, 2: true}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		rowNum := i + 1
		if len(row) < 6 || strings.TrimSpace(row[0]) == "" {
			continue
		}

		noPeserta := strings.TrimSpace(row[0])
		namaPeserta := strings.TrimSpace(row[1])
		kelamin := strings.ToUpper(strings.TrimSpace(row[2]))
		asalSekolah := strings.TrimSpace(row[3])
		jalur, jalurErr := strconv.Atoi(strings.TrimSpace(row[4]))
		laporDiri, laporErr := strconv.Atoi(strings.TrimSpace(row[5]))

		switch {
		case noPeserta == "":
			errorsList = append(errorsList, importError{rowNum, "no_peserta wajib diisi"})
		case namaPeserta == "":
			errorsList = append(errorsList, importError{rowNum, "nama_peserta wajib diisi"})
		case !validKelamin[kelamin]:
			errorsList = append(errorsList, importError{rowNum, "kelamin harus L atau P"})
		case asalSekolah == "":
			errorsList = append(errorsList, importError{rowNum, "asal_sekolah wajib diisi"})
		case jalurErr != nil || !validJalur[jalur]:
			errorsList = append(errorsList, importError{rowNum, "jalur harus angka 1-7"})
		case laporErr != nil || !validLaporDiri[laporDiri]:
			errorsList = append(errorsList, importError{rowNum, "lapor_diri harus angka 0-2"})
		case seenInFile[noPeserta]:
			errorsList = append(errorsList, importError{rowNum, "no_peserta duplikat di dalam file"})
		default:
			seenInFile[noPeserta] = true
			noPesertaList = append(noPesertaList, noPeserta)
			validData = append(validData, SPMB{
				NoPeserta:   noPeserta,
				NamaPeserta: namaPeserta,
				Kelamin:     kelamin,
				AsalSekolah: asalSekolah,
				Jalur:       jalur,
				LaporDiri:   laporDiri,
			})
		}
	}

	if len(errorsList) == 0 && len(noPesertaList) > 0 {
		existing, err := GetExistingNoPeserta(noPesertaList)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		for i, d := range validData {
			if existing[d.NoPeserta] {
				errorsList = append(errorsList, importError{i + 2, "no_peserta " + d.NoPeserta + " sudah terdaftar"})
			}
		}
	}

	if len(errorsList) > 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Import gagal, terdapat kesalahan pada data",
			"errors":  errorsList,
		})
	}

	if err := ImportSPMB(validData); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": strconv.Itoa(len(validData)) + " data peserta berhasil diimpor",
	})
}

func ExportHandler(c *fiber.Ctx) error {
	search := c.Query("search", "")
	jalur, _ := strconv.Atoi(c.Query("jalur", "0"))
	laporDiri, err := strconv.Atoi(c.Query("lapor_diri", "-1"))
	if err != nil {
		laporDiri = -1
	}

	data, err := GetAllForExport(search, jalur, laporDiri)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	f := excelize.NewFile()
	defer f.Close()
	sheet := "SPMB"
	f.SetSheetName("Sheet1", sheet)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"2563EB"}, Pattern: 1},
	})

	for i, h := range spmbExcelHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	for i, d := range data {
		row := i + 2
		values := []interface{}{d.NoPeserta, d.NamaPeserta, d.Kelamin, d.AsalSekolah, d.Jalur, d.LaporDiri}
		cell, _ := excelize.CoordinatesToCellName(1, row)
		f.SetSheetRow(sheet, cell, &values)
	}

	f.SetActiveSheet(0)

	filename := fmt.Sprintf("data-spmb-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return f.Write(c.Response().BodyWriter())
}

type bulkLaporDiriPayload struct {
	NoPeserta []string `json:"no_peserta"`
	LaporDiri int      `json:"lapor_diri"`
}

func BulkUpdateLaporDiriHandler(c *fiber.Ctx) error {
	var payload bulkLaporDiriPayload
	if err := c.BodyParser(&payload); err != nil || len(payload.NoPeserta) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak valid",
		})
	}

	if err := BulkUpdateLaporDiri(payload.NoPeserta, payload.LaporDiri); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Status lapor diri berhasil diperbarui",
	})
}

type bulkDeletePayload struct {
	NoPeserta []string `json:"no_peserta"`
}

func BulkDeleteHandler(c *fiber.Ctx) error {
	var payload bulkDeletePayload
	if err := c.BodyParser(&payload); err != nil || len(payload.NoPeserta) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak valid",
		})
	}

	if err := BulkDelete(payload.NoPeserta); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data peserta berhasil dihapus",
	})
}
