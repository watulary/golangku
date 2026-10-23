package fotoekskul

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

func GetLatest(limit int) ([]FotoEkskul, error) {

	var data []FotoEkskul

	err := database.DB.Select(
		&data,
		`SELECT
			id,
			nama,
			kegiatan,
			foto,
			waktu
		FROM fotoekskul
		ORDER BY waktu DESC
		LIMIT ?`,
		limit,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetAll() ([]FotoEkskul, error) {

	var data []FotoEkskul

	err := database.DB.Select(
		&data,
		`SELECT
			id,
			nama,
			kegiatan,
			foto,
			waktu
		FROM fotoekskul
		ORDER BY waktu DESC`,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetAllAdmin(search, idPengguna, sortBy, sortDir string, limit, offset int) ([]FotoEkskul, error) {

	var data []FotoEkskul

	query := `
		SELECT
			b.id,
			b.nama,
			b.kegiatan,
			b.foto,
			b.waktu,
			b.id_pengguna,
			p.nama_pengguna
		FROM fotoekskul b
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
		FROM fotoekskul b
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

func GetByID(id string) (*FotoEkskul, error) {

	var data FotoEkskul

	err := database.DB.Get(
		&data,
		`SELECT
			b.id,
			b.nama,
			b.kegiatan,
			b.foto,
			b.waktu,
			b.id_pengguna,
			p.nama_pengguna
		FROM fotoekskul b
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

func Create(data FotoEkskul) error {

	_, err := database.DB.Exec(
		`INSERT INTO fotoekskul
			(id, nama, kegiatan, foto, id_pengguna)
		VALUES
			(?, ?, ?, ?, ?)`,
		data.ID,
		data.Nama,
		data.Kegiatan,
		data.Foto,
		data.IDPengguna,
	)

	return err
}

func Update(id string, data FotoEkskul) error {

	_, err := database.DB.Exec(
		`UPDATE fotoekskul SET
			nama = ?,
			kegiatan = ?
		WHERE id = ?`,
		data.Nama,
		data.Kegiatan,
		id,
	)

	return err
}

func UpdateFoto(id, foto string) error {

	_, err := database.DB.Exec(
		`UPDATE fotoekskul SET foto = ? WHERE id = ?`,
		foto,
		id,
	)

	return err
}

func Delete(id string) error {

	_, err := database.DB.Exec(
		`DELETE FROM fotoekskul WHERE id = ?`,
		id,
	)

	return err
}
