package nilai

import (
	"strings"

	"sekolah/internal/database"
)

// EkskulSel = isi 1 sel (1 ekskul, 1 tingkat kelas, 1 semester) di tabel Ekstrakurikuler Buku Induk.
type EkskulSel struct {
	Nilai      string
	Keterangan string
}

// EkskulRiwayat = riwayat 1 ekskul milik 1 siswa untuk kelas 7, 8, 9 (key = tingkat) x semester.
type EkskulRiwayat struct {
	Nama string
	Sel  map[int]map[string]EkskulSel
}

// TingkatKelas menerjemahkan nama kelas ("7", "VII", "7 (Tujuh)", "Kelas 8", "IX") menjadi 7 / 8 / 9.
// 0 = tidak dikenali.
func TingkatKelas(nama string) int {
	n := strings.ToUpper(strings.TrimSpace(nama))

	// cek VIII sebelum VII (VII adalah awalan VIII)
	switch {
	case strings.Contains(n, "VIII"):
		return 8
	case strings.Contains(n, "VII"):
		return 7
	case strings.Contains(n, "IX"):
		return 9
	}

	for _, r := range n {
		switch r {
		case '7':
			return 7
		case '8':
			return 8
		case '9':
			return 9
		}
	}

	return 0
}

// GetRiwayatEkskul = seluruh nilai ekskul 1 siswa (semua tahun pelajaran), dikelompokkan per ekskul.
// Kelas dibaca dari nilai_ekskul.id_kelas, jadi tidak tergantung tahun pelajaran yang dicetak.
// Kalau siswa tinggal kelas (tingkat sama, tahun berbeda), tahun yang lebih baru yang dipakai.
func GetRiwayatEkskul(nis string) ([]EkskulRiwayat, error) {
	rows := []struct {
		Nama       string `db:"nama_ekskul"`
		NamaKelas  string `db:"nama_kelas"`
		Semester   string `db:"semester"`
		Nilai      string `db:"nilai"`
		Keterangan string `db:"keterangan"`
	}{}

	err := database.DB.Select(&rows, `
		SELECT e.nama_ekskul AS nama_ekskul,
			COALESCE(k.nama_kelas, '') AS nama_kelas,
			n.semester AS semester,
			COALESCE(n.nilai, '') AS nilai,
			COALESCE(n.keterangan, '') AS keterangan
		FROM nilai_ekskul n
		JOIN tbl_ekskul e ON e.id_ekskul = n.id_ekskul
		LEFT JOIN tbl_kelas k ON k.id_kelas = n.id_kelas
		WHERE n.nis = ?
		ORDER BY e.nama_ekskul ASC, n.id_tahunpelajaran ASC
	`, nis)
	if err != nil {
		return nil, err
	}

	hasil := []EkskulRiwayat{}
	index := map[string]int{}

	for _, r := range rows {
		tingkat := TingkatKelas(r.NamaKelas)
		if tingkat == 0 {
			continue
		}

		pos, ok := index[r.Nama]
		if !ok {
			pos = len(hasil)
			index[r.Nama] = pos
			hasil = append(hasil, EkskulRiwayat{Nama: r.Nama, Sel: map[int]map[string]EkskulSel{}})
		}

		if hasil[pos].Sel[tingkat] == nil {
			hasil[pos].Sel[tingkat] = map[string]EkskulSel{}
		}

		hasil[pos].Sel[tingkat][r.Semester] = EkskulSel{Nilai: r.Nilai, Keterangan: r.Keterangan}
	}

	return hasil, nil
}
