package nilai

import (
	"strings"

	"sekolah/internal/database"
)

// Dimensi Profil Pelajar Pancasila (urutan sama dengan tabel A di PDF, 1..6).
var DimensiP5 = []string{
	"Beriman, bertakwa kepada Tuhan Yang Maha Esa dan berakhlak mulia",
	"Berkebinekaan Global",
	"Bergotong-royong",
	"Mandiri",
	"Bernalar Kritis",
	"Kreatif",
}

const (
	maxElemenP5 = 255
	maxTargetP5 = 500
)

// isNilaiP5Valid: BB / MB / BSH / SB (atau kosong).
func isNilaiP5Valid(n string) bool {
	return n == "" || n == "BB" || n == "MB" || n == "BSH" || n == "SB"
}

func isDimensiP5Valid(d int) bool { return d >= 1 && d <= len(DimensiP5) }

// P5Siswa = satu baris layar input: siswa 1 rombel + isian P5 untuk 1 dimensi / semester / tahun.
type P5Siswa struct {
	Nis        string `db:"nis" json:"nis"`
	NamaSiswa  string `db:"nama_siswa" json:"nama_siswa"`
	NamaKelas  string `db:"nama_kelas" json:"nama_kelas"`
	NamaRombel string `db:"nama_rombel" json:"nama_rombel"`
	Nilai      string `db:"nilai" json:"nilai"`
	Elemen     string `db:"elemen" json:"elemen"`
	SubElemen  string `db:"sub_elemen" json:"sub_elemen"`
	Target     string `db:"target" json:"target"`
}

type P5BulkItem struct {
	Nis       string `json:"nis"`
	Nilai     string `json:"nilai"`
	Elemen    string `json:"elemen"`
	SubElemen string `json:"sub_elemen"`
	Target    string `json:"target"`
}

type P5BulkRequest struct {
	IDTahunPelajaran int          `json:"id_tahunpelajaran"`
	Semester         string       `json:"semester"`
	IDRombel         int          `json:"id_rombel"`
	Dimensi          int          `json:"dimensi"`
	IDKelasCatat     int          `json:"id_kelas_catat"`
	Items            []P5BulkItem `json:"items"`
}

func GetDaftarP5(idTP int, semester string, idRombel, dimensi int) ([]P5Siswa, error) {
	list := []P5Siswa{}

	err := database.DB.Select(&list, `
		SELECT
			p.nis AS nis,
			p.nama_siswa AS nama_siswa,
			COALESCE(k.nama_kelas, '') AS nama_kelas,
			COALESCE(r.nama_rombel, '') AS nama_rombel,
			COALESCE(n.nilai, '') AS nilai,
			COALESCE(n.elemen, '') AS elemen,
			COALESCE(n.sub_elemen, '') AS sub_elemen,
			COALESCE(n.target, '') AS target
		FROM pesertadidik p
		LEFT JOIN p5_nilai n
			ON n.nis = p.nis
			AND n.id_tahunpelajaran = ?
			AND n.semester = ?
			AND n.dimensi = ?
		LEFT JOIN tbl_kelas k ON k.id_kelas = p.id_kelas
		LEFT JOIN tbl_rombel r ON r.id_rombel = p.id_rombel
		WHERE p.id_rombel = ?
		ORDER BY p.nama_siswa ASC
	`, idTP, semester, dimensi, idRombel)

	return list, err
}

func isEmptyP5(it P5BulkItem) bool {
	return it.Nilai == "" && strings.TrimSpace(it.Elemen) == "" &&
		strings.TrimSpace(it.SubElemen) == "" && strings.TrimSpace(it.Target) == ""
}

// BulkSaveP5: baris berisi -> upsert (unik nis+tahun+semester+dimensi); baris kosong -> hapus isian lama.
func BulkSaveP5(idTP int, semester string, idKelas, dimensi int, items []P5BulkItem) (saved int, cleared int, err error) {
	tx, err := database.DB.Beginx()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	for _, it := range items {
		if isEmptyP5(it) {
			if _, err := tx.Exec(
				`DELETE FROM p5_nilai WHERE nis = ? AND id_tahunpelajaran = ? AND semester = ? AND dimensi = ?`,
				it.Nis, idTP, semester, dimensi,
			); err != nil {
				return 0, 0, err
			}

			cleared++

			continue
		}

		if _, err := tx.Exec(
			`INSERT INTO p5_nilai
				(nis, id_tahunpelajaran, semester, id_kelas, dimensi, nilai, elemen, sub_elemen, target)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				id_kelas = VALUES(id_kelas),
				nilai = VALUES(nilai),
				elemen = VALUES(elemen),
				sub_elemen = VALUES(sub_elemen),
				target = VALUES(target)`,
			it.Nis, idTP, semester, idKelas, dimensi,
			nullStr(it.Nilai), nullStr(it.Elemen), nullStr(it.SubElemen), nullStr(it.Target),
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

// ---------- untuk PDF Buku Induk (tabel A halaman Data Nilai) ----------

// P5Sel = isian 1 dimensi pada 1 semester.
type P5Sel struct {
	Nilai     string
	Elemen    string
	SubElemen string
	Target    string
}

// P5Riwayat[tingkat 7/8/9][dimensi 1..6][semester "Ganjil"/"Genap"].
type P5Riwayat map[int]map[int]map[string]P5Sel

// GetRiwayatP5 = seluruh isian P5 1 siswa, dikelompokkan per tingkat kelas (dari id_kelas yang dicatat).
// Kalau ada beberapa tahun di tingkat yang sama (tinggal kelas), tahun terbaru menimpa.
func GetRiwayatP5(nis string) (P5Riwayat, error) {
	rows := []struct {
		NamaKelas string `db:"nama_kelas"`
		Semester  string `db:"semester"`
		Dimensi   int    `db:"dimensi"`
		Nilai     string `db:"nilai"`
		Elemen    string `db:"elemen"`
		SubElemen string `db:"sub_elemen"`
		Target    string `db:"target"`
	}{}

	err := database.DB.Select(&rows, `
		SELECT COALESCE(k.nama_kelas, '') AS nama_kelas,
			n.semester AS semester,
			n.dimensi AS dimensi,
			COALESCE(n.nilai, '') AS nilai,
			COALESCE(n.elemen, '') AS elemen,
			COALESCE(n.sub_elemen, '') AS sub_elemen,
			COALESCE(n.target, '') AS target
		FROM p5_nilai n
		JOIN tbl_tahunpelajaran tp ON tp.id_tahunpelajaran = n.id_tahunpelajaran
		LEFT JOIN tbl_kelas k ON k.id_kelas = n.id_kelas
		WHERE n.nis = ?
		ORDER BY tp.nama_tahunpelajaran ASC
	`, nis)
	if err != nil {
		return nil, err
	}

	hasil := P5Riwayat{}

	for _, r := range rows {
		t := TingkatKelas(r.NamaKelas)
		if t == 0 {
			continue
		}

		if hasil[t] == nil {
			hasil[t] = map[int]map[string]P5Sel{}
		}

		if hasil[t][r.Dimensi] == nil {
			hasil[t][r.Dimensi] = map[string]P5Sel{}
		}

		hasil[t][r.Dimensi][r.Semester] = P5Sel{Nilai: r.Nilai, Elemen: r.Elemen, SubElemen: r.SubElemen, Target: r.Target}
	}

	return hasil, nil
}
