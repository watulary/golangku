package strukturkomite

import "sekolah/internal/database"

func Get() (*StrukturKomite, error) {

	var data StrukturKomite

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
		FROM strukturkomite s
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
		`UPDATE strukturkomite
		SET hits = hits + 1
		ORDER BY id DESC
		LIMIT 1`,
	)

	return err
}

func UpdateDashboard(data *StrukturKomite) error {

	_, err := database.DB.Exec(
		`UPDATE strukturkomite
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
