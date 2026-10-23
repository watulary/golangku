package walikelas

import (
	"sekolah/internal/database"
	"strings"
)

var allowedSortColumns = map[string]string{
	"kelas":   "k.nama_kelas",
	"rombel":  "r.nama_rombel",
	"pegawai": "p.nama_pegawai",
}

func GetAllAdmin(search string, idKelas int, sortBy, sortDir string, limit, offset int) ([]WaliKelas, error) {
	var data []WaliKelas

	query := `
		SELECT
			w.id_walikelas,
			w.id_kelas,
			w.id_rombel,
			w.nuptk,
			k.nama_kelas,
			r.nama_rombel,
			p.nama_pegawai
		FROM tbl_walikelas w
		LEFT JOIN tbl_kelas k ON k.id_kelas = w.id_kelas
		LEFT JOIN tbl_rombel r ON r.id_rombel = w.id_rombel
		LEFT JOIN pegawai p ON p.nuptk = w.nuptk
		WHERE (
			k.nama_kelas LIKE ?
			OR r.nama_rombel LIKE ?
			OR p.nama_pegawai LIKE ?
			OR w.nuptk LIKE ?
		)`

	keyword := "%" + search + "%"

	args := []interface{}{keyword, keyword, keyword, keyword}

	if idKelas > 0 {
		query += ` AND w.id_kelas = ?`
		args = append(args, idKelas)
	}

	column, ok := allowedSortColumns[sortBy]
	if !ok {
		column = "k.nama_kelas"
	}

	direction := "ASC"
	if strings.ToLower(sortDir) == "desc" {
		direction = "DESC"
	}

	query += ` ORDER BY ` + column + ` ` + direction + `, r.nama_rombel ASC, w.id_walikelas ASC LIMIT ? OFFSET ?`

	args = append(args, limit, offset)

	err := database.DB.Select(&data, query, args...)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func CountAdmin(search string, idKelas int) (int, error) {
	var total int

	query := `
		SELECT COUNT(*)
		FROM tbl_walikelas w
		LEFT JOIN tbl_kelas k ON k.id_kelas = w.id_kelas
		LEFT JOIN tbl_rombel r ON r.id_rombel = w.id_rombel
		LEFT JOIN pegawai p ON p.nuptk = w.nuptk
		WHERE (
			k.nama_kelas LIKE ?
			OR r.nama_rombel LIKE ?
			OR p.nama_pegawai LIKE ?
			OR w.nuptk LIKE ?
		)`

	keyword := "%" + search + "%"

	args := []interface{}{keyword, keyword, keyword, keyword}

	if idKelas > 0 {
		query += ` AND w.id_kelas = ?`
		args = append(args, idKelas)
	}

	err := database.DB.Get(&total, query, args...)

	return total, err
}

func GetByID(id int) (*WaliKelas, error) {
	var data WaliKelas

	err := database.DB.Get(
		&data,
		`SELECT
			w.id_walikelas,
			w.id_kelas,
			w.id_rombel,
			w.nuptk,
			k.nama_kelas,
			r.nama_rombel,
			p.nama_pegawai
		FROM tbl_walikelas w
		LEFT JOIN tbl_kelas k ON k.id_kelas = w.id_kelas
		LEFT JOIN tbl_rombel r ON r.id_rombel = w.id_rombel
		LEFT JOIN pegawai p ON p.nuptk = w.nuptk
		WHERE w.id_walikelas = ?
		LIMIT 1`,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func KelasExists(idKelas int) (bool, error) {
	var total int

	err := database.DB.Get(
		&total,
		`SELECT COUNT(*) FROM tbl_kelas WHERE id_kelas = ?`,
		idKelas,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

func RombelBelongsToKelas(idRombel, idKelas int) (bool, error) {
	var total int

	err := database.DB.Get(
		&total,
		`SELECT COUNT(*)
		FROM tbl_rombel
		WHERE id_rombel = ?
		AND id_kelas = ?`,
		idRombel,
		idKelas,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

func PegawaiExists(nuptk string) (bool, error) {
	var total int

	err := database.DB.Get(
		&total,
		`SELECT COUNT(*) FROM pegawai WHERE nuptk = ?`,
		nuptk,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

func ExistsByRombel(idRombel int, excludeID int) (bool, error) {
	var total int

	err := database.DB.Get(
		&total,
		`SELECT COUNT(*)
		FROM tbl_walikelas
		WHERE id_rombel = ?
		AND id_walikelas != ?`,
		idRombel,
		excludeID,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

func ExistsByPegawai(nuptk string, excludeID int) (bool, error) {
	var total int

	err := database.DB.Get(
		&total,
		`SELECT COUNT(*)
		FROM tbl_walikelas
		WHERE nuptk = ?
		AND id_walikelas != ?`,
		nuptk,
		excludeID,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

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

func Create(data WaliKelas) (int64, error) {
	result, err := database.DB.Exec(
		`INSERT INTO tbl_walikelas
			(id_kelas, id_rombel, nuptk)
		VALUES
			(?, ?, ?)`,
		data.IDKelas,
		data.IDRombel,
		data.NUPTK,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func Update(id int, data WaliKelas) error {
	_, err := database.DB.Exec(
		`UPDATE tbl_walikelas SET
			id_kelas = ?,
			id_rombel = ?,
			nuptk = ?
		WHERE id_walikelas = ?`,
		data.IDKelas,
		data.IDRombel,
		data.NUPTK,
		id,
	)

	return err
}

func Delete(id int) error {
	_, err := database.DB.Exec(
		`DELETE FROM tbl_walikelas WHERE id_walikelas = ?`,
		id,
	)

	return err
}
