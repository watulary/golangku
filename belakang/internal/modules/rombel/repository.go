package rombel

import (
	"sekolah/internal/database"
	"strings"
)

var allowedSortColumns = map[string]string{
	"nama":  "r.nama_rombel",
	"kelas": "k.nama_kelas",
	"aktif": "r.aktif",
}

func GetAllAdmin(search string, idKelas int, sortBy, sortDir string, limit, offset int) ([]Rombel, error) {
	var data []Rombel

	query := `
		SELECT
			r.id_rombel,
			r.id_kelas,
			r.nama_rombel,
			r.aktif,
			k.nama_kelas
		FROM tbl_rombel r
		LEFT JOIN tbl_kelas k ON k.id_kelas = r.id_kelas
		WHERE (
			r.nama_rombel LIKE ?
			OR k.nama_kelas LIKE ?
		)`

	keyword := "%" + search + "%"

	args := []interface{}{keyword, keyword}

	if idKelas > 0 {
		query += ` AND r.id_kelas = ?`
		args = append(args, idKelas)
	}

	column, ok := allowedSortColumns[sortBy]
	if !ok {
		column = "r.nama_rombel"
	}

	direction := "ASC"
	if strings.ToLower(sortDir) == "desc" {
		direction = "DESC"
	}

	query += ` ORDER BY ` + column + ` ` + direction + `, r.nama_rombel ASC, r.id_rombel ASC LIMIT ? OFFSET ?`

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
		FROM tbl_rombel r
		LEFT JOIN tbl_kelas k ON k.id_kelas = r.id_kelas
		WHERE (
			r.nama_rombel LIKE ?
			OR k.nama_kelas LIKE ?
		)`

	keyword := "%" + search + "%"

	args := []interface{}{keyword, keyword}

	if idKelas > 0 {
		query += ` AND r.id_kelas = ?`
		args = append(args, idKelas)
	}

	err := database.DB.Get(&total, query, args...)

	return total, err
}

func GetByID(id int) (*Rombel, error) {
	var data Rombel

	err := database.DB.Get(
		&data,
		`SELECT
			r.id_rombel,
			r.id_kelas,
			r.nama_rombel,
			r.aktif,
			k.nama_kelas
		FROM tbl_rombel r
		LEFT JOIN tbl_kelas k ON k.id_kelas = r.id_kelas
		WHERE r.id_rombel = ?
		LIMIT 1`,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func ExistsByNama(nama string, idKelas int, excludeID int) (bool, error) {
	var total int

	err := database.DB.Get(
		&total,
		`SELECT COUNT(*)
		FROM tbl_rombel
		WHERE nama_rombel = ?
		AND id_kelas = ?
		AND id_rombel != ?`,
		nama,
		idKelas,
		excludeID,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
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

func Create(data Rombel) (int64, error) {
	result, err := database.DB.Exec(
		`INSERT INTO tbl_rombel
			(id_kelas, nama_rombel, aktif)
		VALUES
			(?, ?, ?)`,
		data.IDKelas,
		data.Nama,
		data.Aktif,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func Update(id int, data Rombel) error {
	_, err := database.DB.Exec(
		`UPDATE tbl_rombel SET
			id_kelas = ?,
			nama_rombel = ?
		WHERE id_rombel = ?`,
		data.IDKelas,
		data.Nama,
		id,
	)

	return err
}

func Delete(id int) error {
	_, err := database.DB.Exec(
		`DELETE FROM tbl_rombel WHERE id_rombel = ?`,
		id,
	)

	return err
}

func ToggleStatus(id int, aktif int) error {
	_, err := database.DB.Exec(
		`UPDATE tbl_rombel SET aktif = ? WHERE id_rombel = ?`,
		aktif,
		id,
	)

	return err
}
