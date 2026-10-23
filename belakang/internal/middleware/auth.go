package middleware

import (
	"strings"

	"sekolah/internal/auth"
	"sekolah/internal/modules/pengguna"

	"github.com/gofiber/fiber/v2"
)

func Auth() fiber.Handler {

	return func(c *fiber.Ctx) error {

		// Ambil Authorization Header
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token tidak ditemukan",
			})
		}

		// Format harus:
		// Authorization: Bearer xxxxxxxxx
		const bearer = "Bearer "

		if !strings.HasPrefix(authHeader, bearer) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Format token tidak valid",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, bearer)

		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token tidak valid",
			})
		}

		// =========================================
		// Cek apakah akun masih aktif
		// =========================================

		user, err := pengguna.GetByID(claims.IDPengguna)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Pengguna tidak ditemukan",
			})
		}

		if user.Aktif == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"code":    "ACCOUNT_DISABLED",
				"message": "Akun Anda telah dinonaktifkan.",
			})
		}

		// Simpan Claims agar bisa dipakai handler berikutnya
		c.Locals("claims", claims)

		return c.Next()
	}
}
