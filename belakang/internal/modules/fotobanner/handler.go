package fotobanner

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"sekolah/internal/auth"
	"sekolah/internal/config"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Get(c *fiber.Ctx) error {

	data, err := GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	for i := range data {
		data[i].GambarURL = buildGambarURL(data[i].Gambar)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func GetAllHandler(c *fiber.Ctx) error {

	data, err := GetAll()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	for i := range data {
		data[i].GambarURL = buildGambarURL(data[i].Gambar)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func GetAdminListHandler(c *fiber.Ctx) error {

	search := c.Query("search", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sort_by", "waktu")
	sortDir := c.Query("sort_dir", "desc")

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid.",
		})
	}

	idPengguna := ""

	if claims.Level != "Admin" {
		idPengguna = claims.IDPengguna
	}

	data, err := GetAllAdmin(search, idPengguna, sortBy, sortDir, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	total, err := CountAdmin(search, idPengguna)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	for i := range data {
		data[i].GambarURL = buildGambarURL(data[i].Gambar)
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

	id := c.Params("id")

	data, err := GetByID(id)
	if err != nil {

		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Banner tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	claims, allowed := isOwnerOrAdmin(c, data.IDPengguna)
	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid.",
		})
	}

	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Anda tidak memiliki akses ke banner ini.",
		})
	}

	data.GambarURL = buildGambarURL(data.Gambar)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func CreateHandler(c *fiber.Ctx) error {

	var input FotoBanner

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	if strings.TrimSpace(input.Nama) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nama wajib diisi.",
		})
	}

	if strings.TrimSpace(input.Kegiatan) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Kegiatan wajib diisi.",
		})
	}

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid.",
		})
	}

	id, err := uuid.NewV7()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	input.ID = id.String()
	input.IDPengguna = claims.IDPengguna

	file, err := c.FormFile("gambar")
	if err != nil || file == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Gambar wajib diunggah.",
		})
	}

	if msg := validateGambarBanner(file); msg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": msg,
		})
	}

	filename := input.ID + filepath.Ext(file.Filename)

	if err := c.SaveFile(file, "./storage/fotobanner/"+filename); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengunggah gambar.",
		})
	}

	input.Gambar = filename

	if err := Create(input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	input.GambarURL = buildGambarURL(input.Gambar)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Banner berhasil ditambahkan.",
		"data":    input,
	})
}

func UpdateHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	var input FotoBanner

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	if strings.TrimSpace(input.Nama) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nama wajib diisi.",
		})
	}

	if strings.TrimSpace(input.Kegiatan) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Kegiatan wajib diisi.",
		})
	}

	oldData, err := GetByID(id)
	if err != nil {

		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Banner tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	claims, allowed := isOwnerOrAdmin(c, oldData.IDPengguna)
	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid.",
		})
	}

	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Anda tidak memiliki akses untuk mengubah banner ini.",
		})
	}

	changed := oldData.Nama != input.Nama ||
		oldData.Kegiatan != input.Kegiatan

	file, ferr := c.FormFile("gambar")
	hasNewGambar := ferr == nil && file != nil

	if hasNewGambar {
		if msg := validateGambarBanner(file); msg != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": msg,
			})
		}
	}

	if !changed && !hasNewGambar {
		return c.JSON(fiber.Map{
			"success": true,
			"changed": false,
			"message": "Tidak ada perubahan data.",
			"data":    oldData,
		})
	}

	if changed {
		if err := Update(id, input); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
	}

	if hasNewGambar {

		newID, err := uuid.NewV7()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		filename := newID.String() + filepath.Ext(file.Filename)

		if err := c.SaveFile(file, "./storage/fotobanner/"+filename); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal mengunggah gambar.",
			})
		}

		if err := UpdateGambar(id, filename); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		go deleteGambarFile(oldData.Gambar)
	}

	data, err := GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	data.GambarURL = buildGambarURL(data.Gambar)

	return c.JSON(fiber.Map{
		"success": true,
		"changed": true,
		"message": "Banner berhasil diperbarui.",
		"data":    data,
	})
}

func DeleteHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	data, err := GetByID(id)
	if err != nil {

		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Banner tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	claims, allowed := isOwnerOrAdmin(c, data.IDPengguna)
	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid.",
		})
	}

	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Anda tidak memiliki akses untuk menghapus banner ini.",
		})
	}

	if err := Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	go deleteGambarFile(data.Gambar)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Banner berhasil dihapus.",
	})
}

func buildGambarURL(gambar string) string {

	if gambar == "" {
		return ""
	}

	return config.App.AppURL + "/storage/fotobanner/" + gambar
}

var allowedGambarBannerExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

const maxGambarBannerSize = 100 * 1024 // 100 KB

// validateGambarBanner mengecek ekstensi & ukuran file gambar banner.
// Mengembalikan pesan error kalau tidak valid, string kosong kalau valid.
func validateGambarBanner(file *multipart.FileHeader) string {

	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !allowedGambarBannerExt[ext] {
		return "Gambar harus berformat JPG, JPEG, PNG atau WEBP."
	}

	if file.Size > maxGambarBannerSize {
		return "Ukuran gambar maksimal 100 KB."
	}

	return ""
}

func deleteGambarFile(gambar string) {

	if gambar == "" {
		return
	}

	path := "./storage/fotobanner/" + gambar

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
			fmt.Println("Gagal menghapus file gambar lama:", err)
		}

		time.Sleep(retryDelay)
	}
}

func isOwnerOrAdmin(c *fiber.Ctx, idPenggunaGambar string) (*auth.Claims, bool) {

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return nil, false
	}

	if claims.Level == "Admin" {
		return claims, true
	}

	return claims, claims.IDPengguna == idPenggunaGambar
}
