package rombel

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetAdminListHandler(c *fiber.Ctx) error {

	search := c.Query("search", "")
	idKelas, _ := strconv.Atoi(c.Query("id_kelas", "0"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sort_by", "nama")
	sortDir := c.Query("sort_dir", "asc")

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	data, err := GetAllAdmin(search, idKelas, sortBy, sortDir, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	total, err := CountAdmin(search, idKelas)
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
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid.",
		})
	}

	data, err := GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Rombel tidak ditemukan.",
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

func CreateHandler(c *fiber.Ctx) error {
	var input Rombel

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	input.Nama = strings.TrimSpace(input.Nama)

	if input.Nama == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nama rombel wajib diisi.",
		})
	}

	if len(input.Nama) > 50 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nama rombel maksimal 50 karakter.",
		})
	}

	if input.IDKelas <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Kelas wajib dipilih.",
		})
	}

	kelasExists, err := KelasExists(input.IDKelas)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if !kelasExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Kelas tidak ditemukan.",
		})
	}

	exists, err := ExistsByNama(input.Nama, input.IDKelas, 0)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"message": "Rombel dengan nama tersebut sudah ada pada kelas ini.",
		})
	}

	input.Aktif = 1

	newID, err := Create(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	data, err := GetByID(int(newID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Rombel berhasil ditambahkan.",
		"data":    data,
	})
}

func UpdateHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid.",
		})
	}

	var input Rombel

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	input.Nama = strings.TrimSpace(input.Nama)

	if input.Nama == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nama rombel wajib diisi.",
		})
	}

	if len(input.Nama) > 50 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nama rombel maksimal 50 karakter.",
		})
	}

	if input.IDKelas <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Kelas wajib dipilih.",
		})
	}

	oldData, err := GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Rombel tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if oldData.Nama == input.Nama && oldData.IDKelas == input.IDKelas {
		return c.JSON(fiber.Map{
			"success": true,
			"changed": false,
			"message": "Tidak ada perubahan data.",
			"data":    oldData,
		})
	}

	kelasExists, err := KelasExists(input.IDKelas)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if !kelasExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Kelas tidak ditemukan.",
		})
	}

	exists, err := ExistsByNama(input.Nama, input.IDKelas, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"message": "Rombel dengan nama tersebut sudah ada pada kelas ini.",
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

	return c.JSON(fiber.Map{
		"success": true,
		"changed": true,
		"message": "Rombel berhasil diperbarui.",
		"data":    data,
	})
}

func DeleteHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid.",
		})
	}

	if _, err := GetByID(id); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Rombel tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Rombel tidak dapat dihapus, kemungkinan masih digunakan oleh data lain.",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Rombel berhasil dihapus.",
	})
}

func ToggleStatusHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid.",
		})
	}

	oldData, err := GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Rombel tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	newStatus := 1
	if oldData.Aktif == 1 {
		newStatus = 0
	}

	if err := ToggleStatus(id, newStatus); err != nil {
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

	message := "Rombel berhasil dinonaktifkan."
	if newStatus == 1 {
		message = "Rombel berhasil diaktifkan."
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}
