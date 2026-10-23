package tahunpelajaran

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var tahunAjaranRegex = regexp.MustCompile(`^\d{4}/\d{4}$`)

func isSemesterValid(s string) bool {
	return s == "Ganjil" || s == "Genap"
}

func GetAdminListHandler(c *fiber.Ctx) error {

	search := c.Query("search", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sort_by", "nama")
	sortDir := c.Query("sort_dir", "desc")

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	data, err := GetAllAdmin(search, sortBy, sortDir, limit, offset)
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
				"message": "Tahun pelajaran tidak ditemukan.",
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
	var input TahunPelajaran

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
			"message": "Nama tahun pelajaran wajib diisi.",
		})
	}

	if !tahunAjaranRegex.MatchString(input.Nama) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format tahun pelajaran tidak valid. Gunakan format YYYY/YYYY, contoh: 2025/2026.",
		})
	}

	input.Semester = strings.TrimSpace(input.Semester)

	if !isSemesterValid(input.Semester) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Semester wajib dipilih: Ganjil atau Genap.",
		})
	}

	exists, err := ExistsByNama(input.Nama, 0)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"message": "Tahun pelajaran sudah ada.",
		})
	}

	input.Aktif = 0

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
		"message": "Tahun pelajaran berhasil ditambahkan.",
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

	var input TahunPelajaran

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
			"message": "Nama tahun pelajaran wajib diisi.",
		})
	}

	if !tahunAjaranRegex.MatchString(input.Nama) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format tahun pelajaran tidak valid. Gunakan format YYYY/YYYY, contoh: 2025/2026.",
		})
	}

	input.Semester = strings.TrimSpace(input.Semester)

	if !isSemesterValid(input.Semester) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Semester wajib dipilih: Ganjil atau Genap.",
		})
	}

	oldData, err := GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Tahun pelajaran tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if oldData.Nama == input.Nama && oldData.Semester == input.Semester {
		return c.JSON(fiber.Map{
			"success": true,
			"changed": false,
			"message": "Tidak ada perubahan data.",
			"data":    oldData,
		})
	}

	// cek duplikat nama hanya kalau namanya memang berubah
	if oldData.Nama != input.Nama {
		exists, err := ExistsByNama(input.Nama, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		if exists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false,
				"message": "Tahun pelajaran sudah ada.",
			})
		}
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
		"message": "Tahun pelajaran berhasil diperbarui.",
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

	data, err := GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Tahun pelajaran tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if data.Aktif == 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Tahun pelajaran yang sedang aktif tidak dapat dihapus.",
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
		"message": "Tahun pelajaran berhasil dihapus.",
	})
}

func SetActiveHandler(c *fiber.Ctx) error {
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
				"message": "Tahun pelajaran tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := SetActive(id); err != nil {
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
		"message": "Tahun pelajaran berhasil diaktifkan.",
		"data":    data,
	})
}
