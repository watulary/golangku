package nilai

import "sekolah/internal/database"

// GetTahunPerTingkat = tahun pelajaran (nama, mis. "2024/2025") tempat 1 siswa punya nilai mapel,
// dikelompokkan per tingkat kelas 7 / 8 / 9 (dibaca dari nilai_mapel.id_kelas). Kalau siswa tinggal
// kelas (tingkat sama di dua tahun), tahun yang lebih baru dipakai. Tingkat tanpa nilai tidak ada di map.
func GetTahunPerTingkat(nis string) (map[int]string, error) {
	rows := []struct {
		NamaTahun string `db:"nama_tahunpelajaran"`
		NamaKelas string `db:"nama_kelas"`
	}{}

	err := database.DB.Select(&rows, `
		SELECT DISTINCT tp.nama_tahunpelajaran AS nama_tahunpelajaran,
			COALESCE(k.nama_kelas, '') AS nama_kelas
		FROM nilai_mapel n
		JOIN tbl_tahunpelajaran tp ON tp.id_tahunpelajaran = n.id_tahunpelajaran
		LEFT JOIN tbl_kelas k ON k.id_kelas = n.id_kelas
		WHERE n.nis = ?
		ORDER BY tp.nama_tahunpelajaran ASC
	`, nis)
	if err != nil {
		return nil, err
	}

	hasil := map[int]string{}

	for _, r := range rows {
		if t := TingkatKelas(r.NamaKelas); t != 0 {
			hasil[t] = r.NamaTahun // urut naik: tahun terbaru menimpa
		}
	}

	return hasil, nil
}
