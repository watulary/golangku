package bukuinduk

import (
	"fmt"
	"strings"

	"github.com/go-pdf/fpdf"
)

func renderNomorUrutPage(pdf *fpdf.Fpdf, siswaList []BukuIndukData, kepsek KepalaSekolahInfo, namaPetugas, nipPetugas, kotaSekolah string, waliKelasNama, waliKelasNIP string) {
	pdf.AddPage()
	centerText(pdf, "II. NOMOR URUT PESERTA DIDIK DALAM BUKU INDUK REGISTER", 12, true)
	pdf.Ln(4)

	headers := []string{"No.\nUrut", "Nomor\nInduk Siswa", "Nomor Induk\nSiswa Nasional", "Nama Peserta Didik", "L/P", "Rombel"}
	widths := []float64{12, 26, 30, 68, 12, 38}
	headerHeight := 10.0
	dataFontSize := 8.0
	lineH := 3.8

	drawNomorUrutHeader := func() {
		pdf.SetFillColor(224, 224, 224)
		pdf.SetFont("Arial", "B", 8)
		y := pdf.GetY()
		x := marginL
		for i, h := range headers {
			pdf.Rect(x, y, widths[i], headerHeight, "DF")
			lines := strings.Split(h, "\n")
			cellLineH := headerHeight / float64(len(lines))
			for li, line := range lines {
				pdf.SetXY(x, y+float64(li)*cellLineH)
				pdf.CellFormat(widths[i], cellLineH, line, "", 0, "C", false, 0, "")
			}
			x += widths[i]
		}
		pdf.SetXY(marginL, y+headerHeight)
	}

	drawNomorUrutHeader()
	pdf.SetFont("Arial", "", dataFontSize)

	for i, s := range siswaList {
		values := []string{
			fmt.Sprintf("%d.", i+1),
			s.NIS, s.NISN, s.NamaLengkap, jenisKelaminSingkat(s.JenisKelamin), s.Kelas,
		}

		maxLines := 1
		for j, v := range values {
			lines := pdf.SplitLines([]byte(v), widths[j]-1.5)
			if len(lines) > maxLines {
				maxLines = len(lines)
			}
		}
		rowHeight := float64(maxLines) * lineH
		if rowHeight < 6 {
			rowHeight = 6
		}

		if pdf.GetY()+rowHeight > pageHeight-marginB-40 {
			pdf.AddPage()
			drawNomorUrutHeader()
			pdf.SetFont("Arial", "", dataFontSize)
		}

		x := marginL
		y := pdf.GetY()
		for j, v := range values {
			pdf.Rect(x, y, widths[j], rowHeight, "D")
			align := "C"
			if j == 3 { // Nama Peserta Didik rata kiri
				align = "L"
			}
			cellLines := pdf.SplitLines([]byte(v), widths[j]-1.5)
			nCellLines := len(cellLines)
			if nCellLines < 1 {
				nCellLines = 1
			}
			textY := y + (rowHeight-float64(nCellLines)*lineH)/2
			pdf.SetXY(x, textY)
			pdf.MultiCell(widths[j], lineH, v, "", align, false)
			x += widths[j]
		}
		pdf.SetXY(marginL, y+rowHeight)
	}

	renderSignatureBlock3(pdf, kepsek, waliKelasNama, waliKelasNIP, namaPetugas, nipPetugas, kotaSekolah)
}

func jenisKelaminSingkat(jk string) string {
	if len(jk) > 0 {
		return string(jk[0])
	}
	return ""
}
