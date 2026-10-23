package statistik

import (
	"sekolah/internal/auth"

	"github.com/gofiber/fiber/v2"
)

// GetPublicHandler dipakai halaman publik (route /api/public/statistik,
// TANPA middleware.Auth()) — selalu tampilkan total keseluruhan tanpa
// filter, karena visitor anonim tidak punya claims sama sekali.
func GetPublicHandler(c *fiber.Ctx) error {
	data, err := GetStatistik("", true)

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
		"data":    data,
	})
}

// GetHandler dipakai dashboard (route /api/statistik di grup "api" yg
// sudah middleware.Auth()) — statistik konten website difilter per
// id_pengguna pembuat, kecuali Super Admin yg lihat semua.
func GetHandler(c *fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Tidak terautentikasi.",
		})
	}
	isSuperAdmin := claims.Level == "Admin"

	data, err := GetStatistik(claims.IDPengguna, isSuperAdmin)

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
		"data":    data,
	})
}

func GetRekapRombelHandler(c *fiber.Ctx) error {
	data, err := GetRekapRombel()
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
