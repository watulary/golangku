package pengumuman

import (
	"sekolah/internal/database"
	"strings"
)

const sqlTanggalIndonesia = `
	CONCAT(
		DAY(b.tanggal_publish), ' ',
		CASE MONTH(b.tanggal_publish)
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
		' ', YEAR(b.tanggal_publish)
	)
`
const sqlTanggalUbahIndonesia = `
	CONCAT(
		DAY(b.tanggal), ' ',
		CASE MONTH(b.tanggal)
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
		' ', YEAR(b.tanggal)
	)
`
const sqlHitsFormatted = `
	CASE
		WHEN b.hits < 1000 THEN CAST(b.hits AS CHAR)
		WHEN b.hits < 1000000 THEN CONCAT(
			TRIM(TRAILING '.0' FROM CAST(ROUND(b.hits / 1000, 1) AS CHAR)), 'K'
		)
		WHEN b.hits < 1000000000 THEN CONCAT(
			TRIM(TRAILING '.0' FROM CAST(ROUND(b.hits / 1000000, 1) AS CHAR)), 'M'
		)
		WHEN b.hits < 1000000000000 THEN CONCAT(
			TRIM(TRAILING '.0' FROM CAST(ROUND(b.hits / 1000000000, 1) AS CHAR)), 'B'
		)
		ELSE CONCAT(
			TRIM(TRAILING '.0' FROM CAST(ROUND(b.hits / 1000000000000, 1) AS CHAR)), 'T'
		)
	END
`

var allowedSortColumns = map[string]string{
	"judul":        "b.judul_pengumuman",
	"pengguna":     "p.nama_pengguna",
	"status":       "b.status_pengumuman",
	"hits":         "b.hits",
	"tanggal":      "b.tanggal_publish",
	"tanggal_ubah": "b.tanggal",
}

func GetLatest(limit int) ([]Pengumuman, error) {

	var data []Pengumuman

	err := database.DB.Select(
		&data,
		`SELECT
			id_pengumuman,
			judul_pengumuman,
			slug_pengumuman,
			ringkasan,
			foto,
			hits,
			tanggal_publish
		FROM pengumuman
		WHERE status_pengumuman = 'Publish'
		ORDER BY tanggal_publish DESC
		LIMIT ?`,
		limit,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetAll() ([]Pengumuman, error) {

	var data []Pengumuman

	err := database.DB.Select(
		&data,
		`SELECT
			id_pengumuman,
			judul_pengumuman,
			slug_pengumuman,
			ringkasan,
			foto,
			hits,
			tanggal_publish
		FROM pengumuman
		WHERE status_pengumuman='Publish'
		ORDER BY tanggal_publish DESC`,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetBySlug(slug string) (*Pengumuman, error) {
	var data Pengumuman

	err := database.DB.Get(
		&data,
		`SELECT
    b.id_pengumuman,
    b.judul_pengumuman,
    b.slug_pengumuman,
    b.ringkasan,
    b.isi_pengumuman,
    b.foto,
    b.hits,
    b.tanggal_publish,
    p.nama_pengguna
FROM pengumuman b
LEFT JOIN pengguna p
    ON p.id_pengguna = b.id_pengguna
WHERE b.slug_pengumuman = ?
AND b.status_pengumuman = 'Publish'
LIMIT 1`,
		slug,
	)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func IncrementHits(idPengumuman string) error {

	_, err := database.DB.Exec(
		`UPDATE pengumuman
		SET hits = hits + 1
		WHERE id_pengumuman = ?`,
		idPengumuman,
	)

	return err
}

func GetAllAdmin(search, idPengguna, sortBy, sortDir string, limit, offset int) ([]Pengumuman, error) {

	var data []Pengumuman

	query := `
		SELECT
			b.id_pengumuman,
			b.judul_pengumuman,
			b.slug_pengumuman,
			b.ringkasan,
			b.foto,
			b.hits,
			b.tanggal_publish,
			b.tanggal,
			b.status_pengumuman,
			p.nama_pengguna
		FROM pengumuman b
		LEFT JOIN pengguna p ON p.id_pengguna = b.id_pengguna
		WHERE (
			b.judul_pengumuman LIKE ?
			OR p.nama_pengguna LIKE ?
			OR b.status_pengumuman LIKE ?
			OR CAST(b.hits AS CHAR) LIKE ?
			OR ` + sqlHitsFormatted + ` LIKE ?
			OR ` + sqlTanggalIndonesia + ` LIKE ?
			OR ` + sqlTanggalUbahIndonesia + ` LIKE ?
			OR DATE_FORMAT(b.tanggal_publish, '%d-%m-%Y %H:%i:%s') LIKE ?
			OR DATE_FORMAT(b.tanggal, '%d-%m-%Y %H:%i:%s') LIKE ?
		)`

	keyword := "%" + search + "%"

	args := []interface{}{
		keyword, keyword, keyword, keyword,
		keyword, keyword, keyword, keyword, keyword,
	}

	if idPengguna != "" {
		query += ` AND b.id_pengguna = ?`
		args = append(args, idPengguna)
	}

	column, ok := allowedSortColumns[sortBy]
	if !ok {
		column = "b.tanggal_publish"
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
		FROM pengumuman b
		LEFT JOIN pengguna p ON p.id_pengguna = b.id_pengguna
		WHERE (
			b.judul_pengumuman LIKE ?
			OR p.nama_pengguna LIKE ?
			OR b.status_pengumuman LIKE ?
			OR CAST(b.hits AS CHAR) LIKE ?
			OR ` + sqlHitsFormatted + ` LIKE ?
			OR ` + sqlTanggalIndonesia + ` LIKE ?
			OR ` + sqlTanggalUbahIndonesia + ` LIKE ?
			OR DATE_FORMAT(b.tanggal_publish, '%d-%m-%Y %H:%i:%s') LIKE ?
			OR DATE_FORMAT(b.tanggal, '%d-%m-%Y %H:%i:%s') LIKE ?
		)`

	keyword := "%" + search + "%"

	args := []interface{}{keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword}

	if idPengguna != "" {
		query += ` AND b.id_pengguna = ?`
		args = append(args, idPengguna)
	}

	err := database.DB.Get(&total, query, args...)

	return total, err
}

func GetByID(idPengumuman string) (*Pengumuman, error) {

	var data Pengumuman

	err := database.DB.Get(
		&data,
		`SELECT
			id_pengumuman,
			judul_pengumuman,
			slug_pengumuman,
			isi_pengumuman,
			ringkasan,
			foto,
			hits,
			tanggal_publish,
			status_pengumuman,
			id_pengguna
		FROM pengumuman
		WHERE id_pengumuman = ?
		LIMIT 1`,
		idPengumuman,
	)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func Create(data Pengumuman) error {

	_, err := database.DB.Exec(
		`INSERT INTO pengumuman
			(id_pengumuman, judul_pengumuman, slug_pengumuman, isi_pengumuman, ringkasan,
			 foto, tanggal_publish, status_pengumuman, id_pengguna)
		VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		data.IDPengumuman,
		data.JudulPengumuman,
		data.SlugPengumuman,
		data.IsiPengumuman,
		data.Ringkasan,
		data.Foto,
		data.TanggalPublish,
		data.StatusPengumuman,
		data.IDPengguna,
	)

	return err
}

func Update(idPengumuman string, data Pengumuman) error {

	_, err := database.DB.Exec(
		`UPDATE pengumuman SET
			judul_pengumuman = ?,
			slug_pengumuman = ?,
			isi_pengumuman = ?,
			ringkasan = ?,
			status_pengumuman = ?
		WHERE id_pengumuman = ?`,
		data.JudulPengumuman,
		data.SlugPengumuman,
		data.IsiPengumuman,
		data.Ringkasan,
		data.StatusPengumuman,
		idPengumuman,
	)

	return err
}

func UpdateFoto(idPengumuman, foto string) error {

	_, err := database.DB.Exec(
		`UPDATE pengumuman SET foto = ? WHERE id_pengumuman = ?`,
		foto,
		idPengumuman,
	)

	return err
}

func Delete(idPengumuman string) error {

	_, err := database.DB.Exec(
		`DELETE FROM pengumuman WHERE id_pengumuman = ?`,
		idPengumuman,
	)

	return err
}
