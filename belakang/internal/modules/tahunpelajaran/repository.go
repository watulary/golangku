package tahunpelajaran

import (
	"sekolah/internal/database"
	"strings"
)

var allowedSortColumns = map[string]string{
	"nama":     "nama_tahunpelajaran",
	"semester": "semester",
	"aktif":    "aktif",
}

func GetAllAdmin(search, sortBy, sortDir string, limit, offset int) ([]TahunPelajaran, error) {
	var data []TahunPelajaran

	query := `
			SELECT
			id_tahunpelajaran,
			nama_tahunpelajaran,
			semester,
			aktif
		FROM tbl_tahunpelajaran
		WHERE nama_tahunpelajaran LIKE ?`

	keyword := "%" + search + "%"

	args := []interface{}{keyword}

	column, ok := allowedSortColumns[sortBy]
	if !ok {
		column = "nama_tahunpelajaran"
	}

	direction := "DESC"
	if strings.ToLower(sortDir) == "asc" {
		direction = "ASC"
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
		FROM tbl_tahunpelajaran
		WHERE nama_tahunpelajaran LIKE ?`,
		keyword,
	)

	return total, err
}

func GetByID(id int) (*TahunPelajaran, error) {
	var data TahunPelajaran

	err := database.DB.Get(
		&data,
		`SELECT
			id_tahunpelajaran,
			nama_tahunpelajaran,
			semester,
			aktif
		FROM tbl_tahunpelajaran
		WHERE id_tahunpelajaran = ?
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
		FROM tbl_tahunpelajaran
		WHERE nama_tahunpelajaran = ?
		AND id_tahunpelajaran != ?`,
		nama,
		excludeID,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

func Create(data TahunPelajaran) (int64, error) {
	result, err := database.DB.Exec(
		`INSERT INTO tbl_tahunpelajaran
			(nama_tahunpelajaran, semester, aktif)
		VALUES
			(?, ?, ?)`,
		data.Nama,
		data.Semester,
		data.Aktif,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func Update(id int, data TahunPelajaran) error {
	_, err := database.DB.Exec(
		`UPDATE tbl_tahunpelajaran SET
			nama_tahunpelajaran = ?,
			semester = ?
		WHERE id_tahunpelajaran = ?`,
		data.Nama,
		data.Semester,
		id,
	)

	return err
}

func Delete(id int) error {
	_, err := database.DB.Exec(
		`DELETE FROM tbl_tahunpelajaran WHERE id_tahunpelajaran = ?`,
		id,
	)

	return err
}

func SetActive(id int) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`UPDATE tbl_tahunpelajaran SET aktif = 0`); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(
		`UPDATE tbl_tahunpelajaran SET aktif = 1 WHERE id_tahunpelajaran = ?`,
		id,
	); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
