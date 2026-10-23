package bukuinduk

import (
	"fmt"

	"github.com/go-pdf/fpdf"
)

const (
	eNoColWidth    = 8.0
	eAspekColWidth = 30.0
	eGroupCount    = 6
)

// emptyCellMark dipakai utk sel sempit (Bagian E) yg tidak muat kalau
// pakai dotsIfEmpty (titik-titiknya kepanjangan utk kolom ~12mm).
const emptyCellMark = "-"

func renderBagianE(pdf *fpdf.Fpdf, siswa BukuIndukData, tahunPelajaranSaatIni string) {
	sectionHeading(pdf, "E. Lain-lain")
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 6, "1. Tinggi & berat badan", "", 1, "L", false, 0, "")
	renderTinggiBeratTable(pdf, siswa, tahunPelajaranSaatIni)

	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 6, "2. Kondisi Kesehatan", "", 1, "L", false, 0, "")
	renderKesehatanTable(pdf, tahunPelajaranSaatIni)
}

func renderTinggiBeratTable(pdf *fpdf.Fpdf, siswa BukuIndukData, tahunPelajaranSaatIni string) {
	usableWidth := pageWidth - marginL - marginR
	groupWidth := (usableWidth - eNoColWidth - eAspekColWidth) / eGroupCount
	subWidth := groupWidth / 2

	startX := marginL
	headerY := pdf.GetY()
	rowHeaderH := 10.0
	subH := 7.0

	// No. & Aspek yg dinilai — tinggi digabung (header + sub-header)
	// supaya sejajar dgn kolom Thn Pelajaran di sebelahnya
	pdf.SetFont("Arial", "B", 8)
	fullHeaderH := rowHeaderH + subH
	pdf.Rect(startX, headerY, eNoColWidth, fullHeaderH, "D")
	pdf.SetXY(startX, headerY)
	pdf.CellFormat(eNoColWidth, fullHeaderH, "No.", "", 0, "C", false, 0, "")
	pdf.Rect(startX+eNoColWidth, headerY, eAspekColWidth, fullHeaderH, "D")
	pdf.SetXY(startX+eNoColWidth, headerY)
	pdf.CellFormat(eAspekColWidth, fullHeaderH, "Aspek yg dinilai", "", 0, "C", false, 0, "")

	x := startX + eNoColWidth + eAspekColWidth
	for g := 0; g < eGroupCount; g++ {
		label := "..... / ....."
		if g == 0 {
			label = tahunPelajaranSaatIni
		}
		pdf.SetXY(x, headerY)
		pdf.SetFont("Arial", "B", 7)
		pdf.MultiCell(groupWidth, rowHeaderH/2, "Thn pelajaran\n"+label, "1", "C", false)
		x += groupWidth
	}

	subY := headerY + rowHeaderH
	x = startX + eNoColWidth + eAspekColWidth
	pdf.SetFont("Arial", "B", 6)
	for g := 0; g < eGroupCount; g++ {
		pdf.SetXY(x, subY)
		pdf.MultiCell(subWidth, subH/2, "Sem\nganjil", "1", "C", false)
		pdf.SetXY(x+subWidth, subY)
		pdf.MultiCell(subWidth, subH/2, "Sem\ngenap", "1", "C", false)
		x += groupWidth
	}

	dataY := subY + subH
	dataRowH := 7.0
	pdf.SetFont("Arial", "", 8)
	rows := []struct{ no, label, satuan, nilai string }{
		{"1", "Tinggi badan", " cm", siswa.TinggiBadan},
		{"2", "Berat badan", " kg", siswa.BeratBadan},
	}
	for _, r := range rows {
		pdf.SetXY(startX, dataY)
		pdf.CellFormat(eNoColWidth, dataRowH, r.no, "1", 0, "C", false, 0, "")
		pdf.CellFormat(eAspekColWidth, dataRowH, r.label, "1", 0, "L", false, 0, "")
		x = startX + eNoColWidth + eAspekColWidth
		for g := 0; g < eGroupCount; g++ {
			ganjil := emptyCellMark
			if g == 0 && r.nilai != "" {
				ganjil = r.nilai + r.satuan
			}
			pdf.CellFormat(subWidth, dataRowH, ganjil, "1", 0, "C", false, 0, "")
			pdf.CellFormat(subWidth, dataRowH, emptyCellMark, "1", 0, "C", false, 0, "")
			x += groupWidth
		}
		dataY += dataRowH
	}
	pdf.SetY(dataY)
}

func renderKesehatanTable(pdf *fpdf.Fpdf, tahunPelajaranSaatIni string) {
	usableWidth := pageWidth - marginL - marginR
	groupWidth := (usableWidth - eNoColWidth - eAspekColWidth) / eGroupCount

	startX := marginL
	headerY := pdf.GetY()
	headerH := 15.0

	pdf.SetFont("Arial", "B", 8)
	pdf.Rect(startX, headerY, eNoColWidth, headerH, "D")
	pdf.SetXY(startX, headerY)
	pdf.CellFormat(eNoColWidth, headerH, "No.", "", 0, "C", false, 0, "")
	pdf.Rect(startX+eNoColWidth, headerY, eAspekColWidth, headerH, "D")
	pdf.SetXY(startX+eNoColWidth, headerY)
	pdf.CellFormat(eAspekColWidth, headerH, "Aspek yg dinilai", "", 0, "C", false, 0, "")

	x := startX + eNoColWidth + eAspekColWidth
	for g := 0; g < eGroupCount; g++ {
		label := "..... / ....."
		if g == 0 {
			label = tahunPelajaranSaatIni
		}
		pdf.SetXY(x, headerY)
		pdf.SetFont("Arial", "B", 7)
		pdf.MultiCell(groupWidth, headerH/3, "Thn pelajaran\n"+label+"\nKeterangan", "1", "C", false)
		x += groupWidth
	}

	dataY := headerY + headerH
	dataRowH := 7.0
	pdf.SetFont("Arial", "", 8)
	aspek := []string{"Pendengaran", "Pengelihatan", "Gigi", "-"}
	for i, a := range aspek {
		pdf.SetXY(startX, dataY)
		pdf.CellFormat(eNoColWidth, dataRowH, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(eAspekColWidth, dataRowH, a, "1", 0, "L", false, 0, "")
		x = startX + eNoColWidth + eAspekColWidth
		for g := 0; g < eGroupCount; g++ {
			pdf.CellFormat(groupWidth, dataRowH, emptyCellMark, "1", 0, "C", false, 0, "")
			x += groupWidth
		}
		dataY += dataRowH
	}
	pdf.SetY(dataY)
}
