package spmb

import (
	"sekolah/internal/database"

	"github.com/jmoiron/sqlx"
)

var allowedSortColumns = map[string]bool{
	"no_peserta":   true,
	"nama_peserta": true,
	"kelamin":      true,
	"asal_sekolah": true,
	"jalur":        true,
	"lapor_diri":   true,
}

func Get() ([]SPMB, error) {
	var data []SPMB
	err := database.DB.Select(
		&data,
		`
		SELECT
			no_peserta,
			nama_peserta,
			kelamin,
			asal_sekolah,
			jalur,
			lapor_diri
		FROM ppdb
		ORDER BY nama_peserta ASC
		`,
	)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetAllAdmin(search string, jalur int, laporDiri int, sortBy string, sortDir string, limit int, offset int) ([]SPMB, error) {
	if !allowedSortColumns[sortBy] {
		sortBy = "nama_peserta"
	}
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "asc"
	}

	query := `
		SELECT
			no_peserta,
			nama_peserta,
			kelamin,
			asal_sekolah,
			jalur,
			lapor_diri
		FROM ppdb
		WHERE (no_peserta LIKE ? OR nama_peserta LIKE ? OR asal_sekolah LIKE ?)
	`
	args := []interface{}{"%" + search + "%", "%" + search + "%", "%" + search + "%"}

	if jalur != 0 {
		query += " AND jalur = ?"
		args = append(args, jalur)
	}
	if laporDiri != -1 {
		query += " AND lapor_diri = ?"
		args = append(args, laporDiri)
	}

	query += " ORDER BY " + sortBy + " " + sortDir + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var data []SPMB
	err := database.DB.Select(&data, query, args...)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func CountAdmin(search string, jalur int, laporDiri int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM ppdb
		WHERE (no_peserta LIKE ? OR nama_peserta LIKE ? OR asal_sekolah LIKE ?)
	`
	args := []interface{}{"%" + search + "%", "%" + search + "%", "%" + search + "%"}

	if jalur != 0 {
		query += " AND jalur = ?"
		args = append(args, jalur)
	}
	if laporDiri != -1 {
		query += " AND lapor_diri = ?"
		args = append(args, laporDiri)
	}

	var count int
	err := database.DB.Get(&count, query, args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetByID(noPeserta string) (*SPMB, error) {
	var data SPMB
	err := database.DB.Get(
		&data,
		`
		SELECT
			no_peserta,
			nama_peserta,
			kelamin,
			asal_sekolah,
			jalur,
			lapor_diri
		FROM ppdb
		WHERE no_peserta = ?
		`,
		noPeserta,
	)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func Create(data SPMB) error {
	_, err := database.DB.NamedExec(
		`
		INSERT INTO ppdb (no_peserta, nama_peserta, kelamin, asal_sekolah, jalur, lapor_diri)
		VALUES (:no_peserta, :nama_peserta, :kelamin, :asal_sekolah, :jalur, :lapor_diri)
		`,
		data,
	)
	return err
}

func Update(noPeserta string, data SPMB) error {
	data.NoPeserta = noPeserta
	_, err := database.DB.NamedExec(
		`
		UPDATE ppdb SET
			nama_peserta = :nama_peserta,
			kelamin = :kelamin,
			asal_sekolah = :asal_sekolah,
			jalur = :jalur,
			lapor_diri = :lapor_diri
		WHERE no_peserta = :no_peserta
		`,
		data,
	)
	return err
}

func Delete(noPeserta string) error {
	_, err := database.DB.Exec(`DELETE FROM ppdb WHERE no_peserta = ?`, noPeserta)
	return err
}

func GetAllForExport(search string, jalur int, laporDiri int) ([]SPMB, error) {
	query := `
		SELECT
			no_peserta,
			nama_peserta,
			kelamin,
			asal_sekolah,
			jalur,
			lapor_diri
		FROM ppdb
		WHERE (no_peserta LIKE ? OR nama_peserta LIKE ? OR asal_sekolah LIKE ?)
	`
	args := []interface{}{"%" + search + "%", "%" + search + "%", "%" + search + "%"}

	if jalur != 0 {
		query += " AND jalur = ?"
		args = append(args, jalur)
	}
	if laporDiri != -1 {
		query += " AND lapor_diri = ?"
		args = append(args, laporDiri)
	}

	query += " ORDER BY nama_peserta ASC"

	var data []SPMB
	err := database.DB.Select(&data, query, args...)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetExistingNoPeserta mengembalikan no_peserta yang SUDAH ada di DB dari daftar yang diberikan
func GetExistingNoPeserta(noPesertaList []string) (map[string]bool, error) {
	existing := make(map[string]bool)
	if len(noPesertaList) == 0 {
		return existing, nil
	}

	query, args, err := sqlx.In(
		`SELECT no_peserta FROM ppdb WHERE no_peserta IN (?)`,
		noPesertaList,
	)
	if err != nil {
		return nil, err
	}
	query = database.DB.Rebind(query)

	var found []string
	if err := database.DB.Select(&found, query, args...); err != nil {
		return nil, err
	}
	for _, np := range found {
		existing[np] = true
	}
	return existing, nil
}

// ImportSPMB insert semua data sekaligus dalam 1 transaction (all-or-nothing)
func ImportSPMB(data []SPMB) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, d := range data {
		_, err := tx.NamedExec(
			`
			INSERT INTO ppdb (no_peserta, nama_peserta, kelamin, asal_sekolah, jalur, lapor_diri)
			VALUES (:no_peserta, :nama_peserta, :kelamin, :asal_sekolah, :jalur, :lapor_diri)
			`,
			d,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// BulkUpdateLaporDiri mengubah status lapor_diri utk banyak peserta sekaligus
func BulkUpdateLaporDiri(noPesertaList []string, laporDiri int) error {
	if len(noPesertaList) == 0 {
		return nil
	}
	query, args, err := sqlx.In(
		`UPDATE ppdb SET lapor_diri = ? WHERE no_peserta IN (?)`,
		laporDiri, noPesertaList,
	)
	if err != nil {
		return err
	}
	query = database.DB.Rebind(query)
	_, err = database.DB.Exec(query, args...)
	return err
}

// BulkDelete menghapus banyak peserta sekaligus
func BulkDelete(noPesertaList []string) error {
	if len(noPesertaList) == 0 {
		return nil
	}
	query, args, err := sqlx.In(
		`DELETE FROM ppdb WHERE no_peserta IN (?)`,
		noPesertaList,
	)
	if err != nil {
		return err
	}
	query = database.DB.Rebind(query)
	_, err = database.DB.Exec(query, args...)
	return err
}
