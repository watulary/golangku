package pengumuman

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"sekolah/internal/auth"
	"sekolah/internal/config"
	"sekolah/internal/modules/tiptap"
	"sekolah/internal/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetLatestHandler(c *fiber.Ctx) error {

	data, err := GetLatest(6)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	for i := range data {

		data[i].FotoURL = buildFotoURL(data[i].Foto)
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

		data[i].FotoURL = buildFotoURL(data[i].Foto)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func GetDetailHandler(c *fiber.Ctx) error {

	slug := c.Params("slug")

	data, err := GetBySlug(slug)

	if err != nil {

		if err == sql.ErrNoRows {

			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Pengumuman tidak ditemukan",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if data.Foto != "" {

		data.FotoURL =
			config.App.AppURL +
				"/storage/fotopengumuman/" +
				data.Foto
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func IncrementHitsHandler(c *fiber.Ctx) error {

	idPengumuman := c.Params("id")

	err := IncrementHits(idPengumuman)

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
		"message": "Hits berhasil ditambah",
	})
}

const (
	DefaultFotoPengumuman = "1636812404_b39deec2f158cb789567.jpg"
)

func GetAdminListHandler(c *fiber.Ctx) error {

	search := c.Query("search", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sort_by", "tanggal")
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
		data[i].FotoURL = buildFotoURL(data[i].Foto)
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
				"message": "Pengumuman tidak ditemukan.",
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
			"message": "Anda tidak memiliki akses ke pengumuman ini.",
		})
	}

	if data.Foto != "" {
		data.FotoURL = buildFotoURL(data.Foto)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func CreateHandler(c *fiber.Ctx) error {

	var input Pengumuman

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	if strings.TrimSpace(input.JudulPengumuman) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Judul pengumuman wajib diisi.",
		})
	}

	if strings.TrimSpace(input.IsiPengumuman) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Isi pengumuman wajib diisi.",
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

	input.IDPengumuman = id.String()
	input.SlugPengumuman = utils.Slugify(input.JudulPengumuman)
	input.TanggalPublish = time.Now()
	input.IDPengguna = claims.IDPengguna

	if input.StatusPengumuman == "" {
		input.StatusPengumuman = "Draft"
	}

	file, err := c.FormFile("foto")
	if err == nil && file != nil {

		if msg := validateFotoPengumuman(file); msg != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": msg,
			})
		}

		filename := input.IDPengumuman + filepath.Ext(file.Filename)

		if err := c.SaveFile(file, "./storage/fotopengumuman/"+filename); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal mengunggah foto.",
			})
		}

		input.Foto = filename
	} else {
		input.Foto = DefaultFotoPengumuman
	}

	if err := Create(input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Pengumuman berhasil dibuat.",
		"data":    input,
	})
}

func UpdateHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	var input Pengumuman

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	if strings.TrimSpace(input.JudulPengumuman) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Judul pengumuman wajib diisi.",
		})
	}

	if strings.TrimSpace(input.IsiPengumuman) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Isi pengumuman wajib diisi.",
		})
	}

	oldData, err := GetByID(id)
	if err != nil {

		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Pengumuman tidak ditemukan.",
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
			"message": "Anda tidak memiliki akses untuk mengubah pengumuman ini.",
		})
	}

	input.SlugPengumuman = utils.Slugify(input.JudulPengumuman)

	changed := oldData.JudulPengumuman != input.JudulPengumuman ||
		oldData.IsiPengumuman != input.IsiPengumuman ||
		oldData.Ringkasan != input.Ringkasan ||
		oldData.StatusPengumuman != input.StatusPengumuman

	file, ferr := c.FormFile("foto")
	hasNewFoto := ferr == nil && file != nil
	hapusFoto := c.FormValue("hapus_foto") == "1"

	if hasNewFoto {
		if msg := validateFotoPengumuman(file); msg != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": msg,
			})
		}
	}

	if !changed && !hasNewFoto && !hapusFoto {
		return c.JSON(fiber.Map{
			"success": true,
			"changed": false,
			"message": "Tidak ada perubahan data.",
			"data":    oldData,
		})
	}

	if changed {
		if err := Update(id, input); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		// Isi pengumuman bisa berubah tanpa isi_pengumuman ikut berubah
		// (mis. cuma ganti judul/status), jadi cek dulu supaya tidak salah hapus.
		if oldData.IsiPengumuman != input.IsiPengumuman {
			tiptap.CleanupOrphanImages(oldData.IsiPengumuman, input.IsiPengumuman)
		}
	}

	if hasNewFoto {

		newID, err := uuid.NewV7()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		filename := newID.String() + filepath.Ext(file.Filename)

		if err := c.SaveFile(file, "./storage/fotopengumuman/"+filename); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal mengunggah foto.",
			})
		}

		if err := UpdateFoto(id, filename); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		go deleteFotoFile(oldData.Foto)

	} else if hapusFoto {

		if err := UpdateFoto(id, DefaultFotoPengumuman); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		go deleteFotoFile(oldData.Foto)
	}

	data, err := GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	data.FotoURL = buildFotoURL(data.Foto)

	return c.JSON(fiber.Map{
		"success": true,
		"changed": true,
		"message": "Pengumuman berhasil diperbarui.",
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
				"message": "Pengumuman tidak ditemukan.",
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
			"message": "Anda tidak memiliki akses untuk menghapus pengumuman ini.",
		})
	}

	if err := Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	go deleteFotoFile(data.Foto)

	for _, f := range tiptap.ExtractImages(data.IsiPengumuman) {
		go tiptap.DeleteImageFile(f)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pengumuman berhasil dihapus.",
	})
}

func buildFotoURL(foto string) string {

	if foto == "" {
		foto = DefaultFotoPengumuman
	}

	return config.App.AppURL + "/storage/fotopengumuman/" + foto
}

var allowedFotoPengumumanExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

const maxFotoPengumumanSize = 100 * 1024 // 100 KB

// validateFotoPengumuman mengecek ekstensi & ukuran file foto pengumuman.
// Mengembalikan pesan error kalau tidak valid, string kosong kalau valid.
func validateFotoPengumuman(file *multipart.FileHeader) string {

	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !allowedFotoPengumumanExt[ext] {
		return "Foto harus berformat JPG, JPEG, PNG atau WEBP."
	}

	if file.Size > maxFotoPengumumanSize {
		return "Ukuran foto maksimal 100 KB."
	}

	return ""
}

func deleteFotoFile(foto string) {

	if foto == "" || foto == DefaultFotoPengumuman {
		return
	}

	path := "./storage/fotopengumuman/" + foto

	// Beri waktu proses lain (antivirus, thumbnail cache, cloud sync, dsb)
	// melepas lock file sebelum mulai mencoba menghapus.
	time.Sleep(3 * time.Second)

	const maxAttempt = 30
	const retryDelay = 1 * time.Second

	for attempt := 1; attempt <= maxAttempt; attempt++ {

		err := os.Remove(path)

		if err == nil || os.IsNotExist(err) {
			return
		}

		if attempt == maxAttempt {
			fmt.Println("Gagal menghapus file foto lama:", err)
		}

		time.Sleep(retryDelay)
	}
}

func isOwnerOrAdmin(c *fiber.Ctx, idPenggunaPengumuman string) (*auth.Claims, bool) {

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return nil, false
	}

	if claims.Level == "Admin" {
		return claims, true
	}

	return claims, claims.IDPengguna == idPenggunaPengumuman
}
