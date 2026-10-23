package bukuinduk

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/go-pdf/fpdf"

	"sekolah/internal/modules/nilai"
)

type mapelItem struct {
	No      string
	Nama    string
	SubList []string
}

var daftarMapel = []mapelItem{
	{No: "1", Nama: "Pendidikan Agama dan Budi Pekerti"},
	{No: "2", Nama: "Pendidikan Pancasila"},
	{No: "3", Nama: "Bahasa Indonesia"},
	{No: "4", Nama: "Matematika"},
	{No: "5", Nama: "Ilmu Pengetahuan Alam"},
	{No: "6", Nama: "Ilmu Pengetahuan Sosial"},
	{No: "7", Nama: "Bahasa Inggris"},
	{No: "8", Nama: "PJOK"},
	{No: "9", Nama: "Teknologi Informasi & Komunikasi"},
	{No: "10", Nama: "Seni Budaya & Prakarya", SubList: []string{"Seni", "Seni", "Seni", "Seni", "Prakarya (Budidaya, Kerajinan, Rekayasa, Pengolahan)"}},
	{No: "11", Nama: "Proyek Penguatan Profil Pelajar Pancasila (P5)", SubList: []string{"Proyek", "Proyek", "Proyek", "Proyek"}},
	{No: "12", Nama: "Bimbingan Konseling"},
	{No: "MULOK", Nama: "Muatan Lokal", SubList: []string{"", "", ""}},
}

var dimensiPancasila = []string{
	"Beriman, bertakwa kepada Tuhan Yang Maha Esa dan berakhlak mulia",
	"Berkebinekaan Global",
	"Bergotong-royong",
	"Mandiri",
	"Bernalar Kritis",
	"Kreatif",
}

// renderDataNilaiPerKelas mencetak halaman Data Nilai untuk kelas 7, 8, dan 9 (3 halaman per
// siswa). Tahun pelajaran tiap kelas dicari dari nilai_mapel siswa itu (nilai_mapel.id_kelas
// -> tingkat). Tingkat yang belum punya nilai tetap dicetak (kosong, tahun kosong), kecuali
// tingkat kelas siswa saat ini: dipakai tahun pelajaran yang sedang dicetak, sehingga
// ketidakhadiran semester berjalan tetap tampil walau nilai belum diinput.
//
// tingkatDicetak membatasi tingkat mana yg ikut dicetak (mis. cetak arsip semester berjalan
// yg baru terisi Kelas 7 -> Kelas 8 & 9 dilewati, tidak perlu cetak halaman kosong). nil/kosong
// = cetak semua (7, 8, 9), kompatibel dgn perilaku lama.
func renderDataNilaiPerKelas(pdf *fpdf.Fpdf, siswa BukuIndukData, tahunCetak string, kepsek KepalaSekolahInfo, tingkatDicetak map[int]bool) {
	tahunPerTingkat, err := nilai.GetTahunPerTingkat(siswa.NIS)
	if err != nil {
		log.Printf("buku induk: gagal membaca tahun nilai nis %s: %v", siswa.NIS, err)
		tahunPerTingkat = nil
	}

	tingkatSekarang := nilai.TingkatKelas(siswa.Kelas)

	// wali kelas per tingkat (menu Buku Induk > Wali Kelas). Belum diisi -> kosong (ttd manual),
	// TIDAK jatuh ke wali kelas rombel siswa saat ini.
	waliPerTingkat, errWali := GetWaliKelasPerTingkat(siswa.NIS)
	if errWali != nil {
		log.Printf("buku induk: gagal membaca wali kelas per tingkat nis %s: %v", siswa.NIS, errWali)
		waliPerTingkat = nil
	}

	for tingkat := 7; tingkat <= 9; tingkat++ {
		if len(tingkatDicetak) > 0 && !tingkatDicetak[tingkat] {
			continue
		}

		tahun := tahunPerTingkat[tingkat]
		if tahun == "" && tingkat == tingkatSekarang {
			tahun = tahunCetak
		}

		renderDataNilaiPage1(pdf, siswa, tahun, tingkat, kepsek, waliPerTingkat[tingkat])
	}
}

// renderDataNilaiPage1 — halaman "A. Projek Penguatan Profil Pelajar
// Pancasila" + "B. Pengetahuan dan Keterampilan" + Ketidakhadiran.
// Bagian B (nilai per mapel per semester) dan Ketidakhadiran diisi dari
// nilai_mapel & absensi (nilai.GetDataBukuInduk). Bagian A (P5) masih kosong.
func renderDataNilaiPage1(pdf *fpdf.Fpdf, siswa BukuIndukData, tahunPelajaran string, tingkat int, kepsek KepalaSekolahInfo, wali WaliKelasTingkatInfo) {
	// nilai + ketidakhadiran diambil dari nilai_mapel & absensi utk tahun pelajaran ini.
	// Gagal membaca tidak menggagalkan cetak: halaman tetap tercetak kosong.
	data, err := nilai.GetDataBukuInduk(siswa.NIS, tahunPelajaran)
	if err != nil {
		log.Printf("buku induk: gagal membaca nilai nis %s (%s): %v", siswa.NIS, tahunPelajaran, err)
		data = nil
	}

	// nama kelas di header tabel Ketidakhadiran: dari nilai_mapel kalau ada, kalau tidak angka tingkatnya
	if data != nil && data.NamaKelas == "" && tingkat > 0 {
		data.NamaKelas = strconv.Itoa(tingkat)
	}

	// isian P5 (tabel A) dari p5_nilai; gagal membaca tidak menggagalkan cetak
	p5, err := nilai.GetRiwayatP5(siswa.NIS)
	if err != nil {
		log.Printf("buku induk: gagal membaca P5 nis %s: %v", siswa.NIS, err)
		p5 = nil
	}

	pdf.AddPage()
	centerText(pdf, "IV. DATA NILAI KURIKULUM MERDEKA PESERTA DIDIK SISWA / SISWI SMP", 11, true)
	pdf.Ln(4)

	renderDataNilaiHeader(pdf, siswa, tahunPelajaran)

	pdf.Ln(3)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 6, "A. PROJEK PENGUATAN PROFIL PELAJAR PANCASILA", "", 1, "L", false, 0, "")
	pdf.Ln(2)
	renderPancasilaLegend(pdf)
	renderPancasilaMerdekaTable(pdf, p5)

	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 6, "B. PENGETAHUAN DAN KETERAMPILAN", "", 1, "L", false, 0, "")
	renderMapelNilaiTable(pdf, data)

	pdf.Ln(4)
	renderKetidakhadiranMerdekaTable(pdf, data, kepsek, wali.Nama, wali.NIP, siswa.NIS, tingkat)
}

// p5Teks = teks Elemen / Sub Elemen / Target pertama yang terisi dari isian semester Ganjil lalu Genap
// di tingkat itu untuk 1 dimensi.
func p5Teks(p5 nilai.P5Riwayat, tingkat, dimensi int, ambil func(nilai.P5Sel) string) string {
	for _, sem := range semesterUrut {
		if t := ambil(p5[tingkat][dimensi][sem]); t != "" {
			return t
		}
	}

	return ""
}

// renderPancasilaMerdekaTable = tabel A (Fase D): kolom kiri Elemen + nilai kelas 7, kolom kanan
// Sub Elemen + nilai kelas 8, Target Pencapaian Di Akhir Fase. p5 boleh nil -> tabel kosong.
func renderPancasilaMerdekaTable(pdf *fpdf.Fpdf, p5 nilai.P5Riwayat) {
	fullWidth := pageWidth - marginL - marginR
	wDimensi, wElemen, wNilai, wSub := 35.0, 25.0, 12.0, 35.0
	wTarget := fullWidth - wDimensi - wElemen - wNilai*4 - wSub

	headerH1, headerH2 := 5.0, 7.0
	y := pdf.GetY()
	pdf.SetFont("Arial", "B", 8)

	drawHeaderCell := func(x, w float64, label string) {
		pdf.SetFillColor(224, 224, 224)
		pdf.Rect(x, y, w, headerH1+headerH2, "DF")
		pdf.SetFont("Arial", "B", 8)
		lines := pdf.SplitLines([]byte(label), w-1.5)
		nLines := len(lines)
		if nLines < 1 {
			nLines = 1
		}
		lineH := 4.0
		textH := float64(nLines) * lineH
		startY := y + (headerH1+headerH2-textH)/2
		pdf.SetXY(x, startY)
		pdf.MultiCell(w, lineH, label, "", "C", false)
	}

	x := marginL
	drawHeaderCell(x, wDimensi, "Dimensi")
	x += wDimensi
	drawHeaderCell(x, wElemen, "Elemen")
	x += wElemen

	drawKelasGroup := func(x float64, namaKelas string) {
		pdf.SetFillColor(224, 224, 224)
		pdf.SetFont("Arial", "B", 8)
		pdf.Rect(x, y, wNilai*2, headerH1, "DF")
		pdf.SetXY(x, y)
		pdf.CellFormat(wNilai*2, headerH1, namaKelas, "", 0, "C", false, 0, "")

		pdf.SetFont("Arial", "B", 6)
		pdf.Rect(x, y+headerH1, wNilai, headerH2, "DF")
		pdf.SetXY(x, y+headerH1)
		pdf.MultiCell(wNilai, headerH2/2, "Nilai\nSem.1", "", "C", false)
		pdf.Rect(x+wNilai, y+headerH1, wNilai, headerH2, "DF")
		pdf.SetXY(x+wNilai, y+headerH1)
		pdf.MultiCell(wNilai, headerH2/2, "Nilai\nSem.2", "", "C", false)
	}
	drawKelasGroup(x, "Kelas 7")
	x += wNilai * 2

	drawHeaderCell(x, wSub, "Sub Elemen")
	x += wSub
	drawKelasGroup(x, "Kelas 8")
	x += wNilai * 2

	drawHeaderCell(x, wTarget, "Target Pencapaian Di Akhir Fase")

	pdf.SetY(y + headerH1 + headerH2)

	pdf.SetAutoPageBreak(false, 0) // Rect() tidak dipantau auto-pagebreak; dijaga manual oleh ensureRow

	for i, dim := range dimensiPancasila {
		d := i + 1
		elemen := p5Teks(p5, 7, d, func(s nilai.P5Sel) string { return s.Elemen })
		sub := p5Teks(p5, 8, d, func(s nilai.P5Sel) string { return s.SubElemen })
		target := p5Teks(p5, 8, d, func(s nilai.P5Sel) string { return s.Target })

		if target == "" {
			target = p5Teks(p5, 7, d, func(s nilai.P5Sel) string { return s.Target })
		}

		if sub == "" {
			sub = p5Teks(p5, 7, d, func(s nilai.P5Sel) string { return s.SubElemen })
		}

		// tinggi baris = kolom yang paling banyak barisnya (Dimensi font 8, lainnya font 6.5)
		pdf.SetFont("Arial", "", 8)
		hDim := float64(jumlahBaris(pdf, fmt.Sprintf("%d. %s", d, dim), wDimensi)) * 4.0

		pdf.SetFont("Arial", "", 6.5)
		hLain := 0.0
		for _, kv := range []struct {
			t string
			w float64
		}{{elemen, wElemen}, {sub, wSub}, {target, wTarget}} {
			if h := float64(jumlahBaris(pdf, kv.t, kv.w)) * 3.0; h > hLain {
				hLain = h
			}
		}

		rowH := hDim
		if hLain > rowH {
			rowH = hLain
		}

		rowH += 2 * nilaiPadY
		if rowH < 7 {
			rowH = 7
		}

		rowY := ensureRow(pdf, pdf.GetY(), rowH)
		x = marginL

		// kolom: Dimensi | Elemen | K7 Sem1 | K7 Sem2 | Sub Elemen | K8 Sem1 | K8 Sem2 | Target
		lebar := []float64{wDimensi, wElemen, wNilai, wNilai, wSub, wNilai, wNilai, wTarget}
		isi := []string{
			fmt.Sprintf("%d. %s", d, dim), elemen,
			p5[7][d]["Ganjil"].Nilai, p5[7][d]["Genap"].Nilai,
			sub,
			p5[8][d]["Ganjil"].Nilai, p5[8][d]["Genap"].Nilai,
			target,
		}

		for ci, w := range lebar {
			pdf.Rect(x, rowY, w, rowH, "D")

			switch ci {
			case 0:
				pdf.SetFont("Arial", "", 8)
				drawTeksTengahV(pdf, x, rowY, w, rowH, 4.0, isi[ci], "L")
			case 2, 3, 5, 6: // nilai BB/MB/BSH/SB: di tengah sel
				pdf.SetFont("Arial", "", 8)
				pdf.SetXY(x, rowY)
				pdf.CellFormat(w, rowH, isi[ci], "", 0, "C", false, 0, "")
			default:
				pdf.SetFont("Arial", "", 6.5)
				drawTeksTengahV(pdf, x, rowY, w, rowH, 3.0, isi[ci], "L")
			}

			x += w
		}

		pdf.SetXY(marginL, rowY+rowH)
	}

	pdf.SetAutoPageBreak(true, marginB)
}

// drawMapelDataRow menggambar satu baris tabel mata-pelajaran (dipakai
// renderMapelNilaiTable & renderRataRataUjianTable). Border digambar
// manual via Rect() (bukan MultiCell border="1"). PENTING: cek ruang
// tersisa & loncat halaman DULU sebelum menggambar apapun pada baris
// ini -- Rect() tidak dipantau auto-pagebreak fpdf (beda dari
// MultiCell), jadi kalau tidak dicek manual di sini, border & teks
// bisa kepisah ke 2 halaman berbeda.
// drawMapelDataRow menggambar baris mapel biasa (No terpisah, Mapel
// terpisah, 6 kotak nilai terpisah — TIDAK di-merge). Dipakai utk
// baris yg tidak punya sub-item.
func drawMapelDataRow(pdf *fpdf.Fpdf, dataY, colNo, colMapel, subWidth float64, no, nama string) float64 {
	nLines := len(pdf.SplitLines([]byte(nama), colMapel-1.5))
	if nLines < 1 {
		nLines = 1
	}
	lineH := 3.5
	rowH := float64(nLines) * lineH
	if rowH < 6 {
		rowH = 6
	}

	if dataY+rowH > pageHeight-marginB {
		pdf.AddPage()
		dataY = marginT
	}

	pdf.Rect(marginL, dataY, colNo, rowH, "D")
	pdf.SetXY(marginL, dataY)
	pdf.CellFormat(colNo, rowH, no, "", 0, "C", false, 0, "")

	pdf.Rect(marginL+colNo, dataY, colMapel, rowH, "D")
	textY := dataY + (rowH-float64(nLines)*lineH)/2
	pdf.SetXY(marginL+colNo, textY)
	pdf.MultiCell(colMapel, lineH, nama, "", "L", false)

	x := marginL + colNo + colMapel
	for i := 0; i < 6; i++ {
		pdf.Rect(x, dataY, subWidth, rowH, "D")
		x += subWidth
	}
	return dataY + rowH
}

// drawMapelParentRow menggambar baris INDUK (mis. "10 Seni Budaya &
// Prakarya") sbg SATU sel utuh merentang penuh No+Mapel+semua kolom
// nilai — tanpa kotak data sama sekali (data ada di baris sub-item
// di bawahnya, bukan di baris induk ini).
func drawMapelParentRow(pdf *fpdf.Fpdf, dataY, colNo, colMapel, subWidth float64, no, nama string) float64 {
	mergedWidth := colMapel + subWidth*6
	nLines := len(pdf.SplitLines([]byte(nama), mergedWidth-1.5))
	if nLines < 1 {
		nLines = 1
	}
	lineH := 3.5
	rowH := float64(nLines) * lineH
	if rowH < 6 {
		rowH = 6
	}

	if dataY+rowH > pageHeight-marginB {
		pdf.AddPage()
		dataY = marginT
	}

	pdf.SetFillColor(224, 224, 224)
	pdf.Rect(marginL, dataY, colNo, rowH, "DF")
	pdf.SetXY(marginL, dataY)
	pdf.CellFormat(colNo, rowH, no, "", 0, "C", false, 0, "")

	pdf.Rect(marginL+colNo, dataY, mergedWidth, rowH, "DF")
	textY := dataY + (rowH-float64(nLines)*lineH)/2
	pdf.SetXY(marginL+colNo, textY)
	pdf.MultiCell(mergedWidth, lineH, nama, "", "L", false)
	return dataY + rowH
}

// drawMapelSubRow menggambar baris sub-item (a./b./c./dst) — kolom No
// kosong, label diberi awalan huruf di kolom Mapel, dan 6 kotak nilai
// TETAP ada (tidak di-merge, siap diisi tangan).
func drawMapelSubRow(pdf *fpdf.Fpdf, dataY, colNo, colMapel, subWidth float64, letter, nama string) float64 {
	label := letter + ". " + nama
	nLines := len(pdf.SplitLines([]byte(label), colMapel-1.5))
	if nLines < 1 {
		nLines = 1
	}
	lineH := 3.5
	rowH := float64(nLines) * lineH
	if rowH < 6 {
		rowH = 6
	}

	if dataY+rowH > pageHeight-marginB {
		pdf.AddPage()
		dataY = marginT
	}

	pdf.Rect(marginL, dataY, colNo, rowH, "D")
	pdf.Rect(marginL+colNo, dataY, colMapel, rowH, "D")
	textY := dataY + (rowH-float64(nLines)*lineH)/2
	pdf.SetXY(marginL+colNo, textY)
	pdf.MultiCell(colMapel, lineH, label, "", "L", false)

	x := marginL + colNo + colMapel
	for i := 0; i < 6; i++ {
		pdf.Rect(x, dataY, subWidth, rowH, "D")
		x += subWidth
	}
	return dataY + rowH
}

// drawMapelRows menggambar seluruh baris daftarMapel: baris tanpa
// sub-item pakai drawMapelDataRow, baris berpunya sub-item pakai
// drawMapelParentRow (merge penuh) + drawMapelSubRow per sub-item
// (huruf a/b/c/dst otomatis).
func drawMapelRows(pdf *fpdf.Fpdf, dataY, colNo, colMapel, subWidth float64) float64 {
	letters := "abcdefgh"
	for _, m := range daftarMapel {
		if m.No == "MULOK" {
			dataY = drawMulokRows(pdf, dataY, colNo, colMapel, subWidth, m.SubList)
			continue
		}
		if len(m.SubList) > 0 {
			dataY = drawMapelParentRow(pdf, dataY, colNo, colMapel, subWidth, m.No, m.Nama)
			for i, sub := range m.SubList {
				letter := "?"
				if i < len(letters) {
					letter = string(letters[i])
				}
				dataY = drawMapelSubRow(pdf, dataY, colNo, colMapel, subWidth, letter, sub)
			}
		} else {
			dataY = drawMapelDataRow(pdf, dataY, colNo, colMapel, subWidth, m.No, m.Nama)
		}
	}
	return dataY
}

// drawMulokRows menggambar baris MULOK: kolom No berisi huruf
// "M/U/L/O/K" ditumpuk vertikal (1 huruf per baris) merentang penuh
// tinggi seluruh sub-item a/b/c — TANPA baris header "MULOK" terpisah
// spt baris induk lain, karena kolom No di sini LANGSUNG jadi label.
func drawMulokRows(pdf *fpdf.Fpdf, dataY, colNo, colMapel, subWidth float64, subList []string) float64 {
	rowH := 6.0
	n := len(subList)
	if n == 0 {
		n = 1
	}
	totalH := rowH * float64(n)

	if dataY+totalH > pageHeight-marginB {
		pdf.AddPage()
		dataY = marginT
	}

	pdf.Rect(marginL, dataY, colNo, totalH, "D")
	mulokLetters := []rune("MULOK")
	letterH := totalH / float64(len(mulokLetters))
	pdf.SetFont("Arial", "B", 7)
	for i, r := range mulokLetters {
		pdf.SetXY(marginL, dataY+float64(i)*letterH)
		pdf.CellFormat(colNo, letterH, string(r), "", 0, "C", false, 0, "")
	}

	letters := "abc"
	pdf.SetFont("Arial", "", 7)
	for i := 0; i < n; i++ {
		rY := dataY + float64(i)*rowH
		letter := "?"
		if i < len(letters) {
			letter = string(letters[i])
		}
		label := letter + "."
		if i < len(subList) && subList[i] != "" {
			label = letter + ". " + subList[i]
		}
		pdf.Rect(marginL+colNo, rY, colMapel, rowH, "D")
		pdf.SetXY(marginL+colNo, rY)
		pdf.CellFormat(colMapel, rowH, label, "", 0, "L", false, 0, "")
		x := marginL + colNo + colMapel
		for j := 0; j < 6; j++ {
			pdf.Rect(x, rY, subWidth, rowH, "D")
			x += subWidth
		}
	}
	return dataY + totalH
}

func renderKetidakhadiranMerdekaTable(pdf *fpdf.Fpdf, data *nilai.DataBukuInduk, kepsek KepalaSekolahInfo, waliKelasNama, waliKelasNIP, nis string, tingkat int) {
	fullWidth := pageWidth - marginL - marginR
	colNo, colKet, colSem, colKenaikan, colSign := 8.0, 26.0, 18.0, 6.0, 36.0
	colKeputusan := fullWidth - colNo - colKet - colSem*2 - colKenaikan - colSign*2

	headerH1, headerH2 := 5.0, 5.0
	rowH := 6.0
	totalTableH := headerH1 + headerH2 + rowH*3

	y := pdf.GetY()
	if y+totalTableH > pageHeight-marginB {
		pdf.AddPage()
		y = marginT
	}
	x := marginL

	pdf.SetFont("Arial", "B", 7)
	pdf.SetFillColor(224, 224, 224)
	pdf.Rect(x, y, colNo+colKet, headerH1+headerH2, "DF")
	pdf.SetXY(x, y+3)
	pdf.CellFormat(colNo+colKet, 4, "KETIDAKHADIRAN", "", 0, "C", false, 0, "")
	x += colNo + colKet

	pdf.Rect(x, y, colSem*2, headerH1, "D")
	pdf.SetXY(x, y)
	pdf.CellFormat(colSem*2, headerH1, labelKelas(data), "", 0, "C", false, 0, "")
	pdf.Rect(x, y+headerH1, colSem, headerH2, "D")
	pdf.SetXY(x, y+headerH1)
	pdf.CellFormat(colSem, headerH2, "SEMESTER 1", "", 0, "C", false, 0, "")
	pdf.Rect(x+colSem, y+headerH1, colSem, headerH2, "D")
	pdf.SetXY(x+colSem, y+headerH1)
	pdf.CellFormat(colSem, headerH2, "SEMESTER 2", "", 0, "C", false, 0, "")
	x += colSem * 2

	kenaikanX := x
	x += colKenaikan

	pdf.Rect(x, y, colKeputusan, headerH1+headerH2, "DF")
	pdf.SetFillColor(224, 224, 224)
	pdf.SetFont("Arial", "B", 6.5)
	pdf.SetXY(x+1, y+1)
	pdf.MultiCell(colKeputusan-2, 3.2, "KEPUTUSAN HASIL SEMESTER 1 DAN 2", "", "C", false)
	x += colKeputusan

	waliX := x
	x += colSign

	kepsekX := x

	dataY := y + headerH1 + headerH2
	labels := []string{"Sakit", "Izin", "Tanpa Keterangan"}

	// label "NAIK KE KELAS"/"TINGGAL DI KELAS" dipisah dari tanda kurungnya & diberi
	// lebar tetap (spt TAB di Word), supaya "(   )" kedua baris sejajar lurus - pola sama
	// dgn LULUS/TIDAK LULUS di kotak Keputusan (Halaman 2).
	keputusanTabW := pdf.GetStringWidth("TINGGAL DI KELAS") + 5
	keputusanLabels := [2]string{"NAIK KE KELAS", "TINGGAL DI KELAS"}
	keputusanTanggalLine := "Tanggal: ......../......../........"

	// tanda "(   )" per baris: utk tingkat 7 & 8 diisi OTOMATIS dgn nomor kelas tujuan
	// (Naik -> kelas berikutnya) / kelas sekarang (Tinggal -> kelas ini lagi), dari menu
	// Keputusan Naik Kelas. Tingkat 9 dibiarkan kosong spt sebelumnya (keputusan akhirnya
	// LULUS/TIDAK LULUS di kotak Keputusan halaman 2, bukan di tabel ini).
	keputusanTanda := [2]string{"(   )", "(   )"}
	if tingkat == 7 || tingkat == 8 {
		status, errNaik := nilai.GetKeputusanNaikKelas(nis, tingkat)
		if errNaik != nil {
			log.Printf("buku induk: gagal membaca keputusan naik kelas nis %s tingkat %d: %v", nis, tingkat, errNaik)
			status = ""
		}

		if status == nilai.StatusNaik {
			keputusanTanda[0] = fmt.Sprintf("( %d )", tingkat+1)
		}

		if status == nilai.StatusTinggal {
			keputusanTanda[1] = fmt.Sprintf("( %d )", tingkat)
		}
	}

	pdf.SetFont("Arial", "", 8)
	for i, label := range labels {
		rY := dataY + float64(i)*rowH
		x = marginL
		pdf.Rect(x, rY, colNo, rowH, "D")
		pdf.SetXY(x, rY)
		pdf.CellFormat(colNo, rowH, fmt.Sprintf("%d", i+1), "", 0, "C", false, 0, "")
		x += colNo

		pdf.Rect(x, rY, colKet, rowH, "D")
		pdf.SetXY(x, rY)
		pdf.CellFormat(colKet, rowH, label, "", 0, "L", false, 0, "")
		x += colKet

		pdf.Rect(x, rY, colSem, rowH, "D")
		pdf.SetXY(x, rY)
		pdf.CellFormat(colSem, rowH, hariTeks(data, "Ganjil", i), "", 0, "C", false, 0, "")
		x += colSem

		pdf.Rect(x, rY, colSem, rowH, "D")
		pdf.SetXY(x, rY)
		pdf.CellFormat(colSem, rowH, hariTeks(data, "Genap", i), "", 0, "C", false, 0, "")
		x += colSem
		x += colKenaikan // kolom KENAIKAN digambar terpisah, merentang penuh header+data

		pdf.Rect(x, rY, colKeputusan, rowH, "D")
		if i < 2 {
			pdf.SetXY(x+1, rY)
			pdf.CellFormat(keputusanTabW, rowH, keputusanLabels[i], "", 0, "L", false, 0, "")
			pdf.CellFormat(colKeputusan-2-keputusanTabW, rowH, keputusanTanda[i], "", 0, "L", false, 0, "")
		} else {
			pdf.SetXY(x+1, rY)
			pdf.CellFormat(colKeputusan-2, rowH, keputusanTanggalLine, "", 0, "L", false, 0, "")
		}
	}

	// kolom "KENAIKAN" — 1 kotak merentang PENUH dari header s/d baris
	// terakhir (sejajar "KELAS ...", bukan cuma area data), huruf
	// K-E-N-A-I-K-A-N ditumpuk vertikal
	pdf.SetFillColor(224, 224, 224)
	pdf.Rect(kenaikanX, y, colKenaikan, totalTableH, "DF")
	kenaikanLetters := []rune("KENAIKAN")
	letterH := totalTableH / float64(len(kenaikanLetters))
	pdf.SetFont("Arial", "B", 6)
	for i, r := range kenaikanLetters {
		pdf.SetXY(kenaikanX, y+float64(i)*letterH)
		pdf.CellFormat(colKenaikan, letterH, string(r), "", 0, "C", false, 0, "")
	}

	// kolom Wali Kelas & Kepala Sekolah — 1 kotak merentang PENUH dari
	// header s/d baris terakhir; judul "Mengetahui,.." di atas, nama+NIP
	// langsung di dalam sel yg sama (bukan lagi teks lepas di bawah tabel)
	drawSignatureColumn := func(colX float64, judul, nama, nip string) {
		pdf.Rect(colX, y, colSign, totalTableH, "D")
		pdf.SetFont("Arial", "", 8)
		pdf.SetXY(colX, y+2)
		pdf.MultiCell(colSign, 3.5, judul, "", "C", false)

		nameY := y + totalTableH - 12
		pdf.SetFont("Arial", "", 7)
		pdf.SetXY(colX+1, nameY)
		pdf.CellFormat(colSign-2, 5, nama, "U", 2, "C", false, 0, "")
		pdf.SetXY(colX+1, nameY+5)
		pdf.CellFormat(colSign-2, 5, nip, "", 0, "C", false, 0, "")
	}

	waliNipText := "NIP."
	if waliKelasNIP != "" {
		waliNipText = formatNIP(waliKelasNIP)
	}
	drawSignatureColumn(waliX, "Mengetahui,\nWali Kelas", waliKelasNama, waliNipText)
	drawSignatureColumn(kepsekX, "Mengetahui,\nKepala Sekolah", kepsek.NamaKepsek, formatNIP(kepsek.NIPKepsek))

	pdf.SetY(y + totalTableH + 4)
	pdf.SetFont("Arial", "B", 7)
	pdf.CellFormat(0, 4, "* KKTP (Kriteria Ketercapaian Tujuan Pembelajaran)", "", 1, "L", false, 0, "")
}

// renderWaliKelasKepsekBlock — blok ttd 2 kolom (Wali Kelas | Kepala
// Sekolah), dipakai khusus di halaman Data Nilai (urutan kolom
// kebalik drpd renderWaliKelasSignatureBlock lama).
func renderWaliKelasKepsekBlock(pdf *fpdf.Fpdf, kepsek KepalaSekolahInfo, waliKelasNama, waliKelasNIP string) {
	pdf.Ln(8)
	colWidth := (pageWidth - marginL - marginR) / 2
	startY := pdf.GetY()
	pdf.SetFont("Arial", "", 10)

	pdf.SetXY(marginL, startY)
	pdf.CellFormat(colWidth, 6, "Mengetahui,", "", 2, "L", false, 0, "")
	pdf.SetX(marginL)
	pdf.CellFormat(colWidth, 6, "Wali Kelas", "", 2, "L", false, 0, "")

	pdf.SetXY(marginL+colWidth, startY)
	pdf.CellFormat(colWidth, 6, "Mengetahui,", "", 2, "L", false, 0, "")
	pdf.SetX(marginL + colWidth)
	pdf.CellFormat(colWidth, 6, "Kepala Sekolah", "", 2, "L", false, 0, "")

	pdf.SetXY(marginL, startY+20)
	pdf.CellFormat(colWidth, 6, waliKelasNama, "U", 2, "L", false, 0, "")
	pdf.SetX(marginL)
	if waliKelasNIP != "" {
		pdf.CellFormat(colWidth, 6, formatNIP(waliKelasNIP), "", 2, "L", false, 0, "")
	} else {
		pdf.CellFormat(colWidth, 6, "NIP.", "", 2, "L", false, 0, "")
	}

	pdf.SetXY(marginL+colWidth, startY+20)
	pdf.CellFormat(colWidth, 6, kepsek.NamaKepsek, "U", 2, "L", false, 0, "")
	pdf.SetX(marginL + colWidth)
	pdf.CellFormat(colWidth, 6, formatNIP(kepsek.NIPKepsek), "", 2, "L", false, 0, "")
}

// renderDataNilaiPage2 — halaman kedua: Nilai Rata-rata Raport / Nilai
// Ujian Sekolah / Nilai Sekolah, lalu ringkasan 4 mapel Ujian Sekolah,
// status Lulus/Tidak Lulus. Seluruhnya kosong (blm ada modul nilai).
func renderDataNilaiPage2(pdf *fpdf.Fpdf, siswa BukuIndukData, tahunPelajaran string, kepsek KepalaSekolahInfo) {
	pdf.AddPage()
	centerText(pdf, "IV. DATA NILAI KURIKULUM MERDEKA PESERTA DIDIK SISWA / SISWI SMP", 11, true)
	pdf.Ln(2)

	renderDataNilaiHeader(pdf, siswa, tahunPelajaran)

	pdf.Ln(1)
	d, err := nilai.GetDataHalaman2(siswa.NIS)
	if err != nil {
		log.Printf("buku induk: gagal membaca nilai rata-rata/ujian nis %s: %v", siswa.NIS, err)
		d = nil
	}

	renderRataRataUjianTable(pdf, d)
	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(0, 5, "Nilai Sekolah = Nilai Rata-rata + Nilai Ujian Sekolah", "", 1, "L", false, 0, "")

	// halaman penutup = akhir Kelas 9, jadi wali kelasnya wali kelas tingkat 9
	waliPerTingkat, errWali := GetWaliKelasPerTingkat(siswa.NIS)
	if errWali != nil {
		log.Printf("buku induk: gagal membaca wali kelas per tingkat nis %s: %v", siswa.NIS, errWali)
	}

	wali9 := waliPerTingkat[9]

	pdf.Ln(1)
	renderEmpatMapelTable(pdf, d, kepsek, wali9.Nama, wali9.NIP)
}

// renderEmpatMapelTable — tabel 4 mapel (Bahasa Indonesia/Matematika/
// IPA/IPS) DIGABUNG dgn kotak "Keputusan" yg merentang penuh dari
// header s/d baris terakhir, lalu diikuti footer 2 kolom: kiri "Wali
// Kelas" (nama+NIP), kanan "STATUS AKHIR"+"LULUS/TIDAK LULUS"+Tanggal
// lalu "Kepala Sekolah" (nama+NIP) di bawahnya.
func renderEmpatMapelTable(pdf *fpdf.Fpdf, d *nilai.DataHalaman2, kepsek KepalaSekolahInfo, waliKelasNama, waliKelasNIP string) {
	fullWidth := pageWidth - marginL - marginR
	colNo, colMapel := 8.0, 32.0
	groupWidth := 28.0 // diperlebar (dari 24) supaya "Nilai Ujian Sekolah Standar Daerah"/"Capaian Pembelajaran" tidak tumpuk
	akhirW := 8.0      // kolom Nilai Akhir sempit, Capaian Pembelajaran lebih lebar
	capW := groupWidth - akhirW
	subW := [6]float64{akhirW, capW, akhirW, capW, akhirW, capW}
	colKeputusan := fullWidth - colNo - colMapel - groupWidth*3

	headerH1, headerH2 := 10.0, 8.0 // diperbesar (dari 6/7) supaya teks 2-3 baris muat penuh
	rowH := 6.0
	dataRows := 4

	// isi 4 mapel: Ujian Standar Daerah (dari input), Nilai Sekolah (dari tabel 1), Nilai Akhir = rata-rata keduanya
	mapelEmpat := []string{"Bahasa Indonesia", "Matematika", "Ilmu Pengetahuan Alam", "Ilmu Pengetahuan Sosial"}
	labelEmpat := []string{"Bahasa Indonesia", "Matematika", "IPA", "IPS"}

	var cellsEmpat [4][6]string

	var kolNilai [3][]float64 // utk baris RATA-RATA

	if d != nil {
		for i, nama := range mapelEmpat {
			id := idMapelByNama(d.Mapel, nama)
			if id == 0 {
				continue
			}

			cellsEmpat[i][0], cellsEmpat[i][1] = fmtUjian(d.Daerah, id)
			cellsEmpat[i][2] = fmtMap(d.Gabung, id)
			cellsEmpat[i][4] = fmtMap(d.Akhir, id)

			if s, ok := d.Daerah[id]; ok && s.Nilai != nil {
				kolNilai[0] = append(kolNilai[0], *s.Nilai)
			}

			if v, ok := d.Gabung[id]; ok {
				kolNilai[1] = append(kolNilai[1], v)
			}

			if v, ok := d.Akhir[id]; ok {
				kolNilai[2] = append(kolNilai[2], v)
			}
		}
	}

	// tinggi tiap baris menyesuaikan Capaian Pembelajaran (teks bisa beberapa baris)
	rowHs := make([]float64, dataRows)
	sumH := 0.0

	for i := 0; i < dataRows; i++ {
		rowHs[i] = rowH

		pdf.SetFont("Arial", "", nilaiFontTeks)

		if cap := cellsEmpat[i][1]; cap != "" {
			if h := float64(jumlahBaris(pdf, cap, capW))*nilaiLineH + nilaiPadY*2; h > rowHs[i] {
				rowHs[i] = h
			}
		}

		sumH += rowHs[i]
	}

	tableH := headerH1 + headerH2 + sumH

	y := pdf.GetY()
	// pengecekan pindah halaman di sini DIHAPUS sesuai permintaan —
	// dipercaya ruang tersisa dari "Nilai Sekolah = ..." ke bawah sudah
	// cukup utk menampung tabel+Keputusan+footer ini sekaligus, jadi
	// digambar langsung di posisi kursor saat ini tanpa cek/loncat halaman

	pdf.SetFillColor(224, 224, 224)
	pdf.SetFont("Arial", "B", 7)
	pdf.Rect(marginL, y, colNo, headerH1+headerH2, "DF")
	pdf.SetXY(marginL, y+(headerH1+headerH2)/2-2)
	pdf.CellFormat(colNo, 4, "No.", "", 0, "C", false, 0, "")
	pdf.Rect(marginL+colNo, y, colMapel, headerH1+headerH2, "DF")
	pdf.SetXY(marginL+colNo, y+(headerH1+headerH2)/2-2)
	pdf.CellFormat(colMapel, 4, "Mata Pelajaran", "", 0, "C", false, 0, "")

	groupLabels := []string{"Nilai Ujian Sekolah Standar Daerah", "Nilai Sekolah", "Nilai Akhir *)"}
	x := marginL + colNo + colMapel
	for _, gl := range groupLabels {
		pdf.Rect(x, y, groupWidth, headerH1, "DF")
		lines := pdf.SplitLines([]byte(gl), groupWidth-1.5)
		nLines := len(lines)
		if nLines < 1 {
			nLines = 1
		}
		lineH := 3.0
		pdf.SetXY(x, y+(headerH1-float64(nLines)*lineH)/2)
		pdf.MultiCell(groupWidth, lineH, gl, "", "C", false)
		x += groupWidth
	}
	x = marginL + colNo + colMapel
	subHeaderTexts := []string{"Nilai Akhir", "Capaian Pembelajaran"}
	for i := 0; i < 3; i++ {
		cx := x
		for j, txt := range subHeaderTexts {
			sw := subW[i*2+j]
			pdf.Rect(cx, y+headerH1, sw, headerH2, "DF")
			pdf.SetFont("Arial", "", 6.5)
			lines := pdf.SplitLines([]byte(txt), sw)
			nLines := len(lines)
			if nLines < 1 {
				nLines = 1
			}
			lineH := 2.6
			pdf.SetXY(cx, y+headerH1+(headerH2-float64(nLines)*lineH)/2)
			pdf.MultiCell(sw, lineH, txt, "", "C", false)
			cx += sw
		}
		x += groupWidth
	}
	keputusanX := x

	dataY := y + headerH1 + headerH2
	pdf.SetFont("Arial", "", 8)
	rY := dataY
	for i, mapel := range labelEmpat {
		rh := rowHs[i]
		pdf.SetFont("Arial", "", 8)
		pdf.SetXY(marginL, rY)
		pdf.CellFormat(colNo, rh, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colMapel, rh, mapel, "1", 0, "L", false, 0, "")

		cx := marginL + colNo + colMapel
		for j := 0; j < 6; j++ {
			pdf.Rect(cx, rY, subW[j], rh, "D")

			if t := cellsEmpat[i][j]; t != "" {
				if j%2 == 0 {
					pdf.SetFont("Arial", "", 8)
					pdf.SetXY(cx, rY)
					pdf.CellFormat(subW[j], rh, t, "", 0, "C", false, 0, "")
				} else {
					pdf.SetFont("Arial", "", nilaiFontTeks)
					drawTeksTengahV(pdf, cx, rY, subW[j], rh, nilaiLineH, t, "L")
				}
			}

			cx += subW[j]
		}

		rY += rh
	}

	// kolom "Keputusan" — isi dulu semua teksnya, BARU gambar kotaknya
	// belakangan pakai tinggi sebenarnya yg terpakai (bukan asumsi
	// tableH) — supaya tidak tumpang tindih dgn "Kepala Sekolah" yg
	// menyatu di bagian bawah kotak yg sama (sesuai acuan asli, BUKAN
	// kolom terpisah di footer)
	pdf.SetFont("Arial", "B", 7.5)
	pdf.SetXY(keputusanX+2, y+2)
	pdf.CellFormat(colKeputusan-4, 3.5, "Keputusan:", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.SetX(keputusanX + 2)
	pdf.MultiCell(colKeputusan-4, 3.3, "Berdasarkan hasil belajar yang dicapai pada Semester 1 dan 2 Siswa / Siswi ditetapkan :", "", "L", false)
	pdf.Ln(1)
	// keputusan akhir (diinput manual lewat Buku Induk > Keputusan Akhir); kosong = kotak ( ) dibiarkan kosong
	var kep nilai.Keputusan
	if d != nil && d.Keputusan != nil {
		kep = *d.Keputusan
	}

	// label dipisah dari tanda kurung dan diberi lebar tetap (seperti TAB di Word),
	// supaya "( )" LULUS dan TIDAK LULUS sejajar lurus
	pdf.SetFont("Arial", "B", 8)

	tabW := pdf.GetStringWidth("TIDAK LULUS") + 6

	for _, opsi := range []struct {
		label string
		aktif bool
	}{
		{"LULUS", kep.Status == nilai.StatusLulus},
		{"TIDAK LULUS", kep.Status == nilai.StatusTidakLulus},
	} {
		tanda := "( )"
		if opsi.aktif {
			tanda = "(*)"
		}

		pdf.SetX(keputusanX + 2)
		pdf.CellFormat(tabW, 4.5, opsi.label, "", 0, "L", false, 0, "")
		pdf.CellFormat(colKeputusan-4-tabW, 4.5, tanda, "", 1, "L", false, 0, "")
	}

	// titik dua ketiga baris ini sejajar satu sama lain (lebar label = label terpanjang, "No. Ijazah");
	// tidak terikat ke tanda kurung di atas supaya No. Ijazah yang panjang tetap muat
	pdf.SetFont("Arial", "", 7.5)

	labelW := pdf.GetStringWidth("No. Ijazah") + 3

	for _, b := range []struct{ label, isi string }{
		{"No. Ijazah", kep.NoIjazah},
		{"No. SKHUS", kep.NoSKHUS},
		{"Tanggal", kep.Tanggal},
	} {
		pdf.SetX(keputusanX + 2)
		pdf.CellFormat(labelW, 4.5, b.label, "", 0, "L", false, 0, "")
		pdf.CellFormat(colKeputusan-4-labelW, 4.5, ": "+b.isi, "", 1, "L", false, 0, "")
	}

	textEndY := pdf.GetY()

	// "Kepala Sekolah" + Nama + NIP — NIP-nya disejajarkan persis dgn
	// baris "Tanggal:" di footer Status Akhir (pakai rumus footerY yg
	// sama), bukan sekadar mengikuti alur teks di atasnya
	rataYForAlign := dataY + sumH
	footerRowH := 6.0
	footerYForAlign := rataYForAlign + rowH                  // dipakukan, sesuai footerY yg sekarang tanpa jarak
	targetNipY := footerYForAlign + footerRowH*2 + 1.5 - 0.5 // diturunkan tipis lagi (0.5mm) dari sebelumnya (-1 -> -0.5)

	kepsekBlockY := targetNipY - 16
	if kepsekBlockY < textEndY+1 {
		kepsekBlockY = textEndY + 1 // jaga2 spy tidak tabrakan sama teks Keputusan di atasnya
	}
	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(keputusanX, kepsekBlockY)
	pdf.CellFormat(colKeputusan, 5, "Kepala Sekolah", "", 0, "C", false, 0, "")
	pdf.SetFont("Arial", "", 8)
	pdf.SetXY(keputusanX+2, kepsekBlockY+11)
	pdf.CellFormat(colKeputusan-4, 5, kepsek.NamaKepsek, "U", 2, "C", false, 0, "")
	pdf.SetX(keputusanX + 2)
	pdf.CellFormat(colKeputusan-4, 5, formatNIP(kepsek.NIPKepsek), "", 0, "C", false, 0, "")

	// kotak Keputusan digambar TERAKHIR, pakai tinggi maksimum antara
	// tableH (biar sejajar tabel di kiri) dan tinggi konten sebenarnya
	keputusanContentH := pdf.GetY() + 5 - y
	keputusanBoxH := tableH
	if keputusanContentH > keputusanBoxH {
		keputusanBoxH = keputusanContentH
	}
	pdf.Rect(keputusanX, y, colKeputusan, keputusanBoxH, "D")

	// baris RATA-RATA — full width No+Mapel+6 kolom nilai SAJA, tanpa
	// kolom Keputusan (kotaknya sudah berakhir tepat di atas baris ini)
	rataY := dataY + sumH
	pdf.SetFillColor(224, 224, 224)
	pdf.SetFont("Arial", "", 8)
	pdf.SetXY(marginL, rataY)
	pdf.CellFormat(colNo+colMapel, rowH, "RATA-RATA", "1", 0, "C", true, 0, "")
	for j := 0; j < 6; j++ {
		teks := ""
		if j%2 == 0 && len(kolNilai[j/2]) > 0 {
			jml := 0.0
			for _, v := range kolNilai[j/2] {
				jml += v
			}

			teks = fmtAngka(round2f(jml / float64(len(kolNilai[j/2]))))
		}

		pdf.SetFont("Arial", "", 8)
		pdf.CellFormat(subW[j], rowH, teks, "1", 0, "C", true, 0, "")
	}

	// ================= footer: 2 kolom (Wali Kelas | Status Akhir) =================
	// kolom kiri "Wali Kelas" selebar No+Mapel+1 grup (Standar Daerah
	// saja), kolom kanan "Status Akhir" selebar 2 grup (Nilai Sekolah +
	// Nilai Akhir) — berhenti tepat di keputusanX, TIDAK menyentuh kolom
	// Keputusan (Kepala Sekolah sudah menyatu di kotak Keputusan di atas)
	leftW := colNo + colMapel + groupWidth
	rightX := marginL + leftW
	rightW := groupWidth * 2

	rowA, rowB, rowC := 6.0, 6.0, 6.0
	totalH := rowA + rowB + rowC

	// naikkan footer supaya sejajar dgn baris "Kepala Sekolah" di kotak
	// dipakukan ke rataY (posisi stabil, sudah pasti muat di halaman yg
	// sama) — bukan kepsekBlockY lagi krn tingginya bisa berubah2
	// tergantung isi teks Keputusan & sempat memicu pindah halaman
	// nempel LANGSUNG ke RATA-RATA (tanpa jarak sama sekali) — jadi satu
	// tabel menyatu utuh dari header s/d Wali Kelas/Status Akhir, tanpa
	// terputus. TIDAK ADA pengecekan pindah halaman kedua di sini lagi
	// (itu justru penyebab footer sendirian nyasar ke halaman baru
	// sementara tabel di atasnya sudah kadung tercetak) — pengecekan
	// ruang SUDAH dilakukan sekali di awal fungsi (estimatedBlockH),
	// jadi cukup dipercaya di situ saja.
	footerY := rataY + rowH

	// --- kolom kiri: Wali Kelas (satu kotak menyatu, tanpa garis
	// pemisah label vs nama/NIP — sesuai revisi) ---
	pdf.Rect(marginL, footerY, leftW, rowA+rowB+rowC, "D")
	pdf.SetFont("Arial", "", 8)
	pdf.SetXY(marginL, footerY+1.5)
	pdf.CellFormat(leftW, 4, "Wali Kelas", "", 0, "C", false, 0, "")

	pdf.SetFont("Arial", "", 7)
	pdf.SetXY(marginL+3, footerY+rowA+4) // jarak lebih utk ruang tanda tangan
	pdf.CellFormat(leftW-6, 4.5, waliKelasNama, "U", 2, "C", false, 0, "")
	pdf.SetX(marginL + 3)
	if waliKelasNIP != "" {
		pdf.CellFormat(leftW-6, 4.5, formatNIP(waliKelasNIP), "", 0, "C", false, 0, "")
	} else {
		pdf.CellFormat(leftW-6, 4.5, "NIP.", "", 0, "C", false, 0, "")
	}

	// --- kolom kanan: Status Akhir (3 baris terpisah) ---
	pdf.SetFillColor(224, 224, 224)
	pdf.Rect(rightX, footerY, rightW, rowA, "DF")
	pdf.SetFont("Arial", "B", 7.5)
	pdf.SetXY(rightX, footerY+1.5)
	pdf.CellFormat(rightW, 4, "STATUS AKHIR", "", 0, "C", false, 0, "")

	pdf.Rect(rightX, footerY+rowA, rightW, rowB, "D")
	pdf.SetXY(rightX, footerY+rowA+1.5)
	statusTeks := "LULUS / TIDAK LULUS"
	if kep.Status != "" {
		statusTeks = strings.ToUpper(kep.Status)
	}

	pdf.SetFont("Arial", "B", 7.5)
	pdf.CellFormat(rightW, 4, statusTeks, "", 0, "C", false, 0, "")

	pdf.Rect(rightX, footerY+rowA+rowB, rightW, rowC, "D")
	pdf.SetFont("Arial", "", 7)
	pdf.SetXY(rightX+2, footerY+rowA+rowB+1.5)
	pdf.CellFormat(rightW-4, 4, "Tanggal: "+kep.Tanggal, "", 0, "L", false, 0, "")

	pdf.SetY(footerY + totalH + 4)
}

// renderKepsekWaliKelasBlock — kolom kebalik dari renderWaliKelasKepsekBlock
// (Kepala Sekolah dulu baru Wali Kelas), sesuai urutan di halaman 2.
func renderKepsekWaliKelasBlock(pdf *fpdf.Fpdf, kepsek KepalaSekolahInfo, waliKelasNama, waliKelasNIP string) {
	pdf.Ln(8)
	colWidth := (pageWidth - marginL - marginR) / 2
	startY := pdf.GetY()
	pdf.SetFont("Arial", "", 10)

	pdf.SetXY(marginL, startY)
	pdf.CellFormat(colWidth, 6, "Kepala Sekolah", "", 2, "L", false, 0, "")
	pdf.SetXY(marginL+colWidth, startY)
	pdf.CellFormat(colWidth, 6, "Wali Kelas", "", 2, "L", false, 0, "")

	pdf.SetXY(marginL, startY+20)
	pdf.CellFormat(colWidth, 6, kepsek.NamaKepsek, "U", 2, "L", false, 0, "")
	pdf.SetX(marginL)
	pdf.CellFormat(colWidth, 6, formatNIP(kepsek.NIPKepsek), "", 2, "L", false, 0, "")

	pdf.SetXY(marginL+colWidth, startY+20)
	pdf.CellFormat(colWidth, 6, waliKelasNama, "U", 2, "L", false, 0, "")
	pdf.SetX(marginL + colWidth)
	if waliKelasNIP != "" {
		pdf.CellFormat(colWidth, 6, formatNIP(waliKelasNIP), "", 2, "L", false, 0, "")
	} else {
		pdf.CellFormat(colWidth, 6, "NIP.", "", 2, "L", false, 0, "")
	}
}

// renderDataNilaiHeader — header identitas 2 kolom (Nama Siswa|NISN,
// Nomor Induk|Tahun Pelajaran), posisi X/Y dikontrol eksplisit per sel
// supaya tidak salah tumpuk baris (bug sebelumnya).
func renderDataNilaiHeader(pdf *fpdf.Fpdf, siswa BukuIndukData, tahunPelajaran string) {
	pdf.SetFont("Arial", "", 10)
	leftX := marginL
	rightX := pageWidth/2 + 10
	startY := pdf.GetY()

	pdf.SetXY(leftX, startY)
	pdf.CellFormat(30, 6, "Nama Siswa", "", 0, "L", false, 0, "")
	pdf.CellFormat(4, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(rightX-leftX-34, 6, siswa.NamaLengkap, "", 0, "L", false, 0, "")

	pdf.SetXY(rightX, startY)
	pdf.CellFormat(30, 6, "NISN", "", 0, "L", false, 0, "")
	pdf.CellFormat(4, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, siswa.NISN, "", 1, "L", false, 0, "")

	pdf.SetXY(leftX, startY+7)
	pdf.CellFormat(30, 6, "NIS", "", 0, "L", false, 0, "")
	pdf.CellFormat(4, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(rightX-leftX-34, 6, siswa.NIS, "", 0, "L", false, 0, "")

	pdf.SetXY(rightX, startY+7)
	pdf.CellFormat(30, 6, "Tahun Pelajaran", "", 0, "L", false, 0, "")
	pdf.CellFormat(4, 6, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, tahunPelajaran, "", 1, "L", false, 0, "")

	pdf.SetY(startY + 13)
}

func renderPancasilaLegend(pdf *fpdf.Fpdf) {
	fullWidth := pageWidth - marginL - marginR
	colW := fullWidth / 4
	items := []struct{ label, desc string }{
		{"BB *Belum Berkembang", "Siswa masih membutuhkan bimbingan dalam mengembangkan kemampuan"},
		{"MB *Mulai Berkembang", "Siswa mulai mengembangkan kemampuan namun masih belum ajek"},
		{"BSH *Berkembang Sesuai Harapan", "Siswa telah mengembangkan kemampuan hingga berada dalam tahap ajek"},
		{"SB * Sangat Berkembang", "Siswa telah mengembangkan kemampuan hingga berada dalam tahap ajek"},
	}
	y := pdf.GetY()
	x := marginL
	maxH := 0.0
	for _, it := range items {
		pdf.SetXY(x, y)
		pdf.SetFont("Arial", "B", 8)
		pdf.MultiCell(colW-2, 4, it.label, "", "L", false)
		pdf.SetXY(x, pdf.GetY())
		pdf.SetFont("Arial", "", 7)
		pdf.MultiCell(colW-2, 3.5, it.desc, "", "L", false)
		if h := pdf.GetY() - y; h > maxH {
			maxH = h
		}
		x += colW
	}
	pdf.SetY(y + maxH + 2)
}

// idMapelByNama = id mapel daun (tanpa sub-mapel) dengan nama persis itu di pohon mapel; 0 kalau tidak ada.
func idMapelByNama(pohon []nilai.MapelPohon, nama string) int {
	for _, m := range pohon {
		if len(m.Subs) == 0 && strings.EqualFold(strings.TrimSpace(m.Nama), nama) {
			return m.ID
		}
	}

	return 0
}

func round2f(v float64) float64 {
	return math.Round(v*100) / 100
}
