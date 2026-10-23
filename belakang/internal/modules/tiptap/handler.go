package tiptap

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UploadImage(c *fiber.Ctx) error {

	// =====================================================
	// Ambil File
	// =====================================================

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Gambar wajib dipilih.",
		})
	}

	// =====================================================
	// Maksimal 100 KB
	// =====================================================

	const maxSize = 100 * 1024

	if file.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Ukuran gambar maksimal 100 KB.",
		})
	}

	// =====================================================
	// Validasi Ekstensi
	// =====================================================

	ext := strings.ToLower(filepath.Ext(file.Filename))

	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format gambar harus JPG, JPEG, PNG atau WEBP.",
		})
	}

	// =====================================================
	// Buka File
	// =====================================================

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuka gambar.",
		})
	}
	defer src.Close()

	// =====================================================
	// Validasi MIME Type (Magic Bytes)
	// =====================================================

	buffer := make([]byte, 512)

	if _, err := src.Read(buffer); err != nil && err != io.EOF {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membaca gambar.",
		})
	}

	contentType := http.DetectContentType(buffer)

	switch contentType {
	case "image/jpeg",
		"image/png",
		"image/webp":
		// valid
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Tipe gambar tidak valid.",
		})
	}

	// Kembali ke awal file
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membaca gambar.",
		})
	}

	// =====================================================
	// Pastikan Folder Ada
	// =====================================================

	if err := os.MkdirAll(filepath.Join("storage", "tiptap"), 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuat folder penyimpanan.",
		})
	}

	// =====================================================
	// Generate Nama File (UUID v7)
	// =====================================================

	id, err := uuid.NewV7()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuat nama file.",
		})
	}

	filename := id.String() + ext

	dst := filepath.Join("storage", "tiptap", filename)

	// =====================================================
	// Simpan File
	// =====================================================

	out, err := os.Create(dst)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menyimpan gambar.",
		})
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {

		_ = os.Remove(dst)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menyimpan gambar.",
		})
	}

	// =====================================================
	// Response
	// =====================================================

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Gambar berhasil diunggah.",
		"url":     "/storage/tiptap/" + filename,
	})
}
