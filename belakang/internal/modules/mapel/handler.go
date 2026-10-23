package mapel

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func errorJSON(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}

// isForeignKeyError: MySQL error 1451 (delete parent row gagal karena masih
// dirujuk baris lain) / 1216 (varian lama). Dipakai supaya pesan hapus gagal
// karena mapel sudah punya nilai tidak menampilkan error SQL mentah ke user.
func isForeignKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError

	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1451 || mysqlErr.Number == 1216
	}

	return false
}

func sameIntPtr(a, b *int) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}

	return *a == *b
}

// validateInput merapikan & memvalidasi input. selfID = 0 untuk tambah baru.
// Mengembalikan nil kalau valid.
func validateInput(input *Mapel, selfID int) *fiber.Error {
	input.Nama = strings.TrimSpace(input.Nama)

	if input.Nama == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Nama mapel wajib diisi.")
	}

	if utf8.RuneCountInString(input.Nama) > 100 {
		return fiber.NewError(fiber.StatusBadRequest, "Nama mapel maksimal 100 karakter.")
	}

	if input.Urutan < 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Urutan tidak boleh negatif.")
	}

	if input.IDInduk != nil {
		if selfID != 0 && *input.IDInduk == selfID {
			return fiber.NewError(fiber.StatusBadRequest, "Mapel tidak boleh menjadi induk bagi dirinya sendiri.")
		}

		parent, err := GetByID(*input.IDInduk)
		if err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusBadRequest, "Mapel induk tidak ditemukan.")
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		// maksimal 2 tingkat: induk harus mapel utama, bukan sub-mapel
		if parent.IDInduk != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Induk harus berupa mapel utama, bukan sub-mapel.")
		}

		if selfID != 0 {
			hasChildren, err := HasChildren(selfID)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			if hasChildren {
				return fiber.NewError(fiber.StatusBadRequest, "Mapel ini masih punya sub-mapel, sehingga tidak bisa dijadikan sub-mapel.")
			}
		}
	}

	exists, err := ExistsByNama(input.Nama, input.IDInduk, selfID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if exists {
		return fiber.NewError(fiber.StatusConflict, "Mapel dengan nama tersebut sudah ada.")
	}

	return nil
}

func GetAdminListHandler(c *fiber.Ctx) error {
	search := c.Query("search", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sort_by", "urutan")
	sortDir := c.Query("sort_dir", "asc")

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 1000 {
		limit = 10
	}

	offset := (page - 1) * limit

	data, err := GetAllAdmin(search, sortBy, sortDir, limit, offset)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	total, err := CountAdmin(search)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func GetIndukHandler(c *fiber.Ctx) error {
	data, err := GetInduk()
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// GetDinilaiHandler = daftar mapel yang bisa diisi nilainya (mapel tanpa sub-mapel).
func GetDinilaiHandler(c *fiber.Ctx) error {
	data, err := GetDinilai()
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func GetByIDHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "ID tidak valid.")
	}

	data, err := GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errorJSON(c, fiber.StatusNotFound, "Mapel tidak ditemukan.")
		}

		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func CreateHandler(c *fiber.Ctx) error {
	var input Mapel

	if err := c.BodyParser(&input); err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "Payload tidak valid.")
	}

	if ferr := validateInput(&input, 0); ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	input.Aktif = 1

	newID, err := Create(input)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	data, err := GetByID(int(newID))
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Mapel berhasil ditambahkan.",
		"data":    data,
	})
}

func UpdateHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "ID tidak valid.")
	}

	var input Mapel

	if err := c.BodyParser(&input); err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "Payload tidak valid.")
	}

	oldData, err := GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errorJSON(c, fiber.StatusNotFound, "Mapel tidak ditemukan.")
		}

		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	if ferr := validateInput(&input, id); ferr != nil {
		return errorJSON(c, ferr.Code, ferr.Message)
	}

	if oldData.Nama == input.Nama &&
		oldData.Urutan == input.Urutan &&
		sameIntPtr(oldData.IDInduk, input.IDInduk) {
		return c.JSON(fiber.Map{
			"success": true,
			"changed": false,
			"message": "Tidak ada perubahan data.",
			"data":    oldData,
		})
	}

	if err := Update(id, input); err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	data, err := GetByID(id)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"changed": true,
		"message": "Mapel berhasil diperbarui.",
		"data":    data,
	})
}

func DeleteHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "ID tidak valid.")
	}

	if _, err := GetByID(id); err != nil {
		if err == sql.ErrNoRows {
			return errorJSON(c, fiber.StatusNotFound, "Mapel tidak ditemukan.")
		}

		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	hasChildren, err := HasChildren(id)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	if hasChildren {
		return errorJSON(c, fiber.StatusBadRequest, "Mapel ini masih punya sub-mapel. Hapus sub-mapelnya terlebih dahulu.")
	}

	if err := Delete(id); err != nil {
		if isForeignKeyError(err) {
			return errorJSON(c, fiber.StatusBadRequest,
				"Mapel ini sudah punya nilai yang tercatat, sehingga tidak bisa dihapus. Nonaktifkan saja mapel ini jika sudah tidak dipakai.")
		}

		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Mapel berhasil dihapus.",
	})
}

func ToggleStatusHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "ID tidak valid.")
	}

	if _, err := GetByID(id); err != nil {
		if err == sql.ErrNoRows {
			return errorJSON(c, fiber.StatusNotFound, "Mapel tidak ditemukan.")
		}

		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := ToggleStatus(id); err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	data, err := GetByID(id)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	message := "Mapel berhasil dinonaktifkan."
	if data.Aktif == 1 {
		message = "Mapel berhasil diaktifkan."
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}
