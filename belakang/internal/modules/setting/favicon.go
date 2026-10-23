package setting

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"

	_ "golang.org/x/image/webp"

	xdraw "golang.org/x/image/draw"

	"github.com/gofiber/fiber/v2"
)

// =====================================
// Favicon (persegi, proporsi logo terjaga)
// =====================================

// Favicon membaca logo_sekolah, menaruhnya di kanvas PNG persegi 64x64
// transparan tanpa mengubah proporsi, lalu mengirimnya sebagai gambar.
// Route publik: GET /api/public/favicon
func Favicon(c *fiber.Ctx) error {
	data, err := GetSetting()
	if err != nil || data.LogoSekolah == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	// filepath.Base mencegah path traversal dari nilai di DB.
	path := filepath.Join("storage", "logo", filepath.Base(data.LogoSekolah))

	f, err := os.Open(path)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	defer f.Close()

	src, _, err := image.Decode(f)
	if err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	const size = 64

	b := src.Bounds()
	if b.Dx() == 0 || b.Dy() == 0 {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	scale := math.Min(float64(size)/float64(b.Dx()), float64(size)/float64(b.Dy()))
	w := int(math.Round(float64(b.Dx()) * scale))
	h := int(math.Round(float64(b.Dy()) * scale))

	x0 := (size - w) / 2
	y0 := (size - h) / 2

	dst := image.NewNRGBA(image.Rect(0, 0, size, size)) // transparan
	xdraw.CatmullRom.Scale(
		dst,
		image.Rect(x0, y0, x0+w, y0+h),
		src,
		b,
		draw.Over,
		nil,
	)

	var buf bytes.Buffer
	if err := png.Encode(&buf, dst); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Set("Content-Type", "image/png")
	// Nama file logo berubah tiap upload dan dipakai sebagai ?v= di frontend,
	// jadi aman di-cache.
	c.Set("Cache-Control", "public, max-age=86400")

	return c.Send(buf.Bytes())
}
