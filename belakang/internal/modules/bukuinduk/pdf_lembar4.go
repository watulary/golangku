package bukuinduk

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-pdf/fpdf"

	"sekolah/internal/modules/nilai"
)

// extraRightMargin menyisakan ruang kosong di sisi kanan halaman
// selama baris-baris tertentu dicetak (dipakai supaya teks tidak
// tertimpa kotak pas foto yg digambar terpisah di sisi kanan).
var extraRightMargin float64

// dotsIfEmpty menampilkan titik-titik ASCII kalau value kosong (bukan
// karakter ellipsis Unicode "…" — itu menyebabkan mojibake di font
// bawaan fpdf yg cuma dukung encoding satu-byte).
func dotsIfEmpty(value string) string {
	if value == "" {
		return "........................"
	}
	return value
}

// formItem menulis satu baris item bernomor "1. Label : value", value
// WRAP otomatis (MultiCell) supaya kalimat panjang (mis. alamat) tidak
// terpotong/hilang di luar halaman.
func formItem(pdf *fpdf.Fpdf, no, label, value string, labelWidth float64) {
	const numWidth, colonWidth = 8.0, 4.0
	fullWidth := pageWidth - marginL - marginR - extraRightMargin
	valueWidth := fullWidth - numWidth - labelWidth - colonWidth

	pdf.SetX(marginL)
	pdf.CellFormat(numWidth, 6, no, "", 0, "L", false, 0, "")
	pdf.CellFormat(labelWidth, 6, label, "", 0, "L", false, 0, "")
	pdf.CellFormat(colonWidth, 6, ":", "", 0, "L", false, 0, "")
	pdf.MultiCell(valueWidth, 6, value, "", "L", false)
}

// formHeaderItem sama seperti formItem TAPI TANPA titik dua & value —
// dipakai utk baris header yang cuma menaungi sub-item di bawahnya
// (mis. "1. Nama Siswa"), bukan field data sendiri.
func formHeaderItem(pdf *fpdf.Fpdf, no, label string) {
	pdf.SetX(marginL)
	pdf.CellFormat(8, 6, no, "", 0, "L", false, 0, "")
	width := (pageWidth - marginL - marginR - extraRightMargin) - 8
	pdf.CellFormat(width, 6, label, "", 1, "L", false, 0, "")
}

// formSubItem sama seperti formItem tapi indentasi sub-poin ("a.","b.").
func formSubItem(pdf *fpdf.Fpdf, letter, label, value string, labelWidth float64) {
	const indent, numWidth, colonWidth = 10.0, 8.0, 4.0
	fullWidth := pageWidth - marginL - marginR
	valueWidth := fullWidth - numWidth - labelWidth - colonWidth

	pdf.SetX(marginL + indent)
	pdf.CellFormat(numWidth, 6, letter, "", 0, "L", false, 0, "")
	pdf.CellFormat(labelWidth-indent, 6, label, "", 0, "L", false, 0, "")
	pdf.CellFormat(colonWidth, 6, ":", "", 0, "L", false, 0, "")
	pdf.MultiCell(valueWidth, 6, value, "", "L", false)
}

// formSubHeaderItem sama seperti formHeaderItem tapi utk level "a."/"b."
// (dipakai utk "a. Masuk menjadi siswa kls 1" & "b. Pindahan dari
// sekolah lain" — sama-sama header tanpa data sendiri, konsisten dgn
// perlakuan formHeaderItem, meski tidak eksplisit disebut).
func formSubHeaderItem(pdf *fpdf.Fpdf, letter, label string) {
	pdf.SetX(marginL + 10)
	pdf.CellFormat(8, 6, letter, "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, label, "", 1, "L", false, 0, "")
}

// formSubSubItem utk item bernomor level ketiga (1./2./3. di bawah
// a./b.), labelWidth TETAP sama utk semua panggilan supaya kolom titik
// dua sejajar rapi (sebelumnya tiap baris pakai lebar label berbeda2).
func formSubSubItem(pdf *fpdf.Fpdf, no, label, value string, labelWidth float64) {
	const indent, numWidth, colonWidth = 18.0, 8.0, 4.0
	availableWidth := (pageWidth - marginL - marginR) - indent
	valueWidth := availableWidth - numWidth - labelWidth - colonWidth

	pdf.SetX(marginL + indent)
	pdf.CellFormat(numWidth, 6, no, "", 0, "L", false, 0, "")
	pdf.CellFormat(labelWidth, 6, label, "", 0, "L", false, 0, "")
	pdf.CellFormat(colonWidth, 6, ":", "", 0, "L", false, 0, "")
	pdf.MultiCell(valueWidth, 6, value, "", "L", false)
}

func sectionHeading(pdf *fpdf.Fpdf, text string) {
	pdf.Ln(3)
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 7, text, "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
}

const formLabelWidth = 65.0
const subSubLabelWidth = 47.0

func renderLembarIVPages(pdf *fpdf.Fpdf, siswa BukuIndukData, _ string) {
	pdf.AddPage()
	centerText(pdf, "IV. LEMBAR BUKU INDUK SISWA", 13, true)
	pdf.Ln(4)

	pdf.SetFont("Arial", "", 10)
	pdf.SetX(marginL)
	pdf.CellFormat(50, 6, "Nomor Induk Siswa", "", 0, "L", false, 0, "")
	pdf.CellFormat(4, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, siswa.NIS, "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, "Lembar ke: 1", "", 1, "L", false, 0, "")

	pdf.SetX(marginL)
	pdf.CellFormat(50, 6, "Nomor Induk Siswa Nasional", "", 0, "L", false, 0, "")
	pdf.CellFormat(4, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, siswa.NISN, "", 1, "L", false, 0, "")

	pdf.Ln(2)

	// kotak foto 4x6, sejajar item 2-10 Bagian A (bukan di header atas)
	fotoW, fotoH := 26.0, 36.0
	fotoX := pageWidth - marginR - fotoW
	fotoStartY := pdf.GetY() + 12 // turun sedikit, sejajar mulai item "2. Jenis Kelamin"
	renderFotoBox(pdf, siswa.FotoURL, siswa.JenisKelamin, fotoX, fotoStartY, fotoW, fotoH)
	extraRightMargin = fotoW + 4

	sectionHeading(pdf, "A. Keterangan Tentang Diri Siswa")
	formHeaderItem(pdf, "1.", "Nama Siswa")
	formSubItem(pdf, "a.", "Nama Lengkap", siswa.NamaLengkap, formLabelWidth)
	formSubItem(pdf, "b.", "Nama Panggilan", siswa.NamaPanggilan, formLabelWidth)
	formItem(pdf, "2.", "Jenis Kelamin", siswa.JenisKelamin, formLabelWidth)
	formItem(pdf, "3.", "Tempat dan Tanggal Lahir", tempatTanggalLahir(siswa), formLabelWidth)
	formItem(pdf, "4.", "Agama", siswa.Agama, formLabelWidth)
	formItem(pdf, "5.", "Kewarganegaraan", siswa.Kewarganegaraan, formLabelWidth)
	formItem(pdf, "6.", "Anak Keberapa", siswa.AnakKe, formLabelWidth)
	formItem(pdf, "7.", "Jumlah Saudara Kandung", siswa.JumlahSaudara, formLabelWidth)
	formItem(pdf, "8.", "Jumlah Saudara Tiri", siswa.SaudaraTiri, formLabelWidth)
	formItem(pdf, "9.", "Jumlah Saudara Angkat", siswa.SaudaraAngkat, formLabelWidth)
	formItem(pdf, "10.", "Anak yatim / Yatim Piatu", siswa.YatimPiatu, formLabelWidth)
	formItem(pdf, "11.", "Bahasa Sehari-hari Dirumah", siswa.BahasaSehariHari, formLabelWidth)
	// ================= B. Keterangan Tempat Tinggal =================
	sectionHeading(pdf, "B. Keterangan Tempat Tinggal")
	formItem(pdf, "12.", "Alamat", siswa.Alamat, formLabelWidth)
	formItem(pdf, "13.", "Nomor Telepon", siswa.NoTelpHp, formLabelWidth)
	formItem(pdf, "14.", "Tinggal Dengan Orang Tua / Saudara", siswa.BertempatTinggalPada, formLabelWidth)
	formItem(pdf, "", "Di Asrama / Kost", "", formLabelWidth)
	formItem(pdf, "15.", "Jarak Tempat Tinggal ke Sekolah", jarakKmLabel(siswa.JarakSekolah), formLabelWidth)

	extraRightMargin = 0 // Bagian B dst sudah pasti di bawah foto (foto cuma setinggi Bagian A)

	// ================= C. Keterangan Kesehatan =================
	sectionHeading(pdf, "C. Keterangan Kesehatan")
	formItem(pdf, "16.", "Golongan Darah", siswa.GolonganDarah, formLabelWidth)
	formItem(pdf, "17.", "TBC / Cacar / Malaria dan lain-lain", siswa.PernahSakit, formLabelWidth)
	formItem(pdf, "18.", "Kelainan Jasmani", siswa.KebutuhanKhusus, formLabelWidth)
	formItem(pdf, "19.", "Tinggi dan Berat Badan", tinggiBeratLabel(siswa), formLabelWidth)

	// ================= D. Keterangan Pendidikan =================
	sectionHeading(pdf, "D. Keterangan Pendidikan")
	ensureSpace(pdf, 26)
	formHeaderItem(pdf, "20.", "Pendidikan Sebelumnya")
	formSubItem(pdf, "a.", "Lulus Dari", siswa.AsalSekolahMasuk, formLabelWidth)
	formSubItem(pdf, "b.", "Tanggal dan Nomor STTB", tanggalNomorIjazah(siswa), formLabelWidth)
	formSubItem(pdf, "c.", "Lama Belajar", siswa.LamaBelajar, formLabelWidth)
	ensureSpace(pdf, 20)
	formHeaderItem(pdf, "21.", "Pindahan")
	formSubItem(pdf, "a.", "Dari Sekolah", siswa.AsalSekolahPindahan, formLabelWidth)
	formSubItem(pdf, "b.", "Alasan", "", formLabelWidth)
	ensureSpace(pdf, 32)
	formHeaderItem(pdf, "22.", "Diterima Disekolah ini")
	formSubItem(pdf, "a.", "Ditingkat", "", formLabelWidth)
	formSubItem(pdf, "b.", "Kelompok", "", formLabelWidth)
	formSubItem(pdf, "c.", "Jurusan", "", formLabelWidth)
	formSubItem(pdf, "d.", "Tanggal", siswa.DiterimaTanggalPindahan, formLabelWidth)

	// ================= E. Keterangan Ayah Kandung =================
	sectionHeading(pdf, "E. Keterangan Ayah Kandung")
	formItem(pdf, "23.", "Nama", siswa.NamaAyah, formLabelWidth)
	formItem(pdf, "24.", "Tempat dan Tanggal Lahir", tempatTanggal(siswa.TempatLahirAyah, siswa.TanggalLahirAyah), formLabelWidth)
	formItem(pdf, "25.", "Agama", siswa.AgamaAyah, formLabelWidth)
	formItem(pdf, "26.", "Kewarganegaraan", siswa.WargaAyah, formLabelWidth)
	formItem(pdf, "27.", "Pendidikan", siswa.PendidikanAyah, formLabelWidth)
	formItem(pdf, "28.", "Pekerjaan", siswa.PekerjaanAyah, formLabelWidth)
	formItem(pdf, "29.", "Penghasilan Perbulan", siswa.GajiAyah, formLabelWidth)
	formItem(pdf, "30.", "Alamat Rumah / Nomor Telepon", alamatTelp(siswa.AlamatAyah, siswa.TelpAyah), formLabelWidth)
	formItem(pdf, "31.", "Masih Hidup / Meninggal", siswa.StatusAyah, formLabelWidth)

	// ================= F. Keterangan Ibu Kandung =================
	sectionHeading(pdf, "F. Keterangan Ibu Kandung")
	formItem(pdf, "32.", "Nama", siswa.NamaIbu, formLabelWidth)
	formItem(pdf, "33.", "Tempat dan Tanggal Lahir", tempatTanggal(siswa.TempatLahirIbu, siswa.TanggalLahirIbu), formLabelWidth)
	formItem(pdf, "34.", "Agama", siswa.AgamaIbu, formLabelWidth)
	formItem(pdf, "35.", "Kewarganegaraan", siswa.WargaIbu, formLabelWidth)
	formItem(pdf, "36.", "Pendidikan", siswa.PendidikanIbu, formLabelWidth)
	formItem(pdf, "37.", "Pekerjaan", siswa.PekerjaanIbu, formLabelWidth)
	formItem(pdf, "38.", "Penghasilan Perbulan", siswa.GajiIbu, formLabelWidth)
	formItem(pdf, "39.", "Alamat Rumah / Nomor Telepon", alamatTelp(siswa.AlamatIbu, siswa.TelpIbu), formLabelWidth)
	formItem(pdf, "40.", "Masih Hidup / Meninggal", siswa.StatusIbu, formLabelWidth)

	// ================= G. Keterangan Wali =================
	sectionHeading(pdf, "G. Keterangan Wali")
	formItem(pdf, "41.", "Nama", siswa.NamaWali, formLabelWidth)
	formItem(pdf, "42.", "Tempat dan Tanggal Lahir", tempatTanggal(siswa.TempatLahirWali, siswa.TanggalLahirWali), formLabelWidth)
	formItem(pdf, "43.", "Agama", siswa.AgamaWali, formLabelWidth)
	formItem(pdf, "44.", "Kewarganegaraan", siswa.WargaWali, formLabelWidth)
	formItem(pdf, "45.", "Pendidikan", siswa.PendidikanWali, formLabelWidth)
	formItem(pdf, "46.", "Pekerjaan", siswa.PekerjaanWali, formLabelWidth)
	formItem(pdf, "47.", "Penghasilan Perbulan", siswa.GajiWali, formLabelWidth)
	formItem(pdf, "48.", "Alamat Rumah / Nomor Telepon", alamatTelp(siswa.AlamatWali, siswa.TelpWali), formLabelWidth)

	// ================= H. Kegemaran Siswa =================
	sectionHeading(pdf, "H. Kegemaran Siswa")
	formItem(pdf, "49.", "Kesenian", siswa.HobiSeni, formLabelWidth)
	formItem(pdf, "50.", "Olahraga", siswa.HobiOrkes, formLabelWidth)
	formItem(pdf, "51.", "Kemasyarakatan / Organisasi", siswa.HobiOrganisasi, formLabelWidth)
	formItem(pdf, "52.", "Lain-lain", siswa.HobiLain, formLabelWidth)

	// ================= I. Keterangan Perkembangan Siswa =================
	sectionHeading(pdf, "I. Keterangan Perkembangan Siswa")
	formItem(pdf, "53.", "Menerima Beasiswa", "", formLabelWidth)
	ensureSpace(pdf, 20)
	formHeaderItem(pdf, "54.", "Meninggalkan Sekolah")
	formSubItem(pdf, "a.", "Tanggal meninggalkan sekolah", siswa.KeluarTanggal, formLabelWidth)
	formSubItem(pdf, "b.", "Alasan", siswa.KeluarAlasan, formLabelWidth)
	ensureSpace(pdf, 20)
	formHeaderItem(pdf, "55.", "Akhir Pendidikan")
	formSubItem(pdf, "a.", "Tamat Belajar", siswa.LulusTahun, formLabelWidth)
	formSubItem(pdf, "b.", "Nomor STTB", siswa.LulusNomorIjazah, formLabelWidth)

	// ================= J. Keterangan Selesai Pendidikan =================
	sectionHeading(pdf, "J. Keterangan Selesai Pendidikan")
	formItem(pdf, "56.", "Melanjutkan di", siswa.LulusMelanjutkanKe, formLabelWidth)
	formHeaderItem(pdf, "57.", "Bekerja")
	formSubItem(pdf, "a.", "Tanggal Mulai Kerja", "", formLabelWidth)
	formSubItem(pdf, "b.", "Nama Perusahaan / Lembaga", "", formLabelWidth)
	formSubItem(pdf, "c.", "Penghasilan", "", formLabelWidth)
	renderBagianK(pdf, siswa)
	extraRightMargin = 0 // reset, baris2 berikutnya (Bagian B dst) sudah di bawah foto
}

func renderBagianK(pdf *fpdf.Fpdf, siswa BukuIndukData) {
	sectionHeading(pdf, "K. Catatan Prestasi")
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(0, 6, "Catatan Yang Penting", "", 1, "L", false, 0, "")

	fullWidth := pageWidth - marginL - marginR
	colKelas, colPrestasi := 15.0, 65.0
	colSem := (fullWidth - colKelas - colPrestasi) / 2

	y := pdf.GetY()
	pdf.SetFillColor(224, 224, 224)
	pdf.SetFont("Arial", "B", 8)
	pdf.Rect(marginL, y, colKelas, 12, "DF")
	pdf.SetXY(marginL, y+3)
	pdf.CellFormat(colKelas, 6, "Kelas", "", 0, "C", false, 0, "")
	pdf.Rect(marginL+colKelas, y, colPrestasi, 12, "DF")
	pdf.SetXY(marginL+colKelas, y+3)
	pdf.CellFormat(colPrestasi, 6, "Prestasi yang diraih", "", 0, "C", false, 0, "")
	pdf.Rect(marginL+colKelas+colPrestasi, y, colSem, 6, "DF")
	pdf.SetXY(marginL+colKelas+colPrestasi, y)
	pdf.CellFormat(colSem, 6, "Semester Ganjil", "", 0, "C", false, 0, "")
	pdf.Rect(marginL+colKelas+colPrestasi+colSem, y, colSem, 6, "DF")
	pdf.SetXY(marginL+colKelas+colPrestasi+colSem, y)
	pdf.CellFormat(colSem, 6, "Semester Genap", "", 0, "C", false, 0, "")

	subY := y + 6
	pNilaiW := colSem * 0.4 // Nilai (pendek) lebih sempit dari Keterangan
	pKetW := colSem - pNilaiW
	x := marginL + colKelas + colPrestasi
	for i := 0; i < 2; i++ {
		pdf.Rect(x, subY, pNilaiW, 6, "DF")
		pdf.SetXY(x, subY)
		pdf.CellFormat(pNilaiW, 6, "Nilai", "", 0, "C", false, 0, "")
		pdf.Rect(x+pNilaiW, subY, pKetW, 6, "DF")
		pdf.SetXY(x+pNilaiW, subY)
		pdf.CellFormat(pKetW, 6, "Keterangan", "", 0, "C", false, 0, "")
		x += colSem
	}

	dataY := subY + 6

	// Prestasi per tingkat kelas dari prestasi_siswa. Setiap prestasi = 1 baris; Nilai/Keterangan
	// tampil di kolom semester tempat prestasi itu dicatat. Kelas tanpa prestasi = 1 baris kosong.
	prestasiKelas, errP := nilai.GetRiwayatPrestasi(siswa.NIS)
	if errP != nil {
		log.Printf("buku induk: gagal membaca prestasi nis %s: %v", siswa.NIS, errP)
		prestasiKelas = nil
	}

	pdf.SetAutoPageBreak(false, 0) // page-break dikerjakan manual oleh ensureRow
	for tingkat := 7; tingkat <= 9; tingkat++ {
		entri := prestasiKelas[tingkat]
		if len(entri) == 0 {
			entri = []nilai.PrestasiSel{{}}
		}

		// kolom Kelas = 1 sel gabungan untuk seluruh prestasi kelas itu, nomor di tengah.
		// Kalau kelompok terpotong pindah halaman, tiap potongan digambar sendiri.
		segStart := dataY
		gambarKelas := func(yAkhir float64) {
			if yAkhir <= segStart {
				return
			}
			pdf.Rect(marginL, segStart, colKelas, yAkhir-segStart, "D")
			pdf.SetFont("Arial", "", 8)
			pdf.SetXY(marginL, segStart)
			pdf.CellFormat(colKelas, yAkhir-segStart, strconv.Itoa(tingkat), "", 0, "C", false, 0, "")
		}

		for _, e := range entri {
			pdf.SetFont("Arial", "", 8)
			maxH := float64(jumlahBaris(pdf, e.Prestasi, colPrestasi)) * 3.6

			pdf.SetFont("Arial", "", 6.5)
			if h := float64(jumlahBaris(pdf, e.Nilai, pNilaiW)) * 3.0; h > maxH {
				maxH = h
			}
			if h := float64(jumlahBaris(pdf, e.Keterangan, pKetW)) * 3.0; h > maxH {
				maxH = h
			}

			rowH := maxH + 1.6
			if rowH < 7 {
				rowH = 7
			}

			yBaru := ensureRow(pdf, dataY, rowH)
			if yBaru != dataY { // pindah halaman: tutup potongan sebelumnya
				gambarKelas(dataY)
				dataY = yBaru
				segStart = dataY
			}

			pdf.SetFont("Arial", "", 8)
			pdf.Rect(marginL+colKelas, dataY, colPrestasi, rowH, "D")
			drawTeksTengahV(pdf, marginL+colKelas, dataY, colPrestasi, rowH, 3.6, e.Prestasi, "L")

			x = marginL + colKelas + colPrestasi
			for _, sem := range semesterUrut {
				pdf.Rect(x, dataY, pNilaiW, rowH, "D")
				pdf.Rect(x+pNilaiW, dataY, pKetW, rowH, "D")
				if e.Semester == sem {
					pdf.SetFont("Arial", "", 6.5)
					drawTeksTengahV(pdf, x, dataY, pNilaiW, rowH, 3.0, e.Nilai, "C")
					drawTeksTengahV(pdf, x+pNilaiW, dataY, pKetW, rowH, 3.0, e.Keterangan, "C")
				}
				x += colSem
			}

			dataY += rowH
		}

		gambarKelas(dataY)
	}
	pdf.SetAutoPageBreak(true, marginB)
	pdf.SetY(dataY + 4)

	// ================= Ekstrakurikuler =================
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(0, 6, "Ekstrakurikuler", "", 1, "L", false, 0, "")

	colEkskul := 30.0
	kelasGroupWidth := (fullWidth - colEkskul) / 3
	semWidth := kelasGroupWidth / 2
	nilaiW := semWidth * 0.4 // Nilai (pendek) lebih sempit dari Keterangan
	ketW := semWidth - nilaiW

	y = pdf.GetY()
	pdf.SetFillColor(224, 224, 224)
	pdf.SetFont("Arial", "B", 7)
	pdf.Rect(marginL, y, colEkskul, 15, "DF")
	pdf.SetXY(marginL, y+5)
	pdf.CellFormat(colEkskul, 5, "Ekstrakurikuler", "", 0, "C", false, 0, "")

	kelasLabels := []string{"Kelas: 7 (Tujuh)", "Kelas: 8 (Delapan)", "Kelas: 9 (Sembilan)"}
	x = marginL + colEkskul
	for _, kl := range kelasLabels {
		pdf.Rect(x, y, kelasGroupWidth, 5, "DF")
		pdf.SetXY(x, y)
		pdf.CellFormat(kelasGroupWidth, 5, kl, "", 0, "C", false, 0, "")
		x += kelasGroupWidth
	}

	x = marginL + colEkskul
	for i := 0; i < 3; i++ {
		pdf.Rect(x, y+5, semWidth, 5, "DF")
		pdf.SetXY(x, y+5)
		pdf.CellFormat(semWidth, 5, "Semester Ganjil", "", 0, "C", false, 0, "")
		pdf.Rect(x+semWidth, y+5, semWidth, 5, "DF")
		pdf.SetXY(x+semWidth, y+5)
		pdf.CellFormat(semWidth, 5, "Semester Genap", "", 0, "C", false, 0, "")
		x += kelasGroupWidth
	}

	x = marginL + colEkskul
	pdf.SetFont("Arial", "B", 6)
	for i := 0; i < 6; i++ {
		pdf.Rect(x, y+10, nilaiW, 5, "DF")
		pdf.SetXY(x, y+10)
		pdf.CellFormat(nilaiW, 5, "Nilai", "", 0, "C", false, 0, "")
		pdf.Rect(x+nilaiW, y+10, ketW, 5, "DF")
		pdf.SetXY(x+nilaiW, y+10)
		pdf.CellFormat(ketW, 5, "Keterangan", "", 0, "C", false, 0, "")
		x += semWidth
	}

	dataY = y + 15

	// Baris = ekskul yang diikuti sekarang (EkskulList) + ekskul yang pernah diikuti menurut
	// keanggotaan per periode (ekskul_anggota, walau sudah bukan anggota lagi sekarang) + ekskul
	// lain yang punya nilai tercatat. Nilai/Keterangan per kelas & semester diambil dari nilai_ekskul.
	riwayat, err := nilai.GetRiwayatEkskul(siswa.NIS)
	if err != nil {
		log.Printf("buku induk: gagal membaca nilai ekskul nis %s: %v", siswa.NIS, err)
		riwayat = nil
	}

	namaPeriode, err := nilai.GetNamaEkskulPeriode(siswa.NIS)
	if err != nil {
		log.Printf("buku induk: gagal membaca keanggotaan ekskul per periode nis %s: %v", siswa.NIS, err)
		namaPeriode = nil
	}

	namaEkskul := append([]string{}, siswa.EkskulList...)
	namaEkskul = append(namaEkskul, namaPeriode...)

	baris := []nilai.EkskulRiwayat{}
	sudah := map[string]bool{}
	for _, nama := range namaEkskul {
		if sudah[nama] {
			continue
		}

		r := nilai.EkskulRiwayat{Nama: nama}
		for _, rw := range riwayat {
			if rw.Nama == nama {
				r.Sel = rw.Sel
				break
			}
		}

		baris = append(baris, r)
		sudah[nama] = true
	}

	for _, rw := range riwayat {
		if !sudah[rw.Nama] {
			baris = append(baris, rw)
		}
	}

	if len(baris) == 0 {
		baris = []nilai.EkskulRiwayat{{}} // minimal 1 baris kosong siap isi tangan
	}

	cellW := [2]float64{nilaiW, ketW}
	pdf.SetAutoPageBreak(false, 0) // page-break dikerjakan manual oleh ensureRow
	for _, r := range baris {
		// tinggi baris = sel dengan teks terpanjang
		pdf.SetFont("Arial", "", 8)
		maxH := float64(jumlahBaris(pdf, r.Nama, colEkskul)) * 3.6

		pdf.SetFont("Arial", "", 6.5)
		for tingkat := 7; tingkat <= 9; tingkat++ {
			for _, sem := range semesterUrut {
				sel := r.Sel[tingkat][sem]
				if h := float64(jumlahBaris(pdf, sel.Nilai, cellW[0])) * 3.0; h > maxH {
					maxH = h
				}
				if h := float64(jumlahBaris(pdf, sel.Keterangan, cellW[1])) * 3.0; h > maxH {
					maxH = h
				}
			}
		}

		rowH := maxH + 1.6
		if rowH < 7 {
			rowH = 7
		}

		dataY = ensureRow(pdf, dataY, rowH)

		pdf.Rect(marginL, dataY, colEkskul, rowH, "D")
		pdf.SetFont("Arial", "", 8)
		drawTeksTengahV(pdf, marginL, dataY, colEkskul, rowH, 3.6, r.Nama, "L")

		x = marginL + colEkskul
		pdf.SetFont("Arial", "", 6.5)
		for tingkat := 7; tingkat <= 9; tingkat++ {
			for _, sem := range semesterUrut {
				sel := r.Sel[tingkat][sem]

				pdf.Rect(x, dataY, nilaiW, rowH, "D")
				drawTeksTengahV(pdf, x, dataY, nilaiW, rowH, 3.0, sel.Nilai, "C")
				x += nilaiW

				pdf.Rect(x, dataY, ketW, rowH, "D")
				drawTeksTengahV(pdf, x, dataY, ketW, rowH, 3.0, sel.Keterangan, "C")
				x += ketW
			}
		}

		dataY += rowH
	}
	pdf.SetAutoPageBreak(true, marginB)
	pdf.SetXY(marginL, dataY)
}

func tempatTanggalLahir(s BukuIndukData) string {
	tempat := strings.TrimSpace(s.TempatLahir)
	tanggal := strings.TrimSpace(s.TanggalLahir)
	switch {
	case tempat != "" && tanggal != "":
		return tempat + ", " + tanggal
	case tempat != "":
		return tempat
	default:
		return tanggal
	}
}

func tempatTanggal(tempat, tanggal string) string {
	tempat = strings.TrimSpace(tempat)
	tanggal = strings.TrimSpace(tanggal)
	switch {
	case tempat != "" && tanggal != "":
		return tempat + ", " + tanggal
	case tempat != "":
		return tempat
	default:
		return tanggal
	}
}

func alamatTelp(alamat, telp string) string {
	switch {
	case alamat != "" && telp != "":
		return alamat + " / " + telp
	case alamat != "":
		return alamat
	default:
		return telp
	}
}

func jumlahSaudaraLabel(v string) string {
	if v == "" {
		return ""
	}
	return v + " orang"
}

func jarakLabel(v string) string {
	if v == "" {
		return ""
	}
	return v + " kilometer"
}

func jarakKmLabel(v string) string {
	if v == "" {
		return ""
	}
	return v + " km"
}

func tinggiBeratLabel(s BukuIndukData) string {
	if s.TinggiBadan == "" && s.BeratBadan == "" {
		return ""
	}
	return s.TinggiBadan + " cm / " + s.BeratBadan + " kg"
}

func alamatDenganTelp(s BukuIndukData) string {
	switch {
	case s.Alamat != "" && s.NoTelpHp != "":
		return s.Alamat + " Telp/HP: " + s.NoTelpHp
	case s.Alamat != "":
		return s.Alamat
	default:
		return s.NoTelpHp
	}
}

func tanggalNomorIjazah(s BukuIndukData) string {
	switch {
	case s.TanggalIjazahMasuk != "" && s.NomorIjazahMasuk != "":
		return s.TanggalIjazahMasuk + " / No. " + s.NomorIjazahMasuk
	case s.NomorIjazahMasuk != "":
		return "No. " + s.NomorIjazahMasuk
	default:
		return s.TanggalIjazahMasuk
	}
}

func renderBagianD(pdf *fpdf.Fpdf, siswa BukuIndukData) {
	formHeaderItem(pdf, "17.", "Tamat belajar / lulus")
	if siswa.StatusKeluar == StatusKeluarLulus {
		formSubItem(pdf, "a.", "Tahun", siswa.LulusTahun, formLabelWidth)
		formSubItem(pdf, "b.", "Nomor ijazah", siswa.LulusNomorIjazah, formLabelWidth)
		formSubItem(pdf, "c.", "Melanjutkan ke", siswa.LulusMelanjutkanKe, formLabelWidth)
	} else {
		formSubItem(pdf, "a.", "Tahun", "", formLabelWidth)
		formSubItem(pdf, "b.", "Nomor ijazah", "", formLabelWidth)
		formSubItem(pdf, "c.", "Melanjutkan ke", "", formLabelWidth)
	}

	formHeaderItem(pdf, "18.", "Pindah sekolah")
	if siswa.StatusKeluar == StatusKeluarPindah {
		formSubItem(pdf, "a.", "Tingkat/kelas", siswa.PindahTingkatKelas, formLabelWidth)
		formSubItem(pdf, "b.", "Ke sekolah", siswa.PindahKeSekolah, formLabelWidth)
		formSubItem(pdf, "c.", "Ke tingkat/kelas", siswa.PindahKeTingkatKelas, formLabelWidth)
	} else {
		formSubItem(pdf, "a.", "Tingkat/kelas", "", formLabelWidth)
		formSubItem(pdf, "b.", "Ke sekolah", "", formLabelWidth)
		formSubItem(pdf, "c.", "Ke tingkat/kelas", "", formLabelWidth)
	}

	formHeaderItem(pdf, "19.", "Keluar sekolah")
	if siswa.StatusKeluar == StatusKeluarKeluar {
		formSubItem(pdf, "a.", "Alasan keluar sekolah", siswa.KeluarAlasan, formLabelWidth)
		formSubItem(pdf, "b.", "Pada tanggal", siswa.KeluarTanggal, formLabelWidth)
	} else {
		formSubItem(pdf, "a.", "Alasan keluar sekolah", "", formLabelWidth)
		formSubItem(pdf, "b.", "Pada tanggal", "", formLabelWidth)
	}
}

func renderFotoBox(pdf *fpdf.Fpdf, fotoURL string, jenisKelamin string, x, y, w, h float64) {
	var physicalPath string
	if fotoURL != "" {
		physicalPath = "." + fotoURL
	} else {
		// kolom foto kosong/NULL -> pakai foto default sesuai jenis
		// kelamin dari storage/pesertadidik/L.jpg atau P.jpg
		defaultFile := "P.jpg"
		if jenisKelaminSingkat(jenisKelamin) == "L" {
			defaultFile = "L.jpg"
		}
		physicalPath = "./storage/fotosiswa/" + defaultFile
	}

	data, err := os.ReadFile(physicalPath)
	if err != nil {
		pdf.SetFillColor(235, 235, 235)
		pdf.Rect(x, y, w, h, "F")
		pdf.SetFont("Arial", "", 8)
		pdf.SetXY(x, y+h/2-2)
		pdf.CellFormat(w, 4, "Pas Foto", "", 0, "C", false, 0, "")
		return
	}

	ext := strings.ToUpper(strings.TrimPrefix(filepath.Ext(physicalPath), "."))
	if ext == "JPG" {
		ext = "JPEG"
	}
	imgName := "foto-" + filepath.Base(physicalPath)
	pdf.RegisterImageOptionsReader(imgName, fpdf.ImageOptions{ImageType: ext}, bytes.NewReader(data))
	pdf.ImageOptions(imgName, x, y, w, h, false, fpdf.ImageOptions{ImageType: ext}, 0, "")
}
