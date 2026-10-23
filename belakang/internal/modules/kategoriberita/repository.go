package kategoriberita

import (
	"database/sql"

	"sekolah/internal/database"
)

// GetAll — publik, dipakai halaman depan
func GetAll() ([]KategoriBerita, error) {

	var data []KategoriBerita

	err := database.DB.Select(
		&data,
		`SELECT
			id_kategoriberita,
			nama_kategoriberita,
			slug_kategoriberita
		FROM kategori_berita
		ORDER BY nama_kategoriberita ASC`,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetAllAdmin — untuk halaman admin, dengan pencarian & pagination, plus nama_pengguna hasil JOIN
func GetAllAdmin(limit, offset int, search string) ([]KategoriBerita, error) {

	var data []KategoriBerita

	err := database.DB.Select(
		&data,
		`SELECT
			kb.id_kategoriberita,
			kb.id_pengguna,
			kb.nama_kategoriberita,
			kb.slug_kategoriberita,
			kb.urutan,
			kb.hits,
			kb.tanggal,
			p.nama_pengguna
		FROM kategori_berita kb
		JOIN pengguna p ON p.id_pengguna = kb.id_pengguna
		WHERE kb.nama_kategoriberita LIKE ?
		ORDER BY kb.urutan ASC, kb.tanggal DESC
		LIMIT ? OFFSET ?`,
		"%"+search+"%", limit, offset,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

// CountAdmin — total baris untuk pagination
func CountAdmin(search string) (int, error) {

	var total int

	err := database.DB.Get(
		&total,
		`SELECT COUNT(*) FROM kategori_berita WHERE nama_kategoriberita LIKE ?`,
		"%"+search+"%",
	)

	if err != nil {
		return 0, err
	}

	return total, nil
}

// GetByID — satu kategori berita untuk admin (dengan nama_pengguna)
func GetByID(id string) (*KategoriBerita, error) {

	var data KategoriBerita

	err := database.DB.Get(
		&data,
		`SELECT
			kb.id_kategoriberita,
			kb.id_pengguna,
			kb.nama_kategoriberita,
			kb.slug_kategoriberita,
			kb.urutan,
			kb.hits,
			kb.tanggal,
			p.nama_pengguna
		FROM kategori_berita kb
		JOIN pengguna p ON p.id_pengguna = kb.id_pengguna
		WHERE kb.id_kategoriberita = ?`,
		id,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &data, nil
}

// IsSlugExist — cek duplikasi slug, excludeID dikosongkan saat create
func IsSlugExist(slugValue, excludeID string) (bool, error) {

	var count int

	err := database.DB.Get(
		&count,
		`SELECT COUNT(*) FROM kategori_berita WHERE slug_kategoriberita = ? AND id_kategoriberita != ?`,
		slugValue, excludeID,
	)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Create — tambah kategori berita baru
func Create(k *KategoriBerita) error {

	_, err := database.DB.Exec(
		`INSERT INTO kategori_berita
			(id_kategoriberita, id_pengguna, nama_kategoriberita, slug_kategoriberita, urutan, hits)
		VALUES (?, ?, ?, ?, ?, 0)`,
		k.IDKategoriBerita, k.IDPengguna, k.NamaKategoriBerita, k.SlugKategoriBerita, k.Urutan,
	)

	return err
}

// Update — ubah kategori berita
func Update(k *KategoriBerita) error {

	_, err := database.DB.Exec(
		`UPDATE kategori_berita
		SET nama_kategoriberita = ?, slug_kategoriberita = ?, urutan = ?
		WHERE id_kategoriberita = ?`,
		k.NamaKategoriBerita, k.SlugKategoriBerita, k.Urutan, k.IDKategoriBerita,
	)

	return err
}

// IncrementHits — naikkan hits saat kategori diakses di halaman publik
func IncrementHits(id string) error {

	_, err := database.DB.Exec(
		`UPDATE kategori_berita SET hits = hits + 1 WHERE id_kategoriberita = ?`,
		id,
	)

	return err
}

// Delete — hapus kategori berita
func Delete(id string) error {

	_, err := database.DB.Exec(
		`DELETE FROM kategori_berita WHERE id_kategoriberita = ?`,
		id,
	)

	return err
}

// IsUsedByBerita — cek apakah kategori masih dipakai tabel berita
func IsUsedByBerita(id string) (bool, error) {

	var count int

	err := database.DB.Get(
		&count,
		`SELECT COUNT(*) FROM berita WHERE id_kategoriberita = ?`,
		id,
	)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
