package suratmasuk

import (
	"fmt"
	"os"
	"path/filepath"
	"sekolah/internal/auth"
	"sekolah/internal/config"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const dokumenMaxSize = 2 * 1024 * 1024 // 2 MB

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
		data[i].DokumenURL = buildDokumenURL(data[i].Dokumen)
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
			"message": "Surat masuk tidak ditemukan.",
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
			"message": "Anda tidak memiliki akses ke surat masuk ini.",
		})
	}

	data.DokumenURL = buildDokumenURL(data.Dokumen)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func CreateHandler(c *fiber.Ctx) error {

	var input SuratMasuk

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	if strings.TrimSpace(input.Nama) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Judul surat wajib diisi.",
		})
	}

	if strings.TrimSpace(input.Nomor) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nomor surat wajib diisi.",
		})
	}

	if strings.TrimSpace(input.Tanggal) == "" || strings.TrimSpace(input.TanggalTerima) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Tanggal surat dan tanggal terima wajib diisi.",
		})
	}

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid.",
		})
	}

	file, err := c.FormFile("dokumen")
	if err != nil || file == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Dokumen surat wajib diunggah.",
		})
	}

	if ext := strings.ToLower(filepath.Ext(file.Filename)); ext != ".pdf" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Dokumen harus berformat PDF.",
		})
	}
	if file.Size > dokumenMaxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Ukuran dokumen maksimal 2MB.",
		})
	}

	id, err := uuid.NewV7()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	filename := id.String() + filepath.Ext(file.Filename)

	if err := c.SaveFile(file, "./storage/suratmasuk/"+filename); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengunggah dokumen.",
		})
	}

	input.ID = id.String()
	input.Dokumen = filename
	input.IDPengguna = claims.IDPengguna

	if err := Create(input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	input.DokumenURL = buildDokumenURL(input.Dokumen)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Surat masuk berhasil ditambahkan.",
		"data":    input,
	})
}

func UpdateHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	var input SuratMasuk

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	if strings.TrimSpace(input.Nama) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nama/judul surat wajib diisi.",
		})
	}

	oldData, err := GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	if oldData == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Surat masuk tidak ditemukan.",
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
			"message": "Anda tidak memiliki akses untuk mengubah surat masuk ini.",
		})
	}

	changed := oldData.Nama != input.Nama ||
		oldData.Nomor != input.Nomor ||
		oldData.Tanggal != input.Tanggal ||
		oldData.TanggalTerima != input.TanggalTerima ||
		oldData.Perihal != input.Perihal ||
		oldData.Asal != input.Asal ||
		oldData.Tujuan != input.Tujuan

	file, ferr := c.FormFile("dokumen")
	hasNewDokumen := ferr == nil && file != nil

	if hasNewDokumen {
		if ext := strings.ToLower(filepath.Ext(file.Filename)); ext != ".pdf" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Dokumen harus berformat PDF.",
			})
		}
		if file.Size > dokumenMaxSize {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Ukuran dokumen maksimal 2MB.",
			})
		}
	}

	if !changed && !hasNewDokumen {
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
	}

	if hasNewDokumen {

		newID, err := uuid.NewV7()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		filename := newID.String() + filepath.Ext(file.Filename)

		if err := c.SaveFile(file, "./storage/suratmasuk/"+filename); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal mengunggah dokumen.",
			})
		}

		if err := UpdateDokumen(id, filename); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		go deleteDokumenFile(oldData.Dokumen)
	}

	data, err := GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	data.DokumenURL = buildDokumenURL(data.Dokumen)

	return c.JSON(fiber.Map{
		"success": true,
		"changed": true,
		"message": "Surat masuk berhasil diperbarui.",
		"data":    data,
	})
}

func DeleteHandler(c *fiber.Ctx) error {

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
			"message": "Surat masuk tidak ditemukan.",
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
			"message": "Anda tidak memiliki akses untuk menghapus surat masuk ini.",
		})
	}

	if err := Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	go deleteDokumenFile(data.Dokumen)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Surat masuk berhasil dihapus.",
	})
}

func buildDokumenURL(dokumen string) string {
	if dokumen == "" {
		return ""
	}
	return config.App.AppURL + "/storage/suratmasuk/" + dokumen
}

func deleteDokumenFile(dokumen string) {

	if dokumen == "" {
		return
	}

	path := "./storage/suratmasuk/" + dokumen

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
			fmt.Println("Gagal menghapus file dokumen lama:", err)
		}

		time.Sleep(retryDelay)
	}
}

func isOwnerOrAdmin(c *fiber.Ctx, idPenggunaSuratMasuk string) (*auth.Claims, bool) {

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return nil, false
	}

	if claims.Level == "Admin" {
		return claims, true
	}

	return claims, claims.IDPengguna == idPenggunaSuratMasuk
}
