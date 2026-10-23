package bukuinduk

import (
	"strconv"
	"strings"

	"github.com/go-pdf/fpdf"

	"sekolah/internal/modules/nilai"
)

// nilaiCols = lebar kolom tabel "B. Pengetahuan dan Keterampilan".
// Nilai Akhir dibuat sempit, Capaian Pembelajaran & KKTP dilebarkan karena berisi teks.
type nilaiCols struct {
	No, Mapel            float64
	Akhir, Capaian, Kktp float64
}

func (c nilaiCols) sem() float64    { return c.Akhir + c.Capaian + c.Kktp }
func (c nilaiCols) merged() float64 { return c.Mapel + 2*c.sem() }

const (
	nilaiFontMapel  = 7.5 // nama mapel, No, Nilai Akhir
	nilaiFontTeks   = 6   // Capaian Pembelajaran & KKTP (teks panjang)
	nilaiLineHMapel = 3.4
	nilaiLineH      = 2.8 // tinggi baris teks Capaian/KKTP
	nilaiMinRowH    = 6.0
	nilaiPadY       = 0.8
)

var semesterUrut = []string{"Ganjil", "Genap"}

// fmtAngka: 85 -> "85", 85.5 -> "85,5" (koma desimal), maksimal 2 desimal.
func fmtAngka(v float64) string {
	return strings.ReplaceAll(strconv.FormatFloat(v, 'f', -1, 64), ".", ",")
}

// lebarTeks = lebar yang diberikan ke SplitLines. fpdf.SplitLines SUDAH mengurangi margin sel
// (2 x 1 mm) sendiri, jadi lebar sel diberikan apa adanya. Mengurangi lagi di sini membuat
// lebar efektif terlalu sempit (nama panjang terpecah padahal masih muat).
func lebarTeks(w float64) float64 { return w }

func jumlahBaris(pdf *fpdf.Fpdf, text string, w float64) int {
	if text == "" {
		return 1
	}

	n := len(pdf.SplitLines([]byte(text), lebarTeks(w)))
	if n < 1 {
		n = 1
	}

	return n
}

// drawTeksTengahV menulis teks rata kiri tapi terpusat vertikal di sel (x, y, w, h).
// Pemotongan baris & penggambaran memakai SplitLines yang sama (bukan MultiCell), jadi jumlah
// baris yang dipakai utk menghitung posisi PASTI sama dgn yang tergambar.
func drawTeksTengahV(pdf *fpdf.Fpdf, x, y, w, h, lineH float64, text, align string) {
	lines := pdf.SplitLines([]byte(text), lebarTeks(w))
	if len(lines) == 0 {
		return
	}

	startY := y + (h-float64(len(lines))*lineH)/2
	for i, ln := range lines {
		pdf.SetXY(x, startY+float64(i)*lineH)
		pdf.CellFormat(w, lineH, string(ln), "", 0, align, false, 0, "")
	}
}

// mapelPohonStatis = daftar mapel bawaan (daftarMapel) dalam bentuk pohon; dipakai kalau tbl_mapel kosong / gagal dibaca.
func mapelPohonStatis() []nilai.MapelPohon {
	out := []nilai.MapelPohon{}

	for _, m := range daftarMapel {
		if m.No == "MULOK" {
			continue
		}

		p := nilai.MapelPohon{No: m.No, Nama: m.Nama}
		for _, s := range m.SubList {
			p.Subs = append(p.Subs, nilai.MapelPohon{Nama: s})
		}

		out = append(out, p)
	}

	return out
}

// selNilaiTeks = 6 teks sel utk 1 mapel: [Ganjil: akhir, capaian, kktp, Genap: akhir, capaian, kktp].
func selNilaiTeks(data *nilai.DataBukuInduk, idMapel int) [6]string {
	var out [6]string
	if data == nil || idMapel == 0 {
		return out
	}

	for i, sem := range semesterUrut {
		s, ok := data.Nilai[idMapel][sem]
		if !ok {
			continue
		}

		if s.Akhir != nil {
			out[i*3] = fmtAngka(*s.Akhir)
		}

		out[i*3+1] = s.Capaian
		out[i*3+2] = s.Kktp
	}

	return out
}

// ensureRow pindah halaman kalau baris setinggi rowH tidak muat. Rect() tidak dipantau
// auto-pagebreak fpdf, jadi pengecekan manual ini wajib sebelum menggambar apa pun.
func ensureRow(pdf *fpdf.Fpdf, y, rowH float64) float64 {
	if y+rowH > pageHeight-marginB {
		pdf.AddPage()
		return marginT
	}

	return y
}

// drawNilaiDataRow menggambar 1 baris mapel (No, Mapel, 6 sel nilai). Tinggi baris menyesuaikan
// teks terpanjang (Capaian/KKTP bisa beberapa baris).
func drawNilaiDataRow(pdf *fpdf.Fpdf, y float64, c nilaiCols, no, label string, cells [6]string) float64 {
	widths := [6]float64{c.Akhir, c.Capaian, c.Kktp, c.Akhir, c.Capaian, c.Kktp}

	// tinggi baris = yang tertinggi antara nama mapel (font besar) dan Capaian/KKTP (font kecil)
	pdf.SetFont("Arial", "", nilaiFontMapel)
	tinggiLabel := float64(jumlahBaris(pdf, label, c.Mapel)) * nilaiLineHMapel

	pdf.SetFont("Arial", "", nilaiFontTeks)
	tinggiTeks := 0.0
	for i, t := range cells {
		if i%3 == 0 {
			continue // Nilai Akhir = angka pendek, satu baris
		}

		if h := float64(jumlahBaris(pdf, t, widths[i])) * nilaiLineH; h > tinggiTeks {
			tinggiTeks = h
		}
	}

	rowH := tinggiLabel
	if tinggiTeks > rowH {
		rowH = tinggiTeks
	}

	rowH += nilaiPadY * 2
	if rowH < nilaiMinRowH {
		rowH = nilaiMinRowH
	}

	y = ensureRow(pdf, y, rowH)

	pdf.Rect(marginL, y, c.No, rowH, "D")
	pdf.SetFont("Arial", "", nilaiFontMapel)
	pdf.SetXY(marginL, y)
	pdf.CellFormat(c.No, rowH, no, "", 0, "C", false, 0, "")

	pdf.Rect(marginL+c.No, y, c.Mapel, rowH, "D")
	// rata kiri, tapi terpusat vertikal supaya sejajar dengan No dan Nilai Akhir
	drawTeksTengahV(pdf, marginL+c.No, y, c.Mapel, rowH, nilaiLineHMapel, label, "L")

	x := marginL + c.No + c.Mapel
	for i := 0; i < 6; i++ {
		pdf.Rect(x, y, widths[i], rowH, "D")

		if cells[i] != "" {
			if i%3 == 0 {
				pdf.SetFont("Arial", "", nilaiFontMapel)
				pdf.SetXY(x, y)
				pdf.CellFormat(widths[i], rowH, cells[i], "", 0, "C", false, 0, "")
			} else {
				pdf.SetFont("Arial", "", nilaiFontTeks)
				drawTeksTengahV(pdf, x, y, widths[i], rowH, nilaiLineH, cells[i], "L")
			}
		}

		x += widths[i]
	}

	return y + rowH
}

// drawNilaiParentRow = baris induk (mis. "10 Seni Budaya & Prakarya"), 1 sel utuh tanpa kotak nilai.
func drawNilaiParentRow(pdf *fpdf.Fpdf, y float64, c nilaiCols, no, nama string) float64 {
	pdf.SetFont("Arial", "", nilaiFontMapel)

	n := jumlahBaris(pdf, nama, c.merged())
	rowH := float64(n) * 3.5
	if rowH < nilaiMinRowH {
		rowH = nilaiMinRowH
	}

	y = ensureRow(pdf, y, rowH)

	pdf.SetFillColor(224, 224, 224)
	pdf.Rect(marginL, y, c.No, rowH, "DF")
	pdf.SetXY(marginL, y)
	pdf.CellFormat(c.No, rowH, no, "", 0, "C", false, 0, "")

	pdf.Rect(marginL+c.No, y, c.merged(), rowH, "DF")
	drawTeksTengahV(pdf, marginL+c.No, y, c.merged(), rowH, 3.5, nama, "L")

	return y + rowH
}

// renderMapelNilaiTable = tabel "B. Pengetahuan dan Keterampilan". data boleh nil -> tabel kosong (siap isi tangan).
func renderMapelNilaiTable(pdf *fpdf.Fpdf, data *nilai.DataBukuInduk) {
	fullWidth := pageWidth - marginL - marginR

	c := nilaiCols{No: 8, Mapel: 46, Akhir: 11}
	sisa := (fullWidth - c.No - c.Mapel) / 2 // lebar 1 semester
	c.Capaian = (sisa - c.Akhir) / 2
	c.Kktp = c.Capaian

	y := pdf.GetY()
	pdf.SetFillColor(224, 224, 224)
	pdf.SetFont("Arial", "B", 7)

	pdf.Rect(marginL, y, c.No, 14, "DF")
	pdf.SetXY(marginL, y+5)
	pdf.CellFormat(c.No, 5, "No.", "", 0, "C", false, 0, "")
	pdf.Rect(marginL+c.No, y, c.Mapel, 14, "DF")
	pdf.SetXY(marginL+c.No, y+5)
	pdf.CellFormat(c.Mapel, 5, "Mata Pelajaran", "", 0, "C", false, 0, "")

	x := marginL + c.No + c.Mapel
	pdf.Rect(x, y, c.sem(), 7, "DF")
	pdf.SetXY(x, y)
	pdf.CellFormat(c.sem(), 7, "Semester 1 (Ganjil)", "", 0, "C", false, 0, "")
	pdf.Rect(x+c.sem(), y, c.sem(), 7, "DF")
	pdf.SetXY(x+c.sem(), y)
	pdf.CellFormat(c.sem(), 7, "Semester 2 (Genap)", "", 0, "C", false, 0, "")

	subHeaders := []string{"Nilai Akhir", "Capaian Pembelajaran", "KKTP"}
	subWidths := []float64{c.Akhir, c.Capaian, c.Kktp}
	pdf.SetFont("Arial", "B", 6.5)
	x = marginL + c.No + c.Mapel
	for s := 0; s < 2; s++ {
		for i, sh := range subHeaders {
			w := subWidths[i]
			pdf.Rect(x, y+7, w, 7, "DF")
			n := jumlahBaris(pdf, sh, w)
			const lineH = 3.2
			pdf.SetXY(x, y+7+(7-float64(n)*lineH)/2)
			pdf.MultiCell(w, lineH, sh, "", "C", false)
			x += w
		}
	}

	pohon := mapelPohonStatis()
	if data != nil && len(data.Mapel) > 0 {
		pohon = data.Mapel
	}

	dataY := y + 14
	pdf.SetAutoPageBreak(false, 0) // page-break dikerjakan manual oleh ensureRow
	letters := "abcdefgh"

	for _, m := range pohon {
		if len(m.Subs) == 0 {
			dataY = drawNilaiDataRow(pdf, dataY, c, m.No, m.Nama, selNilaiTeks(data, m.ID))
			continue
		}

		dataY = drawNilaiParentRow(pdf, dataY, c, m.No, m.Nama)
		for i, sub := range m.Subs {
			letter := "?"
			if i < len(letters) {
				letter = string(letters[i])
			}

			dataY = drawNilaiDataRow(pdf, dataY, c, "", letter+". "+sub.Nama, selNilaiTeks(data, sub.ID))
		}
	}

	pdf.SetAutoPageBreak(true, marginB)
	pdf.SetXY(marginL, dataY) // Rect() tidak menggerakkan kursor; sinkronkan dulu sebelum ensureSpace
	ensureSpace(pdf, 12)      // "Jumlah Nilai" + "Rata-rata" harus tetap 1 halaman
	dataY = pdf.GetY()

	pdf.SetFillColor(224, 224, 224)
	pdf.SetFont("Arial", "B", 8)

	barisTotal := []struct {
		label string
		nilai map[string]float64
	}{
		{"Jumlah Nilai", nil},
		{"Rata-rata", nil},
	}
	if data != nil {
		barisTotal[0].nilai = data.Jumlah
		barisTotal[1].nilai = data.Rata
	}

	for _, b := range barisTotal {
		// hanya label (No + Mata Pelajaran) yang di-shading
		pdf.SetXY(marginL, dataY)
		pdf.CellFormat(c.No+c.Mapel, 6, b.label, "1", 0, "C", true, 0, "")

		// area nilai: tanpa shading, tapi tetap bergaris per kolom (Nilai Akhir / Capaian / KKTP)
		// sehingga garis kolom tabel di atasnya menerus sampai baris Rata-rata
		for _, sem := range semesterUrut {
			teks := ""
			if v, ok := b.nilai[sem]; ok {
				teks = fmtAngka(v)
			}

			pdf.CellFormat(c.Akhir, 6, teks, "1", 0, "C", false, 0, "")
			pdf.CellFormat(c.Capaian, 6, "", "1", 0, "C", false, 0, "")
			pdf.CellFormat(c.Kktp, 6, "", "1", 0, "C", false, 0, "")
		}

		dataY += 6
	}

	pdf.SetY(dataY)
}

// hariTeks = isi sel ketidakhadiran ("3 hari"). Kalau absensi semester itu belum pernah diinput,
// dibiarkan "hari" kosong seperti template lama (bukan "0 hari" yang menyesatkan).
func hariTeks(data *nilai.DataBukuInduk, semester string, baris int) string {
	if data == nil {
		return "hari"
	}

	k, ok := data.Ketidakhadiran[semester]
	if !ok || k.Tercatat == 0 {
		return "hari"
	}

	jumlah := [3]int{k.Sakit, k.Izin, k.Alpa}

	return strconv.Itoa(jumlah[baris]) + " hari"
}

func labelKelas(data *nilai.DataBukuInduk) string {
	if data != nil && data.NamaKelas != "" {
		return "KELAS " + strings.ToUpper(data.NamaKelas)
	}

	return "KELAS ..."
}
