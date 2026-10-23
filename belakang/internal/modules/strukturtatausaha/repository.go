package strukturtatausaha

import "sekolah/internal/database"

func Get() (*StrukturTatausaha, error) {

	var data StrukturTatausaha

	err := database.DB.Get(
		&data,
		`SELECT
			s.id,
			s.id_pengguna,
			p.nama_pengguna,
			s.judul,
			s.gambar,
			s.hits,
			s.waktu
		FROM strukturtatausaha s
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
		`UPDATE strukturtatausaha
		SET hits = hits + 1
		ORDER BY id DESC
		LIMIT 1`,
	)

	return err
}

func UpdateDashboard(data *StrukturTatausaha) error {

	_, err := database.DB.Exec(
		`UPDATE strukturtatausaha
		SET
			id_pengguna = ?,
			judul       = ?,
			gambar      = ?,
			waktu       = NOW()
		WHERE id = ?`,
		data.IDPengguna,
		data.Judul,
		data.Gambar,
		data.ID,
	)

	return err
}
