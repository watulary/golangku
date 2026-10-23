package strukturakademik

import "sekolah/internal/database"

func Get() (*StrukturAkademik, error) {

	var data StrukturAkademik

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
		FROM strukturakademik s
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
		`UPDATE strukturakademik
		SET hits = hits + 1
		ORDER BY id DESC
		LIMIT 1`,
	)

	return err
}

func UpdateDashboard(data *StrukturAkademik) error {

	_, err := database.DB.Exec(
		`UPDATE strukturakademik
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
