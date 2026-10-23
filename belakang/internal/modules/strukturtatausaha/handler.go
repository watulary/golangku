package strukturtatausaha

import (
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"sekolah/internal/auth"
	"sekolah/internal/config"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetHandler(c *fiber.Ctx) error {

	data, err := Get()

	if err != nil {

		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"success": false,
					"message": "Data struktur tata usaha tidak ditemukan",
				},
			)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	if data.Gambar != "" {
		data.FotoURL =
			config.App.AppURL +
				"/storage/fotostruktur/" +
				data.Gambar
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func IncrementHitsHandler(c *fiber.Ctx) error {

	err := IncrementHits()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func GetDashboardHandler(c *fiber.Ctx) error {

	data, err := Get()

	if err != nil {

		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Data struktur tata usaha tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if data.Gambar != "" {
		data.FotoURL = config.App.AppURL + "/storage/fotostruktur/" + data.Gambar
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func UpdateDashboardHandler(c *fiber.Ctx) error {

	const defaultGambar = "1636812404_b39deec2f158cb789567.jpg"

	judul := strings.TrimSpace(c.FormValue("judul"))

	if judul == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Judul struktur tata usaha wajib diisi.",
		})
	}

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid.",
		})
	}

	oldData, err := Get()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membaca data struktur tata usaha.",
		})
	}

	updated := StrukturTatausaha{
		ID:         oldData.ID,
		IDPengguna: claims.IDPengguna,
		Judul:      judul,
		Gambar:     oldData.Gambar,
	}

	var newGambarPath string
	var oldGambarPath string

	// Hapus gambar: kembalikan ke gambar default (file default TIDAK pernah dihapus).
	// File lama dihapus setelah DB berhasil diupdate. Bila pada request yang sama
	// ada gambar baru, blok upload di bawah menimpanya.
	if c.FormValue("hapus_foto") == "1" && oldData.Gambar != defaultGambar {
		updated.Gambar = defaultGambar

		if oldData.Gambar != "" {
			oldGambarPath = filepath.Join(
				"storage",
				"fotostruktur",
				filepath.Base(oldData.Gambar),
			)
		}
	}

	gambarFile, err := c.FormFile("gambar")

	if err == nil && gambarFile != nil {

		const maxSize = 100 * 1024

		if gambarFile.Size > maxSize {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Ukuran gambar maksimal 100 KB.",
			})
		}

		ext := strings.ToLower(filepath.Ext(gambarFile.Filename))

		switch ext {
		case ".jpg", ".jpeg", ".png", ".webp":
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Gambar harus berformat JPG, JPEG, PNG atau WEBP.",
			})
		}

		if err := os.MkdirAll(filepath.Join("storage", "fotostruktur"), 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal menyiapkan folder penyimpanan gambar.",
			})
		}

		id, err := uuid.NewV7()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal membuat nama file.",
			})
		}

		filename := id.String() + ext

		newGambarPath = filepath.Join("storage", "fotostruktur", filename)

		src, err := gambarFile.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal membuka file gambar.",
			})
		}
		defer src.Close()

		out, err := os.Create(newGambarPath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal menyimpan gambar.",
			})
		}

		if _, err := io.Copy(out, src); err != nil {
			out.Close()
			_ = os.Remove(newGambarPath)

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal menyimpan gambar.",
			})
		}

		if err := out.Close(); err != nil {
			_ = os.Remove(newGambarPath)

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal menyimpan gambar.",
			})
		}

		updated.Gambar = filename

		if oldData.Gambar != "" && oldData.Gambar != defaultGambar {
			oldGambarPath = filepath.Join("storage", "fotostruktur", oldData.Gambar)
		}
	}

	// =====================================
	// Cek perubahan
	// =====================================

	changed := oldData.Judul != updated.Judul ||
		oldData.Gambar != updated.Gambar

	if !changed {
		// Tidak ada perubahan: hapus file baru yang barusan diunggah.
		if newGambarPath != "" {
			_ = os.Remove(newGambarPath)
		}

		if oldData.Gambar != "" {
			oldData.FotoURL = config.App.AppURL + "/storage/fotostruktur/" + oldData.Gambar
		}

		return c.JSON(fiber.Map{
			"success": true,
			"changed": false,
			"message": "Tidak ada perubahan data.",
			"data":    oldData,
		})
	}

	if err := UpdateDashboard(&updated); err != nil {

		if newGambarPath != "" {
			_ = os.Remove(newGambarPath)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if oldGambarPath != "" && oldGambarPath != newGambarPath {
		_ = os.Remove(oldGambarPath)
	}

	data, err := Get()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if data.Gambar != "" {
		data.FotoURL = config.App.AppURL + "/storage/fotostruktur/" + data.Gambar
	}

	return c.JSON(fiber.Map{
		"success": true,
		"changed": true,
		"message": "Struktur tata usaha berhasil diperbarui.",
		"data":    data,
	})
}
