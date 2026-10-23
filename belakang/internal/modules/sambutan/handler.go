package sambutan

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sekolah/internal/auth"
	"sekolah/internal/config"
	"sekolah/internal/modules/tiptap"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Get(c *fiber.Ctx) error {

	data, err := GetSambutan()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if data.FotoKepsek != "" {
		data.FotoURL =
			config.App.AppURL +
				"/storage/logo/" +
				data.FotoKepsek
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func IncrementHits(c *fiber.Ctx) error {

	err := UpdateHits()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

// =====================================
// Dashboard
// =====================================

func GetDashboardHandler(c *fiber.Ctx) error {

	data, err := GetSambutan()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if data.FotoKepsek != "" {
		data.FotoURL =
			config.App.AppURL +
				"/storage/logo/" +
				data.FotoKepsek
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func UpdateDashboardHandler(c *fiber.Ctx) error {

	const defaultFoto = "1636812404_b39deec2f158cb789567.jpg"

	sambutan := parseSambutanForm(c)

	if strings.TrimSpace(sambutan.JudulSambutan) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Judul sambutan wajib diisi.",
		})
	}

	if strings.TrimSpace(sambutan.IsiSambutan) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Isi sambutan wajib diisi.",
		})
	}

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid.",
		})
	}

	sambutan.IDPengguna = claims.IDPengguna

	// Ambil data lama
	oldSambutan, err := GetSambutan()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membaca data sambutan.",
		})
	}

	sambutan.IDSambutan = oldSambutan.IDSambutan
	sambutan.FotoKepsek = oldSambutan.FotoKepsek

	var newFotoPath string
	var oldFotoPath string

	// Hapus foto: kembalikan ke foto default (file default TIDAK pernah dihapus).
	// File lama dihapus setelah DB berhasil diupdate. Bila pada request yang sama
	// ada foto baru, blok upload di bawah menimpanya.
	if c.FormValue("hapus_foto") == "1" && oldSambutan.FotoKepsek != defaultFoto {
		sambutan.FotoKepsek = defaultFoto

		if oldSambutan.FotoKepsek != "" {
			oldFotoPath = filepath.Join(
				"storage",
				"logo",
				filepath.Base(oldSambutan.FotoKepsek),
			)
		}
	}

	fotoFile, err := c.FormFile("foto")

	if err == nil && fotoFile != nil {

		const maxSize = 100 * 1024

		if fotoFile.Size > maxSize {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Ukuran foto maksimal 100 KB.",
			})
		}

		ext := strings.ToLower(filepath.Ext(fotoFile.Filename))

		switch ext {
		case ".jpg", ".jpeg", ".png", ".webp":
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Foto harus berformat JPG, JPEG, PNG atau WEBP.",
			})
		}

		if err := os.MkdirAll(filepath.Join("storage", "logo"), 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal menyiapkan folder penyimpanan foto.",
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

		newFotoPath = filepath.Join("storage", "logo", filename)

		src, err := fotoFile.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal membuka file foto.",
			})
		}
		defer src.Close()

		out, err := os.Create(newFotoPath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal menyimpan foto.",
			})
		}

		if _, err := io.Copy(out, src); err != nil {
			out.Close()
			_ = os.Remove(newFotoPath)

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal menyimpan foto.",
			})
		}

		if err := out.Close(); err != nil {
			_ = os.Remove(newFotoPath)

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal menyimpan foto.",
			})
		}

		sambutan.FotoKepsek = filename

		// Foto default tidak boleh ikut terhapus saat diganti foto baru.
		if oldSambutan.FotoKepsek != "" && oldSambutan.FotoKepsek != defaultFoto {
			oldFotoPath = filepath.Join(
				"storage",
				"logo",
				filepath.Base(oldSambutan.FotoKepsek),
			)
		}
	}

	// =====================================
	// Cek perubahan
	// =====================================

	changed := oldSambutan.JudulSambutan != sambutan.JudulSambutan ||
		oldSambutan.IsiSambutan != sambutan.IsiSambutan ||
		oldSambutan.FotoKepsek != sambutan.FotoKepsek

	if !changed {
		// Tidak ada perubahan: hapus file baru yang barusan diunggah.
		if newFotoPath != "" {
			_ = os.Remove(newFotoPath)
		}

		if oldSambutan.FotoKepsek != "" {
			oldSambutan.FotoURL = config.App.AppURL + "/storage/logo/" + oldSambutan.FotoKepsek
		}

		return c.JSON(fiber.Map{
			"success": true,
			"changed": false,
			"message": "Tidak ada perubahan data.",
			"data":    oldSambutan,
		})
	}

	err = UpdateDashboard(&sambutan)
	if err != nil {

		if newFotoPath != "" {
			_ = os.Remove(newFotoPath)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if oldFotoPath != "" && oldFotoPath != newFotoPath {
		removeFileEventually(oldFotoPath)
	}

	if oldSambutan.IsiSambutan != sambutan.IsiSambutan {
		tiptap.CleanupOrphanImages(oldSambutan.IsiSambutan, sambutan.IsiSambutan)
	}

	data, err := GetSambutan()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if data.FotoKepsek != "" {
		data.FotoURL = config.App.AppURL + "/storage/logo/" + data.FotoKepsek
	}

	return c.JSON(fiber.Map{
		"success": true,
		"changed": true,
		"message": "Sambutan berhasil diperbarui.",
		"data":    data,
	})
}

func parseSambutanForm(c *fiber.Ctx) Sambutan {
	return Sambutan{
		JudulSambutan: c.FormValue("judul"),
		IsiSambutan:   c.FormValue("isi"),
	}
}

func removeFileEventually(path string) {
	go func() {
		time.Sleep(3 * time.Second)

		const maxAttempt = 30
		const retryDelay = time.Second

		for attempt := 1; attempt <= maxAttempt; attempt++ {
			err := os.Remove(path)

			if err == nil || os.IsNotExist(err) {
				return
			}

			time.Sleep(retryDelay)
		}

		fmt.Println("Gagal menghapus foto lama:", path)
	}()
}
