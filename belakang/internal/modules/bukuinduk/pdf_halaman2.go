package bukuinduk

import (
	"github.com/go-pdf/fpdf"

	"sekolah/internal/modules/nilai"
)

// Lebar kolom "Nilai Akhir" (angka pendek) & "Capaian Pembelajaran" (teks) di dalam 1 grup tabel halaman 2.
const (
	uAkhirW = 12.0
)

// fmtPtr: nil -> "" (sel dibiarkan kosong).
func fmtPtr(v *float64) string {
	if v == nil {
		return ""
	}

	return fmtAngka(*v)
}

func fmtMap(m map[int]float64, id int) string {
	if v, ok := m[id]; ok {
		return fmtAngka(v)
	}

	return ""
}

func fmtUjian(m map[int]nilai.SelUjian, id int) (string, string) {
	s, ok := m[id]
	if !ok {
		return "", ""
	}

	return fmtPtr(s.Nilai), s.Capaian
}

// ujianCells = 6 teks sel utk 1 mapel:
// [Raport akhir, Raport capaian, Ujian akhir, Ujian capaian, Nilai Sekolah akhir, Nilai Sekolah capaian].
// Capaian Raport & Nilai Sekolah sengaja kosong (tidak punya sumber data).
func ujianCells(d *nilai.DataHalaman2, id int) [6]string {
	var out [6]string
	if d == nil || id == 0 {
		return out
	}

	out[0] = fmtMap(d.Raport, id)
	out[2], out[3] = fmtUjian(d.Sekolah, id)
	out[4] = fmtMap(d.Gabung, id)

	return out
}

// drawUjianDataRow = 1 baris mapel: No, Mapel, 6 sel (3 grup x [Nilai Akhir, Capaian]).
func drawUjianDataRow(pdf *fpdf.Fpdf, y, colNo, colMapel float64, widths [6]float64, no, label string, cells [6]string) float64 {
	pdf.SetFont("Arial", "", nilaiFontMapel)
	tinggiLabel := float64(jumlahBaris(pdf, label, colMapel)) * nilaiLineHMapel

	pdf.SetFont("Arial", "", nilaiFontTeks)

	tinggiTeks := 0.0

	for i, t := range cells {
		if i%2 == 0 {
			continue
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

	pdf.Rect(marginL, y, colNo, rowH, "D")
	pdf.SetFont("Arial", "", nilaiFontMapel)
	pdf.SetXY(marginL, y)
	pdf.CellFormat(colNo, rowH, no, "", 0, "C", false, 0, "")

	pdf.Rect(marginL+colNo, y, colMapel, rowH, "D")
	drawTeksTengahV(pdf, marginL+colNo, y, colMapel, rowH, nilaiLineHMapel, label, "L")

	x := marginL + colNo + colMapel

	for i := 0; i < 6; i++ {
		pdf.Rect(x, y, widths[i], rowH, "D")

		if cells[i] != "" {
			if i%2 == 0 {
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

// drawUjianParentRow = baris induk (Seni Budaya & Prakarya, P5): 1 sel utuh, tanpa kotak nilai.
func drawUjianParentRow(pdf *fpdf.Fpdf, y, colNo, mergedW float64, no, nama string) float64 {
	pdf.SetFont("Arial", "", nilaiFontMapel)

	rowH := float64(jumlahBaris(pdf, nama, mergedW)) * 3.5
	if rowH < nilaiMinRowH {
		rowH = nilaiMinRowH
	}

	y = ensureRow(pdf, y, rowH)

	pdf.SetFillColor(224, 224, 224)
	pdf.Rect(marginL, y, colNo, rowH, "DF")
	pdf.SetXY(marginL, y)
	pdf.CellFormat(colNo, rowH, no, "", 0, "C", false, 0, "")
	pdf.Rect(marginL+colNo, y, mergedW, rowH, "DF")
	drawTeksTengahV(pdf, marginL+colNo, y, mergedW, rowH, 3.5, nama, "L")

	return y + rowH
}

// renderRataRataUjianTable = tabel 1 halaman 2: Nilai Rata-rata Raport / Nilai Ujian Sekolah / Nilai Sekolah.
// Rata-rata raport dihitung otomatis dari nilai_mapel, ujian dari nilai_ujian (jenis Sekolah),
// Nilai Sekolah = bobot * raport + bobot * ujian (lihat nilai/ujian.go). d boleh nil -> tabel kosong.
func renderRataRataUjianTable(pdf *fpdf.Fpdf, d *nilai.DataHalaman2) {
	fullWidth := pageWidth - marginL - marginR
	colNo, colMapel := 10.0, 50.0
	groupWidth := (fullWidth - colNo - colMapel) / 3
	capW := groupWidth - uAkhirW
	widths := [6]float64{uAkhirW, capW, uAkhirW, capW, uAkhirW, capW}

	y := pdf.GetY()
	pdf.SetFillColor(224, 224, 224)
	pdf.SetFont("Arial", "B", 7)
	pdf.Rect(marginL, y, colNo, 12, "DF")
	pdf.SetXY(marginL, y+3)
	pdf.CellFormat(colNo, 6, "No.", "", 0, "C", false, 0, "")
	pdf.Rect(marginL+colNo, y, colMapel, 12, "DF")
	pdf.SetXY(marginL+colNo, y+3)
	pdf.CellFormat(colMapel, 6, "Mata Pelajaran", "", 0, "C", false, 0, "")

	x := marginL + colNo + colMapel
	for _, gl := range []string{"Nilai Rata-rata Raport", "Nilai Ujian Sekolah", "Nilai Sekolah"} {
		pdf.Rect(x, y, groupWidth, 6, "DF")
		pdf.SetXY(x, y)
		pdf.CellFormat(groupWidth, 6, gl, "", 0, "C", false, 0, "")
		x += groupWidth
	}

	x = marginL + colNo + colMapel

	for i, txt := range []string{"Nilai Akhir", "Capaian Pembelajaran", "Nilai Akhir", "Capaian Pembelajaran", "Nilai Akhir", "Capaian Pembelajaran"} {
		w := widths[i]
		pdf.Rect(x, y+6, w, 6, "DF")
		pdf.SetFont("Arial", "B", 6.5)

		n := jumlahBaris(pdf, txt, w)
		lineH := 2.6

		pdf.SetXY(x, y+6+(6-float64(n)*lineH)/2)
		pdf.MultiCell(w, lineH, txt, "", "C", false)

		x += w
	}

	pohon := mapelPohonStatis()
	if d != nil && len(d.Mapel) > 0 {
		pohon = d.Mapel
	}

	mergedW := colMapel + groupWidth*3
	dataY := y + 12

	pdf.SetAutoPageBreak(false, 0) // page-break dikerjakan manual oleh ensureRow

	letters := "abcdefgh"

	for _, m := range pohon {
		if len(m.Subs) == 0 {
			dataY = drawUjianDataRow(pdf, dataY, colNo, colMapel, widths, m.No, m.Nama, ujianCells(d, m.ID))
			continue
		}

		dataY = drawUjianParentRow(pdf, dataY, colNo, mergedW, m.No, m.Nama)

		for i, sub := range m.Subs {
			letter := "?"
			if i < len(letters) {
				letter = string(letters[i])
			}

			dataY = drawUjianDataRow(pdf, dataY, colNo, colMapel, widths, "", letter+". "+sub.Nama, ujianCells(d, sub.ID))
		}
	}

	pdf.SetAutoPageBreak(true, marginB)
	pdf.SetXY(marginL, dataY) // Rect() tidak menggerakkan kursor; sinkronkan dulu sebelum ensureSpace
	ensureSpace(pdf, 12)      // "Jumlah Nilai" + "Rata-rata" harus tetap 1 halaman
	dataY = pdf.GetY()

	var jumlah, rata [3]*float64
	if d != nil {
		jumlah = [3]*float64{d.JumlahRaport, d.JumlahSekolah, d.JumlahGabung}
		rata = [3]*float64{d.RataRaport, d.RataSekolah, d.RataGabung}
	}

	pdf.SetFillColor(224, 224, 224)

	for bi, b := range []struct {
		label string
		nilai [3]*float64
	}{{"Jumlah Nilai", jumlah}, {"Rata-rata", rata}} {
		_ = bi

		pdf.SetFont("Arial", "B", 8)
		pdf.SetXY(marginL, dataY)
		pdf.CellFormat(colNo+colMapel, 6, b.label, "1", 0, "C", true, 0, "") // hanya label yang di-shading

		for g := 0; g < 3; g++ {
			pdf.SetFont("Arial", "", 7.5)
			pdf.CellFormat(uAkhirW, 6, fmtPtr(b.nilai[g]), "1", 0, "C", false, 0, "")
			pdf.CellFormat(capW, 6, "", "1", 0, "C", false, 0, "")
		}

		dataY += 6
	}

	pdf.SetY(dataY)
}
