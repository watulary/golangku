package berita

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
				"message": "Berita tidak ditemukan",
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
				"/storage/fotoberita/" +
				data.Foto
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func IncrementHitsHandler(c *fiber.Ctx) error {

	idBerita := c.Params("id")

	err := IncrementHits(idBerita)

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
	KategoriInformasiDokumen    = "f7349f0cc8bc2651f48f909d07fdefca9157cd37"
	KategoriInformasiBerkala    = "00d1b344e92719e582139de020395ea3ea8be04d"
	KategoriInformasiSetiapSaat = "ecae2b1b6f40ca6d5e5d1c6df94a8efed2dc46a1"
	KategoriInformasiSertaMerta = "e434198b1815fcd7c0e96f8b13d657f788b9fc1f"
	KategoriInformasiKontak     = "5ffb08da9b0d35c9a8753e16872cf744a9462b1a"
	DefaultFotoBerita           = "1636812404_b39deec2f158cb789567.jpg"
)

func getBeritaKategori(
	c *fiber.Ctx,
	idKategori string,
) error {

	data, err := GetByKategori(
		idKategori,
		100,
	)

	if err != nil {
		return c.Status(500).JSON(
			fiber.Map{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	for i := range data {

		data[i].FotoURL = buildFotoURL(data[i].Foto)
	}

	return c.JSON(
		fiber.Map{
			"success": true,
			"data":    data,
		},
	)
}

func GetInformasiDokumenHandler(c *fiber.Ctx) error {
	return getBeritaKategori(
		c,
		KategoriInformasiDokumen,
	)
}

func GetInformasiBerkalaHandler(c *fiber.Ctx) error {
	return getBeritaKategori(
		c,
		KategoriInformasiBerkala,
	)
}

func GetInformasiSetiapSaatHandler(c *fiber.Ctx) error {
	return getBeritaKategori(
		c,
		KategoriInformasiSetiapSaat,
	)
}

func GetInformasiSertaMertaHandler(c *fiber.Ctx) error {
	return getBeritaKategori(
		c,
		KategoriInformasiSertaMerta,
	)
}

func GetInformasiKontakHandler(c *fiber.Ctx) error {
	return getBeritaKategori(
		c,
		KategoriInformasiKontak,
	)
}

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
				"message": "Berita tidak ditemukan.",
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
			"message": "Anda tidak memiliki akses ke berita ini.",
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

	var input Berita

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	if strings.TrimSpace(input.JudulBerita) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Judul berita wajib diisi.",
		})
	}

	if strings.TrimSpace(input.IsiBerita) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Isi berita wajib diisi.",
		})
	}

	if strings.TrimSpace(input.IDKategoriBerita) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Kategori berita wajib dipilih.",
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

	input.IDBerita = id.String()
	input.SlugBerita = utils.Slugify(input.JudulBerita)
	input.TanggalPublish = time.Now()
	input.IDPengguna = claims.IDPengguna

	if input.StatusBerita == "" {
		input.StatusBerita = "Draft"
	}

	file, err := c.FormFile("foto")
	if err == nil && file != nil {

		if msg := validateFotoBerita(file); msg != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": msg,
			})
		}

		filename := input.IDBerita + filepath.Ext(file.Filename)

		if err := c.SaveFile(file, "./storage/fotoberita/"+filename); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal mengunggah foto.",
			})
		}

		input.Foto = filename
	} else {
		input.Foto = DefaultFotoBerita
	}

	if err := Create(input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Berita berhasil dibuat.",
		"data":    input,
	})
}

func UpdateHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	var input Berita

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	if strings.TrimSpace(input.JudulBerita) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Judul berita wajib diisi.",
		})
	}

	if strings.TrimSpace(input.IsiBerita) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Isi berita wajib diisi.",
		})
	}

	oldData, err := GetByID(id)
	if err != nil {

		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Berita tidak ditemukan.",
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
			"message": "Anda tidak memiliki akses untuk mengubah berita ini.",
		})
	}

	input.SlugBerita = utils.Slugify(input.JudulBerita)

	changed := oldData.JudulBerita != input.JudulBerita ||
		oldData.IsiBerita != input.IsiBerita ||
		oldData.Ringkasan != input.Ringkasan ||
		oldData.StatusBerita != input.StatusBerita ||
		oldData.IDKategoriBerita != input.IDKategoriBerita

	file, ferr := c.FormFile("foto")
	hasNewFoto := ferr == nil && file != nil
	hapusFoto := c.FormValue("hapus_foto") == "1"

	if hasNewFoto {
		if msg := validateFotoBerita(file); msg != "" {
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

		// Isi berita bisa berubah tanpa isi_berita ikut berubah (mis. cuma
		// ganti judul/kategori), jadi cek dulu supaya tidak salah hapus.
		if oldData.IsiBerita != input.IsiBerita {
			tiptap.CleanupOrphanImages(oldData.IsiBerita, input.IsiBerita)
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

		if err := c.SaveFile(file, "./storage/fotoberita/"+filename); err != nil {
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

		if err := UpdateFoto(id, DefaultFotoBerita); err != nil {
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
		"message": "Berita berhasil diperbarui.",
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
				"message": "Berita tidak ditemukan.",
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
			"message": "Anda tidak memiliki akses untuk menghapus berita ini.",
		})
	}

	if err := Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	go deleteFotoFile(data.Foto)

	for _, f := range tiptap.ExtractImages(data.IsiBerita) {
		go tiptap.DeleteImageFile(f)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Berita berhasil dihapus.",
	})
}

func buildFotoURL(foto string) string {

	if foto == "" {
		foto = DefaultFotoBerita
	}

	return config.App.AppURL + "/storage/fotoberita/" + foto
}

var allowedFotoBeritaExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

const maxFotoBeritaSize = 100 * 1024 // 100 KB

// validateFotoBerita mengecek ekstensi & ukuran file foto berita.
// Mengembalikan pesan error kalau tidak valid, string kosong kalau valid.
func validateFotoBerita(file *multipart.FileHeader) string {

	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !allowedFotoBeritaExt[ext] {
		return "Foto harus berformat JPG, JPEG, PNG atau WEBP."
	}

	if file.Size > maxFotoBeritaSize {
		return "Ukuran foto maksimal 100 KB."
	}

	return ""
}

func deleteFotoFile(foto string) {

	if foto == "" || foto == DefaultFotoBerita {
		return
	}

	path := "./storage/fotoberita/" + foto

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

func isOwnerOrAdmin(c *fiber.Ctx, idPenggunaBerita string) (*auth.Claims, bool) {

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return nil, false
	}

	if claims.Level == "Admin" {
		return claims, true
	}

	return claims, claims.IDPengguna == idPenggunaBerita
}
