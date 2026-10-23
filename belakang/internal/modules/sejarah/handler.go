package sejarah

import (
	"database/sql"
	"sekolah/internal/auth"
	"sekolah/internal/modules/tiptap"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetHandler(c *fiber.Ctx) error {

	data, err := Get()

	if err != nil {

		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Data sejarah tidak ditemukan",
			})
		}

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

func IncrementHitsHandler(c *fiber.Ctx) error {

	err := IncrementHits()

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

func GetDashboardHandler(c *fiber.Ctx) error {

	data, err := Get()

	if err != nil {

		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Data sejarah tidak ditemukan.",
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

func UpdateDashboardHandler(c *fiber.Ctx) error {

	var sejarah Sejarah

	if err := c.BodyParser(&sejarah); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	if strings.TrimSpace(sejarah.Judul) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Judul sejarah wajib diisi.",
		})
	}

	if strings.TrimSpace(sejarah.Keterangan) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Keterangan sejarah wajib diisi.",
		})
	}

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid.",
		})
	}

	sejarah.IDPengguna = claims.IDPengguna

	// Ambil data lama untuk preserve ID
	oldSejarah, err := Get()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membaca data sejarah.",
		})
	}

	sejarah.ID = oldSejarah.ID

	// =====================================
	// Cek perubahan
	// =====================================

	changed := oldSejarah.Judul != sejarah.Judul ||
		oldSejarah.Keterangan != sejarah.Keterangan

	if !changed {
		return c.JSON(fiber.Map{
			"success": true,
			"changed": false,
			"message": "Tidak ada perubahan data.",
			"data":    oldSejarah,
		})
	}

	if err := UpdateDashboard(&sejarah); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if oldSejarah.Keterangan != sejarah.Keterangan {
		tiptap.CleanupOrphanImages(oldSejarah.Keterangan, sejarah.Keterangan)
	}

	data, err := Get()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"changed": true,
		"message": "Sejarah berhasil diperbarui.",
		"data":    data,
	})
}
