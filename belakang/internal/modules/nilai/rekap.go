package nilai

import (
	"fmt"
	"math"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"

	"sekolah/internal/database"
)

// kolomRekap = satu mapel yang bisa dinilai (aktif, tanpa sub-mapel) = satu kolom di rekap.
type kolomRekap struct {
	ID        int     `db:"id_mapel"`
	IDInduk   *int    `db:"id_induk"`
	Nama      string  `db:"nama_mapel"`
	NamaInduk *string `db:"nama_induk"`
}

type siswaRekap struct {
	Nis  string `db:"nis"`
	Nama string `db:"nama_siswa"`
}

type nilaiRekap struct {
	Nis     string  `db:"nis"`
	IDMapel int     `db:"id_mapel"`
	Nilai   float64 `db:"nilai_akhir"`
}

func getKolomRekap() ([]kolomRekap, error) {
	list := []kolomRekap{}

	err := database.DB.Select(&list, `
		SELECT m.id_mapel, m.id_induk, m.nama_mapel, p.nama_mapel AS nama_induk
		FROM tbl_mapel m
		LEFT JOIN tbl_mapel p ON p.id_mapel = m.id_induk
		WHERE m.aktif = 1
		AND NOT EXISTS (SELECT 1 FROM tbl_mapel c WHERE c.id_induk = m.id_mapel)
		ORDER BY COALESCE(p.urutan, m.urutan), COALESCE(m.id_induk, m.id_mapel), (m.id_induk IS NOT NULL), m.urutan, m.id_mapel
	`)

	return list, err
}

func getSiswaRekap(idRombel int) ([]siswaRekap, error) {
	list := []siswaRekap{}

	err := database.DB.Select(&list,
		`SELECT nis, nama_siswa FROM pesertadidik WHERE id_rombel = ? ORDER BY nama_siswa ASC`,
		idRombel,
	)

	return list, err
}

func getNilaiRekap(idTP int, semester string, idRombel int) ([]nilaiRekap, error) {
	list := []nilaiRekap{}

	err := database.DB.Select(&list, `
		SELECT n.nis, n.id_mapel, n.nilai_akhir
		FROM nilai_mapel n
		JOIN pesertadidik p ON p.nis = n.nis
		WHERE n.id_tahunpelajaran = ? AND n.semester = ? AND p.id_rombel = ?
		AND n.nilai_akhir IS NOT NULL
	`, idTP, semester, idRombel)

	return list, err
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

// RekapHandler = unduh rekap nilai akhir 1 rombel: siswa x mapel, + Jumlah & Rata-rata.
// Mapel yang punya sub-mapel (mis. Seni Budaya & Prakarya) dihitung SATU kali dalam
// Jumlah/Rata-rata, yaitu rata-rata nilai sub-mapelnya yang terisi.
func RekapHandler(c *fiber.Ctx) error {
	idTP := c.QueryInt("id_tahunpelajaran", 0)
	semester := c.Query("semester")
	idRombel := c.QueryInt("id_rombel", 0)

	if idTP <= 0 {
		return errorJSON(c, fiber.StatusBadRequest, "Tahun pelajaran wajib dipilih.")
	}
	if !isSemesterValid(semester) {
		return errorJSON(c, fiber.StatusBadRequest, "Semester wajib dipilih: Ganjil atau Genap.")
	}
	if idRombel <= 0 {
		return errorJSON(c, fiber.StatusBadRequest, "Rombel wajib dipilih.")
	}

	namaTahun, err := GetNamaTahunPelajaran(idTP)
	if err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "Tahun pelajaran tidak ditemukan.")
	}

	rombel, err := GetRombel(idRombel)
	if err != nil {
		return errorJSON(c, fiber.StatusBadRequest, "Rombel tidak ditemukan.")
	}

	kolom, err := getKolomRekap()
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	siswa, err := getSiswaRekap(idRombel)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	rows, err := getNilaiRekap(idTP, semester, idRombel)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	hadir, err := GetKetidakhadiranRombel(idRombel, namaTahun, semester)
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	nilai := map[string]map[int]float64{}
	for _, r := range rows {
		if nilai[r.Nis] == nil {
			nilai[r.Nis] = map[int]float64{}
		}
		nilai[r.Nis][r.IDMapel] = r.Nilai
	}

	// kunci grup: mapel induk (kalau sub-mapel) atau mapel itu sendiri
	groupOf := func(k kolomRekap) int {
		if k.IDInduk != nil {
			return *k.IDInduk
		}

		return k.ID
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Rekap"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "REKAP NILAI AKHIR")
	f.SetCellValue(sheet, "A2", fmt.Sprintf("Tahun Pelajaran %s - Semester %s", namaTahun, semester))
	f.SetCellValue(sheet, "A3", fmt.Sprintf("Kelas %s - Rombel %s", rombel.NamaKelas, rombel.NamaRombel))

	const headerRow = 5

	header := []string{"No", "NIS", "Nama Peserta Didik"}
	for _, k := range kolom {
		if k.NamaInduk != nil {
			header = append(header, *k.NamaInduk+" - "+k.Nama)
		} else {
			header = append(header, k.Nama)
		}
	}
	header = append(header, "Jumlah", "Rata-rata", "Sakit", "Izin", "Alpa")

	colMapel := 4 // kolom mapel pertama (D)
	colJumlah := colMapel + len(kolom)
	colRata := colJumlah + 1

	for i, h := range header {
		cell, _ := excelize.CoordinatesToCellName(i+1, headerRow)
		f.SetCellValue(sheet, cell, h)
	}

	sumKolom := make([]float64, len(kolom))
	cntKolom := make([]int, len(kolom))

	for i, s := range siswa {
		row := headerRow + 1 + i

		put := func(col int, v interface{}) {
			cell, _ := excelize.CoordinatesToCellName(col, row)
			f.SetCellValue(sheet, cell, v)
		}

		put(1, i+1)
		put(2, s.Nis)
		put(3, s.Nama)

		perGrup := map[int][]float64{}
		urutGrup := []int{}

		for j, k := range kolom {
			v, ok := nilai[s.Nis][k.ID]
			if !ok {
				continue
			}

			put(colMapel+j, v)
			sumKolom[j] += v
			cntKolom[j]++

			g := groupOf(k)
			if _, seen := perGrup[g]; !seen {
				urutGrup = append(urutGrup, g)
			}
			perGrup[g] = append(perGrup[g], v)
		}

		kh := hadir[s.Nis]
		put(colRata+1, kh.Sakit)
		put(colRata+2, kh.Izin)
		put(colRata+3, kh.Alpa)

		if len(urutGrup) > 0 {
			jumlah := 0.0

			for _, g := range urutGrup {
				total := 0.0
				for _, v := range perGrup[g] {
					total += v
				}
				jumlah += total / float64(len(perGrup[g]))
			}

			put(colJumlah, round2(jumlah))
			put(colRata, round2(jumlah/float64(len(urutGrup))))
		}
	}

	// baris rata-rata kelas per mapel
	if len(siswa) > 0 {
		row := headerRow + 1 + len(siswa)
		cell, _ := excelize.CoordinatesToCellName(3, row)
		f.SetCellValue(sheet, cell, "Rata-rata Kelas")

		for j := range kolom {
			if cntKolom[j] == 0 {
				continue
			}

			cell, _ := excelize.CoordinatesToCellName(colMapel+j, row)
			f.SetCellValue(sheet, cell, round2(sumKolom[j]/float64(cntKolom[j])))
		}
	}

	if headerStyle, errStyle := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"2563EB"}, Pattern: 1},
		Alignment: &excelize.Alignment{Vertical: "center", WrapText: true},
	}); errStyle == nil {
		first, _ := excelize.CoordinatesToCellName(1, headerRow)
		last, _ := excelize.CoordinatesToCellName(len(header), headerRow)
		f.SetCellStyle(sheet, first, last, headerStyle)
	}

	if titleStyle, errStyle := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}}); errStyle == nil {
		f.SetCellStyle(sheet, "A1", "A1", titleStyle)
	}

	f.SetColWidth(sheet, "A", "A", 6)
	f.SetColWidth(sheet, "B", "B", 14)
	f.SetColWidth(sheet, "C", "C", 34)

	first, _ := excelize.CoordinatesToCellName(colMapel, 1)
	last, _ := excelize.CoordinatesToCellName(colRata+3, 1)
	f.SetColWidth(sheet, strings.TrimRight(first, "0123456789"), strings.TrimRight(last, "0123456789"), 16)

	f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      3,
		YSplit:      headerRow,
		TopLeftCell: "D6",
		ActivePane:  "bottomRight",
	})

	buf, err := f.WriteToBuffer()
	if err != nil {
		return errorJSON(c, fiber.StatusInternalServerError, "Gagal membuat rekap nilai.")
	}

	filename := fmt.Sprintf("rekap-nilai-%s-%s.xlsx", strings.ReplaceAll(rombel.NamaRombel, " ", "-"), semester)
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return c.Send(buf.Bytes())
}
