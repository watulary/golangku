package kelulusan

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

// imgDir - folder gambar tetap (logo, ttd) dipakai generator PDF.
const imgDir = "./storage/fotokelulusan"

var bulanIndo = map[string]string{
	"01": "Januari", "02": "Februari", "03": "Maret", "04": "April",
	"05": "Mei", "06": "Juni", "07": "Juli", "08": "Agustus",
	"09": "September", "10": "Oktober", "11": "November", "12": "Desember",
}

func formatTanggalIndo(raw string) string {
	// raw format "YYYY-MM-DD"
	parts := strings.Split(raw, "-")
	if len(parts) != 3 {
		return raw
	}
	thn, bln, tglRaw := parts[0], parts[1], parts[2]
	tgl, err := strconv.Atoi(tglRaw)
	if err != nil {
		return raw
	}
	nama, ok := bulanIndo[bln]
	if !ok {
		return raw
	}
	return fmt.Sprintf("%d %s %s", tgl, nama, thn)
}

// cell - helper ringkas utk CellFormat, mirip signature FPDF PHP: Cell(w,h,txt,border,ln,align)
func cell(pdf *gofpdf.Fpdf, w, h float64, txt string, border string, ln int, align string) {
	pdf.CellFormat(w, h, txt, border, ln, align, false, 0, "")
}

// multiCell - helper ringkas mirip MultiCell(w,h,txt,border,align) FPDF PHP
func multiCell(pdf *gofpdf.Fpdf, w, h float64, txt string, border string, align string) {
	pdf.MultiCell(w, h, txt, border, align, false)
}

func imgIfExists(pdf *gofpdf.Fpdf, filename string, x, y, w, h float64) {
	path := filepath.Join(imgDir, filename)
	if _, err := os.Stat(path); err == nil {
		pdf.ImageOptions(path, x, y, w, h, false, gofpdf.ImageOptions{ImageType: ""}, 0, "")
	}
}

// GenerateSKL - meniru 1:1 layout printpdf() controller CI lama (FPDF, ukuran
// custom 21x33cm), diterjemahkan ke gofpdf. cetak = data siswa+nilai, skl =
// konten surat (kop*/isi*/ttd) dari kelulusan_waktu.
func GenerateSKL(cetak *CetakData, skl *KelulusanWaktu) ([]byte, error) {
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		OrientationStr: "P",
		UnitStr:        "cm",
		SizeStr:        "",
		Size:           gofpdf.SizeType{Wd: 21, Ht: 33},
	})
	pdf.AddPage()
	pdf.SetFont("Times", "", 12)

	imgIfExists(pdf, "logodki.png", 1.5, 1, 2.3, 2.5)

	pdf.SetX(3)
	multiCell(pdf, 15.5, 0.5, skl.Kop1SKL, "", "C")
	pdf.SetX(3)
	multiCell(pdf, 15.5, 0.5, skl.Kop2SKL, "", "C")
	pdf.SetFont("Times", "B", 12.5)
	pdf.SetX(3)
	multiCell(pdf, 15.5, 0.5, skl.Kop3SKL, "", "C")
	pdf.SetFont("Times", "", 12)
	pdf.SetX(3)
	multiCell(pdf, 15.5, 0.5, skl.Kop4SKL, "", "C")
	pdf.SetX(3)
	multiCell(pdf, 15.5, 0.5, skl.Kop5SKL, "", "C")
	pdf.SetX(3)
	multiCell(pdf, 15.5, 0.5, skl.Kop6SKL, "", "C")
	pdf.SetX(3)
	multiCell(pdf, 15.5, 0.5, skl.Kop7SKL, "", "R")
	pdf.Line(1.1, 4.5, 20, 4.5)
	pdf.Ln(0.5)

	pdf.SetFont("Times", "BU", 12)
	pdf.SetX(3)
	multiCell(pdf, 15.5, 0.5, skl.Isi1SKL, "", "C")
	pdf.SetFont("Times", "", 12)
	pdf.SetX(3)
	multiCell(pdf, 15.5, 0.5, skl.Isi2SKL, "", "C")
	pdf.Ln(0.5)
	pdf.SetX(2.5)
	multiCell(pdf, 15.5, 0.5, skl.Isi3SKL, "", "J")
	pdf.Ln(0.5)

	pdf.SetX(3)
	cell(pdf, 5, 0.5, "Nama Peserta Didik", "", 0, "J")
	pdf.SetX(8.6)
	cell(pdf, 1, 0.5, ":", "", 0, "J")
	pdf.SetX(8.9)
	cell(pdf, 15, 0.5, cetak.NamaSiswa, "", 0, "J")
	pdf.Ln(0.5)

	pdf.SetX(3)
	cell(pdf, 5, 0.5, "Tempat, Tanggal Lahir", "", 0, "J")
	pdf.SetX(8.6)
	cell(pdf, 1, 0.5, ":", "", 0, "J")
	pdf.SetX(8.9)
	cell(pdf, 15, 0.5, cetak.TempatLahirSiswa+", "+formatTanggalIndo(cetak.TanggalLahirSiswa), "", 0, "J")
	pdf.Ln(0.5)

	pdf.SetX(3)
	cell(pdf, 5, 0.5, "Nomor Induk Siswa Nasional", "", 0, "J")
	pdf.SetX(8.6)
	cell(pdf, 1, 0.5, ":", "", 0, "J")
	pdf.SetX(8.9)
	cell(pdf, 15, 0.5, cetak.NISN, "", 0, "J")
	pdf.Ln(0.5)

	pdf.SetX(3)
	cell(pdf, 5, 0.5, "Satuan Pendidikan", "", 0, "J")
	pdf.SetX(8.6)
	cell(pdf, 1, 0.5, ":", "", 0, "J")
	pdf.SetX(8.9)
	cell(pdf, 15, 0.5, cetak.NamaSekolah, "", 0, "J")
	pdf.Ln(0.5)

	pdf.SetX(3)
	cell(pdf, 5, 0.5, "Nomor Pokok Sekolah Nasional", "", 0, "J")
	pdf.SetX(8.6)
	cell(pdf, 1, 0.5, ":", "", 0, "J")
	pdf.SetX(8.9)
	cell(pdf, 15, 0.5, cetak.NPSN, "", 0, "J")
	pdf.Ln(0.5)

	pdf.SetX(3)
	cell(pdf, 5, 0.5, "Tanggal Kelulusan", "", 0, "J")
	pdf.SetX(8.6)
	cell(pdf, 1, 0.5, ":", "", 0, "J")
	pdf.SetX(8.9)
	// TODO: field lulus statis "2 Juni 2026" di CI lama — sebaiknya dijadikan
	// field baru di kelulusan_waktu (mis. tgl_lulus) daripada di-hardcode.
	cell(pdf, 15, 0.5, skl.TglSKL, "", 0, "J")
	pdf.Ln(0.5)

	pdf.SetX(3)
	cell(pdf, 5, 0.5, "Dinyatakan", "", 0, "J")
	pdf.SetX(8.6)
	cell(pdf, 1, 0.5, ":", "", 0, "J")
	pdf.SetFont("Times", "B", 12)
	pdf.SetX(8.9)
	cell(pdf, 15, 0.5, cetak.Ket, "", 0, "J")
	pdf.Ln(0.7)

	pdf.SetFont("Times", "", 12)
	pdf.SetX(2.5)
	multiCell(pdf, 16.5, 0.5, skl.Isi4SKL, "", "J")

	pdf.Ln(0.4)
	pdf.SetX(2.5)
	pdf.SetFont("Times", "B", 12)
	cell(pdf, 1.7, 0.7, "No", "1", 0, "C")
	cell(pdf, 11, 0.7, "Mata Pelajaran", "1", 0, "C")
	cell(pdf, 4, 0.7, "Nilai", "1", 0, "C")
	pdf.Ln(-1)

	type mapel struct {
		no, nama, nilai string
	}
	daftarMapel := []mapel{
		{"1", "Pendidikan Agama dan Budi Pekerti", cetak.PaS},
		{"2", "Pendidikan Pancasila dan Kewarganegaraan", cetak.PpknS},
		{"3", "Bahasa Indonesia", cetak.BindS},
		{"4", "Matematika", cetak.MtkS},
		{"5", "Ilmu Pengetahuan Alam", cetak.IpaS},
		{"6", "Ilmu Pengetahuan Sosial", cetak.IpsS},
		{"7", "Bahasa Inggris", cetak.BingS},
	}
	pdf.SetFont("Times", "", 12)
	for _, m := range daftarMapel {
		pdf.SetX(2.5)
		cell(pdf, 1.7, 0.7, m.no, "1", 0, "C")
		cell(pdf, 11, 0.7, m.nama, "1", 0, "J")
		cell(pdf, 4, 0.7, m.nilai, "1", 0, "C")
		pdf.Ln(-1)
	}

	pdf.SetX(2.5)
	cell(pdf, 1.7, 0.7, "8", "1", 0, "C")
	cell(pdf, 11, 0.7, "Seni Budaya", "1", 0, "J")
	cell(pdf, 4, 0.7, cetak.SenbudS, "1", 0, "C")
	pdf.Ln(-1)

	pdf.SetX(2.5)
	cell(pdf, 1.7, 0.7, "9", "1", 0, "C")
	cell(pdf, 11, 0.7, "Pendidikan Jasmani, Olahraga dan Kesehatan", "1", 0, "J")
	cell(pdf, 4, 0.7, cetak.PenjasS, "1", 0, "C")
	pdf.Ln(-1)

	pdf.SetX(2.5)
	cell(pdf, 1.7, 0.7, "10", "1", 0, "C")
	cell(pdf, 11, 0.7, "Informatika", "1", 0, "J")
	cell(pdf, 4, 0.7, cetak.PrakaryaS, "1", 0, "C")
	pdf.Ln(-1)

	pdf.SetX(2.5)
	pdf.SetFont("Times", "B", 12)
	cell(pdf, 1.7, 0.7, "", "1", 0, "C")
	cell(pdf, 11, 0.7, "Rata-rata", "1", 0, "J")
	cell(pdf, 4, 0.7, cetak.RataS, "1", 0, "C")
	pdf.Ln(1)

	pdf.SetFont("Times", "", 12)
	pdf.SetX(2.5)
	multiCell(pdf, 16.5, 0.5, skl.Isi5SKL, "", "J")
	pdf.Ln(-1)

	pdf.SetX(2.5)
	cell(pdf, 0.7, 0.7, "", "", 0, "C")
	cell(pdf, 10, 0.7, "", "", 0, "J")
	cell(pdf, 4, 0.7, skl.TglSKL, "", 0, "J")
	cell(pdf, 4, 0.7, "", "", 0, "C")
	pdf.Ln(0.5)

	pdf.SetX(2.5)
	cell(pdf, 0.7, 0.7, "", "", 0, "C")
	cell(pdf, 10, 0.7, "", "", 0, "J")
	cell(pdf, 4, 0.7, skl.SekolahSKL, "", 0, "J")
	cell(pdf, 4, 0.7, "", "", 0, "C")
	pdf.Ln(1)

	pdf.SetX(2.5)
	cell(pdf, 0.7, 0.7, "", "", 0, "C")
	cell(pdf, 11, 0.7, "", "", 0, "J")
	cell(pdf, 4, 0.7, skl.TtdSKL, "", 0, "J")
	cell(pdf, 4, 0.7, "", "", 0, "C")
	pdf.Ln(1)

	pdf.SetFont("Times", "U", 12)
	pdf.SetX(2.5)
	cell(pdf, 0.7, 0.7, "", "", 0, "C")
	cell(pdf, 10, 0.7, "", "", 0, "J")
	cell(pdf, 4, 0.7, skl.KepsekSKL, "", 0, "J")
	cell(pdf, 4, 0.7, "", "", 0, "C")
	pdf.Ln(0.5)

	pdf.SetFont("Times", "", 12)
	pdf.SetX(2.5)
	cell(pdf, 0.7, 0.7, "", "", 0, "C")
	cell(pdf, 10, 0.7, "", "", 0, "J")
	cell(pdf, 4, 0.7, skl.NipKepsekSKL, "", 0, "J")
	cell(pdf, 4, 0.7, "", "", 0, "C")

	// heroarea1 -> tampilkan gambar tanda tangan; aturform -> tampilkan foto siswa
	if skl.HeroArea1 != "" && skl.HeroArea1 != "0" {
		imgIfExists(pdf, "ttd.png", 11.3, 23, 7, 3.8)
	}
	if skl.AturForm != "" && skl.AturForm != "0" && cetak.FotoLulus != "" {
		imgIfExists(pdf, cetak.FotoLulus, 9.5, 23, 3, 4)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
