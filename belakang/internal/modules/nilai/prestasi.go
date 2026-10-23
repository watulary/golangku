package nilai

import (
	"strings"

	"sekolah/internal/database"
)

const (
	maxPrestasi           = 255
	maxNilaiPrestasi      = 20
	maxKeteranganPrestasi = 255
)

// SiswaOpsi = siswa 1 rombel (dropdown di layar input prestasi).
type SiswaOpsi struct {
	Nis       string `db:"nis" json:"nis"`
	NamaSiswa string `db:"nama_siswa" json:"nama_siswa"`
}

// PrestasiRow = satu prestasi milik satu siswa pada 1 tahun pelajaran + semester.
// Satu siswa boleh punya banyak prestasi.
type PrestasiRow struct {
	IDPrestasi int    `db:"id_prestasi" json:"id_prestasi"`
	Nis        string `db:"nis" json:"nis"`
	NamaSiswa  string `db:"nama_siswa" json:"nama_siswa"`
	Prestasi   string `db:"prestasi" json:"prestasi"`
	Nilai      string `db:"nilai" json:"nilai"`
	Keterangan string `db:"keterangan" json:"keterangan"`
}

type PrestasiBulkItem struct {
	Nis        string `json:"nis"`
	Prestasi   string `json:"prestasi"`
	Nilai      string `json:"nilai"`
	Keterangan string `json:"keterangan"`
}

type PrestasiBulkRequest struct {
	IDTahunPelajaran int                `json:"id_tahunpelajaran"`
	Semester         string             `json:"semester"`
	IDRombel         int                `json:"id_rombel"`
	IDKelasCatat     int                `json:"id_kelas_catat"`
	Items            []PrestasiBulkItem `json:"items"`
}

func GetSiswaRombel(idRombel int) ([]SiswaOpsi, error) {
	list := []SiswaOpsi{}

	err := database.DB.Select(&list,
		`SELECT nis, nama_siswa FROM pesertadidik WHERE id_rombel = ? ORDER BY nama_siswa ASC`, idRombel)

	return list, err
}

// GetDaftarPrestasi = seluruh prestasi siswa 1 rombel pada 1 tahun pelajaran + semester.
func GetDaftarPrestasi(idTP int, semester string, idRombel int) ([]PrestasiRow, error) {
	list := []PrestasiRow{}

	err := database.DB.Select(&list, `
		SELECT n.id_prestasi AS id_prestasi,
			n.nis AS nis,
			p.nama_siswa AS nama_siswa,
			n.prestasi AS prestasi,
			COALESCE(n.nilai, '') AS nilai,
			COALESCE(n.keterangan, '') AS keterangan
		FROM prestasi_siswa n
		JOIN pesertadidik p ON p.nis = n.nis
		WHERE p.id_rombel = ?
		AND n.id_tahunpelajaran = ?
		AND n.semester = ?
		ORDER BY p.nama_siswa ASC, n.id_prestasi ASC
	`, idRombel, idTP, semester)

	return list, err
}

// SavePrestasi MENGGANTI seluruh prestasi siswa rombel itu pada tahun pelajaran + semester itu
// dengan items (dalam 1 transaksi). items kosong = semua prestasi periode itu dihapus.
func SavePrestasi(idTP int, semester string, idKelas int, nisRombel []string, items []PrestasiBulkItem) (int, error) {
	tx, err := database.DB.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	for _, nis := range nisRombel {
		if _, err := tx.Exec(
			`DELETE FROM prestasi_siswa WHERE nis = ? AND id_tahunpelajaran = ? AND semester = ?`,
			nis, idTP, semester,
		); err != nil {
			return 0, err
		}
	}

	for _, it := range items {
		if _, err := tx.Exec(
			`INSERT INTO prestasi_siswa
				(nis, id_tahunpelajaran, semester, id_kelas, prestasi, nilai, keterangan)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			it.Nis, idTP, semester, idKelas, it.Prestasi, nullStr(it.Nilai), nullStr(it.Keterangan),
		); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return len(items), nil
}

// ---------- untuk PDF Buku Induk ----------

// PrestasiSel = 1 prestasi di tabel Catatan Prestasi Buku Induk.
type PrestasiSel struct {
	Prestasi   string
	Semester   string // Ganjil / Genap
	Nilai      string
	Keterangan string
}

// GetRiwayatPrestasi = seluruh prestasi 1 siswa dikelompokkan per tingkat kelas (7/8/9),
// urut tahun pelajaran, semester (Ganjil dulu), lalu urutan input.
func GetRiwayatPrestasi(nis string) (map[int][]PrestasiSel, error) {
	rows := []struct {
		NamaKelas  string `db:"nama_kelas"`
		Semester   string `db:"semester"`
		Prestasi   string `db:"prestasi"`
		Nilai      string `db:"nilai"`
		Keterangan string `db:"keterangan"`
	}{}

	err := database.DB.Select(&rows, `
		SELECT COALESCE(k.nama_kelas, '') AS nama_kelas,
			n.semester AS semester,
			n.prestasi AS prestasi,
			COALESCE(n.nilai, '') AS nilai,
			COALESCE(n.keterangan, '') AS keterangan
		FROM prestasi_siswa n
		LEFT JOIN tbl_kelas k ON k.id_kelas = n.id_kelas
		WHERE n.nis = ?
		ORDER BY n.id_tahunpelajaran ASC, FIELD(n.semester, 'Ganjil', 'Genap'), n.id_prestasi ASC
	`, nis)
	if err != nil {
		return nil, err
	}

	hasil := map[int][]PrestasiSel{}

	for _, r := range rows {
		tingkat := TingkatKelas(r.NamaKelas)
		if tingkat == 0 {
			continue
		}

		hasil[tingkat] = append(hasil[tingkat], PrestasiSel{
			Prestasi:   strings.TrimSpace(r.Prestasi),
			Semester:   r.Semester,
			Nilai:      r.Nilai,
			Keterangan: r.Keterangan,
		})
	}

	return hasil, nil
}
