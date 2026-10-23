package kelas

import (
	"sekolah/internal/database"
	"strings"
)

var allowedSortColumns = map[string]string{
	"nama":  "nama_kelas",
	"aktif": "aktif",
}

func GetAllAdmin(search, sortBy, sortDir string, limit, offset int) ([]Kelas, error) {
	var data []Kelas

	query := `
		SELECT
			id_kelas,
			nama_kelas,
			aktif
		FROM tbl_kelas
		WHERE nama_kelas LIKE ?`

	keyword := "%" + search + "%"

	args := []interface{}{keyword}

	column, ok := allowedSortColumns[sortBy]
	if !ok {
		column = "nama_kelas"
	}

	direction := "ASC"
	if strings.ToLower(sortDir) == "desc" {
		direction = "DESC"
	}

	query += ` ORDER BY ` + column + ` ` + direction + ` LIMIT ? OFFSET ?`

	args = append(args, limit, offset)

	err := database.DB.Select(&data, query, args...)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func CountAdmin(search string) (int, error) {
	var total int

	keyword := "%" + search + "%"

	err := database.DB.Get(
		&total,
		`SELECT COUNT(*)
		FROM tbl_kelas
		WHERE nama_kelas LIKE ?`,
		keyword,
	)

	return total, err
}

func GetByID(id int) (*Kelas, error) {
	var data Kelas

	err := database.DB.Get(
		&data,
		`SELECT
			id_kelas,
			nama_kelas,
			aktif
		FROM tbl_kelas
		WHERE id_kelas = ?
		LIMIT 1`,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func ExistsByNama(nama string, excludeID int) (bool, error) {
	var total int

	err := database.DB.Get(
		&total,
		`SELECT COUNT(*)
		FROM tbl_kelas
		WHERE nama_kelas = ?
		AND id_kelas != ?`,
		nama,
		excludeID,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

func Create(data Kelas) (int64, error) {
	result, err := database.DB.Exec(
		`INSERT INTO tbl_kelas
			(nama_kelas, aktif)
		VALUES
			(?, ?)`,
		data.Nama,
		data.Aktif,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func Update(id int, data Kelas) error {
	_, err := database.DB.Exec(
		`UPDATE tbl_kelas SET
			nama_kelas = ?
		WHERE id_kelas = ?`,
		data.Nama,
		id,
	)

	return err
}

func Delete(id int) error {
	_, err := database.DB.Exec(
		`DELETE FROM tbl_kelas WHERE id_kelas = ?`,
		id,
	)

	return err
}

func ToggleStatus(id int, aktif int) error {
	_, err := database.DB.Exec(
		`UPDATE tbl_kelas SET aktif = ? WHERE id_kelas = ?`,
		aktif,
		id,
	)

	return err
}
