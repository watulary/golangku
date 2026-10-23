package unggahdokumen

import (
	"fmt"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"sekolah/internal/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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

var allowedDokumenExt = map[string]bool{
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".ppt":  true,
	".pptx": true,
}

const maxDokumenSize = 2 * 1024 * 1024 // 2MB

func isOwnerOrAdmin(c *fiber.Ctx, idPenggunaData string) (*auth.Claims, bool) {
	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return nil, false
	}
	if claims.Level == "Admin" {
		return claims, true
	}
	return claims, claims.IDPengguna == idPenggunaData
}

func GetAdminListHandler(c *fiber.Ctx) error {

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	search := c.Query("q", "")
	sortBy := c.Query("sortBy", "waktu")
	sortDir := c.Query("sortDir", "desc")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	idPengguna := ""
	if claims.Level != "Admin" {
		idPengguna = claims.IDPengguna
	}

	data, err := GetAdmin(search, idPengguna, sortBy, sortDir, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	total, err := CountAdmin(search, idPengguna)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
		"pagination": fiber.Map{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_page": totalPage,
		},
	})
}

// GetByIDHandler mengambil satu dokumen berdasarkan id
func GetByIDHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	data, err := GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Dokumen tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// CreateHandler menyimpan dokumen baru (nama + file)
func CreateHandler(c *fiber.Ctx) error {

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	nama := c.FormValue("nama")
	if nama == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Nama dokumen wajib diisi",
		})
	}

	file, err := c.FormFile("dokumen")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "File dokumen wajib diunggah",
		})
	}

	filename, err := saveDokumenFile(c, file)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	data := &Unggahdokumen{
		ID:         uuid.Must(uuid.NewV7()).String(),
		IDPengguna: claims.IDPengguna,
		Nama:       nama,
		Dokumen:    filename,
	}

	if err := Create(data); err != nil {
		go deleteDokumenFile(filename)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Dokumen berhasil ditambahkan",
		"data":    data,
	})
}

// UpdateHandler mengubah nama dokumen, dan file (opsional)
func UpdateHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	existing, err := GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Dokumen tidak ditemukan",
		})
	}

	if _, ok := isOwnerOrAdmin(c, existing.IDPengguna); !ok {
		return c.Status(403).JSON(fiber.Map{
			"success": false,
			"message": "Tidak memiliki akses",
		})
	}

	nama := c.FormValue("nama")
	changed := false

	if nama != "" && nama != existing.Nama {
		if err := Update(id, nama); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		changed = true
	}

	if file, err := c.FormFile("dokumen"); err == nil {
		filename, err := saveDokumenFile(c, file)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		if err := UpdateDokumen(id, filename); err != nil {
			go deleteDokumenFile(filename)
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		go deleteDokumenFile(existing.Dokumen)
		changed = true
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Dokumen berhasil diperbarui",
		"changed": changed,
	})
}

// DeleteHandler menghapus dokumen
func DeleteHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	existing, err := GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Dokumen tidak ditemukan",
		})
	}

	if _, ok := isOwnerOrAdmin(c, existing.IDPengguna); !ok {
		return c.Status(403).JSON(fiber.Map{
			"success": false,
			"message": "Tidak memiliki akses",
		})
	}

	if err := Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	go deleteDokumenFile(existing.Dokumen)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Dokumen berhasil dihapus",
	})
}

// saveDokumenFile memvalidasi tipe & ukuran, lalu menyimpan file dengan nama uuid v7 baru
func saveDokumenFile(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedDokumenExt[ext] {
		return "", fmt.Errorf("tipe file tidak didukung (%s)", ext)
	}

	if file.Size > maxDokumenSize {
		return "", fmt.Errorf("ukuran file melebihi 2MB")
	}

	filename := uuid.Must(uuid.NewV7()).String() + ext
	dest := filepath.Join("storage", "dokumen", filename)

	if err := os.MkdirAll(filepath.Join("storage", "dokumen"), 0755); err != nil {
		return "", err
	}

	if err := c.SaveFile(file, dest); err != nil {
		return "", err
	}

	return filename, nil
}

// deleteDokumenFile menghapus file lama (retry sabar, dipanggil via goroutine)
func deleteDokumenFile(filename string) {

	if filename == "" {
		return
	}

	path := filepath.Join("storage", "dokumen", filename)

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
			fmt.Println("Gagal menghapus file dokumen lama:", err)
		}

		time.Sleep(retryDelay)
	}
}
