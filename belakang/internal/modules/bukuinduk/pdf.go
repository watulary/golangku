package bukuinduk

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

// BukuIndukPDFInput adalah seluruh data yg dibutuhkan utk merender 1
// file PDF Buku Induk (baik cetak satu siswa maupun massal).
type BukuIndukPDFInput struct {
	SiswaList       []BukuIndukData
	Profil          SekolahProfil
	Kepsek          KepalaSekolahInfo
	TahunPelajaran  string
	NamaPetugas     string
	NIPPetugas      string
	RentangNISAwal  string
	RentangNISAkhir string
	Opsi            CetakOpsi // lembar mana yg dicetak; zero-value -> DefaultCetakOpsi() (cetak semua)
}

// ukuran halaman & margin standar A4, dipakai konsisten di semua halaman
const (
	pageWidth  = 210.0
	pageHeight = 297.0
	marginL    = 12.0
	marginR    = 12.0
	marginT    = 15.0
	marginB    = 15.0
)

// GenerateBukuIndukPDF adalah entry point tunggal — menyusun semua
// halaman berurutan sesuai struktur asli (I-VI), lalu mengembalikan
// bytes PDF jadi.
func GenerateBukuIndukPDF(in BukuIndukPDFInput) ([]byte, error) {
	opsi := in.Opsi
	if !opsi.adaLembarDicetak() {
		opsi = DefaultCetakOpsi() // Opsi tidak diisi (caller lama) -> perilaku lama, cetak semua
	}

	kepsek := opsi.terapkanKepsek(in.Kepsek)

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginL, marginT, marginR)
	pdf.SetAutoPageBreak(true, marginB)
	registerLogos(pdf)
	isFirstPage := true // halaman pertama yg BENAR2 dicetak (apa pun lembarnya) tidak diberi watermark
	pdf.SetHeaderFunc(func() {
		if isFirstPage {
			isFirstPage = false
			return
		}
		renderWatermark(pdf)
	})

	if opsi.Sampul {
		renderCoverPage(pdf, in.Profil, in.TahunPelajaran)
	}

	waliKelasNama, waliKelasNIP := "", ""
	if len(in.SiswaList) > 0 {
		waliKelasNama = in.SiswaList[0].WaliKelasNama
		waliKelasNIP = in.SiswaList[0].WaliKelasNIPNRK
	}
	if opsi.WaliKelasNamaOverride != "" {
		waliKelasNama = opsi.WaliKelasNamaOverride
	}
	if opsi.WaliKelasNIPOverride != "" {
		waliKelasNIP = opsi.WaliKelasNIPOverride
	}

	if opsi.DaftarIsi {
		renderDaftarIsiPage(pdf, kepsek, in.NamaPetugas, in.NIPPetugas, in.Profil.KabupatenSekolah, in.RentangNISAwal, in.RentangNISAkhir, in.SiswaList, waliKelasNama, waliKelasNIP)
	}

	if opsi.Pemeriksaan {
		renderPemeriksaanPage(pdf)
	}

	if opsi.NomorUrut {
		renderNomorUrutPage(pdf, in.SiswaList, kepsek, in.NamaPetugas, in.NIPPetugas, in.Profil.KabupatenSekolah, waliKelasNama, waliKelasNIP)
	}

	if opsi.Petunjuk {
		renderPetunjukPages(pdf)
	}

	tingkatDicetak := opsi.tingkatMap()
	for _, s := range in.SiswaList {
		if opsi.LembarIV {
			renderLembarIVPages(pdf, s, in.TahunPelajaran)
		}

		renderDataNilaiPerKelas(pdf, s, in.TahunPelajaran, kepsek, tingkatDicetak)

		if opsi.HalamanAkhir {
			renderDataNilaiPage2(pdf, s, in.TahunPelajaran, kepsek)
		}
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("gagal generate PDF: %w", err)
	}
	return buf.Bytes(), nil
}

// --- helper umum, dipakai lintas halaman ---

// drawDigitBoxes menggambar kotak-kotak kosong bergaya "kotak isian
// digit" (mis. NO. STATISTIK SEKOLAH), sejumlah n kotak mulai dari x,y.
func drawDigitBoxes(pdf *fpdf.Fpdf, x, y, boxSize float64, n int) {
	for i := 0; i < n; i++ {
		pdf.Rect(x+float64(i)*boxSize, y, boxSize, boxSize, "D")
	}
}

// labeledLine menulis "Label : nilai" dgn label rata kiri lebar tetap
// (dipakai di banyak halaman isian: Cover, Lembar IV, dst).
func labeledLine(pdf *fpdf.Fpdf, label, value string, labelWidth, lineHeight float64) {
	x := pdf.GetX()
	pdf.CellFormat(labelWidth, lineHeight, label, "", 0, "L", false, 0, "")
	pdf.CellFormat(4, lineHeight, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, lineHeight, value, "", 1, "L", false, 0, "")
	pdf.SetX(x)
}

// centerText menulis teks rata tengah horizontal di posisi Y saat ini.
func centerText(pdf *fpdf.Fpdf, text string, fontSize float64, bold bool) {
	style := ""
	if bold {
		style = "B"
	}
	pdf.SetFont("Arial", style, fontSize)
	pdf.CellFormat(0, fontSize/2, text, "", 1, "C", false, 0, "")
}

// renderCoverPage — Halaman Sampul (data diri sekolah kosong siap isi
// no. statistik/kode sekolah tulis tangan, judul, kotak identitas
// sekolah terisi otomatis dari SekolahProfil + tahun pelajaran pilihan).
// renderCoverPage — versi baru "BUKU INDUK REGISTER PESERTA DIDIK
// KURIKULUM MERDEKA": logo Merdeka Belajar pojok kiri, 3 baris kotak
// digit (NSS/NPSN/Kode Sekolah masing2 12 kotak) di kanan atas, logo
// TutWuri di tengah, judul 4 baris, kotak identitas 6 baris (tanpa
// Kecamatan — form baru tidak menampilkannya).
func renderCoverPage(pdf *fpdf.Fpdf, profil SekolahProfil, tahunPelajaran string) {
	pdf.AddPage()

	// logo Merdeka Belajar, pojok kiri atas (rasio asli 501x345 px)
	merdekaW := 35.0
	if logoAvailable["logo-merdeka"] {
		pdf.ImageOptions("logo-merdeka", marginL, marginT, merdekaW, 0, false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	}

	// 3 baris kotak digit di kanan atas: NSS, NPSN, Kode Sekolah — 12 kotak/baris
	const headerBoxSize, headerFontSize = 6.0, 11.0
	boxGroupWidth := 12 * headerBoxSize
	blockX := pageWidth - marginR - boxGroupWidth

	pdf.SetFont("Arial", "B", 9)
	y := marginT
	pdf.SetXY(blockX, y)
	pdf.CellFormat(boxGroupWidth, 5, "NOMOR STATISTIK SEKOLAH", "", 0, "L", false, 0, "")
	drawDigitBoxesWithValue(pdf, blockX, y+6, headerBoxSize, headerFontSize, 12, profil.NSS)

	y += 15
	pdf.SetFont("Arial", "B", 9)
	pdf.SetXY(blockX, y)
	pdf.CellFormat(boxGroupWidth, 5, "NOMOR POKOK SEKOLAH NASIONAL", "", 0, "L", false, 0, "")
	drawDigitBoxesWithValue(pdf, blockX, y+6, headerBoxSize, headerFontSize, 12, profil.NPSN)

	y += 15
	pdf.SetFont("Arial", "B", 9)
	pdf.SetXY(blockX, y)
	pdf.CellFormat(boxGroupWidth, 5, "KODE SEKOLAH", "", 0, "L", false, 0, "")
	drawDigitBoxesWithValue(pdf, blockX, y+6, headerBoxSize, headerFontSize, 12, profil.NSP)

	// logo Tut Wuri Handayani, tengah
	tutWuriW := 55.0
	if logoAvailable["logo-tutwuri"] {
		pdf.ImageOptions("logo-tutwuri", (pageWidth-tutWuriW)/2, marginT+53, tutWuriW, 0, false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	}

	pdf.SetY(marginT + 115)
	pdf.SetTextColor(30, 125, 191)
	centerText(pdf, "BUKU INDUK REGISTER", 24, true)
	centerText(pdf, "PESERTA DIDIK", 24, true)
	centerText(pdf, "KURIKULUM MERDEKA", 24, true)
	centerText(pdf, "SEKOLAH MENENGAH PERTAMA (SMP)", 16, true)
	pdf.SetTextColor(0, 0, 0)

	// kotak identitas sekolah — 6 baris (tanpa Kecamatan, sesuai form baru)
	boxY := pdf.GetY() + 12
	interiorWidth := (pageWidth - marginL - marginR) - 16
	valueWidth := interiorWidth - 45 - 4

	pdf.SetXY(marginL+8, boxY+8)
	pdf.SetFont("Arial", "", 11)
	rows := []struct{ label, value string }{
		{"Nama Sekolah", profil.NamaSekolah},
		{"Alamat Sekolah", profil.AlamatSekolah},
		{"Desa/Kelurahan", profil.KelurahanSekolah},
		{"Kota/Kabupaten", profil.KabupatenSekolah},
		{"Provinsi", profil.ProvinsiSekolah},
		{"Tahun Pelajaran", tahunPelajaran},
	}
	for _, r := range rows {
		pdf.SetX(marginL + 8)
		labeledLineBoxed(pdf, r.label, r.value, 45, valueWidth, 8)
	}

	boxHeight := (pdf.GetY() + 8) - boxY
	pdf.Rect(marginL, boxY, pageWidth-marginL-marginR, boxHeight, "D")
}

// labeledLineBoxed sama seperti labeledLine, tapi nilai (value) dibatasi
// lebar tertentu dan WRAP otomatis ke baris baru kalau teksnya panjang
// (mis. alamat sekolah) — supaya tidak meluber keluar kotak pembatas.
func labeledLineBoxed(pdf *fpdf.Fpdf, label, value string, labelWidth, valueWidth, lineHeight float64) {
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(labelWidth, lineHeight, label, "", 0, "L", false, 0, "")
	pdf.CellFormat(4, lineHeight, ":", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.MultiCell(valueWidth, lineHeight, value, "", "L", false)
}

// renderSignatureBlock menggambar blok tanda tangan dua kolom yg
// berulang di banyak halaman: kiri "Mengetahui, Kepala Sekolah" (data
// otomatis dari KepalaSekolahInfo), kanan "[Kota], [tanggal] / Petugas"
// (nama otomatis dari admin yg login, tanggal = tanggal cetak).
//
// CATATAN ASUMSI: form asli menulis literal "Jakarta" sbg kota pada
// baris kanan (kemungkinan itu cuma placeholder bawaan template, bukan
// data sungguhan) — saya ganti jadi kabupaten/kota sekolah dari profil
// supaya sesuai sekolah masing-masing. Tolong koreksi kalau ternyata
// memang harus tetap literal "Jakarta".
func renderSignatureBlock(pdf *fpdf.Fpdf, kepsek KepalaSekolahInfo, namaPetugas string, nipPetugas string, kotaSekolah string) {
	pdf.Ln(10)
	colWidth := (pageWidth - marginL - marginR) / 2
	startY := pdf.GetY()

	pdf.SetFont("Arial", "", 11)

	// kolom kiri: Mengetahui, Kepala Sekolah
	pdf.SetXY(marginL, startY)
	pdf.CellFormat(colWidth, 6, "Mengetahui,", "", 2, "L", false, 0, "")
	pdf.SetX(marginL)
	pdf.CellFormat(colWidth, 6, "Kepala Sekolah,", "", 2, "L", false, 0, "")

	pdf.SetXY(marginL, startY+25)
	pdf.CellFormat(colWidth, 6, kepsek.NamaKepsek, "U", 2, "L", false, 0, "")
	pdf.SetX(marginL)
	pdf.CellFormat(colWidth, 6, formatNIP(kepsek.NIPKepsek), "", 2, "L", false, 0, "")

	// kolom kanan: [Kota], [tanggal cetak] / Petugas
	tanggalCetak := formatTanggalIndo(time.Now())
	pdf.SetXY(marginL+colWidth, startY)
	pdf.CellFormat(colWidth, 6, kotaSekolah+", "+tanggalCetak, "", 2, "L", false, 0, "")
	pdf.SetX(marginL + colWidth)
	pdf.CellFormat(colWidth, 6, "Petugas,", "", 2, "L", false, 0, "")

	pdf.SetXY(marginL+colWidth, startY+25)
	pdf.CellFormat(colWidth, 6, namaPetugas, "U", 2, "L", false, 0, "")
	pdf.SetX(marginL + colWidth)
	if nipPetugas != "" {
		pdf.CellFormat(colWidth, 6, formatNIP(nipPetugas), "", 2, "L", false, 0, "")
	} else {
		pdf.CellFormat(colWidth, 6, "NIP.", "", 2, "L", false, 0, "")
	}
}

// formatTanggalIndo memformat tanggal ke gaya "5 September 2026".
func formatTanggalIndo(t time.Time) string {
	bulan := [...]string{"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	return fmt.Sprintf("%d %s %d", t.Day(), bulan[t.Month()-1], t.Year())
}

// renderDaftarIsiPage — Halaman Daftar Isi (statis, sama utk semua cetak).
func renderDaftarIsiPage(pdf *fpdf.Fpdf, kepsek KepalaSekolahInfo, namaPetugas, nipPetugas, kotaSekolah string, nisAwal, nisAkhir string, siswaList []BukuIndukData, waliKelasNama, waliKelasNIP string) {
	pdf.AddPage()
	centerText(pdf, "DAFTAR ISI", 14, true)
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 11)
	items := []string{
		"I.    PEMERIKSAAN BUKU INDUK SISWA",
		"II.   NOMOR URUT MURID DALAM BUKU INDUK SISWA",
		"III.  PETUNJUK PENGISIAN BUKU INDUK SISWA",
		"      1. PETUNJUK UMUM",
		"      2. PETUNJUK KHUSUS",
		"IV.   DATA NILAI SISWA",
	}
	for _, item := range items {
		pdf.CellFormat(0, 7, item, "", 1, "L", false, 0, "")
	}

	pdf.Ln(4)
	jumlahL, jumlahP := computeJumlahLP(siswaList)
	pdf.CellFormat(0, 7, fmt.Sprintf("JUMLAH MURID LAKI-LAKI %d SISWA", jumlahL), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 7, fmt.Sprintf("JUMLAH MURID WANITA %d SISWI", jumlahP), "", 1, "L", false, 0, "")

	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 7, "KETERANGAN:", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 7, "Halaman .......................................................................... gagal", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 7, "Halaman .......................................................................... tidak terpakai", "", 1, "L", false, 0, "")
	pdf.Ln(7)
	rentangText := fmt.Sprintf("Nomor induk siswa (NIS) mulai dari nomor %s s/d %s", nisAwal, nisAkhir)
	pdf.CellFormat(0, 7, rentangText, "", 1, "C", false, 0, "")
	pdf.Ln(7)

	renderSignatureBlock3(pdf, kepsek, waliKelasNama, waliKelasNIP, namaPetugas, nipPetugas, kotaSekolah)
}

// renderPemeriksaanPage — Halaman I: tabel pemeriksaan buku induk
// register (kosong, diisi tangan tiap kali ada pemeriksaan berkala).
func renderPemeriksaanPage(pdf *fpdf.Fpdf) {
	pdf.AddPage()
	centerText(pdf, "I. PEMERIKSAAN BUKU INDUK REGISTER", 13, true)
	pdf.Ln(4)

	colWidths := []float64{15, 30, 55, 30, 56}
	headers := []string{"No. Urut", "Hari/Tanggal", "Nama Petugas Pemeriksa", "Jabatan", "Catatan"}
	headerHeight := 8.0
	rowHeight := 8.0

	drawPemeriksaanHeader := func() {
		pdf.SetFillColor(224, 224, 224)
		pdf.SetFont("Arial", "B", 9)
		y := pdf.GetY()
		x := marginL
		for i, h := range headers {
			pdf.Rect(x, y, colWidths[i], headerHeight, "DF")
			pdf.SetXY(x, y)
			pdf.CellFormat(colWidths[i], headerHeight, h, "", 0, "C", false, 0, "")
			x += colWidths[i]
		}
		pdf.SetXY(marginL, y+headerHeight)
	}

	drawPemeriksaanHeader()
	pdf.SetFont("Arial", "", 9)

	for row := 0; row < 30; row++ {
		if pdf.GetY()+rowHeight > pageHeight-marginB {
			pdf.AddPage()
			drawPemeriksaanHeader()
			pdf.SetFont("Arial", "", 9)
		}
		y := pdf.GetY()
		x := marginL
		for _, w := range colWidths {
			pdf.Rect(x, y, w, rowHeight, "D")
			x += w
		}
		pdf.SetXY(marginL, y+rowHeight)
	}
}

func formatNIP(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if strings.HasPrefix(strings.ToUpper(trimmed), "NIP") {
		return trimmed
	}
	return "NIP. " + trimmed
}

// logoAvailable melacak logo mana yg berhasil dimuat, supaya
// renderCoverPage/renderWatermark tidak memanggil ImageOptions dgn
// nama gambar yg gagal ter-register (bisa bikin pdf.Err() kotor).
var logoAvailable = map[string]bool{}

func registerLogos(pdf *fpdf.Fpdf) {
	registerOne := func(name, path string) {
		data, err := os.ReadFile(path)
		if err != nil {
			logoAvailable[name] = false
			return
		}
		pdf.RegisterImageOptionsReader(name, fpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(data))
		logoAvailable[name] = true
	}
	registerOne("logo-tutwuri", "storage/bukuinduk/tutwuri.png")
	registerOne("logo-orang", "storage/bukuinduk/orang.png")
	registerOne("logo-merdeka", "storage/bukuinduk/merdeka.png")
	registerOne("logo-254", "storage/bukuinduk/logo254.png")
}

// drawDigitBoxesWithValue seperti drawDigitBoxes, tapi tiap kotak diisi
// satu karakter dari value (rata kiri; sisa kotak kosong kalau value
// lebih pendek dari n, karakter kelebihan dipotong kalau lebih panjang).
func drawDigitBoxesWithValue(pdf *fpdf.Fpdf, x, y, boxSize, fontSize float64, n int, value string) {
	runes := []rune(value)
	for i := 0; i < n; i++ {
		bx := x + float64(i)*boxSize
		pdf.Rect(bx, y, boxSize, boxSize, "D")
		if i < len(runes) {
			pdf.SetFont("Arial", "", fontSize)
			pdf.SetXY(bx, y)
			pdf.CellFormat(boxSize, boxSize, string(runes[i]), "", 0, "C", false, 0, "")
		}
	}
}

// renderWatermark menggambar logo Tut Wuri Handayani transparan besar
// di tengah halaman sbg watermark latar belakang.
func renderWatermark(pdf *fpdf.Fpdf) {
	wmSize := 130.0
	x := (pageWidth - wmSize) / 2
	y := (pageHeight-wmSize)/2 - 20 // digeser naik ~2cm
	if logoAvailable["logo-254"] {
		pdf.SetAlpha(0.08, "Normal")
		pdf.ImageOptions("logo-254", x, y, wmSize, 0, false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")
		pdf.SetAlpha(1, "Normal")
	}
}

// renderSignatureBlock3 — versi 3 kolom (Kepala Sekolah | Wali Kelas |
// Petugas), dipakai di halaman umum (Daftar Isi/Bag.I/Bag.II) pada
// format baru. Wali Kelas & NIP-nya diambil dari SISWA PERTAMA di
// daftar yg dicetak (tepat utk cetak 1 siswa/1 rombel; utk bulk lintas
// rombel, itu jadi representasi siswa pertama saja).
func renderSignatureBlock3(pdf *fpdf.Fpdf, kepsek KepalaSekolahInfo, waliKelasNama, waliKelasNIP string, namaPetugas, nipPetugas string, kotaSekolah string) {
	pdf.Ln(10)
	colWidth := (pageWidth - marginL - marginR) / 3
	startY := pdf.GetY()
	pdf.SetFont("Arial", "", 10)

	// kolom 1: Mengetahui, Kepala Sekolah
	pdf.SetXY(marginL, startY)
	pdf.CellFormat(colWidth, 6, "Mengetahui,", "", 2, "L", false, 0, "")
	pdf.SetX(marginL)
	pdf.CellFormat(colWidth, 6, "Kepala Sekolah", "", 2, "L", false, 0, "")

	// kolom 2: Wali Kelas
	pdf.SetXY(marginL+colWidth, startY+6)
	pdf.CellFormat(colWidth, 6, "Wali Kelas", "", 2, "L", false, 0, "")

	// kolom 3: [Kota], [tanggal] / Petugas
	tanggalCetak := formatTanggalIndo(time.Now())
	pdf.SetXY(marginL+2*colWidth, startY)
	pdf.CellFormat(colWidth, 6, kotaSekolah+", "+tanggalCetak, "", 2, "L", false, 0, "")
	pdf.SetX(marginL + 2*colWidth)
	pdf.CellFormat(colWidth, 6, "Petugas", "", 2, "L", false, 0, "")

	// baris nama & NIP, sejajar 25mm di bawah baris label
	nameY := startY + 25
	pdf.SetXY(marginL, nameY)
	pdf.CellFormat(colWidth, 6, kepsek.NamaKepsek, "U", 2, "L", false, 0, "")
	pdf.SetX(marginL)
	pdf.CellFormat(colWidth, 6, formatNIP(kepsek.NIPKepsek), "", 2, "L", false, 0, "")

	pdf.SetXY(marginL+colWidth, nameY)
	pdf.CellFormat(colWidth, 6, waliKelasNama, "U", 2, "L", false, 0, "")
	pdf.SetX(marginL + colWidth)
	if waliKelasNIP != "" {
		pdf.CellFormat(colWidth, 6, formatNIP(waliKelasNIP), "", 2, "L", false, 0, "")
	} else {
		pdf.CellFormat(colWidth, 6, "NIP.", "", 2, "L", false, 0, "")
	}

	pdf.SetXY(marginL+2*colWidth, nameY)
	pdf.CellFormat(colWidth, 6, namaPetugas, "U", 2, "L", false, 0, "")
	pdf.SetX(marginL + 2*colWidth)
	if nipPetugas != "" {
		pdf.CellFormat(colWidth, 6, formatNIP(nipPetugas), "", 2, "L", false, 0, "")
	} else {
		pdf.CellFormat(colWidth, 6, "NIP.", "", 2, "L", false, 0, "")
	}
}

// computeJumlahLP menghitung jumlah siswa laki-laki & perempuan dari
// daftar yg SEDANG DICETAK (bukan seluruh sekolah).
func computeJumlahLP(siswaList []BukuIndukData) (jumlahL, jumlahP int) {
	for _, s := range siswaList {
		switch jenisKelaminSingkat(s.JenisKelamin) {
		case "L":
			jumlahL++
		case "P":
			jumlahP++
		}
	}
	return
}

// ensureSpace memaksa pindah halaman kalau ruang tersisa di halaman
// ini kurang dari neededHeight — dipakai supaya grup butir bernomor
// atau baris ringkasan (Jumlah Nilai/Rata-rata) tidak terpotong
// page-break otomatis di tengah-tengah.
func ensureSpace(pdf *fpdf.Fpdf, neededHeight float64) {
	if pdf.GetY()+neededHeight > pageHeight-marginB {
		pdf.AddPage()
	}
}
