package mapel

import (
	"sekolah/internal/database"
	"strings"
)

const selectBase = `
	SELECT
		m.id_mapel,
		m.id_induk,
		m.nama_mapel,
		m.urutan,
		m.aktif,
		p.nama_mapel AS nama_induk
	FROM tbl_mapel m
	LEFT JOIN tbl_mapel p ON p.id_mapel = m.id_induk`

// Urutan "pohon": mapel utama dulu, lalu sub-mapelnya, sesuai kolom urutan.
const treeOrder = `COALESCE(p.urutan, m.urutan), COALESCE(m.id_induk, m.id_mapel), (m.id_induk IS NOT NULL), m.urutan, m.id_mapel`

// sort_by "urutan" (default) = urutan pohon di atas, arah sort diabaikan.
var allowedSortColumns = map[string]string{
	"nama":  "m.nama_mapel",
	"aktif": "m.aktif",
}

const searchWhere = ` WHERE (m.nama_mapel LIKE ? OR p.nama_mapel LIKE ?)`

func GetAllAdmin(search, sortBy, sortDir string, limit, offset int) ([]Mapel, error) {
	data := []Mapel{}

	keyword := "%" + search + "%"

	direction := "ASC"
	if strings.ToLower(sortDir) == "desc" {
		direction = "DESC"
	}

	orderBy := treeOrder
	if column, ok := allowedSortColumns[sortBy]; ok {
		orderBy = column + " " + direction + ", " + treeOrder
	}

	query := selectBase + searchWhere + ` ORDER BY ` + orderBy + ` LIMIT ? OFFSET ?`

	err := database.DB.Select(&data, query, keyword, keyword, limit, offset)
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
		FROM tbl_mapel m
		LEFT JOIN tbl_mapel p ON p.id_mapel = m.id_induk`+searchWhere,
		keyword,
		keyword,
	)

	return total, err
}

// GetInduk = daftar mapel utama (id_induk NULL) untuk dropdown "Induk" di form.
func GetInduk() ([]Mapel, error) {
	data := []Mapel{}

	err := database.DB.Select(
		&data,
		selectBase+` WHERE m.id_induk IS NULL ORDER BY m.urutan, m.nama_mapel`,
	)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetDinilai = mapel aktif yang BISA diisi nilainya, yaitu mapel "daun"
// (tidak punya sub-mapel). Mapel induk seperti "Seni Budaya & Prakarya"
// hanya judul di Buku Induk, nilainya ada di sub-mapel.
// Dipakai untuk dropdown pilih mapel di layar input nilai.
func GetDinilai() ([]Mapel, error) {
	data := []Mapel{}

	err := database.DB.Select(
		&data,
		selectBase+`
		WHERE m.aktif = 1
		AND NOT EXISTS (
			SELECT 1 FROM tbl_mapel c WHERE c.id_induk = m.id_mapel
		)
		ORDER BY `+treeOrder,
	)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetByID(id int) (*Mapel, error) {
	var data Mapel

	err := database.DB.Get(
		&data,
		selectBase+` WHERE m.id_mapel = ? LIMIT 1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// ExistsByNama: nama harus unik di dalam induk yang sama
// (`<=>` = sama dengan yang aman untuk NULL, jadi berlaku juga utk mapel utama).
func ExistsByNama(nama string, idInduk *int, excludeID int) (bool, error) {
	var total int

	err := database.DB.Get(
		&total,
		`SELECT COUNT(*)
		FROM tbl_mapel
		WHERE nama_mapel = ?
		AND id_induk <=> ?
		AND id_mapel != ?`,
		nama,
		idInduk,
		excludeID,
	)
	if err != nil {
		return false, err
	}

	return total > 0, nil
}

func HasChildren(id int) (bool, error) {
	var total int

	err := database.DB.Get(
		&total,
		`SELECT COUNT(*) FROM tbl_mapel WHERE id_induk = ?`,
		id,
	)
	if err != nil {
		return false, err
	}

	return total > 0, nil
}

func Create(data Mapel) (int64, error) {
	result, err := database.DB.Exec(
		`INSERT INTO tbl_mapel
			(id_induk, nama_mapel, urutan, aktif)
		VALUES
			(?, ?, ?, ?)`,
		data.IDInduk,
		data.Nama,
		data.Urutan,
		data.Aktif,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func Update(id int, data Mapel) error {
	_, err := database.DB.Exec(
		`UPDATE tbl_mapel SET
			id_induk = ?,
			nama_mapel = ?,
			urutan = ?
		WHERE id_mapel = ?`,
		data.IDInduk,
		data.Nama,
		data.Urutan,
		id,
	)

	return err
}

func Delete(id int) error {
	_, err := database.DB.Exec(
		`DELETE FROM tbl_mapel WHERE id_mapel = ?`,
		id,
	)

	return err
}

func ToggleStatus(id int) error {
	_, err := database.DB.Exec(
		`UPDATE tbl_mapel SET aktif = 1 - aktif WHERE id_mapel = ?`,
		id,
	)

	return err
}
