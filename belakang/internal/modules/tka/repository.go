package tka

import (
	"fmt"
	"sekolah/internal/database"

	"github.com/jmoiron/sqlx"
)

var allowedSortColumns = map[string]bool{
	"nis": true, "nopes": true, "nama_siswa": true,
	"nama_kelas": true, "nama_rombel": true, "mtk": true, "bind": true,
}

func GetAllAdmin(search, sortBy, sortDir string, limit, offset int) ([]TkaAdmin, error) {
	if !allowedSortColumns[sortBy] {
		sortBy = "nama_rombel"
	}
	orderExpr := sortBy
	if sortBy == "mtk" || sortBy == "bind" {
		orderExpr = "CAST(" + sortBy + " AS UNSIGNED)"
	}
	query := fmt.Sprintf(`
		SELECT t.id_tahunpelajaran, t.id_kelas, t.id_rombel, t.nis, t.nopes,
			t.mtk, t.ket_mtk, t.bind, t.ket_bind,
			DATE_FORMAT(t.waktu,'%%Y-%%m-%%d %%H:%%i:%%s') AS waktu,
			tp.nama_tahunpelajaran, k.nama_kelas, r.nama_rombel,
			COALESCE(pd.nama_siswa,'') AS nama_siswa
		FROM tka t
		LEFT JOIN tbl_tahunpelajaran tp ON tp.id_tahunpelajaran = t.id_tahunpelajaran
		LEFT JOIN tbl_kelas k ON k.id_kelas = t.id_kelas
		LEFT JOIN tbl_rombel r ON r.id_rombel = t.id_rombel
		LEFT JOIN pesertadidik pd ON pd.nis = t.nis
		WHERE t.nis LIKE ? OR t.nopes LIKE ? OR pd.nama_siswa LIKE ?
		ORDER BY %s %s, nama_siswa ASC LIMIT ? OFFSET ?`, orderExpr, sortDir)

	var rows []TkaAdmin
	like := "%" + search + "%"
	err := database.DB.Select(&rows, query, like, like, like, limit, offset)
	return rows, err
}

func CountAdmin(search string) (int, error) {
	var count int
	like := "%" + search + "%"
	err := database.DB.Get(&count, `
		SELECT COUNT(*) FROM tka t
		LEFT JOIN pesertadidik pd ON pd.nis = t.nis
		WHERE t.nis LIKE ? OR t.nopes LIKE ? OR pd.nama_siswa LIKE ?`, like, like, like)
	return count, err
}

func GetAllForExport(search string) ([]TkaExport, error) {
	like := "%" + search + "%"
	var rows []TkaExport
	err := database.DB.Select(&rows, `
		SELECT t.id_tahunpelajaran, t.id_kelas, t.id_rombel, t.nis, t.nopes,
			t.mtk, t.ket_mtk, t.bind, t.ket_bind,
			DATE_FORMAT(t.waktu,'%Y-%m-%d %H:%i:%s') AS waktu,
			tp.nama_tahunpelajaran, k.nama_kelas, r.nama_rombel,
			COALESCE(pd.nama_siswa,'') AS nama_siswa,
			COALESCE(pd.nisn,'') AS nisn,
			COALESCE(pd.jenis_kelamin,'') AS jenis_kelamin,
			COALESCE(pd.tempat_lahirsiswa,'') AS tempat_lahirsiswa,
			COALESCE(pd.tanggal_lahirsiswa,'') AS tanggal_lahirsiswa
		FROM tka t
		LEFT JOIN tbl_tahunpelajaran tp ON tp.id_tahunpelajaran = t.id_tahunpelajaran
		LEFT JOIN tbl_kelas k ON k.id_kelas = t.id_kelas
		LEFT JOIN tbl_rombel r ON r.id_rombel = t.id_rombel
		LEFT JOIN pesertadidik pd ON pd.nis = t.nis
		WHERE t.nis LIKE ? OR t.nopes LIKE ? OR pd.nama_siswa LIKE ?
		ORDER BY r.id_rombel ASC, pd.nama_siswa ASC`, like, like, like)
	return rows, err
}

func ExistsByNIS(nis string) (bool, error) {
	var count int
	err := database.DB.Get(&count, `SELECT COUNT(*) FROM tka WHERE nis = ?`, nis)
	return count > 0, err
}

func ExistsByNopes(nopes string) (bool, error) {
	var count int
	err := database.DB.Get(&count, `SELECT COUNT(*) FROM tka WHERE nopes = ?`, nopes)
	return count > 0, err
}

func ExistsPesertaDidik(nis string) (bool, error) {
	var count int
	err := database.DB.Get(&count, `SELECT COUNT(*) FROM pesertadidik WHERE nis = ?`, nis)
	return count > 0, err
}

func Insert(tx *sqlx.Tx, data Tka) error {
	_, err := tx.NamedExec(`
		INSERT INTO tka (id_tahunpelajaran, id_kelas, id_rombel, nis, nopes, mtk, ket_mtk, bind, ket_bind)
		VALUES (:id_tahunpelajaran, :id_kelas, :id_rombel, :nis, :nopes, :mtk, :ket_mtk, :bind, :ket_bind)`, data)
	return err
}

func BulkDelete(nisList []string) error {
	query, args, err := sqlx.In(`DELETE FROM tka WHERE nis IN (?)`, nisList)
	if err != nil {
		return err
	}
	query = database.DB.Rebind(query)
	_, err = database.DB.Exec(query, args...)
	return err
}

func GetWaktu() (TkaWaktu, error) {
	var w TkaWaktu
	err := database.DB.Get(&w, `
		SELECT tka_waktu.id_waktu, tka_waktu.id_pengguna,
			DATE_FORMAT(tka_waktu.aturwaktu,'%Y-%m-%d %H:%i:%s') AS aturwaktu,
			tka_waktu.aturform,
			tka_waktu.judul, tka_waktu.nomor, tka_waktu.pengantar, tka_waktu.penutup, tka_waktu.akhir,
			tka_waktu.tanggal, tka_waktu.mengetahui, tka_waktu.kepsek, tka_waktu.nipkepsek,
			COALESCE(tka_waktu.ttd_skl,'') AS ttd_skl,
			COALESCE(tka_waktu.tahun_pelajaran,'') AS tahun_pelajaran,
			DATE_FORMAT(tka_waktu.waktu,'%Y-%m-%d %H:%i:%s') AS waktu,
			COALESCE(pengguna.nama_pengguna,'') AS nama_pengguna
		FROM tka_waktu
		LEFT JOIN pengguna ON pengguna.id_pengguna = tka_waktu.id_pengguna
		WHERE tka_waktu.id_waktu = 1`)
	return w, err
}

func UpdateWaktu(idPengguna string, data TkaWaktu) error {
	data.IDPengguna = idPengguna
	_, err := database.DB.NamedExec(`
		UPDATE tka_waktu SET id_pengguna=:id_pengguna, aturwaktu=:aturwaktu, aturform=:aturform,
			judul=:judul, nomor=:nomor, pengantar=:pengantar, penutup=:penutup, akhir=:akhir,
			tanggal=:tanggal, mengetahui=:mengetahui, kepsek=:kepsek, nipkepsek=:nipkepsek,
			ttd_skl=:ttd_skl, tahun_pelajaran=:tahun_pelajaran
		WHERE id_waktu = 1`, data)
	return err
}

// GetByNIS dipakai handler /tka/me (sisi Peserta Didik, reuse fungsi admin)
func GetByNIS(nis string) (*TkaAdmin, error) {
	var row TkaAdmin
	err := database.DB.Get(&row, `
		SELECT t.id_tahunpelajaran, t.id_kelas, t.id_rombel, t.nis, t.nopes,
			t.mtk, t.ket_mtk, t.bind, t.ket_bind,
			DATE_FORMAT(t.waktu,'%Y-%m-%d %H:%i:%s') AS waktu,
			COALESCE(tp.nama_tahunpelajaran,'') AS nama_tahunpelajaran,
			COALESCE(k.nama_kelas,'') AS nama_kelas,
			COALESCE(r.nama_rombel,'') AS nama_rombel,
			COALESCE(pd.nama_siswa,'') AS nama_siswa
		FROM tka t
		LEFT JOIN tbl_tahunpelajaran tp ON tp.id_tahunpelajaran = t.id_tahunpelajaran
		LEFT JOIN tbl_kelas k ON k.id_kelas = t.id_kelas
		LEFT JOIN tbl_rombel r ON r.id_rombel = t.id_rombel
		LEFT JOIN pesertadidik pd ON pd.nis = t.nis
		WHERE t.nis = ?`, nis)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

type refItem struct {
	ID   int    `db:"id"`
	Nama string `db:"nama"`
}

func getTahunPelajaranAktifList() ([]refItem, error) {
	var rows []refItem
	err := database.DB.Select(&rows, `
		SELECT id_tahunpelajaran AS id, nama_tahunpelajaran AS nama
		FROM tbl_tahunpelajaran WHERE aktif = 1 ORDER BY id_tahunpelajaran ASC`)
	return rows, err
}

func getKelasAktifList() ([]refItem, error) {
	var rows []refItem
	err := database.DB.Select(&rows, `
		SELECT id_kelas AS id, nama_kelas AS nama
		FROM tbl_kelas WHERE aktif = 1 ORDER BY id_kelas ASC`)
	return rows, err
}

func getRombelAktifList() ([]refItem, error) {
	var rows []refItem
	err := database.DB.Select(&rows, `
		SELECT id_rombel AS id, nama_rombel AS nama
		FROM tbl_rombel WHERE aktif = 1 ORDER BY id_rombel ASC`)
	return rows, err
}
