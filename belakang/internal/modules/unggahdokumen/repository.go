package unggahdokumen

import (
	"sekolah/internal/config"
	"sekolah/internal/database"
	"strings"
)

func Get() ([]Unggahdokumen, error) {

	var data []Unggahdokumen

	err := database.DB.Select(
		&data,
		`
		SELECT
			u.id,
			u.id_pengguna,
			p.nama_pengguna,
			u.nama,
			u.dokumen,
			u.waktu
		FROM unggahdokumen u
		LEFT JOIN pengguna p
			ON p.id_pengguna = u.id_pengguna
		ORDER BY u.waktu DESC
		`,
	)

	if err != nil {
		return nil, err
	}

	for i := range data {
		if data[i].Dokumen != "" {
			data[i].DokumenURL =
				config.App.AppURL +
					"/storage/dokumen/" +
					data[i].Dokumen
		}
	}

	return data, nil
}

// GetByID mengambil satu dokumen berdasarkan id
func GetByID(id string) (*Unggahdokumen, error) {

	var data Unggahdokumen

	err := database.DB.Get(
		&data,
		`
		SELECT
			u.id,
			u.id_pengguna,
			p.nama_pengguna,
			u.nama,
			u.dokumen,
			u.waktu
		FROM unggahdokumen u
		LEFT JOIN pengguna p
			ON p.id_pengguna = u.id_pengguna
		WHERE u.id = ?
		`,
		id,
	)

	if err != nil {
		return nil, err
	}

	if data.Dokumen != "" {
		data.DokumenURL = config.App.AppURL + "/storage/dokumen/" + data.Dokumen
	}

	return &data, nil
}

var allowedSortColumnsUnggahdokumen = map[string]string{
	"nama":          "u.nama",
	"waktu":         "u.waktu",
	"nama_pengguna": "p.nama_pengguna",
}

// GetAdmin mengambil daftar dokumen untuk admin (pagination + search + sort)
func GetAdmin(search, idPengguna, sortBy, sortDir string, limit, offset int) ([]Unggahdokumen, error) {

	var data []Unggahdokumen

	query := `
		SELECT
			u.id,
			u.id_pengguna,
			p.nama_pengguna,
			u.nama,
			u.dokumen,
			u.waktu
		FROM unggahdokumen u
		LEFT JOIN pengguna p
			ON p.id_pengguna = u.id_pengguna
	`

	args := []interface{}{}
	var conditions []string

	if search != "" {
		conditions = append(conditions, "(u.nama LIKE ? OR p.nama_pengguna LIKE ?)")
		term := "%" + search + "%"
		args = append(args, term, term)
	}

	if idPengguna != "" {
		conditions = append(conditions, "u.id_pengguna = ?")
		args = append(args, idPengguna)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	sortColumn, ok := allowedSortColumnsUnggahdokumen[sortBy]
	if !ok {
		sortColumn = "u.waktu"
	}

	dir := "DESC"
	if strings.EqualFold(sortDir, "asc") {
		dir = "ASC"
	}

	query += " ORDER BY " + sortColumn + " " + dir
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	err := database.DB.Select(&data, query, args...)

	if err != nil {
		return nil, err
	}

	for i := range data {
		if data[i].Dokumen != "" {
			data[i].DokumenURL = config.App.AppURL + "/storage/dokumen/" + data[i].Dokumen
		}
	}

	return data, nil
}

// CountAdmin menghitung total dokumen untuk admin (dipakai pagination)
func CountAdmin(search, idPengguna string) (int, error) {

	var count int

	query := `
		SELECT COUNT(*)
		FROM unggahdokumen u
		LEFT JOIN pengguna p
			ON p.id_pengguna = u.id_pengguna
	`

	args := []interface{}{}
	var conditions []string

	if search != "" {
		conditions = append(conditions, "(u.nama LIKE ? OR p.nama_pengguna LIKE ?)")
		term := "%" + search + "%"
		args = append(args, term, term)
	}

	if idPengguna != "" {
		conditions = append(conditions, "u.id_pengguna = ?")
		args = append(args, idPengguna)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	err := database.DB.Get(&count, query, args...)

	return count, err
}

// Create menyimpan dokumen baru
func Create(data *Unggahdokumen) error {

	_, err := database.DB.Exec(
		`
		INSERT INTO unggahdokumen (id, id_pengguna, nama, dokumen)
		VALUES (?, ?, ?, ?)
		`,
		data.ID,
		data.IDPengguna,
		data.Nama,
		data.Dokumen,
	)

	return err
}

// Update mengubah nama dokumen (tanpa mengubah file)
func Update(id, nama string) error {

	_, err := database.DB.Exec(
		`UPDATE unggahdokumen SET nama = ? WHERE id = ?`,
		nama,
		id,
	)

	return err
}

// UpdateDokumen mengganti file dokumen
func UpdateDokumen(id, dokumen string) error {

	_, err := database.DB.Exec(
		`UPDATE unggahdokumen SET dokumen = ? WHERE id = ?`,
		dokumen,
		id,
	)

	return err
}

// Delete menghapus dokumen
func Delete(id string) error {

	_, err := database.DB.Exec(
		`DELETE FROM unggahdokumen WHERE id = ?`,
		id,
	)

	return err
}
