package auth

import (
	"errors"
	"sekolah/internal/security"
	"strconv"

	"time"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {

	var req LoginRequest

	// Parse request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(LoginResponse{
			Success: false,
			Message: "Request tidak valid.",
		})
	}

	// Cek apakah login sedang diblok
	blocked, retryAfter := security.CheckLoginLimit(c, req.Email)

	if blocked {

		c.Set("Retry-After", strconv.Itoa(retryAfter))

		return c.Status(fiber.StatusTooManyRequests).JSON(LoginResponse{
			Success: false,
			Message: "Login bermasalah, silahkan nanti login kembali.",
		})
	}

	// Cari user berdasarkan email
	user, err := GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {

			// Tambah counter gagal login
			security.IncrementLoginFail(c, req.Email)

			return c.Status(fiber.StatusUnauthorized).JSON(LoginResponse{
				Success: false,
				Message: "Email atau password salah.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(LoginResponse{
			Success: false,
			Message: "Terjadi kesalahan pada server.",
		})
	}

	// Cek akun aktif
	if user.Aktif == 0 {
		return c.Status(fiber.StatusForbidden).JSON(LoginResponse{
			Success: false,
			Message: "Akun anda tidak aktif, hubungi administrator.",
		})
	}

	// Cek password
	if err := CheckPassword(user.Password, req.Password); err != nil {

		// Tambah counter gagal login
		security.IncrementLoginFail(c, req.Email)

		return c.Status(fiber.StatusUnauthorized).JSON(LoginResponse{
			Success: false,
			Message: "Email atau password salah.",
		})
	}

	// Login berhasil, reset counter
	security.ResetLoginFail(c, req.Email)

	// Update last login
	if err := UpdateLastLogin(user.IDPengguna, time.Now()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(LoginResponse{
			Success: false,
			Message: "Gagal memperbarui waktu login.",
		})
	}

	// Generate JWT
	token, err := GenerateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(LoginResponse{
			Success: false,
			Message: "Gagal membuat token.",
		})
	}

	// Response
	return c.Status(fiber.StatusOK).JSON(LoginResponse{
		Success: true,
		Message: "Login berhasil.",
		Token:   token,
		User:    user,
	})
}

func Me(c *fiber.Ctx) error {

	claims, ok := c.Locals("claims").(*Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized.",
		})
	}

	user, err := GetByID(claims.IDPengguna)
	if err != nil {

		if errors.Is(err, ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Pengguna tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Terjadi kesalahan pada server.",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}
