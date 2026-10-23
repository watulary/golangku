package walikelas

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
	sortBy := c.Query("sort_by", "kelas")
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
				"message": "Wali kelas tidak ditemukan.",
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

func SearchPegawaiHandler(c *fiber.Ctx) error {
	search := c.Query("q", "")

	data, err := SearchPegawai(search, 15)
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

func validateRelasi(idKelas, idRombel int, nuptk string, excludeID int) (int, string, bool) {
	kelasExists, err := KelasExists(idKelas)
	if err != nil {
		return fiber.StatusInternalServerError, err.Error(), false
	}
	if !kelasExists {
		return fiber.StatusBadRequest, "Kelas tidak ditemukan.", false
	}

	rombelValid, err := RombelBelongsToKelas(idRombel, idKelas)
	if err != nil {
		return fiber.StatusInternalServerError, err.Error(), false
	}
	if !rombelValid {
		return fiber.StatusBadRequest, "Rombel tidak sesuai dengan kelas yang dipilih.", false
	}

	pegawaiExists, err := PegawaiExists(nuptk)
	if err != nil {
		return fiber.StatusInternalServerError, err.Error(), false
	}
	if !pegawaiExists {
		return fiber.StatusBadRequest, "Pegawai tidak ditemukan.", false
	}

	rombelTaken, err := ExistsByRombel(idRombel, excludeID)
	if err != nil {
		return fiber.StatusInternalServerError, err.Error(), false
	}
	if rombelTaken {
		return fiber.StatusConflict, "Rombel ini sudah memiliki wali kelas.", false
	}

	pegawaiTaken, err := ExistsByPegawai(nuptk, excludeID)
	if err != nil {
		return fiber.StatusInternalServerError, err.Error(), false
	}
	if pegawaiTaken {
		return fiber.StatusConflict, "Pegawai ini sudah menjadi wali kelas di rombel lain.", false
	}

	return 0, "", true
}

func CreateHandler(c *fiber.Ctx) error {
	var input WaliKelas

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	input.NUPTK = strings.TrimSpace(input.NUPTK)

	if input.IDKelas <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Kelas wajib dipilih.",
		})
	}

	if input.IDRombel <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Rombel wajib dipilih.",
		})
	}

	if input.NUPTK == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Wali kelas (pegawai) wajib dipilih.",
		})
	}

	if status, msg, ok := validateRelasi(input.IDKelas, input.IDRombel, input.NUPTK, 0); !ok {
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"message": msg,
		})
	}

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
		"message": "Wali kelas berhasil ditambahkan.",
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

	var input WaliKelas

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	input.NUPTK = strings.TrimSpace(input.NUPTK)

	if input.IDKelas <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Kelas wajib dipilih.",
		})
	}

	if input.IDRombel <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Rombel wajib dipilih.",
		})
	}

	if input.NUPTK == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Wali kelas (pegawai) wajib dipilih.",
		})
	}

	oldData, err := GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Wali kelas tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if oldData.IDKelas == input.IDKelas &&
		oldData.IDRombel == input.IDRombel &&
		oldData.NUPTK == input.NUPTK {
		return c.JSON(fiber.Map{
			"success": true,
			"changed": false,
			"message": "Tidak ada perubahan data.",
			"data":    oldData,
		})
	}

	if status, msg, ok := validateRelasi(input.IDKelas, input.IDRombel, input.NUPTK, id); !ok {
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"message": msg,
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
		"message": "Wali kelas berhasil diperbarui.",
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
				"message": "Wali kelas tidak ditemukan.",
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
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Wali kelas berhasil dihapus.",
	})
}
