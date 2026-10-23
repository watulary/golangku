package nilai

import (
	"database/sql"
	"fmt"

	"sekolah/internal/database"
)

// SelNilai = isi 1 sel nilai (1 mapel, 1 semester) di halaman Data Nilai Buku Induk.
type SelNilai struct {
	Akhir   *float64
	Capaian string
	Kktp    string
}

// MapelPohon = 1 baris mapel di PDF. Subs terisi kalau mapel itu punya sub-mapel
// (baris induk tidak punya nilai sendiri).
type MapelPohon struct {
	ID   int
	No   string
	Nama string
	Subs []MapelPohon
}

// DataBukuInduk = semua data yang dibutuhkan halaman Data Nilai Buku Induk
// untuk 1 siswa pada 1 tahun pelajaran (semester Ganjil + Genap).
type DataBukuInduk struct {
	NamaKelas string // dari nilai_mapel.id_kelas; "" kalau belum ada nilai

	Mapel []MapelPohon
	Nilai map[int]map[string]SelNilai // id_mapel -> semester -> sel

	Jumlah map[string]float64 // semester -> jumlah nilai (hanya ada kalau semester itu punya nilai)
	Rata   map[string]float64 // semester -> rata-rata

	Ketidakhadiran map[string]Ketidakhadiran // semester -> rekap dari absensi
}

// GetPohonMapel = mapel aktif berbentuk pohon (mapel utama + sub-mapelnya), nomor urut 1..n untuk mapel utama.
func GetPohonMapel() ([]MapelPohon, error) {
	rows := []struct {
		ID      int    `db:"id_mapel"`
		IDInduk *int   `db:"id_induk"`
		Nama    string `db:"nama_mapel"`
	}{}

	err := database.DB.Select(&rows, `
		SELECT m.id_mapel, m.id_induk, m.nama_mapel
		FROM tbl_mapel m
		LEFT JOIN tbl_mapel p ON p.id_mapel = m.id_induk
		WHERE m.aktif = 1
		ORDER BY COALESCE(p.urutan, m.urutan), COALESCE(m.id_induk, m.id_mapel), (m.id_induk IS NOT NULL), m.urutan, m.id_mapel
	`)
	if err != nil {
		return nil, err
	}

	pohon := []MapelPohon{}
	index := map[int]int{} // id_mapel induk -> posisi di pohon

	for _, r := range rows {
		if r.IDInduk == nil {
			index[r.ID] = len(pohon)
			pohon = append(pohon, MapelPohon{ID: r.ID, No: fmt.Sprintf("%d", len(pohon)+1), Nama: r.Nama})
			continue
		}

		if pos, ok := index[*r.IDInduk]; ok {
			pohon[pos].Subs = append(pohon[pos].Subs, MapelPohon{ID: r.ID, Nama: r.Nama})
		}
	}

	return pohon, nil
}

// GetDataBukuInduk mengumpulkan nilai + ketidakhadiran 1 siswa utk tahun pelajaran (nama, mis. "2024/2025").
// Tahun pelajaran tidak ditemukan / belum ada nilai bukan error: hasilnya tetap berisi pohon mapel dan nilai kosong.
func GetDataBukuInduk(nis, namaTahun string) (*DataBukuInduk, error) {
	pohon, err := GetPohonMapel()
	if err != nil {
		return nil, err
	}

	data := &DataBukuInduk{
		Mapel:          pohon,
		Nilai:          map[int]map[string]SelNilai{},
		Jumlah:         map[string]float64{},
		Rata:           map[string]float64{},
		Ketidakhadiran: map[string]Ketidakhadiran{},
	}

	var idTP int
	err = database.DB.Get(&idTP, `SELECT id_tahunpelajaran FROM tbl_tahunpelajaran WHERE nama_tahunpelajaran = ? LIMIT 1`, namaTahun)
	if err == sql.ErrNoRows {
		return data, nil
	}
	if err != nil {
		return nil, err
	}

	rows := []struct {
		IDMapel  int      `db:"id_mapel"`
		Semester string   `db:"semester"`
		Akhir    *float64 `db:"nilai_akhir"`
		Capaian  string   `db:"capaian_pembelajaran"`
		Kktp     string   `db:"kktp"`
	}{}

	err = database.DB.Select(&rows, `
		SELECT id_mapel, semester, nilai_akhir,
			COALESCE(capaian_pembelajaran, '') AS capaian_pembelajaran,
			COALESCE(kktp, '') AS kktp
		FROM nilai_mapel
		WHERE nis = ? AND id_tahunpelajaran = ?
	`, nis, idTP)
	if err != nil {
		return nil, err
	}

	for _, r := range rows {
		if data.Nilai[r.IDMapel] == nil {
			data.Nilai[r.IDMapel] = map[string]SelNilai{}
		}

		data.Nilai[r.IDMapel][r.Semester] = SelNilai{Akhir: r.Akhir, Capaian: r.Capaian, Kktp: r.Kktp}
	}

	// nama kelas saat nilai diinput (nilai_mapel.id_kelas)
	var namaKelas string
	err = database.DB.Get(&namaKelas, `
		SELECT COALESCE(k.nama_kelas, '')
		FROM nilai_mapel n
		JOIN tbl_kelas k ON k.id_kelas = n.id_kelas
		WHERE n.nis = ? AND n.id_tahunpelajaran = ?
		LIMIT 1
	`, nis, idTP)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	data.NamaKelas = namaKelas

	// Jumlah & rata-rata per semester. Mapel dengan sub-mapel dihitung SATU kali
	// (rata-rata sub-mapelnya yang terisi), sama seperti rekap nilai.
	for _, sem := range []string{"Ganjil", "Genap"} {
		jumlah, n := 0.0, 0

		for _, m := range pohon {
			if len(m.Subs) == 0 {
				if s, ok := data.Nilai[m.ID][sem]; ok && s.Akhir != nil {
					jumlah += *s.Akhir
					n++
				}

				continue
			}

			total, cnt := 0.0, 0
			for _, sub := range m.Subs {
				if s, ok := data.Nilai[sub.ID][sem]; ok && s.Akhir != nil {
					total += *s.Akhir
					cnt++
				}
			}

			if cnt > 0 {
				jumlah += total / float64(cnt)
				n++
			}
		}

		if n > 0 {
			data.Jumlah[sem] = round2(jumlah)
			data.Rata[sem] = round2(jumlah / float64(n))
		}

		// ketidakhadiran: format nama tahun tak dikenali -> dilewati (kolom tetap kosong)
		if kh, errKh := GetKetidakhadiran(nis, namaTahun, sem); errKh == nil {
			data.Ketidakhadiran[sem] = kh
		}
	}

	return data, nil
}
