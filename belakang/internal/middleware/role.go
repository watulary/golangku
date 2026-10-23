package middleware

import (
	"strings"

	"sekolah/internal/auth"

	"github.com/gofiber/fiber/v2"
)

func RequireRole(roles ...string) fiber.Handler {

	return func(c *fiber.Ctx) error {

		claims, ok := c.Locals("claims").(*auth.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized.",
			})
		}

		for _, role := range roles {

			if strings.EqualFold(claims.Level, role) {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Anda tidak memiliki hak akses.",
		})
	}
}
