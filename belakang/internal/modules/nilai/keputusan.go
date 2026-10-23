package nilai

import (
	"strings"
	"time"

	"sekolah/internal/database"
)

const (
	StatusLulus      = "Lulus"
	StatusTidakLulus = "Tidak Lulus"

	maxNoDokumen = 50
)

func isStatusKeputusanValid(s string) bool {
	return s == "" || s == StatusLulus || s == StatusTidakLulus
}

// KeputusanSiswa = satu baris layar input keputusan akhir (halaman 2 Data Nilai): 1 siswa = 1 keputusan.
// Tanggal berformat yyyy-mm-dd ("" kalau kosong).
type KeputusanSiswa struct {
	Nis        string `db:"nis" json:"nis"`
	NamaSiswa  string `db:"nama_siswa" json:"nama_siswa"`
	NamaKelas  string `db:"nama_kelas" json:"nama_kelas"`
	NamaRombel string `db:"nama_rombel" json:"nama_rombel"`
	Status     string `db:"status" json:"status"`
	NoIjazah   string `db:"no_ijazah" json:"no_ijazah"`
	NoSKHUS    string `db:"no_skhus" json:"no_skhus"`
	Tanggal    string `db:"tanggal" json:"tanggal"`
}

type KeputusanBulkItem struct {
	Nis      string `json:"nis"`
	Status   string `json:"status"`
	NoIjazah string `json:"no_ijazah"`
	NoSKHUS  string `json:"no_skhus"`
	Tanggal  string `json:"tanggal"`
}

type KeputusanBulkRequest struct {
	IDTahunPelajaran int                 `json:"id_tahunpelajaran"`
	IDRombel         int                 `json:"id_rombel"`
	Items            []KeputusanBulkItem `json:"items"`
}

// GetDaftarKeputusan: siswa yg belum punya keputusan (a.status NULL) diisi default StatusLulus,
// jadi admin di layar Keputusan Akhir hanya perlu mengubah yg Tidak Lulus (efisiensi input -
// mayoritas siswa memang lulus). Disimpan sungguhan sbg baris "Lulus" begitu admin klik Simpan.
func GetDaftarKeputusan(idRombel int) ([]KeputusanSiswa, error) {
	list := []KeputusanSiswa{}

	err := database.DB.Select(&list, `
		SELECT
			p.nis AS nis,
			p.nama_siswa AS nama_siswa,
			COALESCE(k.nama_kelas, '') AS nama_kelas,
			COALESCE(r.nama_rombel, '') AS nama_rombel,
			COALESCE(a.status, ?) AS status,
			COALESCE(a.no_ijazah, '') AS no_ijazah,
			COALESCE(a.no_skhus, '') AS no_skhus,
			COALESCE(DATE_FORMAT(a.tanggal, '%Y-%m-%d'), '') AS tanggal
		FROM pesertadidik p
		LEFT JOIN keputusan_akhir a ON a.nis = p.nis
		LEFT JOIN tbl_kelas k ON k.id_kelas = p.id_kelas
		LEFT JOIN tbl_rombel r ON r.id_rombel = p.id_rombel
		WHERE p.id_rombel = ?
		ORDER BY p.nama_siswa ASC
	`, StatusLulus, idRombel)

	return list, err
}

func isEmptyKeputusan(it KeputusanBulkItem) bool {
	return it.Status == "" && strings.TrimSpace(it.NoIjazah) == "" &&
		strings.TrimSpace(it.NoSKHUS) == "" && strings.TrimSpace(it.Tanggal) == ""
}

// parseTanggalKeputusan menerima yyyy-mm-dd, dd-mm-yyyy, atau dd/mm/yyyy; hasil yyyy-mm-dd.
func parseTanggalKeputusan(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", true
	}

	for _, layout := range []string{"2006-01-02", "02-01-2006", "2-1-2006", "02/01/2006", "2/1/2006"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.Format("2006-01-02"), true
		}
	}

	return "", false
}

// BulkSaveKeputusan: baris berisi -> upsert (unik per NIS); baris kosong -> hapus keputusan lama siswa itu.
func BulkSaveKeputusan(idTP int, items []KeputusanBulkItem) (saved int, cleared int, err error) {
	tx, err := database.DB.Beginx()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	for _, it := range items {
		if isEmptyKeputusan(it) {
			if _, err := tx.Exec(`DELETE FROM keputusan_akhir WHERE nis = ?`, it.Nis); err != nil {
				return 0, 0, err
			}

			cleared++

			continue
		}

		if _, err := tx.Exec(
			`INSERT INTO keputusan_akhir (nis, id_tahunpelajaran, status, no_ijazah, no_skhus, tanggal)
			VALUES (?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				id_tahunpelajaran = VALUES(id_tahunpelajaran),
				status = VALUES(status),
				no_ijazah = VALUES(no_ijazah),
				no_skhus = VALUES(no_skhus),
				tanggal = VALUES(tanggal)`,
			it.Nis, idTP, nullStr(it.Status), nullStr(it.NoIjazah), nullStr(it.NoSKHUS), nullStr(it.Tanggal),
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

// Keputusan = isi kotak Keputusan halaman 2 PDF (tanggal sudah dd-mm-yyyy).
type Keputusan struct {
	Status   string
	NoIjazah string
	NoSKHUS  string
	Tanggal  string
}

// GetKeputusan: nil kalau siswa belum punya keputusan.
func GetKeputusan(nis string) (*Keputusan, error) {
	rows := []struct {
		Status   string `db:"status"`
		NoIjazah string `db:"no_ijazah"`
		NoSKHUS  string `db:"no_skhus"`
		Tanggal  string `db:"tanggal"`
	}{}

	err := database.DB.Select(&rows, `
		SELECT COALESCE(status, '') AS status,
			COALESCE(no_ijazah, '') AS no_ijazah,
			COALESCE(no_skhus, '') AS no_skhus,
			COALESCE(DATE_FORMAT(tanggal, '%d-%m-%Y'), '') AS tanggal
		FROM keputusan_akhir WHERE nis = ? LIMIT 1`, nis)
	if err != nil || len(rows) == 0 {
		return nil, err
	}

	r := rows[0]

	return &Keputusan{Status: r.Status, NoIjazah: r.NoIjazah, NoSKHUS: r.NoSKHUS, Tanggal: r.Tanggal}, nil
}
