package tiptap

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
)

const uploadDir = "./storage/tiptap"

func SaveImage(c *fiber.Ctx, field string, filename string) error {
	file, err := c.FormFile(field)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return err
	}

	dst := filepath.Join(uploadDir, filename)

	if err := c.SaveFile(file, dst); err != nil {
		return err
	}

	return nil
}

func ImageURL(appURL, filename string) string {
	return fmt.Sprintf("%s/storage/tiptap/%s", appURL, filename)
}

// DeleteImageFile menghapus satu file gambar tiptap dari disk (retry sabar,
// dipanggil via goroutine oleh modul lain spt berita saat gambar tsb sudah
// tidak dipakai lagi di isi konten). Dipakai lintas modul, makanya exported.
func DeleteImageFile(filename string) {

	if filename == "" {
		return
	}

	path := filepath.Join(uploadDir, filename)

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
			fmt.Println("Gagal menghapus file gambar tiptap lama:", err)
		}

		time.Sleep(retryDelay)
	}
}

// imageRegex menangkap nama file gambar yg disisipkan via tiptap
// (src="/storage/tiptap/<nama-file>.<ext>") di dalam konten HTML apapun
// (isi_berita, isi_sambutan, isi_pengumuman, isi_sejarah, isi_visimisi, dst).
var imageRegex = regexp.MustCompile(`/storage/tiptap/([a-zA-Z0-9\-]+\.(?:jpg|jpeg|png|webp))`)

// ExtractImages mengambil daftar nama file gambar tiptap (unik, tanpa
// duplikat) yg disisipkan di dalam suatu konten HTML. Dipakai bersama oleh
// semua modul yg punya field isi/konten berbasis TiptapEditor.
func ExtractImages(html string) []string {

	matches := imageRegex.FindAllStringSubmatch(html, -1)

	seen := map[string]bool{}
	var filenames []string

	for _, m := range matches {

		if len(m) < 2 || seen[m[1]] {
			continue
		}

		seen[m[1]] = true
		filenames = append(filenames, m[1])
	}

	return filenames
}

// CleanupOrphanImages membandingkan gambar tiptap di konten HTML lama vs
// baru, lalu menghapus (via goroutine) gambar yg ada di konten lama tapi
// sudah tidak dipakai lagi di konten baru. Aman dipanggil meski konten tidak
// berubah sama sekali (hasilnya tidak akan ada yg dihapus). Panggil ini
// SETELAH konten baru berhasil disimpan ke database.
func CleanupOrphanImages(oldHTML, newHTML string) {

	newImages := map[string]bool{}

	for _, f := range ExtractImages(newHTML) {
		newImages[f] = true
	}

	for _, f := range ExtractImages(oldHTML) {

		if !newImages[f] {
			go DeleteImageFile(f)
		}
	}
}
