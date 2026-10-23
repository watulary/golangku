package kenaikan

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"sekolah/internal/database"
)

var allowedSortColumns = map[string]string{
	"nis":                 "kenaikan.nis",
	"nama_siswa":          "pesertadidik.nama_siswa",
	"nama_tahunpelajaran": "tbl_tahunpelajaran.nama_tahunpelajaran",
	"nama_kelas":          "tbl_kelas.nama_kelas",
	"nama_rombel":         "tbl_rombel.nama_rombel",
	"ket":                 "kenaikan.ket",
}

const baseSelectKenaikan = `
	SELECT
		kenaikan.id_tahunpelajaran, kenaikan.id_kelas, kenaikan.id_rombel,
		kenaikan.nis, kenaikan.ket,
		COALESCE(tbl_tahunpelajaran.nama_tahunpelajaran, '') AS nama_tahunpelajaran,
		COALESCE(tbl_kelas.nama_kelas, '') AS nama_kelas,
		COALESCE(tbl_rombel.nama_rombel, '') AS nama_rombel,
		COALESCE(pesertadidik.nama_siswa, '') AS nama_siswa,
		COALESCE(pesertadidik.nisn, '') AS nisn,
		COALESCE(pesertadidik.no_ujiansd, '') AS no_ujiansd,
		COALESCE(pesertadidik.tempat_lahirsiswa, '') AS tempat_lahirsiswa,
		COALESCE(pesertadidik.tanggal_lahirsiswa, '') AS tanggal_lahirsiswa,
		COALESCE(pesertadidik.jenis_kelamin, '') AS jenis_kelamin,
		COALESCE(pesertadidik.foto, '') AS foto
	FROM kenaikan
	LEFT JOIN tbl_tahunpelajaran ON tbl_tahunpelajaran.id_tahunpelajaran = kenaikan.id_tahunpelajaran
	LEFT JOIN tbl_kelas ON tbl_kelas.id_kelas = kenaikan.id_kelas
	LEFT JOIN tbl_rombel ON tbl_rombel.id_rombel = kenaikan.id_rombel
	LEFT JOIN pesertadidik ON pesertadidik.nis = kenaikan.nis
`

func GetAllAdmin(search, sortBy, sortDir string, limit, offset int) ([]KenaikanListItem, error) {
	query := baseSelectKenaikan
	args := []interface{}{}

	if search != "" {
		query += ` WHERE (kenaikan.nis LIKE ? OR pesertadidik.nama_siswa LIKE ? OR pesertadidik.nisn LIKE ?)`
		like := "%" + search + "%"
		args = append(args, like, like, like)
	}

	orderCol, ok := allowedSortColumns[sortBy]
	if !ok {
		orderCol = "tbl_rombel.id_rombel"
	}
	dir := "ASC"
	if strings.ToLower(sortDir) == "desc" {
		dir = "DESC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s, pesertadidik.nama_siswa ASC", orderCol, dir)
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var items []KenaikanListItem
	if err := database.DB.Select(&items, database.DB.Rebind(query), args...); err != nil {
		return nil, err
	}
	return items, nil
}

func CountAdmin(search string) (int, error) {
	query := `SELECT COUNT(*) FROM kenaikan
		LEFT JOIN pesertadidik ON pesertadidik.nis = kenaikan.nis`
	args := []interface{}{}

	if search != "" {
		query += ` WHERE (kenaikan.nis LIKE ? OR pesertadidik.nama_siswa LIKE ? OR pesertadidik.nisn LIKE ?)`
		like := "%" + search + "%"
		args = append(args, like, like, like)
	}

	var total int
	if err := database.DB.Get(&total, database.DB.Rebind(query), args...); err != nil {
		return 0, err
	}
	return total, nil
}

func GetAllForExport(search string) ([]KenaikanListItem, error) {
	query := baseSelectKenaikan
	args := []interface{}{}

	if search != "" {
		query += ` WHERE (kenaikan.nis LIKE ? OR pesertadidik.nama_siswa LIKE ? OR pesertadidik.nisn LIKE ?)`
		like := "%" + search + "%"
		args = append(args, like, like, like)
	}
	query += " ORDER BY tbl_rombel.id_rombel ASC, pesertadidik.nama_siswa ASC"

	var items []KenaikanListItem
	if err := database.DB.Select(&items, database.DB.Rebind(query), args...); err != nil {
		return nil, err
	}
	return items, nil
}

func GetByID(nis string) (*KenaikanListItem, error) {
	query := baseSelectKenaikan + " WHERE kenaikan.nis = ?"
	var item KenaikanListItem
	if err := database.DB.Get(&item, database.DB.Rebind(query), nis); err != nil {
		return nil, err
	}
	return &item, nil
}

func Exists(nis string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM kenaikan WHERE nis = ?`
	if err := database.DB.Get(&count, database.DB.Rebind(query), nis); err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsPesertaDidik mengecek apakah nis benar-benar terdaftar di tabel pesertadidik
// (dipakai validasi form tambah/edit manual & import Excel, biar tidak ada data kenaikan
// untuk siswa yang tidak ada).
func ExistsPesertaDidik(nis string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM pesertadidik WHERE nis = ?`
	if err := database.DB.Get(&count, database.DB.Rebind(query), nis); err != nil {
		return false, err
	}
	return count > 0, nil
}

func Create(data Kenaikan) error {
	query := `INSERT INTO kenaikan (id_tahunpelajaran, id_kelas, id_rombel, nis, ket)
		VALUES (?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(database.DB.Rebind(query),
		data.IDTahunPelajaran, data.IDKelas, data.IDRombel, data.NIS, data.Ket)
	return err
}

func Update(nis string, data Kenaikan) error {
	query := `UPDATE kenaikan SET id_tahunpelajaran = ?, id_kelas = ?, id_rombel = ?, ket = ?
		WHERE nis = ?`
	_, err := database.DB.Exec(database.DB.Rebind(query),
		data.IDTahunPelajaran, data.IDKelas, data.IDRombel, data.Ket, nis)
	return err
}

func Delete(nis string) error {
	query := `DELETE FROM kenaikan WHERE nis = ?`
	_, err := database.DB.Exec(database.DB.Rebind(query), nis)
	return err
}

func BulkDelete(nisList []string) error {
	if len(nisList) == 0 {
		return nil
	}
	query, args, err := sqlx.In(`DELETE FROM kenaikan WHERE nis IN (?)`, nisList)
	if err != nil {
		return err
	}
	_, err = database.DB.Exec(database.DB.Rebind(query), args...)
	return err
}

// --- kenaikan_waktu (pengaturan pengumuman, single-row id_waktu = 1) ---

func GetWaktu() (*KenaikanWaktu, error) {
	query := `SELECT kenaikan_waktu.id_waktu, kenaikan_waktu.id_pengguna,
			COALESCE(kenaikan_waktu.aturwaktu, '') AS aturwaktu,
			kenaikan_waktu.waktu,
			COALESCE(kenaikan_waktu.judul, '') AS judul,
			COALESCE(kenaikan_waktu.nomor, '') AS nomor,
			COALESCE(kenaikan_waktu.pengantar, '') AS pengantar,
			COALESCE(kenaikan_waktu.penutup, '') AS penutup,
			COALESCE(kenaikan_waktu.akhir, '') AS akhir,
			COALESCE(kenaikan_waktu.tanggal, '') AS tanggal,
			COALESCE(kenaikan_waktu.mengetahui, '') AS mengetahui,
			COALESCE(kenaikan_waktu.kepsek, '') AS kepsek,
			COALESCE(kenaikan_waktu.nipkepsek, '') AS nipkepsek,
			COALESCE(kenaikan_waktu.ttd_skl, '') AS ttd_skl,
			COALESCE(kenaikan_waktu.tahun_pelajaran, '') AS tahun_pelajaran,
			COALESCE(pengguna.nama_pengguna, '') AS nama_pengguna
		FROM kenaikan_waktu
		LEFT JOIN pengguna ON pengguna.id_pengguna = kenaikan_waktu.id_pengguna
		WHERE kenaikan_waktu.id_waktu = 1`

	var w KenaikanWaktu
	if err := database.DB.Get(&w, query); err != nil {
		return nil, err
	}
	return &w, nil
}

func UpdateWaktu(idPengguna string, data KenaikanWaktu) error {
	query := `UPDATE kenaikan_waktu SET
			id_pengguna = ?, aturwaktu = ?, judul = ?, nomor = ?, pengantar = ?, penutup = ?,
			akhir = ?, tanggal = ?, mengetahui = ?, kepsek = ?, nipkepsek = ?, ttd_skl = ?, tahun_pelajaran = ?
		WHERE id_waktu = 1`
	_, err := database.DB.Exec(database.DB.Rebind(query),
		idPengguna, data.Aturwaktu, data.Judul, data.Nomor, data.Pengantar, data.Penutup,
		data.Akhir, data.Tanggal, data.Mengetahui, data.Kepsek, data.NIPKepsek, data.TTDSkl, data.TahunPelajaran)
	return err
}
