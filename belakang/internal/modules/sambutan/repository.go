package sambutan

import "sekolah/internal/database"

func GetSambutan() (*Sambutan, error) {

	var data Sambutan

	err := database.DB.Get(
		&data,
		`SELECT
			s.id_sambutan,
			s.id_pengguna,
			p.nama_pengguna,
			s.judul_sambutan,
			s.isi_sambutan,
			s.foto_kepsek,
			s.hits,
			s.waktu
		FROM sambutan s
		LEFT JOIN pengguna p
			ON p.id_pengguna = s.id_pengguna
		LIMIT 1`,
	)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func UpdateHits() error {

	_, err := database.DB.Exec(`
		UPDATE sambutan
		SET hits = hits + 1
	`)

	return err
}

// =====================================
// Dashboard
// =====================================

func UpdateDashboard(data *Sambutan) error {

	_, err := database.DB.Exec(`
		UPDATE sambutan
		SET
			id_pengguna = ?,
			judul_sambutan = ?,
			isi_sambutan = ?,
			foto_kepsek = ?,
			waktu = NOW()
		WHERE id_sambutan = ?
	`,
		data.IDPengguna,
		data.JudulSambutan,
		data.IsiSambutan,
		data.FotoKepsek,
		data.IDSambutan,
	)

	return err
}
