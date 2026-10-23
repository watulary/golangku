package visimisi

import "sekolah/internal/database"

func Get() (*VisiMisi, error) {

	var data VisiMisi

	err := database.DB.Get(
		&data,
		`SELECT
			v.id,
			v.id_pengguna,
			p.nama_pengguna,
			v.judul,
			v.keterangan,
			v.hits,
			v.waktu
		FROM visimisi v
		LEFT JOIN pengguna p
			ON p.id_pengguna = v.id_pengguna
		ORDER BY v.id DESC
		LIMIT 1`,
	)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func IncrementHits() error {

	_, err := database.DB.Exec(
		`UPDATE visimisi
		SET hits = hits + 1
		ORDER BY id DESC
		LIMIT 1`,
	)

	return err
}

func UpdateDashboard(data *VisiMisi) error {

	_, err := database.DB.Exec(
		`UPDATE visimisi
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
