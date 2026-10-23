package nilai

import (
	"strings"

	"sekolah/internal/database"
)

// GetDaftar mengambil siswa 1 rombel + nilainya utk (tahun pelajaran, semester, mapel).
// Siswa yang belum punya nilai tetap tampil (LEFT JOIN) dengan nilai kosong.
func GetDaftar(idTP int, semester string, idRombel, idMapel int) ([]NilaiSiswa, error) {
	list := []NilaiSiswa{}

	query := `
		SELECT
			IFNULL(n.id_nilai, 0) AS id_nilai,
			p.nis AS nis,
			p.nama_siswa AS nama_siswa,
			COALESCE(k.nama_kelas, '') AS nama_kelas,
			COALESCE(r.nama_rombel, '') AS nama_rombel,
			n.nilai_akhir AS nilai_akhir,
			COALESCE(n.capaian_pembelajaran, '') AS capaian_pembelajaran,
			COALESCE(n.kktp, '') AS kktp
		FROM pesertadidik p
		LEFT JOIN nilai_mapel n
			ON n.nis = p.nis
			AND n.id_tahunpelajaran = ?
			AND n.semester = ?
			AND n.id_mapel = ?
		LEFT JOIN tbl_kelas k ON k.id_kelas = p.id_kelas
		LEFT JOIN tbl_rombel r ON r.id_rombel = p.id_rombel
		WHERE p.id_rombel = ?
		ORDER BY p.nama_siswa ASC
	`

	if err := database.DB.Select(&list, query, idTP, semester, idMapel, idRombel); err != nil {
		return nil, err
	}

	return list, nil
}

// GetRombel: sql.ErrNoRows kalau rombel tidak ada.
func GetRombel(idRombel int) (*RombelInfo, error) {
	var info RombelInfo

	err := database.DB.Get(
		&info,
		`SELECT
			r.id_kelas AS id_kelas,
			r.nama_rombel AS nama_rombel,
			COALESCE(k.nama_kelas, '') AS nama_kelas
		FROM tbl_rombel r
		LEFT JOIN tbl_kelas k ON k.id_kelas = r.id_kelas
		WHERE r.id_rombel = ?
		LIMIT 1`,
		idRombel,
	)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// GetNISRombel = himpunan NIS siswa aktif di 1 rombel (utk memastikan yang disimpan memang anggota rombel itu).
func GetNISRombel(idRombel int) (map[string]bool, error) {
	var list []string

	if err := database.DB.Select(&list, `SELECT nis FROM pesertadidik WHERE id_rombel = ?`, idRombel); err != nil {
		return nil, err
	}

	result := make(map[string]bool, len(list))
	for _, nis := range list {
		result[nis] = true
	}

	return result, nil
}

// GetNamaTahunPelajaran: sql.ErrNoRows kalau tahun pelajaran tidak ada.
func GetNamaTahunPelajaran(id int) (string, error) {
	var nama string

	err := database.DB.Get(
		&nama,
		`SELECT COALESCE(nama_tahunpelajaran, '') FROM tbl_tahunpelajaran WHERE id_tahunpelajaran = ? LIMIT 1`,
		id,
	)

	return nama, err
}

// GetNamaMapelDinilai mengembalikan nama mapel (sub-mapel ditulis "Induk - Sub")
// HANYA kalau mapel itu aktif dan tidak punya sub-mapel (bisa dinilai).
// sql.ErrNoRows kalau mapel tidak ada / nonaktif / merupakan mapel induk.
func GetNamaMapelDinilai(idMapel int) (string, error) {
	var nama string

	err := database.DB.Get(
		&nama,
		`SELECT
			CASE WHEN p.nama_mapel IS NULL
				THEN m.nama_mapel
				ELSE CONCAT(p.nama_mapel, ' - ', m.nama_mapel)
			END
		FROM tbl_mapel m
		LEFT JOIN tbl_mapel p ON p.id_mapel = m.id_induk
		WHERE m.id_mapel = ?
		AND m.aktif = 1
		AND NOT EXISTS (SELECT 1 FROM tbl_mapel c WHERE c.id_induk = m.id_mapel)
		LIMIT 1`,
		idMapel,
	)

	return nama, err
}

func nullStr(s string) interface{} {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	return s
}

func isEmptyItem(item BulkItem) bool {
	return item.NilaiAkhir == nil &&
		strings.TrimSpace(item.Capaian) == "" &&
		strings.TrimSpace(item.Kktp) == ""
}

// BulkSave menyimpan nilai banyak siswa dalam 1 transaksi.
//   - baris berisi  -> INSERT ... ON DUPLICATE KEY UPDATE (unik: nis+tahun+semester+mapel)
//   - baris kosong  -> hapus nilai lama siswa itu (kalau ada), supaya "dikosongkan" benar-benar kosong
func BulkSave(idTP int, semester string, idKelas, idMapel int, items []BulkItem) (saved int, cleared int, err error) {
	tx, err := database.DB.Beginx()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	for _, item := range items {
		if isEmptyItem(item) {
			if _, err := tx.Exec(
				`DELETE FROM nilai_mapel
				WHERE nis = ? AND id_tahunpelajaran = ? AND semester = ? AND id_mapel = ?`,
				item.Nis, idTP, semester, idMapel,
			); err != nil {
				return 0, 0, err
			}

			cleared++
			continue
		}

		if _, err := tx.Exec(
			`INSERT INTO nilai_mapel
				(nis, id_tahunpelajaran, semester, id_kelas, id_mapel, nilai_akhir, capaian_pembelajaran, kktp)
			VALUES
				(?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				id_kelas = VALUES(id_kelas),
				nilai_akhir = VALUES(nilai_akhir),
				capaian_pembelajaran = VALUES(capaian_pembelajaran),
				kktp = VALUES(kktp)`,
			item.Nis, idTP, semester, idKelas, idMapel,
			item.NilaiAkhir, nullStr(item.Capaian), nullStr(item.Kktp),
		); err != nil {
			return 0, 0, err
		}

		saved++
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, err
	}

	return saved, cleared, nil
}
