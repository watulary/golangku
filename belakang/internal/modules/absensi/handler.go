package absensi

import (
	"fmt"
	"sekolah/internal/auth"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

var allowedStatus = map[string]bool{
	"Hadir": true,
	"Sakit": true,
	"Izin":  true,
	"Alpa":  true,
}

var absensiExcelHeaders = []string{
	"No", "Kelas", "Rombel", "NIS", "Nama Peserta Didik", "Status", "Keterangan",
}

func absensiExcelHeaderStyle(f *excelize.File) int {
	style, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"1E40AF"}},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "left"},
	})
	if err != nil {
		return -1
	}
	return style
}

// autoFitColumnWidth menghitung lebar kolom berdasarkan teks terpanjang di tiap kolom
// (header + seluruh isi data yang sudah dikumpulkan lewat colWidths), lalu menerapkannya.
// minWidth/maxWidth membatasi supaya kolom tidak terlalu sempit atau terlalu lebar.
func autoFitColumnWidth(f *excelize.File, sheet string, colWidths []float64, minWidth, maxWidth float64) {
	for i, w := range colWidths {
		width := w + 2 // padding kiri-kanan biar teks nggak mepet border sel
		if width < minWidth {
			width = minWidth
		}
		if width > maxWidth {
			width = maxWidth
		}
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheet, colName, colName, width)
	}
}

// GetDaftarHandler mengambil daftar siswa 1 rombel pada tanggal tertentu beserta status absensinya
func GetDaftarHandler(c *fiber.Ctx) error {
	tanggal := c.Query("tanggal")
	idKelas := c.QueryInt("id_kelas", 0)
	idRombel := c.QueryInt("id_rombel", 0)

	if tanggal == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "tanggal wajib diisi",
		})
	}

	list, err := GetDaftarAbsensi(tanggal, idKelas, idRombel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil daftar absensi",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Berhasil mengambil daftar absensi",
		"data":    list,
	})
}

// BulkUpsertHandler menyimpan absensi utk 1 rombel pada 1 tanggal sekaligus
func BulkUpsertHandler(c *fiber.Ctx) error {
	var req AbsensiBulkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid",
		})
	}

	if req.Tanggal == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "tanggal wajib diisi",
		})
	}
	if len(req.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Tidak ada data absensi yang dikirim",
		})
	}

	for _, item := range req.Items {
		if item.Nis == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Ada item tanpa nis",
			})
		}
		if !allowedStatus[item.Status] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Status tidak valid: " + item.Status,
			})
		}
	}

	if err := BulkUpsert(req.Tanggal, req.Items); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menyimpan absensi",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Absensi berhasil disimpan",
	})
}

func ExportHandler(c *fiber.Ctx) error {
	tanggal := c.Query("tanggal")
	idKelas := c.QueryInt("id_kelas", 0)
	idRombel := c.QueryInt("id_rombel", 0)
	status := c.Query("status", "")

	if tanggal == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "tanggal wajib diisi",
		})
	}

	list, err := GetDaftarAbsensi(tanggal, idKelas, idRombel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if status != "" && status != "Semua" {
		filtered := make([]AbsensiSiswa, 0, len(list))
		for _, s := range list {
			if s.Status == status {
				filtered = append(filtered, s)
			}
		}
		list = filtered
	}

	f := excelize.NewFile()
	sheet := "Absensi"
	f.SetSheetName("Sheet1", sheet)

	for i, h := range absensiExcelHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// lebar awal = panjang header, di-update tiap baris data ditulis
	colWidths := make([]float64, len(absensiExcelHeaders))
	for i, h := range absensiExcelHeaders {
		colWidths[i] = float64(len([]rune(h)))
	}

	counts := map[string]int{"Hadir": 0, "Sakit": 0, "Izin": 0, "Alpa": 0}
	for i, s := range list {
		r := i + 2
		values := []interface{}{i + 1, s.NamaKelas, s.NamaRombel, s.Nis, s.NamaSiswa, s.Status, s.Keterangan}
		for colIdx, v := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, r)
			f.SetCellValue(sheet, cell, v)
			if l := float64(len([]rune(fmt.Sprintf("%v", v)))); l > colWidths[colIdx] {
				colWidths[colIdx] = l
			}
		}
		counts[s.Status]++
	}

	startRow := len(list) + 3
	f.SetCellValue(sheet, fmt.Sprintf("A%d", startRow), "Ringkasan")
	labels := []string{"Hadir", "Sakit", "Izin", "Alpa"}
	for i, lbl := range labels {
		r := startRow + 1 + i
		f.SetCellValue(sheet, fmt.Sprintf("A%d", r), lbl)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", r), counts[lbl])
	}

	autoFitColumnWidth(f, sheet, colWidths, 6, 40)

	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	if style := absensiExcelHeaderStyle(f); style != -1 {
		firstCell, _ := excelize.CoordinatesToCellName(1, 1)
		lastHeaderCell, _ := excelize.CoordinatesToCellName(len(absensiExcelHeaders), 1)
		f.SetCellStyle(sheet, firstCell, lastHeaderCell, style)
	}

	filename := fmt.Sprintf("laporan-absensi-%s.xlsx", tanggal)
	if status != "" && status != "Semua" {
		filename = fmt.Sprintf("laporan-absensi-%s-%s.xlsx", tanggal, strings.ToLower(status))
	}
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	return f.Write(c.Response().BodyWriter())
}

func GetRekapBulananHandler(c *fiber.Ctx) error {
	bulan := c.QueryInt("bulan", 0)
	tahun := c.QueryInt("tahun", 0)
	idKelas := c.QueryInt("id_kelas", 0)
	idRombel := c.QueryInt("id_rombel", 0)

	if bulan < 1 || bulan > 12 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "bulan tidak valid (1-12)",
		})
	}
	if tahun < 2000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "tahun wajib diisi",
		})
	}

	list, jumlahHari, err := GetRekapBulanan(bulan, tahun, idKelas, idRombel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Berhasil mengambil rekap bulanan",
		"data": fiber.Map{
			"bulan":       bulan,
			"tahun":       tahun,
			"jumlah_hari": jumlahHari,
			"siswa":       list,
		},
	})
}

func ExportRekapBulananHandler(c *fiber.Ctx) error {
	bulan := c.QueryInt("bulan", 0)
	tahun := c.QueryInt("tahun", 0)
	idKelas := c.QueryInt("id_kelas", 0)
	idRombel := c.QueryInt("id_rombel", 0)

	if bulan < 1 || bulan > 12 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "bulan tidak valid (1-12)",
		})
	}
	if tahun < 2000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "tahun wajib diisi",
		})
	}

	list, jumlahHari, err := GetRekapBulanan(bulan, tahun, idKelas, idRombel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	f := excelize.NewFile()
	sheet := "Rekap Bulanan"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"No", "Kelas", "Rombel", "NIS", "Nama Peserta Didik"}
	for d := 1; d <= jumlahHari; d++ {
		headers = append(headers, strconv.Itoa(d))
	}
	headers = append(headers, "Hadir", "Sakit", "Izin", "Alpa")

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	colWidths := make([]float64, len(headers))
	for i, h := range headers {
		colWidths[i] = float64(len([]rune(h)))
	}

	for i, s := range list {
		r := i + 2
		col := 1
		setNext := func(v interface{}) {
			cell, _ := excelize.CoordinatesToCellName(col, r)
			f.SetCellValue(sheet, cell, v)
			if l := float64(len([]rune(fmt.Sprintf("%v", v)))); l > colWidths[col-1] {
				colWidths[col-1] = l
			}
			col++
		}
		setNext(i + 1)
		setNext(s.NamaKelas)
		setNext(s.NamaRombel)
		setNext(s.Nis)
		setNext(s.NamaSiswa)

		counts := map[string]int{"H": 0, "S": 0, "I": 0, "A": 0}
		for _, kode := range s.Harian {
			setNext(kode)
			if kode != "" {
				counts[kode]++
			}
		}
		setNext(counts["H"])
		setNext(counts["S"])
		setNext(counts["I"])
		setNext(counts["A"])
	}

	autoFitColumnWidth(f, sheet, colWidths, 3, 30)

	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		XSplit:      5,
		TopLeftCell: "F2",
		ActivePane:  "bottomRight",
	})

	if style := absensiExcelHeaderStyle(f); style != -1 {
		firstCell, _ := excelize.CoordinatesToCellName(1, 1)
		lastHeaderCell, _ := excelize.CoordinatesToCellName(len(headers), 1)
		f.SetCellStyle(sheet, firstCell, lastHeaderCell, style)
	}

	filename := fmt.Sprintf("rekap-absensi-%02d-%d.xlsx", bulan, tahun)
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	return f.Write(c.Response().BodyWriter())
}

func GetSayaHandler(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*auth.Claims)
	nis := claims.NIS

	bulan := c.QueryInt("bulan", int(time.Now().Month()))
	tahun := c.QueryInt("tahun", time.Now().Year())

	if bulan < 1 || bulan > 12 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "bulan tidak valid (1-12)",
		})
	}

	list, err := GetAbsensiSaya(nis, bulan, tahun)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Berhasil mengambil absensi",
		"data": fiber.Map{
			"bulan": bulan,
			"tahun": tahun,
			"list":  list,
		},
	})
}
