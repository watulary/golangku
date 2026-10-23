package visimisi

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
			return c.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"success": false,
					"message": "Data visi misi tidak ditemukan",
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
				"message": "Data visi misi tidak ditemukan.",
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

	var visimisi VisiMisi

	if err := c.BodyParser(&visimisi); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	if strings.TrimSpace(visimisi.Judul) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Judul visi misi wajib diisi.",
		})
	}

	if strings.TrimSpace(visimisi.Keterangan) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Keterangan visi misi wajib diisi.",
		})
	}

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid.",
		})
	}

	visimisi.IDPengguna = claims.IDPengguna

	oldVisiMisi, err := Get()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membaca data visi misi.",
		})
	}

	visimisi.ID = oldVisiMisi.ID

	// =====================================
	// Cek perubahan
	// =====================================

	changed := oldVisiMisi.Judul != visimisi.Judul ||
		oldVisiMisi.Keterangan != visimisi.Keterangan

	if !changed {
		return c.JSON(fiber.Map{
			"success": true,
			"changed": false,
			"message": "Tidak ada perubahan data.",
			"data":    oldVisiMisi,
		})
	}

	if err := UpdateDashboard(&visimisi); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if oldVisiMisi.Keterangan != visimisi.Keterangan {
		tiptap.CleanupOrphanImages(oldVisiMisi.Keterangan, visimisi.Keterangan)
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
		"message": "Visi misi berhasil diperbarui.",
		"data":    data,
	})
}
