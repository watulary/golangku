package sejarah

import "sekolah/internal/database"

func Get() (*Sejarah, error) {

	var data Sejarah

	err := database.DB.Get(
		&data,
		`SELECT
			s.id,
			s.id_pengguna,
			p.nama_pengguna,
			s.judul,
			s.keterangan,
			s.hits,
			s.waktu
		FROM sejarah s
		LEFT JOIN pengguna p
			ON p.id_pengguna = s.id_pengguna
		ORDER BY s.id DESC
		LIMIT 1`,
	)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func IncrementHits() error {

	_, err := database.DB.Exec(
		`UPDATE sejarah
		SET hits = hits + 1
		ORDER BY id DESC
		LIMIT 1`,
	)

	return err
}

type UpdateInput struct {
	IDPengguna string `json:"id_pengguna" db:"id_pengguna"`
	Judul      string `json:"judul" db:"judul"`
	Keterangan string `json:"keterangan" db:"keterangan"`
}

func UpdateDashboard(data *Sejarah) error {

	_, err := database.DB.Exec(
		`UPDATE sejarah
		SET
			id_pengguna = ?,
			judul       = ?,
			keterangan  = ?,
			waktu       = NOW()
		WHERE id = ?`,
		data.IDPengguna,
		data.Judul,
		data.Keterangan,
		data.ID,
	)

	return err
}
