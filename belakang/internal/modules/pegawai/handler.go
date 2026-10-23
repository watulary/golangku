package pegawai

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sekolah/internal/auth"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

const fotoPegawaiDir = "./storage/fotopegawai/"

var allowedFotoExt = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}

// pegawaiExcelHeaders dipakai bersama oleh Export & Template.
// Index HARUS sinkron dengan mapping getCell() di ImportHandler (idx 1..53 -> field pegawai).
var pegawaiExcelHeaders = []string{
	"No", "NUPTK", "Nama Pegawai", "Jabatan", "Strata", "Status Kepegawaian",
	"NIKKI", "SIMPKB ID", "Dapodik ID", "NIP Lama", "NIP Baru", "NRK", "NRG",
	"Tempat Lahir", "Tanggal Lahir", "Jenis Kelamin", "Agama", "NIK", "No. KK",
	"Alamat", "RT", "RW", "Kelurahan", "Kecamatan", "Kode Pos",
	"Telp Rumah", "No. HP", "Email",
	"SK Sekolah", "TMT Sekolah", "SK Non KKI", "TMT Non KKI", "SK KKI", "TMT KKI",
	"SK CPNS", "TMT CPNS", "SK PNS", "TMT PNS",
	"TMT KGB", "Pangkat KGB", "Gaji KGB", "SK KGB", "Tanggal SK KGB",
	"Karpeg", "Karis", "Nama Ibu Kandung", "Nama Suami/Istri", "Pekerjaan Suami/Istri",
	"NPWP", "Status Keluar (Kode)", "Keterangan Keluar", "Tanggal Keluar", "Nomor Surat",
	"Foto",
}

// pegawaiExcelHeaderStyle dipakai bersama oleh Export & Template
func pegawaiExcelHeaderStyle(f *excelize.File) int {
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

func GetAdminListHandler(c *fiber.Ctx) error {
	pnsNon := c.Query("pns_non", "")
	search := c.Query("search", "")
	kategori := c.Query("kategori", "")
	status, _ := strconv.Atoi(c.Query("status", "-1"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sort_by", "nama")
	sortDir := c.Query("sort_dir", "asc")

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	data, err := GetAllAdmin(search, kategori, pnsNon, status, sortBy, sortDir, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	total, err := CountAdmin(search, kategori, pnsNon, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func GetByIDHandler(c *fiber.Ctx) error {
	nuptk := c.Params("nuptk")

	data, err := GetByID(nuptk)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Pegawai tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func SearchOptionHandler(c *fiber.Ctx) error {
	search := c.Query("q", "")

	data, err := SearchOption(search, 15)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func GetMeHandler(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*auth.Claims)

	data, err := GetByID(claims.NUPTK)
	if err != nil {
		if err == sql.ErrNoRows {
			// belum ada data pegawai utk akun ini -> kirim null, bukan error
			return c.JSON(fiber.Map{
				"success": true,
				"data":    nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// ================= HELPER MAPPING =================

func mapCreateToModel(input PegawaiCreateRequest) Pegawai {
	return Pegawai{
		Nuptk: input.Nuptk, NamaPegawai: input.NamaPegawai, Jabatan: input.Jabatan,
		Strata: input.Strata, PnsNon: input.PnsNon, Nikki: input.Nikki,
		Simpkbid: input.Simpkbid, Dapodikid: input.Dapodikid,
		NipLama: input.NipLama, NipBaru: input.NipBaru, Nrk: input.Nrk, Nrg: input.Nrg,
		TempatLahir: input.TempatLahir, TanggalLahir: input.TanggalLahir,
		JenisKelamin: input.JenisKelamin, Agama: input.Agama, Nik: input.Nik, NoKk: input.NoKk,
		Alamat: input.Alamat, Rt: input.Rt, Rw: input.Rw, Kel: input.Kel, Kec: input.Kec, Kopos: input.Kopos,
		TelpRmh: input.TelpRmh, TelpHp: input.TelpHp, Email: input.Email,
		SkSekolah: input.SkSekolah, TmtSekolah: input.TmtSekolah,
		SkNonkki: input.SkNonkki, TmtNonkki: input.TmtNonkki,
		SkKki: input.SkKki, TmtKki: input.TmtKki,
		SkCpns: input.SkCpns, TmtCpns: input.TmtCpns,
		SkPns: input.SkPns, TmtPns: input.TmtPns,
		TmtKgb: input.TmtKgb, PangkatKgb: input.PangkatKgb, GajiKgb: input.GajiKgb,
		SkKgb: input.SkKgb, TglskKgb: input.TglskKgb,
		Karpeg: input.Karpeg, Karis: input.Karis,
		NamaIbukdg: input.NamaIbukdg, NmSuamiIstri: input.NmSuamiIstri, PekSuamiIstri: input.PekSuamiIstri,
		Npwp: input.Npwp, StatusKeluar: input.StatusKeluar, KetKeluar: input.KetKeluar,
		TanggalKeluar: input.TanggalKeluar, NomorSurat: input.NomorSurat, Foto: input.Foto,
	}
}

func mapUpdateToModel(nuptk string, input PegawaiUpdateRequest) Pegawai {
	return Pegawai{
		Nuptk: nuptk, NamaPegawai: input.NamaPegawai, Jabatan: input.Jabatan,
		Strata: input.Strata, PnsNon: input.PnsNon, Nikki: input.Nikki,
		Simpkbid: input.Simpkbid, Dapodikid: input.Dapodikid,
		NipLama: input.NipLama, NipBaru: input.NipBaru, Nrk: input.Nrk, Nrg: input.Nrg,
		TempatLahir: input.TempatLahir, TanggalLahir: input.TanggalLahir,
		JenisKelamin: input.JenisKelamin, Agama: input.Agama, Nik: input.Nik, NoKk: input.NoKk,
		Alamat: input.Alamat, Rt: input.Rt, Rw: input.Rw, Kel: input.Kel, Kec: input.Kec, Kopos: input.Kopos,
		TelpRmh: input.TelpRmh, TelpHp: input.TelpHp, Email: input.Email,
		SkSekolah: input.SkSekolah, TmtSekolah: input.TmtSekolah,
		SkNonkki: input.SkNonkki, TmtNonkki: input.TmtNonkki,
		SkKki: input.SkKki, TmtKki: input.TmtKki,
		SkCpns: input.SkCpns, TmtCpns: input.TmtCpns,
		SkPns: input.SkPns, TmtPns: input.TmtPns,
		TmtKgb: input.TmtKgb, PangkatKgb: input.PangkatKgb, GajiKgb: input.GajiKgb,
		SkKgb: input.SkKgb, TglskKgb: input.TglskKgb,
		Karpeg: input.Karpeg, Karis: input.Karis,
		NamaIbukdg: input.NamaIbukdg, NmSuamiIstri: input.NmSuamiIstri, PekSuamiIstri: input.PekSuamiIstri,
		Npwp: input.Npwp, StatusKeluar: input.StatusKeluar, KetKeluar: input.KetKeluar,
		TanggalKeluar: input.TanggalKeluar, NomorSurat: input.NomorSurat, Foto: input.Foto,
	}
}

// ================= CREATE / UPDATE / DELETE =================

func CreateHandler(c *fiber.Ctx) error {
	var input PegawaiCreateRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	input.Nuptk = strings.TrimSpace(input.Nuptk)
	input.NamaPegawai = strings.TrimSpace(input.NamaPegawai)
	input.Jabatan = strings.TrimSpace(input.Jabatan)

	if input.Nuptk == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "NUPTK wajib diisi.",
		})
	}

	if input.NamaPegawai == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nama pegawai wajib diisi.",
		})
	}

	if input.Jabatan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Jabatan wajib diisi.",
		})
	}

	if msg := validateDigitField(input.Nuptk, "NUPTK", 16); msg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": msg,
		})
	}

	if msg := validateDigitField(input.NipBaru, "NIP", 18); msg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": msg,
		})
	}

	if msg := validateDigitField(input.Nik, "NIK", 16); msg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": msg,
		})
	}

	if msg := validateDigitField(input.NoKk, "No. KK", 16); msg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": msg,
		})
	}

	exists, err := ExistsByNuptk(input.Nuptk)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"message": "NUPTK sudah terdaftar.",
		})
	}

	data := mapCreateToModel(input)

	if err := Create(data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Pegawai berhasil ditambahkan.",
		"data":    data,
	})
}

func UpdateHandler(c *fiber.Ctx) error {
	nuptk := c.Params("nuptk")

	var input PegawaiUpdateRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	input.NamaPegawai = strings.TrimSpace(input.NamaPegawai)
	input.Jabatan = strings.TrimSpace(input.Jabatan)

	if input.NamaPegawai == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nama pegawai wajib diisi.",
		})
	}

	if input.Jabatan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Jabatan wajib diisi.",
		})
	}

	if msg := validateDigitField(input.NipBaru, "NIP", 18); msg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": msg,
		})
	}

	if msg := validateDigitField(input.Nik, "NIK", 16); msg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": msg,
		})
	}

	if msg := validateDigitField(input.NoKk, "No. KK", 16); msg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": msg,
		})
	}

	oldData, err := GetByID(nuptk)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Pegawai tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	data := mapUpdateToModel(nuptk, input)
	data.Foto = oldData.Foto
	data.FotoURL = oldData.FotoURL

	// Semua field Pegawai bertipe string sehingga struct bisa dibandingkan langsung dengan ==
	if *oldData == data {
		return c.JSON(fiber.Map{
			"success": true,
			"changed": false,
			"message": "Tidak ada perubahan data.",
			"data":    oldData,
		})
	}

	if err := Update(nuptk, data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"changed": true,
		"message": "Data pegawai berhasil diperbarui.",
		"data":    data,
	})
}

func DeleteHandler(c *fiber.Ctx) error {
	nuptk := c.Params("nuptk")

	data, err := GetByID(nuptk)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Pegawai tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := Delete(nuptk); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Pegawai tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if data.Foto != "" {
		os.Remove(fotoPegawaiDir + data.Foto)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pegawai berhasil dihapus.",
	})
}

// =========================================
// Export Excel: Daftar Pegawai
// =========================================

func ExportHandler(c *fiber.Ctx) error {
	pnsNon := c.Query("pns_non", "")
	search := c.Query("search", "")
	kategori := c.Query("kategori", "")
	status, _ := strconv.Atoi(c.Query("status", "-1"))

	data, err := GetAllForExport(search, kategori, pnsNon, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	f := excelize.NewFile()
	sheet := "Pegawai"
	f.SetSheetName("Sheet1", sheet)

	headers := append(append([]string{}, pegawaiExcelHeaders...), "Status")

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, row := range data {
		r := i + 2
		statusLabel := "Aktif"
		if row.StatusKeluar != "" {
			statusLabel = "Keluar"
		}
		values := []interface{}{
			i + 1,
			row.Nuptk, row.NamaPegawai, row.Jabatan, row.Strata, row.PnsNon,
			row.Nikki, row.Simpkbid, row.Dapodikid, row.NipLama, row.NipBaru, row.Nrk, row.Nrg,
			row.TempatLahir, row.TanggalLahir, row.JenisKelamin, row.Agama, row.Nik, row.NoKk,
			row.Alamat, row.Rt, row.Rw, row.Kel, row.Kec, row.Kopos,
			row.TelpRmh, row.TelpHp, row.Email,
			row.SkSekolah, row.TmtSekolah, row.SkNonkki, row.TmtNonkki, row.SkKki, row.TmtKki,
			row.SkCpns, row.TmtCpns, row.SkPns, row.TmtPns,
			row.TmtKgb, row.PangkatKgb, row.GajiKgb, row.SkKgb, row.TglskKgb,
			row.Karpeg, row.Karis, row.NamaIbukdg, row.NmSuamiIstri, row.PekSuamiIstri,
			row.Npwp, row.StatusKeluar, row.KetKeluar, row.TanggalKeluar, row.NomorSurat,
			row.Foto, statusLabel,
		}
		for colIdx, v := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, r)
			f.SetCellValue(sheet, cell, v)
		}
	}

	// Paksa semua kolom data jadi format Text, biar NIK/NUPTK/NIP/No.KK dll
	// tidak berubah jadi angka/notasi ilmiah saat dibuka atau diedit di Excel
	textStyle, err := f.NewStyle(&excelize.Style{
		CustomNumFmt: func() *string { s := "@"; return &s }(),
	})
	if err == nil {
		firstDataCell, _ := excelize.CoordinatesToCellName(1, 2)
		lastDataCell, _ := excelize.CoordinatesToCellName(len(headers), len(data)+1)
		f.SetCellStyle(sheet, firstDataCell, lastDataCell, textStyle)
	}

	lastCol, _ := excelize.ColumnNumberToName(len(headers))

	// Lebar kolom: No dibuat sempit, sisanya default 18
	f.SetColWidth(sheet, "A", "A", 6)
	secondCol, _ := excelize.ColumnNumberToName(2)
	f.SetColWidth(sheet, secondCol, lastCol, 18)

	firstCell, _ := excelize.CoordinatesToCellName(1, 1)
	lastHeaderCell, _ := excelize.CoordinatesToCellName(len(headers), 1)

	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"1E40AF"},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "left",
		},
	})
	if err == nil {
		f.SetCellStyle(sheet, firstCell, lastHeaderCell, headerStyle)
	}

	filename := fmt.Sprintf(
		"data-pegawai-%s.xlsx",
		time.Now().Format("2006-01-02-15-04-05"),
	)
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	return f.Write(c.Response().BodyWriter())
}

func UpdateFotoHandler(c *fiber.Ctx) error {
	nuptk := c.Params("nuptk")

	oldData, err := GetByID(nuptk)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Pegawai tidak ditemukan.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	file, err := c.FormFile("foto")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File foto wajib diunggah.",
		})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedFotoExt[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format file harus jpg, jpeg, png, atau webp.",
		})
	}

	if file.Size > 100*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Ukuran file maksimal 100 KB.",
		})
	}

	if err := os.MkdirAll(fotoPegawaiDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	newID, err := uuid.NewV7()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	newFilename := newID.String() + ext

	if err := c.SaveFile(file, fotoPegawaiDir+newFilename); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := UpdateFoto(nuptk, newFilename); err != nil {
		os.Remove(fotoPegawaiDir + newFilename)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if oldData.Foto != "" {
		go deleteFotoFileWithRetry(oldData.Foto)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Foto pegawai berhasil diperbarui.",
		"data": fiber.Map{
			"foto": newFilename,
		},
	})
}

// ================= VALIDASI DIGIT =================

func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// validateDigitField mengembalikan pesan error jika value tidak kosong tapi
// panjangnya tidak sesuai atau mengandung karakter selain angka. Kosong dianggap valid
// (biar field opsional seperti nip_baru/nik/no_kk tidak dipaksa selalu diisi).
func validateDigitField(value, label string, length int) string {
	if value == "" {
		return ""
	}
	if len(value) != length || !isAllDigits(value) {
		return fmt.Sprintf("%s harus %d digit angka.", label, length)
	}
	return ""
}

// ================= CEK DUPLIKAT =================

func CheckDuplicateHandler(c *fiber.Ctx) error {
	field := c.Query("field", "")
	value := strings.TrimSpace(c.Query("value", ""))
	exclude := c.Query("exclude", "") // nuptk pegawai yang sedang diedit; kosong saat mode tambah

	if !allowedDuplicateFields[field] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Field tidak valid untuk pengecekan duplikat.",
		})
	}

	if value == "" {
		return c.JSON(fiber.Map{
			"success": true,
			"data":    fiber.Map{"exists": false},
		})
	}

	result, err := FindDuplicate(field, value, exclude)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if result == nil {
		return c.JSON(fiber.Map{
			"success": true,
			"data":    fiber.Map{"exists": false},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"exists":       true,
			"nuptk":        result.Nuptk,
			"nama_pegawai": result.NamaPegawai,
		},
	})
}

// getCell aman mengambil isi kolom ke-idx (0-based); kolom kosong/tidak ada -> ""
func getCell(row []string, idx int) string {
	if idx < len(row) {
		return strings.TrimSpace(row[idx])
	}
	return ""
}

// isRowBlank true kalau seluruh kolom di baris kosong
func isRowBlank(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

func ImportHandler(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File Excel wajib diunggah",
		})
	}

	src, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuka file",
		})
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format file tidak valid, pastikan file .xlsx",
		})
	}

	sheet := f.GetSheetList()[0]
	rows, err := f.GetRows(sheet)
	if err != nil || len(rows) < 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File kosong atau tidak ada data untuk diimport",
		})
	}

	var errorsList []ImportRowError
	var validData []Pegawai
	var allNuptk, allNipBaru, allNik []string
	seenNuptk := map[string]int{}
	seenNipBaru := map[string]int{}
	seenNik := map[string]int{}

	// rows[0] adalah header, mulai dari rows[1]
	for i, row := range rows[1:] {
		rowNum := i + 2 // baris asli di Excel
		if isRowBlank(row) {
			continue
		}

		p := Pegawai{
			Nuptk:         getCell(row, 1),
			NamaPegawai:   getCell(row, 2),
			Jabatan:       getCell(row, 3),
			Strata:        getCell(row, 4),
			PnsNon:        getCell(row, 5),
			Nikki:         getCell(row, 6),
			Simpkbid:      getCell(row, 7),
			Dapodikid:     getCell(row, 8),
			NipLama:       getCell(row, 9),
			NipBaru:       getCell(row, 10),
			Nrk:           getCell(row, 11),
			Nrg:           getCell(row, 12),
			TempatLahir:   getCell(row, 13),
			TanggalLahir:  getCell(row, 14),
			JenisKelamin:  getCell(row, 15),
			Agama:         getCell(row, 16),
			Nik:           getCell(row, 17),
			NoKk:          getCell(row, 18),
			Alamat:        getCell(row, 19),
			Rt:            getCell(row, 20),
			Rw:            getCell(row, 21),
			Kel:           getCell(row, 22),
			Kec:           getCell(row, 23),
			Kopos:         getCell(row, 24),
			TelpRmh:       getCell(row, 25),
			TelpHp:        getCell(row, 26),
			Email:         getCell(row, 27),
			SkSekolah:     getCell(row, 28),
			TmtSekolah:    getCell(row, 29),
			SkNonkki:      getCell(row, 30),
			TmtNonkki:     getCell(row, 31),
			SkKki:         getCell(row, 32),
			TmtKki:        getCell(row, 33),
			SkCpns:        getCell(row, 34),
			TmtCpns:       getCell(row, 35),
			SkPns:         getCell(row, 36),
			TmtPns:        getCell(row, 37),
			TmtKgb:        getCell(row, 38),
			PangkatKgb:    getCell(row, 39),
			GajiKgb:       getCell(row, 40),
			SkKgb:         getCell(row, 41),
			TglskKgb:      getCell(row, 42),
			Karpeg:        getCell(row, 43),
			Karis:         getCell(row, 44),
			NamaIbukdg:    getCell(row, 45),
			NmSuamiIstri:  getCell(row, 46),
			PekSuamiIstri: getCell(row, 47),
			Npwp:          getCell(row, 48),
			StatusKeluar:  getCell(row, 49),
			KetKeluar:     getCell(row, 50),
			TanggalKeluar: getCell(row, 51),
			NomorSurat:    getCell(row, 52),
			Foto:          getCell(row, 53),
		}

		if p.Nuptk == "" {
			errorsList = append(errorsList, ImportRowError{rowNum, "NUPTK wajib diisi"})
			continue
		}

		// cek duplikat di dalam file
		if prev, ok := seenNuptk[p.Nuptk]; ok {
			errorsList = append(errorsList, ImportRowError{rowNum, fmt.Sprintf("NUPTK '%s' duplikat dengan baris %d", p.Nuptk, prev)})
		} else {
			seenNuptk[p.Nuptk] = rowNum
		}
		if p.NipBaru != "" {
			if prev, ok := seenNipBaru[p.NipBaru]; ok {
				errorsList = append(errorsList, ImportRowError{rowNum, fmt.Sprintf("NIP Baru '%s' duplikat dengan baris %d", p.NipBaru, prev)})
			} else {
				seenNipBaru[p.NipBaru] = rowNum
			}
		}
		if p.Nik != "" {
			if prev, ok := seenNik[p.Nik]; ok {
				errorsList = append(errorsList, ImportRowError{rowNum, fmt.Sprintf("NIK '%s' duplikat dengan baris %d", p.Nik, prev)})
			} else {
				seenNik[p.Nik] = rowNum
			}
		}

		validData = append(validData, p)
		allNuptk = append(allNuptk, p.Nuptk)
		if p.NipBaru != "" {
			allNipBaru = append(allNipBaru, p.NipBaru)
		}
		if p.Nik != "" {
			allNik = append(allNik, p.Nik)
		}
	}

	// cek bentrok dengan data yang sudah ada di database
	conflicts, err := GetExistingConflicts(allNuptk, allNipBaru, allNik)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	for _, cf := range conflicts {
		if rowNum, ok := seenNuptk[cf.Nuptk]; ok && cf.Nuptk != "" {
			errorsList = append(errorsList, ImportRowError{rowNum, fmt.Sprintf("NUPTK '%s' sudah terdaftar di database (a.n. %s)", cf.Nuptk, cf.NamaPegawai)})
		}
		if rowNum, ok := seenNipBaru[cf.NipBaru]; ok && cf.NipBaru != "" {
			errorsList = append(errorsList, ImportRowError{rowNum, fmt.Sprintf("NIP Baru '%s' sudah terdaftar di database (a.n. %s)", cf.NipBaru, cf.NamaPegawai)})
		}
		if rowNum, ok := seenNik[cf.Nik]; ok && cf.Nik != "" {
			errorsList = append(errorsList, ImportRowError{rowNum, fmt.Sprintf("NIK '%s' sudah terdaftar di database (a.n. %s)", cf.Nik, cf.NamaPegawai)})
		}
	}

	if len(errorsList) > 0 {
		sort.Slice(errorsList, func(i, j int) bool { return errorsList[i].Row < errorsList[j].Row })
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Ditemukan data bermasalah, seluruh proses import dibatalkan",
			"errors":  errorsList,
		})
	}

	if len(validData) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Tidak ada data valid untuk diimport",
		})
	}

	if err := ImportPegawai(validData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menyimpan data: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("%d data pegawai berhasil diimport", len(validData)),
		"data":    fiber.Map{"imported": len(validData)},
	})
}

func TemplateHandler(c *fiber.Ctx) error {
	f := excelize.NewFile()
	sheet := "Pegawai"
	f.SetSheetName("Sheet1", sheet)

	for i, h := range pegawaiExcelHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	lastCol, _ := excelize.ColumnNumberToName(len(pegawaiExcelHeaders))
	secondCol, _ := excelize.ColumnNumberToName(2)
	f.SetColWidth(sheet, "A", "A", 6)
	f.SetColWidth(sheet, secondCol, lastCol, 18)

	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	if style := pegawaiExcelHeaderStyle(f); style != -1 {
		firstCell, _ := excelize.CoordinatesToCellName(1, 1)
		lastHeaderCell, _ := excelize.CoordinatesToCellName(len(pegawaiExcelHeaders), 1)
		f.SetCellStyle(sheet, firstCell, lastHeaderCell, style)
	}

	// Sheet petunjuk pengisian
	petunjuk := "Petunjuk"
	f.NewSheet(petunjuk)
	instructions := []string{
		"PETUNJUK PENGISIAN TEMPLATE UNGGAH PEGAWAI",
		"",
		"1. Jangan mengubah urutan atau nama kolom pada sheet 'Pegawai'.",
		"2. Kolom 'No' boleh dikosongkan, tidak diproses saat import.",
		"3. Kolom NUPTK wajib diisi, dan tidak boleh sama (duplikat) antar baris maupun dengan data yang sudah ada di sistem.",
		"4. Kolom NIP Baru dan NIK, jika diisi, juga tidak boleh duplikat antar baris maupun dengan data yang sudah ada di sistem.",
		"5. Format tanggal: YYYY-MM-DD (contoh: 1990-05-17).",
		"6. Kolom 'Foto' dikosongkan saja.",
		"7. Jika salah satu baris bermasalah, SELURUH proses import akan dibatalkan (tidak ada data yang tersimpan sebagian).",
		"8. Beberapa kolom hanya menerima nilai tertentu (harus ditulis PERSIS seperti daftar di bawah, huruf besar/kecil dan spasi berpengaruh):",
		"",
		"Kolom 'Jabatan' — pilih salah satu:",
		"Guru | Tatausaha",
		"",
		"Kolom 'Status Kepegawaian' — pilih salah satu:",
		"PNS | CPNS | PPPK | PPPK-PW | KKI | Non KKI | Guru Honor Sekolah | Honor Daerah TK.I Provinsi",
		"",
		"Kolom 'Jenis Kelamin' — pilih salah satu:",
		"L | P ",
		"",
		"Kolom 'Agama' — pilih salah satu:",
		"Islam | Kristen | Katholik | Hindu | Budha | Khonghucu | Kepercayaan kpd Tuhan YME | lainnya",
		"",
		"Kolom 'Strata' — pilih salah satu:",
		"Tidak sekolah | TK / sederajat | SD / sederajat | SMP / sederajat | SMA / sederajat",
		"D1 | D2 | D3 | D4 | S1 | S2 | S2 Terapan | S3 | S3 Terapan | Sp-1 | Sp-2",
		"PAUD | Profesi | Paket A | Paket B | Paket C | Informal | Non formal | Putus SD | Lainnya",
	}
	for i, line := range instructions {
		cell, _ := excelize.CoordinatesToCellName(1, i+1)
		f.SetCellValue(petunjuk, cell, line)
	}
	f.SetColWidth(petunjuk, "A", "A", 90)
	f.SetActiveSheet(0)

	filename := fmt.Sprintf("template-pegawai-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return f.Write(c.Response().BodyWriter())
}

func DeleteFotoHandler(c *fiber.Ctx) error {
	nuptk := c.Params("nuptk")
	oldData, err := GetByID(nuptk)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Pegawai tidak ditemukan.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if oldData.Foto == "" {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Pegawai sudah menggunakan foto default.",
			"data":    fiber.Map{"foto_url": buildFotoURL("", oldData.JenisKelamin)},
		})
	}

	if err := UpdateFoto(nuptk, ""); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	go deleteFotoFileWithRetry(oldData.Foto)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Foto pegawai berhasil dihapus.",
		"data":    fiber.Map{"foto_url": buildFotoURL("", oldData.JenisKelamin)},
	})
}

func deleteFotoFileWithRetry(filename string) {
	path := fotoPegawaiDir + filename

	// Beri waktu proses lain (antivirus, thumbnail cache, cloud sync, dsb)
	// melepas lock file sebelum mulai mencoba menghapus.
	time.Sleep(3 * time.Second)

	const maxAttempt = 30
	const retryDelay = 1 * time.Second

	for attempt := 1; attempt <= maxAttempt; attempt++ {
		err := os.Remove(path)

		if err == nil || os.IsNotExist(err) {
			return
		}

		if attempt == maxAttempt {
			log.Printf("gagal hapus file foto lama %s setelah beberapa percobaan: %v", filename, err)
		}

		time.Sleep(retryDelay)
	}
}
