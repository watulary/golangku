package nilai

import (
	"strings"

	"sekolah/internal/database"
)

// GetEkskulAktif = daftar ekskul aktif untuk dropdown.
func GetEkskulAktif() ([]EkskulOpsi, error) {
	list := []EkskulOpsi{}

	err := database.DB.Select(&list,
		`SELECT id_ekskul, nama_ekskul FROM tbl_ekskul WHERE aktif = 1 ORDER BY nama_ekskul ASC`)

	return list, err
}

// GetNamaEkskul: sql.ErrNoRows kalau ekskul tidak ada / nonaktif.
func GetNamaEkskul(idEkskul int) (string, error) {
	var nama string

	err := database.DB.Get(&nama,
		`SELECT nama_ekskul FROM tbl_ekskul WHERE id_ekskul = ? AND aktif = 1 LIMIT 1`, idEkskul)

	return nama, err
}

// GetDaftarEkskul = siswa 1 rombel yang menjadi anggota ekskul itu (tbl_anggota_ekskul, keanggotaan berjalan),
// ditambah siswa yang tercatat sebagai anggota ekskul itu pada periode (tahun pelajaran + semester) ini
// (ekskul_anggota), ditambah siswa yang sudah punya nilai ekskul itu di periode ini, beserta nilainya.
func GetDaftarEkskul(idTP int, semester string, idRombel, idEkskul int) ([]NilaiEkskulSiswa, error) {
	list := []NilaiEkskulSiswa{}

	err := database.DB.Select(&list, `
		SELECT
			IFNULL(n.id_nilai_ekskul, 0) AS id_nilai_ekskul,
			p.nis AS nis,
			p.nama_siswa AS nama_siswa,
			COALESCE(r.nama_rombel, '') AS nama_rombel,
			COALESCE(n.nilai, '') AS nilai,
			COALESCE(n.keterangan, '') AS keterangan
		FROM pesertadidik p
		LEFT JOIN nilai_ekskul n
			ON n.nis = p.nis
			AND n.id_tahunpelajaran = ?
			AND n.semester = ?
			AND n.id_ekskul = ?
		LEFT JOIN tbl_rombel r ON r.id_rombel = p.id_rombel
		WHERE p.id_rombel = ?
		AND (
			n.id_nilai_ekskul IS NOT NULL
			OR EXISTS (SELECT 1 FROM tbl_anggota_ekskul a WHERE a.id_ekskul = ? AND a.nis = p.nis)
			OR EXISTS (
				SELECT 1 FROM ekskul_anggota ea
				WHERE ea.id_ekskul = ? AND ea.nis = p.nis
				AND ea.id_tahunpelajaran = ? AND ea.semester = ?
			)
		)
		ORDER BY p.nama_siswa ASC
	`, idTP, semester, idEkskul, idRombel, idEkskul, idEkskul, idTP, semester)

	return list, err
}

func isEmptyEkskulItem(item EkskulBulkItem) bool {
	return strings.TrimSpace(item.Nilai) == "" && strings.TrimSpace(item.Keterangan) == ""
}

// BulkSaveEkskul menyimpan nilai ekskul banyak siswa dalam 1 transaksi.
// Baris kosong (nilai & keterangan) menghapus nilai lama siswa itu.
func BulkSaveEkskul(idTP int, semester string, idKelas, idEkskul int, items []EkskulBulkItem) (saved int, cleared int, err error) {
	tx, err := database.DB.Beginx()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	for _, item := range items {
		if isEmptyEkskulItem(item) {
			if _, err := tx.Exec(
				`DELETE FROM nilai_ekskul
				WHERE nis = ? AND id_tahunpelajaran = ? AND semester = ? AND id_ekskul = ?`,
				item.Nis, idTP, semester, idEkskul,
			); err != nil {
				return 0, 0, err
			}

			cleared++
			continue
		}

		if _, err := tx.Exec(
			`INSERT INTO nilai_ekskul
				(nis, id_tahunpelajaran, semester, id_kelas, id_ekskul, nilai, keterangan)
			VALUES
				(?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				id_kelas = VALUES(id_kelas),
				nilai = VALUES(nilai),
				keterangan = VALUES(keterangan)`,
			item.Nis, idTP, semester, idKelas, idEkskul,
			nullStr(item.Nilai), nullStr(item.Keterangan),
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
