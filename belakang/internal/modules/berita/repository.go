package berita

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
	"judul":        "b.judul_berita",
	"kategori":     "kb.nama_kategoriberita",
	"pengguna":     "p.nama_pengguna",
	"status":       "b.status_berita",
	"hits":         "b.hits",
	"tanggal":      "b.tanggal_publish",
	"tanggal_ubah": "b.tanggal",
}

func GetLatest(limit int) ([]Berita, error) {

	var data []Berita

	err := database.DB.Select(
		&data,
		`SELECT
			id_berita,
			judul_berita,
			slug_berita,
			ringkasan,
			foto,
			hits,
			tanggal_publish
		FROM berita
		WHERE status_berita = 'Publish'
		ORDER BY tanggal_publish DESC
		LIMIT ?`,
		limit,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetAll() ([]Berita, error) {

	var data []Berita

	err := database.DB.Select(
		&data,
		`SELECT
			id_berita,
			judul_berita,
			slug_berita,
			ringkasan,
			foto,
			hits,
			tanggal_publish
		FROM berita
		WHERE status_berita='Publish'
		ORDER BY tanggal_publish DESC`,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetBySlug(slug string) (*Berita, error) {
	var data Berita

	err := database.DB.Get(
		&data,
		`SELECT
    b.id_berita,
    b.judul_berita,
    b.slug_berita,
    b.ringkasan,
    b.isi_berita,
    b.foto,
    b.hits,
    b.tanggal_publish,
    p.nama_pengguna
FROM berita b
LEFT JOIN pengguna p
    ON p.id_pengguna = b.id_pengguna
WHERE b.slug_berita = ?
AND b.status_berita = 'Publish'
LIMIT 1`,
		slug,
	)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func IncrementHits(idBerita string) error {

	_, err := database.DB.Exec(
		`UPDATE berita
		SET hits = hits + 1
		WHERE id_berita = ?`,
		idBerita,
	)

	return err
}

func GetByKategori(idKategori string, limit int) ([]Berita, error) {

	var data []Berita

	err := database.DB.Select(
		&data,
		`SELECT
			b.id_berita,
			b.id_kategoriberita,
			kb.nama_kategoriberita,
			kb.slug_kategoriberita,
			b.judul_berita,
			b.slug_berita,
			b.ringkasan,
			b.isi_berita,
			b.foto,
			b.hits,
			b.tanggal_publish,
			p.nama_pengguna
		FROM berita b
		LEFT JOIN pengguna p
			ON p.id_pengguna = b.id_pengguna
		LEFT JOIN kategori_berita kb
			ON kb.id_kategoriberita = b.id_kategoriberita
		WHERE b.status_berita = 'Publish'
		AND b.id_kategoriberita = ?
		ORDER BY b.tanggal_publish DESC
		LIMIT ?`,
		idKategori,
		limit,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetAllAdmin(search, idPengguna, sortBy, sortDir string, limit, offset int) ([]Berita, error) {

	var data []Berita

	query := `
		SELECT
			b.id_berita,
			b.judul_berita,
			b.slug_berita,
			b.ringkasan,
			b.foto,
			b.hits,
			b.tanggal_publish,
			b.tanggal,
			b.status_berita,
			b.id_kategoriberita,
			kb.nama_kategoriberita,
			p.nama_pengguna
		FROM berita b
		LEFT JOIN pengguna p ON p.id_pengguna = b.id_pengguna
		LEFT JOIN kategori_berita kb ON kb.id_kategoriberita = b.id_kategoriberita
		WHERE (
			b.judul_berita LIKE ?
			OR kb.nama_kategoriberita LIKE ?
			OR p.nama_pengguna LIKE ?
			OR b.status_berita LIKE ?
			OR CAST(b.hits AS CHAR) LIKE ?
			OR ` + sqlHitsFormatted + ` LIKE ?
			OR ` + sqlTanggalIndonesia + ` LIKE ?
			OR ` + sqlTanggalUbahIndonesia + ` LIKE ?
			OR DATE_FORMAT(b.tanggal_publish, '%d-%m-%Y %H:%i:%s') LIKE ?
			OR DATE_FORMAT(b.tanggal, '%d-%m-%Y %H:%i:%s') LIKE ?
		)`

	keyword := "%" + search + "%"

	args := []interface{}{
		keyword, keyword, keyword, keyword, keyword,
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
		FROM berita b
		LEFT JOIN pengguna p ON p.id_pengguna = b.id_pengguna
		LEFT JOIN kategori_berita kb ON kb.id_kategoriberita = b.id_kategoriberita
		WHERE (
			b.judul_berita LIKE ?
			OR kb.nama_kategoriberita LIKE ?
			OR p.nama_pengguna LIKE ?
			OR b.status_berita LIKE ?
			OR CAST(b.hits AS CHAR) LIKE ?
			OR ` + sqlHitsFormatted + ` LIKE ?
			OR ` + sqlTanggalIndonesia + ` LIKE ?
			OR ` + sqlTanggalUbahIndonesia + ` LIKE ?
			OR DATE_FORMAT(b.tanggal_publish, '%d-%m-%Y %H:%i:%s') LIKE ?
			OR DATE_FORMAT(b.tanggal, '%d-%m-%Y %H:%i:%s') LIKE ?
		)`

	keyword := "%" + search + "%"

	args := []interface{}{keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword}

	if idPengguna != "" {
		query += ` AND b.id_pengguna = ?`
		args = append(args, idPengguna)
	}

	err := database.DB.Get(&total, query, args...)

	return total, err
}

func GetByID(idBerita string) (*Berita, error) {

	var data Berita

	err := database.DB.Get(
		&data,
		`SELECT
			id_berita,
			judul_berita,
			slug_berita,
			isi_berita,
			ringkasan,
			foto,
			hits,
			tanggal_publish,
			status_berita,
			id_kategoriberita,
			id_pengguna
		FROM berita
		WHERE id_berita = ?
		LIMIT 1`,
		idBerita,
	)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func Create(data Berita) error {

	_, err := database.DB.Exec(
		`INSERT INTO berita
			(id_berita, judul_berita, slug_berita, isi_berita, ringkasan,
			 foto, tanggal_publish, status_berita, id_kategoriberita, id_pengguna)
		VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		data.IDBerita,
		data.JudulBerita,
		data.SlugBerita,
		data.IsiBerita,
		data.Ringkasan,
		data.Foto,
		data.TanggalPublish,
		data.StatusBerita,
		data.IDKategoriBerita,
		data.IDPengguna,
	)

	return err
}

func Update(idBerita string, data Berita) error {

	_, err := database.DB.Exec(
		`UPDATE berita SET
			judul_berita = ?,
			slug_berita = ?,
			isi_berita = ?,
			ringkasan = ?,
			status_berita = ?,
			id_kategoriberita = ?
		WHERE id_berita = ?`,
		data.JudulBerita,
		data.SlugBerita,
		data.IsiBerita,
		data.Ringkasan,
		data.StatusBerita,
		data.IDKategoriBerita,
		idBerita,
	)

	return err
}

func UpdateFoto(idBerita, foto string) error {

	_, err := database.DB.Exec(
		`UPDATE berita SET foto = ? WHERE id_berita = ?`,
		foto,
		idBerita,
	)

	return err
}

func Delete(idBerita string) error {

	_, err := database.DB.Exec(
		`DELETE FROM berita WHERE id_berita = ?`,
		idBerita,
	)

	return err
}
