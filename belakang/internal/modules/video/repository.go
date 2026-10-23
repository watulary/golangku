package video

import (
	"sekolah/internal/database"
	"strings"
)

const sqlTanggalIndonesia = `
	CONCAT(
		DAY(b.waktu), ' ',
		CASE MONTH(b.waktu)
			WHEN 1 THEN 'Januari'
			WHEN 2 THEN 'Februari'
			WHEN 3 THEN 'Maret'
			WHEN 4 THEN 'April'
			WHEN 5 THEN 'Mei'
			WHEN 6 THEN 'Juni'
			WHEN 7 THEN 'Juli'
			WHEN 8 THEN 'Agustus'
			WHEN 9 THEN 'September'
			WHEN 10 THEN 'Oktober'
			WHEN 11 THEN 'November'
			WHEN 12 THEN 'Desember'
		END,
		' ', YEAR(b.waktu)
	)
`

var allowedSortColumns = map[string]string{
	"nama":     "b.nama",
	"kegiatan": "b.kegiatan",
	"pengguna": "p.nama_pengguna",
	"waktu":    "b.waktu",
}

func GetLatest(limit int) ([]Video, error) {

	var data []Video

	err := database.DB.Select(
		&data,
		`SELECT
			id,
			nama,
			kegiatan,
			video,
			waktu
		FROM video
		ORDER BY waktu DESC
		LIMIT ?`,
		limit,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetAll() ([]Video, error) {

	var data []Video

	err := database.DB.Select(
		&data,
		`SELECT
			id,
			nama,
			kegiatan,
			video,
			waktu
		FROM video
		ORDER BY waktu DESC`,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetAllAdmin(search, idPengguna, sortBy, sortDir string, limit, offset int) ([]Video, error) {

	var data []Video

	query := `
		SELECT
			b.id,
			b.nama,
			b.kegiatan,
			b.video,
			b.waktu,
			b.id_pengguna,
			p.nama_pengguna
		FROM video b
		LEFT JOIN pengguna p ON p.id_pengguna = b.id_pengguna
		WHERE (
			b.nama LIKE ?
			OR b.kegiatan LIKE ?
			OR p.nama_pengguna LIKE ?
			OR ` + sqlTanggalIndonesia + ` LIKE ?
			OR DATE_FORMAT(b.waktu, '%d-%m-%Y %H:%i:%s') LIKE ?
		)`

	keyword := "%" + search + "%"

	args := []interface{}{
		keyword, keyword, keyword, keyword, keyword,
	}

	if idPengguna != "" {
		query += ` AND b.id_pengguna = ?`
		args = append(args, idPengguna)
	}

	column, ok := allowedSortColumns[sortBy]
	if !ok {
		column = "b.waktu"
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

func CountAdmin(search, idPengguna string) (int, error) {

	var total int

	query := `
		SELECT COUNT(*)
		FROM video b
		LEFT JOIN pengguna p ON p.id_pengguna = b.id_pengguna
		WHERE (
			b.nama LIKE ?
			OR b.kegiatan LIKE ?
			OR p.nama_pengguna LIKE ?
			OR ` + sqlTanggalIndonesia + ` LIKE ?
			OR DATE_FORMAT(b.waktu, '%d-%m-%Y %H:%i:%s') LIKE ?
		)`

	keyword := "%" + search + "%"

	args := []interface{}{keyword, keyword, keyword, keyword, keyword}

	if idPengguna != "" {
		query += ` AND b.id_pengguna = ?`
		args = append(args, idPengguna)
	}

	err := database.DB.Get(&total, query, args...)

	return total, err
}

func GetByID(id string) (*Video, error) {

	var data Video

	err := database.DB.Get(
		&data,
		`SELECT
			b.id,
			b.nama,
			b.kegiatan,
			b.video,
			b.waktu,
			b.id_pengguna,
			p.nama_pengguna
		FROM video b
		LEFT JOIN pengguna p ON p.id_pengguna = b.id_pengguna
		WHERE b.id = ?
		LIMIT 1`,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func Create(data Video) error {

	_, err := database.DB.Exec(
		`INSERT INTO video
			(id, nama, kegiatan, video, id_pengguna)
		VALUES
			(?, ?, ?, ?, ?)`,
		data.ID,
		data.Nama,
		data.Kegiatan,
		data.Video,
		data.IDPengguna,
	)

	return err
}

func Update(id string, data Video) error {

	_, err := database.DB.Exec(
		`UPDATE video SET
			nama = ?,
			kegiatan = ?,
			video = ?
		WHERE id = ?`,
		data.Nama,
		data.Kegiatan,
		data.Video,
		id,
	)

	return err
}

func Delete(id string) error {

	_, err := database.DB.Exec(
		`DELETE FROM video WHERE id = ?`,
		id,
	)

	return err
}
