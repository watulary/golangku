package suratmasuk

import (
	"database/sql"
	"fmt"
	"strings"

	"sekolah/internal/database"
)

var allowedSortColumns = map[string]string{
	"nama":           "nama",
	"nomor":          "nomor",
	"tanggal":        "tanggal",
	"tanggal_terima": "tanggal_terima",
	"perihal":        "perihal",
	"asal":           "asal",
	"tujuan":         "tujuan",
	"waktu":          "waktu",
}

func GetAllAdmin(search, idPengguna, sortBy, sortDir string, limit, offset int) ([]SuratMasuk, error) {
	var data []SuratMasuk

	query := `
        SELECT sm.*, p.nama_pengguna
        FROM suratmasuk sm
        LEFT JOIN pengguna p ON p.id_pengguna = sm.id_pengguna
        WHERE 1=1`
	args := []interface{}{}

	if search != "" {
		query += ` AND (
        sm.nama LIKE ? OR
        sm.nomor LIKE ? OR
        sm.perihal LIKE ? OR
        sm.asal LIKE ? OR
        sm.tujuan LIKE ? OR
        DATE_FORMAT(sm.tanggal, '%d-%m-%Y') LIKE ? OR
        DATE_FORMAT(sm.tanggal_terima, '%d-%m-%Y') LIKE ? OR
        p.nama_pengguna LIKE ?
    )`
		s := "%" + search + "%"
		args = append(args, s, s, s, s, s, s, s, s)
	}

	if idPengguna != "" {
		query += ` AND sm.id_pengguna = ?`
		args = append(args, idPengguna)
	}

	col, ok := allowedSortColumns[sortBy]
	if !ok {
		col = "waktu"
	}
	dir := "DESC"
	if strings.ToUpper(sortDir) == "ASC" {
		dir = "ASC"
	}
	query += fmt.Sprintf(" ORDER BY sm.%s %s", col, dir)

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	err := database.DB.Select(&data, query, args...)
	return data, err
}

func CountAdmin(search, idPengguna string) (int, error) {
	var total int

	query := `
        SELECT COUNT(*)
        FROM suratmasuk sm
        LEFT JOIN pengguna p ON p.id_pengguna = sm.id_pengguna
        WHERE 1=1`
	args := []interface{}{}

	if search != "" {
		query += ` AND (
        sm.nama LIKE ? OR
        sm.nomor LIKE ? OR
        sm.perihal LIKE ? OR
        sm.asal LIKE ? OR
        sm.tujuan LIKE ? OR
        DATE_FORMAT(sm.tanggal, '%d-%m-%Y') LIKE ? OR
        DATE_FORMAT(sm.tanggal_terima, '%d-%m-%Y') LIKE ? OR
        p.nama_pengguna LIKE ?
    )`
		s := "%" + search + "%"
		args = append(args, s, s, s, s, s, s, s, s)
	}

	if idPengguna != "" {
		query += ` AND sm.id_pengguna = ?`
		args = append(args, idPengguna)
	}

	err := database.DB.Get(&total, query, args...)
	return total, err
}

// GetByID mengambil satu surat masuk berdasarkan id.
func GetByID(id string) (*SuratMasuk, error) {
	var data SuratMasuk
	err := database.DB.Get(&data, `
		SELECT sm.*, p.nama_pengguna
		FROM suratmasuk sm
		LEFT JOIN pengguna p ON p.id_pengguna = sm.id_pengguna
		WHERE sm.id = ?`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

// Create menyimpan data surat masuk baru (termasuk nama file dokumen).
func Create(data SuratMasuk) error {
	_, err := database.DB.Exec(`
		INSERT INTO suratmasuk
			(id, id_pengguna, nama, nomor, tanggal, tanggal_terima, perihal, asal, tujuan, dokumen)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		data.ID, data.IDPengguna, data.Nama, data.Nomor, data.Tanggal,
		data.TanggalTerima, data.Perihal, data.Asal, data.Tujuan, data.Dokumen,
	)
	return err
}

// Update memperbarui data teks surat masuk (tanpa mengubah dokumen).
func Update(id string, data SuratMasuk) error {
	_, err := database.DB.Exec(`
		UPDATE suratmasuk
		SET nama = ?, nomor = ?, tanggal = ?, tanggal_terima = ?, perihal = ?, asal = ?, tujuan = ?
		WHERE id = ?`,
		data.Nama, data.Nomor, data.Tanggal, data.TanggalTerima,
		data.Perihal, data.Asal, data.Tujuan, id,
	)
	return err
}

// UpdateDokumen memperbarui nama file dokumen surat masuk (terpisah dari Update biasa).
func UpdateDokumen(id, dokumen string) error {
	_, err := database.DB.Exec(`UPDATE suratmasuk SET dokumen = ? WHERE id = ?`, dokumen, id)
	return err
}

// Delete menghapus satu surat masuk berdasarkan id.
func Delete(id string) error {
	_, err := database.DB.Exec(`DELETE FROM suratmasuk WHERE id = ?`, id)
	return err
}
