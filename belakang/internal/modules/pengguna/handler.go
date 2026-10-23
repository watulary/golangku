package pengguna

import (
	"fmt"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"sekolah/internal/auth"
	"sekolah/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

// =========================================
// Password Validation
// =========================================

const (
	minPasswordLength = 8
	maxPasswordLength = 64
)

var (
	passwordUpperRegex   = regexp.MustCompile(`[A-Z]`)
	passwordLowerRegex   = regexp.MustCompile(`[a-z]`)
	passwordNumberRegex  = regexp.MustCompile(`[0-9]`)
	passwordSpecialRegex = regexp.MustCompile("[^a-zA-Z0-9]")
)

func validatePassword(password string) string {

	switch {

	case len(password) < minPasswordLength:
		return "Password minimal 8 karakter."

	case len(password) > maxPasswordLength:
		return "Password maksimal 64 karakter."

	case !passwordUpperRegex.MatchString(password):
		return "Password harus mengandung minimal satu huruf besar."

	case !passwordLowerRegex.MatchString(password):
		return "Password harus mengandung minimal satu huruf kecil."

	case !passwordNumberRegex.MatchString(password):
		return "Password harus mengandung minimal satu angka."

	case !passwordSpecialRegex.MatchString(password):
		return "Password harus mengandung minimal satu karakter khusus."

	}

	return ""

}

func ListHandler(c *fiber.Ctx) error {

	// ==========================
	// Query Parameter
	// ==========================

	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit < 1 {
		limit = 10
	}

	if limit > 1000 {
		limit = 1000
	}

	search := c.Query("q")
	level := c.Query("level")
	aktif := c.Query("aktif")

	sort := c.Query("sort", "nama_pengguna")
	order := c.Query("order", "asc")

	// ==========================
	// Ambil Data
	// ==========================

	data, err := List(
		page,
		limit,
		search,
		level,
		aktif,
		sort,
		order,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	// ==========================
	// Total Data
	// ==========================

	total, err := CountList(
		search,
		level,
		aktif,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	// ==========================
	// Foto URL
	// ==========================

	for i := range data {

		data[i].FotoURL = buildFotoURLPengguna(
			data[i].Foto,
			data[i].NUPTK,
			data[i].NIS,
		)
	}

	// ==========================
	// Pagination
	// ==========================

	totalPage := 0

	if total > 0 {
		totalPage = int(math.Ceil(float64(total) / float64(limit)))
	}

	return c.JSON(
		fiber.Map{
			"success": true,
			"message": "Data pengguna berhasil diambil.",
			"data":    data,
			"pagination": fiber.Map{
				"page":       page,
				"limit":      limit,
				"total":      total,
				"total_page": totalPage,
			},
		},
	)
}

// =========================================
// Get Pengguna
// =========================================

func GetPenggunaHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	if id == "" {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "ID pengguna tidak valid.",
			},
		)

	}

	data, err := GetPenggunaByID(id)

	if err != nil {

		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"success": false,
				"message": "Pengguna tidak ditemukan.",
			},
		)

	}

	data.FotoURL = buildFotoURLPengguna(data.Foto, data.NUPTK, data.NIS)

	return c.JSON(
		fiber.Map{
			"success": true,
			"data":    data,
		},
	)

}

// =========================================
// Update Pengguna
// =========================================

func UpdatePenggunaHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	if id == "" {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "ID pengguna tidak valid.",
			},
		)

	}

	var req UpdatePenggunaRequest

	if err := c.BodyParser(&req); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Request tidak valid.",
			},
		)

	}

	req.IDPengguna = id

	changed, err := UpdatePengguna(req)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	if !changed {
		return c.JSON(
			fiber.Map{
				"success": true,
				"changed": false,
				"message": "Tidak ada perubahan data.",
			},
		)
	}

	return c.JSON(
		fiber.Map{
			"success": true,
			"changed": true,
			"message": "Data pengguna berhasil diperbarui.",
		},
	)

}

// =========================================
// Check Duplicate (realtime, Tambah/Ubah)
// =========================================

func CheckDuplicateHandler(c *fiber.Ctx) error {

	field := c.Query("field")
	value := strings.TrimSpace(c.Query("value"))
	exclude := c.Query("exclude")

	if value == "" {
		return c.JSON(fiber.Map{"exists": false})
	}

	var (
		info DuplicateInfo
		err  error
	)

	switch field {

	case "nuptk":
		info, err = CheckDuplicateNUPTK(value, exclude)

	default:

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Field tidak dikenali.",
			},
		)
	}

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	return c.JSON(info)
}

// =========================================
// Create Pengguna
// =========================================

func CreateHandler(c *fiber.Ctx) error {

	var req CreatePenggunaRequest

	// =========================================
	// Parse Body
	// =========================================

	if err := c.BodyParser(&req); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Request tidak valid.",
			},
		)
	}

	// =========================================
	// Simpan
	// =========================================

	idPengguna, err := CreatePengguna(req)

	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	// =========================================
	// Success
	// =========================================

	return c.JSON(
		fiber.Map{
			"success":     true,
			"message":     "Pengguna berhasil ditambahkan.",
			"id_pengguna": idPengguna,
		},
	)
}

// =========================================
// Export Excel
// =========================================

func ExportExcelHandler(c *fiber.Ctx) error {

	// ==========================
	// Ambil Data
	// ==========================

	data, err := ExportExcel()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	// ==========================
	// Buat Workbook
	// ==========================

	f := excelize.NewFile()

	sheet := "Pengguna"

	f.SetSheetName("Sheet1", sheet)

	// ==========================
	// Header
	// ==========================

	headers := []string{
		"No",
		"Nama Pengguna",
		"Email",
		"Level",
		"Status",
		"Password",
		"Login Terakhir",
	}

	for i, h := range headers {

		cell, _ := excelize.CoordinatesToCellName(i+1, 1)

		f.SetCellValue(sheet, cell, h)
	}

	// ==========================
	// Isi Data
	// ==========================

	for i, row := range data {

		r := i + 2

		level := row.Level

		switch row.Level {

		case "Admin":
			level = "Super Admin"

		case "Staf":
			level = "Admin"

		case "Pesertadidik":
			level = "Peserta Didik"
		}

		status := "Tidak Aktif"

		if row.Aktif == 1 {
			status = "Aktif"
		}

		password := "Standar"

		if row.Password {
			password = "Diubah"
		}

		login := "Belum Pernah Masuk"

		if row.LastLogin != nil {

			login = row.LastLogin.
				In(time.FixedZone("WIB", 7*3600)).
				Format("2006-01-02 15:04:05")
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", r), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", r), row.NamaPengguna)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", r), row.Email)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", r), level)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", r), status)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", r), password)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", r), login)
	}

	f.SetColWidth(sheet, "A", "A", 8)
	f.SetColWidth(sheet, "B", "B", 35)
	f.SetColWidth(sheet, "C", "C", 35)
	f.SetColWidth(sheet, "D", "D", 20)
	f.SetColWidth(sheet, "E", "E", 15)
	f.SetColWidth(sheet, "F", "F", 18)
	f.SetColWidth(sheet, "G", "G", 24)

	// ==========================
	// Freeze Header
	// ==========================

	err = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	if err != nil {
		return err
	}

	// ==========================
	// Style Header
	// ==========================

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
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "D1D5DB",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "D1D5DB",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "D1D5DB",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "D1D5DB",
				Style: 1,
			},
		},
	})

	if err == nil {
		f.SetCellStyle(sheet, "A1", "G1", headerStyle)
	}

	// ==========================
	// Nama File
	// ==========================

	filename := fmt.Sprintf(
		"pengguna-%s.xlsx",
		time.Now().Format("2006-01-02-15-04-05"),
	)

	// ==========================
	// Header Download
	// ==========================

	c.Set(
		"Content-Type",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	)

	c.Set(
		"Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, filename),
	)

	c.Set(
		"Access-Control-Expose-Headers",
		"Content-Disposition",
	)

	return f.Write(c.Response().BodyWriter())
}

// =========================================
// Download Template Excel
// =========================================

func DownloadTemplateHandler(c *fiber.Ctx) error {

	f := excelize.NewFile()

	// =========================================
	// Sheet 1 : Data Pengguna
	// =========================================

	dataSheet := "Data Pengguna"

	f.SetSheetName("Sheet1", dataSheet)

	headers := []string{
		"Nama Pengguna",
		"Email",
		"Level",
		"NUPTK",
		"NIS",
	}

	for i, h := range headers {

		cell, _ := excelize.CoordinatesToCellName(i+1, 1)

		f.SetCellValue(dataSheet, cell, h)
	}

	// Contoh Data

	f.SetCellValue(dataSheet, "A2", "Kirana Agnimaya")
	f.SetCellValue(dataSheet, "B2", "kirana@mail.id")
	f.SetCellValue(dataSheet, "C2", "Pegawai")
	f.SetCellValue(dataSheet, "D2", "1987654321001")

	f.SetCellValue(dataSheet, "A3", "Nabawi Raudo")
	f.SetCellValue(dataSheet, "B3", "nabawi@mail.id")
	f.SetCellValue(dataSheet, "C3", "Pesertadidik")
	f.SetCellValue(dataSheet, "E3", "123456")

	// =========================================
	// Sheet 2 : Petunjuk
	// =========================================

	petunjukSheet := "Petunjuk"

	f.NewSheet(petunjukSheet)

	f.SetCellValue(petunjukSheet, "A1", "No")
	f.SetCellValue(petunjukSheet, "B1", "Petunjuk")

	petunjuk := []string{
		"Isi data mulai dari baris ke-2.",
		"Kolom Nama Pengguna, Email, dan Level wajib diisi.",
		"Level hanya boleh: Admin, Staf, Pegawai, atau Pesertadidik.",
		"Sebelum mengisi NUPTK atau NIS, ubah format kolom menjadi Text.",
		"Jika Level = Pegawai maka NUPTK wajib diisi.",
		"Jika Level = Pesertadidik maka NIS wajib diisi.",
		"Email tidak boleh duplikat.",
		"NUPTK tidak boleh duplikat.",
		"NIS tidak boleh duplikat.",
		"Password standar: Abcd1234* (Ubah password ketika pertama kali login).",
		"Status pengguna otomatis menjadi Aktif.",
		"Jika ditemukan satu saja kesalahan maka seluruh proses upload dibatalkan.",
	}

	for i, item := range petunjuk {

		row := i + 2

		f.SetCellValue(petunjukSheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(petunjukSheet, fmt.Sprintf("B%d", row), item)
	}

	// =========================================
	// Style Header
	// =========================================

	headerStyle, _ := f.NewStyle(&excelize.Style{
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

	f.SetCellStyle(dataSheet, "A1", "E1", headerStyle)
	f.SetCellStyle(petunjukSheet, "A1", "B1", headerStyle)

	// =========================================
	// Lebar Kolom
	// =========================================

	f.SetColWidth(dataSheet, "A", "A", 30)
	f.SetColWidth(dataSheet, "B", "B", 35)
	f.SetColWidth(dataSheet, "C", "C", 18)
	f.SetColWidth(dataSheet, "D", "D", 24)
	f.SetColWidth(dataSheet, "E", "E", 20)

	f.SetColWidth(petunjukSheet, "A", "A", 8)
	f.SetColWidth(petunjukSheet, "B", "B", 90)

	// =========================================
	// Freeze Header
	// =========================================

	_ = f.SetPanes(dataSheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	_ = f.SetPanes(petunjukSheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	// Sheet pertama yang dibuka
	f.SetActiveSheet(0)

	filename := "template-pengguna.xlsx"

	c.Set(
		"Content-Type",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	)

	c.Set(
		"Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, filename),
	)

	c.Set(
		"Access-Control-Expose-Headers",
		"Content-Disposition",
	)

	return f.Write(c.Response().BodyWriter())
}

// =========================================
// Bulk Activate
// =========================================

func BulkActivateHandler(c *fiber.Ctx) error {

	var req BulkStatusRequest

	// ==========================
	// Parse Body
	// ==========================

	if err := c.BodyParser(&req); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Request tidak valid.",
			},
		)

	}

	// ==========================
	// Validasi
	// ==========================

	if len(req.IDs) == 0 {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Pilih minimal satu pengguna.",
			},
		)

	}

	// ==========================
	// Update
	// ==========================

	if err := BulkActivate(req.IDs); err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)

	}

	// ==========================
	// Success
	// ==========================

	return c.JSON(
		fiber.Map{
			"success": true,
			"message": fmt.Sprintf(
				"%d Pengguna berhasil diaktifkan.",
				len(req.IDs),
			),
		},
	)

}

// =========================================
// Bulk Deactivate
// =========================================

func BulkDeactivateHandler(c *fiber.Ctx) error {

	var req BulkStatusRequest

	// ==========================
	// Parse Body
	// ==========================

	if err := c.BodyParser(&req); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Request tidak valid.",
			},
		)

	}

	// ==========================
	// Validasi
	// ==========================

	if len(req.IDs) == 0 {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Pilih minimal satu pengguna.",
			},
		)

	}

	// ==========================
	// Update
	// ==========================

	if err := BulkDeactivate(req.IDs); err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)

	}

	// ==========================
	// Success
	// ==========================

	return c.JSON(
		fiber.Map{
			"success": true,
			"message": fmt.Sprintf(
				"%d Pengguna berhasil dinonaktifkan.",
				len(req.IDs),
			),
		},
	)

}

// =========================================
// Bulk Reset Password
// =========================================

func BulkResetPasswordHandler(c *fiber.Ctx) error {

	var req BulkStatusRequest

	// ==========================
	// Parse Body
	// ==========================

	if err := c.BodyParser(&req); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Request tidak valid.",
			},
		)

	}

	// ==========================
	// Validasi
	// ==========================

	if len(req.IDs) == 0 {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Pilih minimal satu pengguna.",
			},
		)

	}

	// ==========================
	// Reset Password
	// ==========================

	if err := BulkResetPassword(req.IDs); err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)

	}

	// ==========================
	// Success
	// ==========================

	return c.JSON(
		fiber.Map{
			"success": true,
			"message": fmt.Sprintf(
				"%d Password pengguna berhasil direset.",
				len(req.IDs),
			),
		},
	)

}

// =========================================
// Change Password
// =========================================

func ChangePasswordHandler(c *fiber.Ctx) error {

	claims := c.Locals("claims").(*auth.Claims)

	var req ChangePasswordRequest

	// ==========================
	// Parse Body
	// ==========================

	if err := c.BodyParser(&req); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Request tidak valid.",
			},
		)

	}

	// ==========================
	// Validasi
	// ==========================

	if req.CurrentPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Password lama wajib diisi.",
		})
	}

	if req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Password baru wajib diisi.",
		})
	}

	if req.ConfirmPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Konfirmasi password wajib diisi.",
		})
	}

	if req.NewPassword != req.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Konfirmasi password tidak sama.",
		})
	}

	if msg := validatePassword(req.NewPassword); msg != "" {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": msg,
			},
		)

	}

	if req.CurrentPassword == req.NewPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Password baru tidak boleh sama dengan password lama.",
		})
	}

	// ==========================
	// Ambil Password Lama
	// ==========================

	user, err := GetPasswordByID(claims.IDPengguna)
	if err != nil {

		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"success": false,
				"message": "Pengguna tidak ditemukan.",
			},
		)

	}

	// ==========================
	// Verifikasi Password Lama
	// ==========================

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.CurrentPassword),
	)

	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Password lama salah.",
			},
		)

	}

	// ==========================
	// Hash Password Baru
	// ==========================

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": "Gagal membuat password.",
			},
		)

	}

	// ==========================
	// Simpan
	// ==========================

	err = ChangePassword(
		claims.IDPengguna,
		string(hashedPassword),
	)

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": "Gagal mengubah password.",
			},
		)

	}

	// ==========================
	// Success
	// ==========================

	return c.JSON(
		fiber.Map{
			"success": true,
			"message": "Password berhasil diubah.",
		},
	)

}

// =========================================
// Foto Pengguna (profil sendiri)
// =========================================

// FotoDefaultPengguna: nama file default LAMA (sebelum ada L.jpg/P.jpg
// per jenis kelamin). Sudah TIDAK PERNAH ditulis ke DB lagi mulai sekarang
// (foto kosong = default gender-based via buildFotoURLPengguna), konstanta
// ini disisakan hanya untuk kompatibilitas data lama yang kadung
// tersimpan foto = FotoDefaultPengguna di DB.
const FotoDefaultPengguna = "1636812404_b39deec2f158cb789567.jpg"

var allowedFotoExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

const maxFotoSize = 100 * 1024 // 100 KB

func deleteFotoFileWithRetry(path string) {

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
			fmt.Println("Gagal menghapus file foto lama:", err)
		}

		time.Sleep(retryDelay)
	}
}

// validateFotoUpload mengecek ekstensi & ukuran file foto.
// Mengembalikan ekstensi (mis. ".jpg") kalau valid, atau pesan error kalau tidak.
func validateFotoUpload(file *multipart.FileHeader) (string, string) {

	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !allowedFotoExt[ext] {
		return "", "Foto harus berformat JPG, JPEG, PNG atau WEBP."
	}

	if file.Size > maxFotoSize {
		return "", "Ukuran foto maksimal 100 KB."
	}

	return ext, ""
}

// saveFotoFile menyimpan file ke storage/fotopengguna dgn nama uuid.NewV7() baru.
func saveFotoFile(c *fiber.Ctx, file *multipart.FileHeader, ext string) (string, error) {

	id, err := uuid.NewV7()

	if err != nil {
		return "", err
	}

	namaFile := id.String() + ext

	if err := c.SaveFile(file, "./storage/fotopengguna/"+namaFile); err != nil {
		return "", err
	}

	return namaFile, nil
}

func UpdateFotoPenggunaHandler(c *fiber.Ctx) error {

	claims := c.Locals("claims").(*auth.Claims)

	// ==========================
	// Ambil File
	// ==========================

	file, err := c.FormFile("foto")

	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "File foto wajib dipilih.",
			},
		)
	}

	ext, msg := validateFotoUpload(file)

	if msg != "" {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": msg,
			},
		)
	}

	// ==========================
	// Foto Lama
	// ==========================

	fotoLama, err := GetFotoPengguna(claims.IDPengguna)

	if err != nil {

		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"success": false,
				"message": "Pengguna tidak ditemukan.",
			},
		)
	}

	// ==========================
	// Simpan File Baru
	// ==========================

	namaFile, err := saveFotoFile(c, file, ext)

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": "Gagal menyimpan file foto.",
			},
		)
	}

	// ==========================
	// Update Database
	// ==========================

	if err := UpdateFotoPengguna(claims.IDPengguna, namaFile); err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": "Gagal memperbarui foto pengguna.",
			},
		)
	}

	// ==========================
	// Hapus File Lama
	// (kalau bukan foto default)
	// ==========================

	if fotoLama != "" && fotoLama != FotoDefaultPengguna {

		go deleteFotoFileWithRetry("./storage/fotopengguna/" + fotoLama)
	}

	// ==========================
	// Success
	// ==========================

	return c.JSON(
		fiber.Map{
			"success": true,
			"message": "Foto profil berhasil diubah.",
			"foto_url": config.App.AppURL +
				"/storage/fotopengguna/" +
				namaFile,
		},
	)
}

func DeleteFotoPenggunaHandler(c *fiber.Ctx) error {

	claims := c.Locals("claims").(*auth.Claims)

	fotoLama, nuptk, nis, err := GetFotoInfoPengguna(claims.IDPengguna)

	if err != nil {

		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"success": false,
				"message": "Pengguna tidak ditemukan.",
			},
		)
	}

	fotoDefaultURL := buildFotoURLPengguna("", nuptk, nis)

	if fotoLama == "" || fotoLama == FotoDefaultPengguna {

		return c.JSON(
			fiber.Map{
				"success":  true,
				"message":  "Foto profil sudah default.",
				"foto_url": fotoDefaultURL,
			},
		)
	}

	// ==========================
	// Kembalikan ke Default
	// (foto dikosongkan, default L.jpg/P.jpg
	// ditentukan saat tampil sesuai jenis kelamin)
	// ==========================

	if err := UpdateFotoPengguna(claims.IDPengguna, ""); err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": "Gagal menghapus foto pengguna.",
			},
		)
	}

	// ==========================
	// Hapus File Lama
	// ==========================

	go deleteFotoFileWithRetry("./storage/fotopengguna/" + fotoLama)

	// ==========================
	// Success
	// ==========================

	return c.JSON(
		fiber.Map{
			"success":  true,
			"message":  "Foto profil berhasil dihapus.",
			"foto_url": fotoDefaultURL,
		},
	)
}

// =========================================
// Foto Pengguna (Admin - kelola pengguna lain,
// dipakai di Tambah.vue & Ubah.vue)
// =========================================

func UpdateFotoPenggunaByIDHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	if id == "" {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "ID pengguna tidak valid.",
			},
		)
	}

	// ==========================
	// Ambil File
	// ==========================

	file, err := c.FormFile("foto")

	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "File foto wajib dipilih.",
			},
		)
	}

	ext, msg := validateFotoUpload(file)

	if msg != "" {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": msg,
			},
		)
	}

	// ==========================
	// Foto Lama
	// ==========================

	fotoLama, err := GetFotoPengguna(id)

	if err != nil {

		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"success": false,
				"message": "Pengguna tidak ditemukan.",
			},
		)
	}

	// ==========================
	// Simpan File Baru
	// ==========================

	namaFile, err := saveFotoFile(c, file, ext)

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": "Gagal menyimpan file foto.",
			},
		)
	}

	// ==========================
	// Update Database
	// ==========================

	if err := UpdateFotoPengguna(id, namaFile); err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": "Gagal memperbarui foto pengguna.",
			},
		)
	}

	// ==========================
	// Hapus File Lama
	// (kalau bukan foto default)
	// ==========================

	if fotoLama != "" && fotoLama != FotoDefaultPengguna {

		go deleteFotoFileWithRetry("./storage/fotopengguna/" + fotoLama)
	}

	// ==========================
	// Success
	// ==========================

	return c.JSON(
		fiber.Map{
			"success": true,
			"message": "Foto pengguna berhasil diubah.",
			"foto_url": config.App.AppURL +
				"/storage/fotopengguna/" +
				namaFile,
		},
	)
}

func DeleteFotoPenggunaByIDHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	if id == "" {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "ID pengguna tidak valid.",
			},
		)
	}

	fotoLama, nuptk, nis, err := GetFotoInfoPengguna(id)

	if err != nil {

		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"success": false,
				"message": "Pengguna tidak ditemukan.",
			},
		)
	}

	fotoDefaultURL := buildFotoURLPengguna("", nuptk, nis)

	if fotoLama == "" || fotoLama == FotoDefaultPengguna {

		return c.JSON(
			fiber.Map{
				"success":  true,
				"message":  "Foto pengguna sudah default.",
				"foto_url": fotoDefaultURL,
			},
		)
	}

	// ==========================
	// Kembalikan ke Default
	// (foto dikosongkan, default L.jpg/P.jpg
	// ditentukan saat tampil sesuai jenis kelamin)
	// ==========================

	if err := UpdateFotoPengguna(id, ""); err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": "Gagal menghapus foto pengguna.",
			},
		)
	}

	// ==========================
	// Hapus File Lama
	// ==========================

	go deleteFotoFileWithRetry("./storage/fotopengguna/" + fotoLama)

	// ==========================
	// Success
	// ==========================

	return c.JSON(
		fiber.Map{
			"success":  true,
			"message":  "Foto pengguna berhasil dihapus.",
			"foto_url": fotoDefaultURL,
		},
	)
}

// =========================================
// Bulk Delete Pengguna
// =========================================

func BulkDeleteHandler(c *fiber.Ctx) error {

	var req BulkStatusRequest

	// ==========================
	// Parse Body
	// ==========================

	if err := c.BodyParser(&req); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Request tidak valid.",
			},
		)

	}

	// ==========================
	// Validasi
	// ==========================

	if len(req.IDs) == 0 {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Pilih minimal satu pengguna.",
			},
		)

	}

	// ==========================
	// Cek Super Admin
	// ==========================

	hasAdmin, err := HasAdmin(req.IDs)

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": "Gagal memvalidasi pengguna.",
			},
		)

	}

	if hasAdmin {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Super Admin tidak dapat dihapus.",
			},
		)

	}

	// ==========================
	// Ambil Foto Lama
	// (sebelum baris pengguna-nya dihapus dari DB)
	// ==========================

	fotoLama, err := GetFotosByIDs(req.IDs)

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": "Gagal membaca data foto pengguna.",
			},
		)

	}

	// ==========================
	// Hapus Pengguna
	// ==========================

	if err := BulkDelete(req.IDs); err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)

	}

	// ==========================
	// Hapus File Foto Lama
	// (kalau bukan foto default)
	// ==========================

	for _, foto := range fotoLama {

		if foto != "" && foto != FotoDefaultPengguna {
			go deleteFotoFileWithRetry("./storage/fotopengguna/" + foto)
		}
	}

	// ==========================
	// Success
	// ==========================

	return c.JSON(
		fiber.Map{
			"success": true,
			"message": fmt.Sprintf(
				"%d Pengguna berhasil dihapus.",
				len(req.IDs),
			),
		},
	)

}

// =========================================
// Upload Excel
// =========================================

func UploadHandler(c *fiber.Ctx) error {

	file, err := c.FormFile("file")

	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "File Excel wajib dipilih.",
			},
		)
	}

	src, err := file.Open()

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": "Gagal membuka file.",
			},
		)
	}

	defer src.Close()

	f, err := excelize.OpenReader(src)

	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Format file tidak valid.",
			},
		)
	}

	defer func() {
		_ = f.Close()
	}()

	// =========================================
	// Baca Excel
	// =========================================

	data, err := ReadUploadExcel(f)

	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	// =========================================
	// Validasi
	// =========================================

	validation, err := ValidateUpload(data)

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	if len(validation) > 0 {

		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "Upload dibatalkan karena terdapat kesalahan.",
				"errors":  validation,
			},
		)
	}

	// =========================================
	// Simpan ke Database
	// =========================================

	err = BulkInsert(data)

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	// =========================================
	// Success
	// =========================================

	return c.JSON(
		fiber.Map{
			"success": true,
			"message": fmt.Sprintf(
				"%d pengguna berhasil diupload.",
				len(data),
			),
		},
	)
}
