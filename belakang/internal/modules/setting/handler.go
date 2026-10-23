package setting

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sekolah/internal/auth"
	"sekolah/internal/config"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// =====================================
// Setting Halaman Publik
// =====================================

func Get(c *fiber.Ctx) error {

	data, err := GetSetting()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if data.LogoSekolah != "" {
		data.LogoURL =
			config.App.AppURL +
				"/storage/logo/" +
				data.LogoSekolah
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// =====================================
// Dashboard Setting
// =====================================

func GetDashboardHandler(c *fiber.Ctx) error {

	data, err := GetDashboardSetting()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if data.LogoSekolah != "" {
		data.LogoURL =
			config.App.AppURL +
				"/storage/logo/" +
				data.LogoSekolah
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// =====================================
// Parse Form
// =====================================

func parseSettingForm(c *fiber.Ctx) Setting {

	return Setting{
		IDSekolah:        1,
		NamaSekolah:      c.FormValue("nama_sekolah"),
		AlamatSekolah:    c.FormValue("alamat_sekolah"),
		NPSN:             c.FormValue("npsn"),
		NSSS:             c.FormValue("nsss"),
		NSP:              c.FormValue("nsp"),
		KelurahanSekolah: c.FormValue("kelurahan_sekolah"),
		KecamatanSekolah: c.FormValue("kecamatan_sekolah"),
		KabupatenSekolah: c.FormValue("kabupaten_sekolah"),
		ProvinsiSekolah:  c.FormValue("provinsi_sekolah"),
		KodePosSekolah:   c.FormValue("kodepos_sekolah"),
		TelpSekolah:      c.FormValue("telp_sekolah"),
		EmailSekolah:     c.FormValue("email_sekolah"),
		WebSekolah:       c.FormValue("web_sekolah"),
		DeskripsiSekolah: c.FormValue("deskripsi_sekolah"),
		VisiSekolah:      c.FormValue("visi_sekolah"),
		FooterSekolah:    c.FormValue("footer_sekolah"),
		ProfilSekolah:    c.FormValue("profil_sekolah"),
		Facebook:         c.FormValue("facebook"),
		Twitter:          c.FormValue("twitter"),
		Instagram:        c.FormValue("instagram"),
		Youtube:          c.FormValue("youtube"),
		PetaSekolah:      c.FormValue("peta_sekolah"),
		HeroArea1:        c.FormValue("heroarea1"),
		HeroArea2:        c.FormValue("heroarea2"),
		HeroArea3:        c.FormValue("heroarea3"),
		HeroArea4:        c.FormValue("heroarea4"),
	}
}

// =====================================
// Update Dashboard Setting
// =====================================

func UpdateDashboardHandler(c *fiber.Ctx) error {
	const defaultLogo = "1636812404_b39deec2f158cb789567.png"

	setting := parseSettingForm(c)

	if strings.TrimSpace(setting.NamaSekolah) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nama sekolah wajib diisi.",
		})
	}

	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid.",
		})
	}

	setting.IDPengguna = claims.IDPengguna

	// Ambil data lama sekali saja.
	oldSetting, err := GetDashboardSetting()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membaca data pengaturan.",
		})
	}

	// Pertahankan nama logo lama bila admin tidak mengunggah logo baru.
	setting.LogoSekolah = oldSetting.LogoSekolah

	var newLogoPath string
	var oldLogoPath string

	// Hapus logo: kembalikan ke logo default (file default TIDAK pernah dihapus).
	// File lama dihapus setelah DB berhasil diupdate. Bila pada request yang sama
	// ada logo baru, blok upload di bawah menimpanya.
	if c.FormValue("hapus_logo") == "1" && oldSetting.LogoSekolah != defaultLogo {
		setting.LogoSekolah = defaultLogo

		if oldSetting.LogoSekolah != "" {
			oldLogoPath = filepath.Join(
				"storage",
				"logo",
				filepath.Base(oldSetting.LogoSekolah),
			)
		}
	}

	logoFile, err := c.FormFile("logo")

	// Tidak ada file logo adalah kondisi normal.
	if err == nil && logoFile != nil {
		const maxSize = 100 * 1024

		if logoFile.Size > maxSize {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Ukuran logo maksimal 100 KB.",
			})
		}

		ext := strings.ToLower(filepath.Ext(logoFile.Filename))

		switch ext {
		case ".jpg", ".jpeg", ".png", ".webp":
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Logo harus berformat JPG, JPEG, PNG, atau WEBP.",
			})
		}

		if err := os.MkdirAll(filepath.Join("storage", "logo"), 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal menyiapkan folder penyimpanan logo.",
			})
		}

		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

		newLogoPath = filepath.Join("storage", "logo", filename)

		src, err := logoFile.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal membuka file logo.",
			})
		}
		defer src.Close()

		out, err := os.Create(newLogoPath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal menyimpan logo.",
			})
		}

		if _, err := io.Copy(out, src); err != nil {
			out.Close()
			_ = os.Remove(newLogoPath)

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal menyimpan logo.",
			})
		}

		if err := out.Close(); err != nil {
			_ = os.Remove(newLogoPath)

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal menyimpan logo.",
			})
		}

		setting.LogoSekolah = filename

		// Logo default tidak boleh ikut terhapus saat diganti logo baru.
		if oldSetting.LogoSekolah != "" && oldSetting.LogoSekolah != defaultLogo {
			oldLogoPath = filepath.Join(
				"storage",
				"logo",
				filepath.Base(oldSetting.LogoSekolah),
			)
		}
	}

	changed, err := UpdateSetting(&setting)
	if err != nil {
		// Database gagal: file baru tidak boleh menjadi sampah.
		if newLogoPath != "" {
			_ = os.Remove(newLogoPath)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if !changed {
		// Tidak ada perubahan: hapus file baru yang barusan diunggah.
		if newLogoPath != "" {
			_ = os.Remove(newLogoPath)
		}

		return c.JSON(fiber.Map{
			"success": true,
			"changed": false,
			"message": "Tidak ada perubahan data.",
		})
	}

	// Database sudah menunjuk logo baru; sekarang aman menghapus logo lama.
	if oldLogoPath != "" && oldLogoPath != newLogoPath {
		removeFileEventually(oldLogoPath)
	}

	data, err := GetDashboardSetting()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if data.LogoSekolah != "" {
		data.LogoURL = config.App.AppURL + "/storage/logo/" + data.LogoSekolah
	}

	return c.JSON(fiber.Map{
		"success": true,
		"changed": true,
		"message": "Website berhasil diperbarui.",
		"data":    data,
	})
}

func removeFileWithRetry(path string) error {
	const maxAttempt = 10
	const retryDelay = 200 * time.Millisecond

	var err error

	for attempt := 1; attempt <= maxAttempt; attempt++ {
		err = os.Remove(path)

		if err == nil || os.IsNotExist(err) {
			return nil
		}

		time.Sleep(retryDelay)
	}

	return err
}

func removeFileEventually(path string) {
	go func() {
		// Beri browser/static server waktu menyelesaikan request logo lama.
		time.Sleep(3 * time.Second)

		const maxAttempt = 30
		const retryDelay = 1 * time.Second

		for attempt := 1; attempt <= maxAttempt; attempt++ {
			err := os.Remove(path)

			if err == nil || os.IsNotExist(err) {
				return
			}

			time.Sleep(retryDelay)
		}

		fmt.Println("Gagal menghapus logo lama setelah beberapa percobaan:", path)
	}()
}
