package video

import (
	"database/sql"
	"regexp"
	"sekolah/internal/auth"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var youtubeIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{11}$`)

func Get(c *fiber.Ctx) error {

	data, err := GetLatest(3)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	for i := range data {
		buildVideoURLs(&data[i])
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func GetAllHandler(c *fiber.Ctx) error {

	data, err := GetAll()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	for i := range data {
		buildVideoURLs(&data[i])
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func GetAdminListHandler(c *fiber.Ctx) error {

	search := c.Query("search", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sort_by", "waktu")
	sortDir := c.Query("sort_dir", "desc")

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid.",
		})
	}

	idPengguna := ""

	if claims.Level != "Admin" {
		idPengguna = claims.IDPengguna
	}

	data, err := GetAllAdmin(search, idPengguna, sortBy, sortDir, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	total, err := CountAdmin(search, idPengguna)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	for i := range data {
		buildVideoURLs(&data[i])
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func GetByIDHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	data, err := GetByID(id)
	if err != nil {

		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Video tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	claims, allowed := isOwnerOrAdmin(c, data.IDPengguna)
	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid.",
		})
	}

	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Anda tidak memiliki akses ke video ini.",
		})
	}

	buildVideoURLs(data)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func CreateHandler(c *fiber.Ctx) error {

	var input Video

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	if strings.TrimSpace(input.Nama) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nama wajib diisi.",
		})
	}

	if strings.TrimSpace(input.Kegiatan) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Kegiatan wajib diisi.",
		})
	}

	input.Video = strings.TrimSpace(input.Video)

	if input.Video == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID video YouTube wajib diisi.",
		})
	}

	if !youtubeIDRegex.MatchString(input.Video) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID video YouTube tidak valid (harus 11 karakter).",
		})
	}

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid.",
		})
	}

	id, err := uuid.NewV7()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	input.ID = id.String()
	input.IDPengguna = claims.IDPengguna

	if err := Create(input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	buildVideoURLs(&input)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Video berhasil ditambahkan.",
		"data":    input,
	})
}

func UpdateHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	var input Video

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	if strings.TrimSpace(input.Nama) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nama wajib diisi.",
		})
	}

	if strings.TrimSpace(input.Kegiatan) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Kegiatan wajib diisi.",
		})
	}

	input.Video = strings.TrimSpace(input.Video)

	if input.Video == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID video YouTube wajib diisi.",
		})
	}

	if !youtubeIDRegex.MatchString(input.Video) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID video YouTube tidak valid (harus 11 karakter).",
		})
	}

	oldData, err := GetByID(id)
	if err != nil {

		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Video tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	claims, allowed := isOwnerOrAdmin(c, oldData.IDPengguna)
	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid.",
		})
	}

	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Anda tidak memiliki akses untuk mengubah video ini.",
		})
	}

	changed := oldData.Nama != input.Nama ||
		oldData.Kegiatan != input.Kegiatan ||
		oldData.Video != input.Video

	if !changed {
		buildVideoURLs(oldData)

		return c.JSON(fiber.Map{
			"success": true,
			"changed": false,
			"message": "Tidak ada perubahan data.",
			"data":    oldData,
		})
	}

	if err := Update(id, input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	data, err := GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	buildVideoURLs(data)

	return c.JSON(fiber.Map{
		"success": true,
		"changed": true,
		"message": "Video berhasil diperbarui.",
		"data":    data,
	})
}

func DeleteHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	data, err := GetByID(id)
	if err != nil {

		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Video tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	claims, allowed := isOwnerOrAdmin(c, data.IDPengguna)
	if claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid.",
		})
	}

	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Anda tidak memiliki akses untuk menghapus video ini.",
		})
	}

	if err := Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Video berhasil dihapus.",
	})
}

func buildVideoURLs(v *Video) {
	if v.Video == "" {
		return
	}
	v.YoutubeURL = "https://www.youtube.com/embed/" + v.Video
	v.Thumbnail = "https://img.youtube.com/vi/" + v.Video + "/mqdefault.jpg"
}

func isOwnerOrAdmin(c *fiber.Ctx, idPenggunaVideo string) (*auth.Claims, bool) {

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return nil, false
	}

	if claims.Level == "Admin" {
		return claims, true
	}

	return claims, claims.IDPengguna == idPenggunaVideo
}
