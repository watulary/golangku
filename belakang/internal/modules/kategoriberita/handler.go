package kategoriberita

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"sekolah/internal/auth"
	"sekolah/internal/utils"
)

func GetAllHandler(c *fiber.Ctx) error {

	data, err := GetAll()

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

func GetAllAdminHandler(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	search := c.Query("search", "")

	data, err := GetAllAdmin(limit, offset, search)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	total, err := CountAdmin(search)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if data == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "kategori berita tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func CreateHandler(c *fiber.Ctx) error {

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "sesi tidak valid, silakan login ulang",
		})
	}

	req := new(KategoriBeritaRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "data tidak valid",
		})
	}

	if req.NamaKategoriBerita == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "nama kategori wajib diisi",
		})
	}

	newSlug := utils.Slugify(req.NamaKategoriBerita)

	exist, err := IsSlugExist(newSlug, "")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if exist {
		newSlug = newSlug + "-" + uuid.New().String()[:6]
	}

	id, err := uuid.NewV7()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	data := &KategoriBerita{
		IDKategoriBerita:   id.String(),
		IDPengguna:         claims.IDPengguna,
		NamaKategoriBerita: req.NamaKategoriBerita,
		SlugKategoriBerita: newSlug,
		Urutan:             req.Urutan,
	}

	if err := Create(data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "kategori berita berhasil ditambahkan",
		"data":    data,
	})
}

func UpdateHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	existing, err := GetByID(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "kategori berita tidak ditemukan",
		})
	}

	req := new(KategoriBeritaRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "data tidak valid",
		})
	}

	if req.NamaKategoriBerita == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "nama kategori wajib diisi",
		})
	}

	newSlug := existing.SlugKategoriBerita

	if req.NamaKategoriBerita != existing.NamaKategoriBerita {

		newSlug = utils.Slugify(req.NamaKategoriBerita)

		exist, err := IsSlugExist(newSlug, id)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		if exist {
			newSlug = newSlug + "-" + uuid.New().String()[:6]
		}
	}

	existing.NamaKategoriBerita = req.NamaKategoriBerita
	existing.SlugKategoriBerita = newSlug
	existing.Urutan = req.Urutan

	if err := Update(existing); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "kategori berita berhasil diperbarui",
		"data":    existing,
	})
}

func DeleteHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	existing, err := GetByID(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if existing == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "kategori berita tidak ditemukan",
		})
	}

	used, err := IsUsedByBerita(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if used {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "kategori masih dipakai oleh berita, tidak bisa dihapus",
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
		"message": "kategori berita berhasil dihapus",
	})
}
