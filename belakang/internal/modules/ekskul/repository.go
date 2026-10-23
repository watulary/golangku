package ekskul

import (
	"database/sql"
	"fmt"
	"sekolah/internal/database"
	"strings"

	"github.com/jmoiron/sqlx"
)

var allowedSortColumns = map[string]string{
	"ekskul":  "e.nama_ekskul",
	"pegawai": "p.nama_pegawai",
	"nuptk":   "e.nuptk",
	"anggota": "jumlah_anggota",
}

var allowedSortColumnsAnggota = map[string]string{
	"rombel": "r.nama_rombel",
	"nis":    "a.nis",
	"siswa":  "pd.nama_siswa",
}

// ================= EKSKUL =================

func GetAllAdmin(search string, statusAktif int, sortBy, sortDir string, limit, offset int) ([]Ekskul, error) {
	var data []Ekskul

	query := `
		SELECT
			e.id_ekskul,
			e.nama_ekskul,
			e.nuptk,
			e.aktif,
			COALESCE(p.nama_pegawai, '') AS nama_pegawai,
			(SELECT COUNT(*) FROM tbl_anggota_ekskul a WHERE a.id_ekskul = e.id_ekskul) AS jumlah_anggota
		FROM tbl_ekskul e
		LEFT JOIN pegawai p ON p.nuptk = e.nuptk
		WHERE (
			e.nama_ekskul LIKE ?
			OR p.nama_pegawai LIKE ?
			OR e.nuptk LIKE ?
		)`

	keyword := "%" + search + "%"
	args := []interface{}{keyword, keyword, keyword}

	if statusAktif >= 0 {
		query += ` AND e.aktif = ?`
		args = append(args, statusAktif)
	}

	column, ok := allowedSortColumns[sortBy]
	if !ok {
		column = "e.nama_ekskul"
	}

	direction := "ASC"
	if strings.ToLower(sortDir) == "desc" {
		direction = "DESC"
	}

	query += ` ORDER BY ` + column + ` ` + direction + `, e.nama_ekskul ASC, e.id_ekskul ASC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	err := database.DB.Select(&data, query, args...)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func CountAdmin(search string, statusAktif int) (int, error) {
	var total int

	query := `
		SELECT COUNT(*)
		FROM tbl_ekskul e
		LEFT JOIN pegawai p ON p.nuptk = e.nuptk
		WHERE (
			e.nama_ekskul LIKE ?
			OR p.nama_pegawai LIKE ?
			OR e.nuptk LIKE ?
		)`

	keyword := "%" + search + "%"
	args := []interface{}{keyword, keyword, keyword}

	if statusAktif >= 0 {
		query += ` AND e.aktif = ?`
		args = append(args, statusAktif)
	}

	err := database.DB.Get(&total, query, args...)

	return total, err
}

func GetByID(id int) (*Ekskul, error) {
	var data Ekskul

	err := database.DB.Get(
		&data,
		`SELECT
			e.id_ekskul,
			e.nama_ekskul,
			e.nuptk,
			e.aktif,
			COALESCE(p.nama_pegawai, '') AS nama_pegawai,
			(SELECT COUNT(*) FROM tbl_anggota_ekskul a WHERE a.id_ekskul = e.id_ekskul) AS jumlah_anggota
		FROM tbl_ekskul e
		LEFT JOIN pegawai p ON p.nuptk = e.nuptk
		WHERE e.id_ekskul = ?
		LIMIT 1`,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GetEkskulBySiswa mengambil semua ekskul yang diikuti seorang siswa (lewat tbl_anggota_ekskul),
// dipakai endpoint /ekskul/me sisi Peserta Didik. Kalau siswa belum ikut ekskul manapun,
// mengembalikan slice kosong (bukan error).
func GetEkskulBySiswa(nis string) ([]Ekskul, error) {
	data := []Ekskul{}

	err := database.DB.Select(
		&data,
		`SELECT
			e.id_ekskul,
			e.nama_ekskul,
			e.nuptk,
			e.aktif,
			COALESCE(p.nama_pegawai, '') AS nama_pegawai
		FROM tbl_anggota_ekskul a
		JOIN tbl_ekskul e ON e.id_ekskul = a.id_ekskul
		LEFT JOIN pegawai p ON p.nuptk = e.nuptk
		WHERE a.nis = ?
		ORDER BY e.nama_ekskul ASC`,
		nis,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

// PembinaValid memastikan nuptk ada di tabel pegawai DAN jabatannya Guru.
func PembinaValid(nuptk string) (bool, error) {
	var total int

	err := database.DB.Get(
		&total,
		`SELECT COUNT(*)
		FROM pegawai
		WHERE nuptk = ?
		AND jabatan = 'Guru'`,
		nuptk,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

func ExistsByNama(nama string, excludeID int) (bool, error) {
	var total int

	err := database.DB.Get(
		&total,
		`SELECT COUNT(*)
		FROM tbl_ekskul
		WHERE nama_ekskul = ?
		AND id_ekskul != ?`,
		nama,
		excludeID,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

// SearchPegawai mencari pegawai dengan jabatan Guru saja, untuk dropdown pembina.
func SearchPegawai(search string, limit int) ([]PegawaiOption, error) {
	var data []PegawaiOption

	err := database.DB.Select(
		&data,
		`SELECT
			nuptk,
			nama_pegawai,
			jabatan
		FROM pegawai
		WHERE (
			nama_pegawai LIKE ?
			OR nuptk LIKE ?
		)
		AND jabatan = 'Guru'
		AND (status_keluar = '' OR status_keluar IS NULL)
		ORDER BY nama_pegawai ASC
		LIMIT ?`,
		"%"+search+"%",
		"%"+search+"%",
		limit,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func Create(data Ekskul) (int64, error) {
	result, err := database.DB.Exec(
		`INSERT INTO tbl_ekskul
			(nama_ekskul, nuptk, aktif)
		VALUES
			(?, ?, ?)`,
		data.NamaEkskul,
		data.NUPTK,
		data.Aktif,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func Update(id int, data Ekskul) error {
	_, err := database.DB.Exec(
		`UPDATE tbl_ekskul SET
			nama_ekskul = ?,
			nuptk = ?,
			aktif = ?
		WHERE id_ekskul = ?`,
		data.NamaEkskul,
		data.NUPTK,
		data.Aktif,
		id,
	)

	return err
}

func Delete(id int) error {
	_, err := database.DB.Exec(
		`DELETE FROM tbl_ekskul WHERE id_ekskul = ?`,
		id,
	)

	return err
}

// ================= ANGGOTA EKSKUL =================

func GetAnggotaByEkskul(idEkskul int, search, sortBy, sortDir string, limit, offset, idKelas, idRombel int) ([]AnggotaEkskul, error) {
	var data []AnggotaEkskul

	column, ok := allowedSortColumnsAnggota[sortBy]
	if !ok {
		column = "pd.nama_siswa"
	}

	direction := "ASC"
	if strings.ToLower(sortDir) == "desc" {
		direction = "DESC"
	}

	query := `
		SELECT
			a.id,
			a.id_ekskul,
			a.nis,
			pd.nama_siswa,
			e.nama_ekskul,
			r.nama_rombel
		FROM tbl_anggota_ekskul a
		JOIN pesertadidik pd ON pd.nis = a.nis
		LEFT JOIN tbl_rombel r ON r.id_rombel = pd.id_rombel
		LEFT JOIN tbl_ekskul e ON e.id_ekskul = a.id_ekskul
		WHERE a.id_ekskul = ?
		AND (
			pd.nama_siswa LIKE ?
			OR a.nis LIKE ?
			OR r.nama_rombel LIKE ?
		)`

	args := []interface{}{idEkskul, "%" + search + "%", "%" + search + "%", "%" + search + "%"}

	if idKelas > 0 {
		query += ` AND r.id_kelas = ?`
		args = append(args, idKelas)
	}
	if idRombel > 0 {
		query += ` AND r.id_rombel = ?`
		args = append(args, idRombel)
	}

	query += fmt.Sprintf(` ORDER BY %s %s, pd.nama_siswa ASC, a.id ASC LIMIT ? OFFSET ?`, column, direction)
	args = append(args, limit, offset)

	err := database.DB.Select(&data, query, args...)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func CountAnggotaByEkskul(idEkskul int, search string, idKelas, idRombel int) (int, error) {
	var total int

	query := `
		SELECT COUNT(*)
		FROM tbl_anggota_ekskul a
		JOIN pesertadidik pd ON pd.nis = a.nis
		LEFT JOIN tbl_rombel r ON r.id_rombel = pd.id_rombel
		WHERE a.id_ekskul = ?
		AND (
			pd.nama_siswa LIKE ?
			OR a.nis LIKE ?
			OR r.nama_rombel LIKE ?
		)`

	args := []interface{}{idEkskul, "%" + search + "%", "%" + search + "%", "%" + search + "%"}

	if idKelas > 0 {
		query += ` AND r.id_kelas = ?`
		args = append(args, idKelas)
	}
	if idRombel > 0 {
		query += ` AND r.id_rombel = ?`
		args = append(args, idRombel)
	}

	err := database.DB.Get(&total, query, args...)

	return total, err
}

// GetAnggotaByEkskulForExport mengambil semua anggota satu ekskul (tanpa paginasi),
// dipakai untuk export Excel per-ekskul dari halaman Anggota.vue.
func GetAnggotaByEkskulForExport(idEkskul int, search string, idKelas, idRombel int) ([]AnggotaEkskul, error) {
	var data []AnggotaEkskul

	query := `
		SELECT
			a.id,
			a.id_ekskul,
			a.nis,
			pd.nama_siswa,
			e.nama_ekskul,
			COALESCE(r.nama_rombel, '') AS nama_rombel
		FROM tbl_anggota_ekskul a
		JOIN pesertadidik pd ON pd.nis = a.nis
		LEFT JOIN tbl_rombel r ON r.id_rombel = pd.id_rombel
		LEFT JOIN tbl_ekskul e ON e.id_ekskul = a.id_ekskul
		WHERE a.id_ekskul = ?
		AND (
			pd.nama_siswa LIKE ?
			OR a.nis LIKE ?
			OR r.nama_rombel LIKE ?
		)`

	args := []interface{}{idEkskul, "%" + search + "%", "%" + search + "%", "%" + search + "%"}

	if idKelas > 0 {
		query += ` AND r.id_kelas = ?`
		args = append(args, idKelas)
	}
	if idRombel > 0 {
		query += ` AND r.id_rombel = ?`
		args = append(args, idRombel)
	}

	query += ` ORDER BY r.nama_rombel ASC, pd.nama_siswa ASC, a.id ASC`

	err := database.DB.Select(&data, query, args...)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func AnggotaExists(idEkskul int, nis string) (bool, error) {
	var total int

	err := database.DB.Get(
		&total,
		`SELECT COUNT(*)
		FROM tbl_anggota_ekskul
		WHERE id_ekskul = ?
		AND nis = ?`,
		idEkskul,
		nis,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

// SiswaAktifValid memastikan nis ada di pesertadidik, belum keluar,
// dan berada di tahun pelajaran yang sedang aktif.
func SiswaAktifValid(nis string) (bool, error) {
	var total int

	err := database.DB.Get(
		&total,
		`SELECT COUNT(*)
		FROM pesertadidik
		WHERE nis = ?
		AND (status_keluar = '' OR status_keluar IS NULL)
		AND id_tahunpelajaran = (
			SELECT id_tahunpelajaran FROM tbl_tahunpelajaran WHERE aktif = 1 LIMIT 1
		)`,
		nis,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

// SearchPesertaDidik mencari siswa aktif tahun berjalan untuk dropdown tambah anggota.
func SearchPesertaDidik(search string, limit int) ([]PesertaDidikOption, error) {
	var data []PesertaDidikOption

	err := database.DB.Select(
		&data,
		`SELECT
			pd.nis,
			pd.nama_siswa,
			COALESCE(r.nama_rombel, '') AS nama_rombel
		FROM pesertadidik pd
		LEFT JOIN tbl_rombel r ON r.id_rombel = pd.id_rombel
		WHERE (
			pd.nama_siswa LIKE ?
			OR pd.nis LIKE ?
		)
		AND (pd.status_keluar = '' OR pd.status_keluar IS NULL)
		AND pd.id_tahunpelajaran = (
			SELECT id_tahunpelajaran FROM tbl_tahunpelajaran WHERE aktif = 1 LIMIT 1
		)
		ORDER BY pd.nama_siswa ASC
		LIMIT ?`,
		"%"+search+"%",
		"%"+search+"%",
		limit,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func AddAnggota(idEkskul int, nis string) (int64, error) {
	result, err := database.DB.Exec(
		`INSERT INTO tbl_anggota_ekskul
			(id_ekskul, nis)
		VALUES
			(?, ?)`,
		idEkskul,
		nis,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func RemoveAnggota(id int) error {
	result, err := database.DB.Exec(
		`DELETE FROM tbl_anggota_ekskul WHERE id = ?`,
		id,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// BulkRemoveAnggota menghapus banyak anggota sekaligus berdasarkan daftar ID.
func BulkRemoveAnggota(ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	query, args, err := sqlx.In(`DELETE FROM tbl_anggota_ekskul WHERE id IN (?)`, ids)
	if err != nil {
		return err
	}

	query = database.DB.Rebind(query)

	_, err = database.DB.Exec(query, args...)

	return err
}

// GetAllForExport mengambil semua ekskul (tanpa paginasi) untuk export Excel.
func GetAllForExport() ([]Ekskul, error) {
	var data []Ekskul

	err := database.DB.Select(
		&data,
		`SELECT
			e.id_ekskul,
			e.nama_ekskul,
			e.nuptk,
			e.aktif,
			COALESCE(p.nama_pegawai, '') AS nama_pegawai,
			(SELECT COUNT(*) FROM tbl_anggota_ekskul a WHERE a.id_ekskul = e.id_ekskul) AS jumlah_anggota
		FROM tbl_ekskul e
		LEFT JOIN pegawai p ON p.nuptk = e.nuptk
		ORDER BY e.nama_ekskul ASC`,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetAllAnggotaForExport mengambil semua anggota dari semua ekskul (tanpa paginasi),
// diurutkan per ekskul untuk keperluan export sheet-per-ekskul.
func GetAllAnggotaForExport() ([]AnggotaEkskul, error) {
	var data []AnggotaEkskul

	err := database.DB.Select(
		&data,
		`SELECT
			a.id,
			a.id_ekskul,
			a.nis,
			COALESCE(pd.nama_siswa, '') AS nama_siswa,
			COALESCE(r.nama_rombel, '') AS nama_rombel,
			e.nama_ekskul
		FROM tbl_anggota_ekskul a
		JOIN pesertadidik pd ON pd.nis = a.nis
		LEFT JOIN tbl_rombel r ON r.id_rombel = pd.id_rombel
		JOIN tbl_ekskul e ON e.id_ekskul = a.id_ekskul
		ORDER BY e.nama_ekskul ASC, pd.nama_siswa ASC`,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}
