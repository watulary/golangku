package nilai

import (
	"database/sql"

	"sekolah/internal/database"
)

const (
	StatusNaik    = "Naik"
	StatusTinggal = "Tinggal"
)

func isStatusNaikKelasValid(s string) bool {
	return s == "" || s == StatusNaik || s == StatusTinggal
}

// NaikKelasSiswa = satu baris layar input Keputusan Naik Kelas (Kelas 7 & 8): 1 siswa = 1 keputusan.
type NaikKelasSiswa struct {
	Nis        string `db:"nis" json:"nis"`
	NamaSiswa  string `db:"nama_siswa" json:"nama_siswa"`
	NamaKelas  string `db:"nama_kelas" json:"nama_kelas"`
	NamaRombel string `db:"nama_rombel" json:"nama_rombel"`
	Status     string `db:"status" json:"status"`
}

type NaikKelasBulkItem struct {
	Nis    string `json:"nis"`
	Status string `json:"status"`
}

type NaikKelasBulkRequest struct {
	IDTahunPelajaran int                 `json:"id_tahunpelajaran"`
	IDRombel         int                 `json:"id_rombel"`
	Items            []NaikKelasBulkItem `json:"items"`
}

// GetDaftarNaikKelas: siswa yg belum punya keputusan (a.status NULL) diisi default StatusNaik,
// jadi admin di layar Keputusan Naik Kelas hanya perlu mengubah yg Tinggal (efisiensi input -
// mayoritas siswa memang naik). Disimpan sungguhan sbg baris "Naik" begitu admin klik Simpan.
func GetDaftarNaikKelas(idRombel, tingkat int) ([]NaikKelasSiswa, error) {
	list := []NaikKelasSiswa{}

	err := database.DB.Select(&list, `
		SELECT
			p.nis AS nis,
			p.nama_siswa AS nama_siswa,
			COALESCE(k.nama_kelas, '') AS nama_kelas,
			COALESCE(r.nama_rombel, '') AS nama_rombel,
			COALESCE(a.status, ?) AS status
		FROM pesertadidik p
		LEFT JOIN keputusan_naikkelas a ON a.nis = p.nis AND a.tingkat = ?
		LEFT JOIN tbl_kelas k ON k.id_kelas = p.id_kelas
		LEFT JOIN tbl_rombel r ON r.id_rombel = p.id_rombel
		WHERE p.id_rombel = ?
		ORDER BY p.nama_siswa ASC
	`, StatusNaik, tingkat, idRombel)

	return list, err
}

func isEmptyNaikKelas(it NaikKelasBulkItem) bool {
	return it.Status == ""
}

// BulkSaveNaikKelas: baris berisi -> upsert (unik per NIS+tingkat); baris kosong -> hapus
// keputusan lama siswa itu utk tingkat ini saja (tingkat lain tidak terpengaruh).
func BulkSaveNaikKelas(idTP, tingkat int, items []NaikKelasBulkItem) (saved int, cleared int, err error) {
	tx, err := database.DB.Beginx()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	for _, it := range items {
		if isEmptyNaikKelas(it) {
			if _, err := tx.Exec(
				`DELETE FROM keputusan_naikkelas WHERE nis = ? AND tingkat = ?`, it.Nis, tingkat,
			); err != nil {
				return 0, 0, err
			}

			cleared++

			continue
		}

		if _, err := tx.Exec(
			`INSERT INTO keputusan_naikkelas (nis, id_tahunpelajaran, tingkat, status)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				id_tahunpelajaran = VALUES(id_tahunpelajaran),
				status = VALUES(status)`,
			it.Nis, idTP, tingkat, it.Status,
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

// GetKeputusanNaikKelas: status Naik/Tinggal siswa utk 1 tingkat ("" = belum diisi). Dipakai PDF.
func GetKeputusanNaikKelas(nis string, tingkat int) (string, error) {
	var status string

	err := database.DB.Get(
		&status,
		`SELECT COALESCE(status, '') FROM keputusan_naikkelas WHERE nis = ? AND tingkat = ? LIMIT 1`,
		nis, tingkat,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}

		return "", err
	}

	return status, nil
}
