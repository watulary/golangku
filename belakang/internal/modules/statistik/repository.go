package statistik

import "sekolah/internal/database"

// GetStatistik mengambil statistik ringkas dashboard. Empat subquery konten
// website (berita/pengumuman/galeri/video) difilter per id_pengguna pembuat
// kecuali isSuperAdmin true (Level == "Admin"), yang tampil semua data tanpa
// memandang siapa pembuatnya. Statistik akademik (murid/rombel/guru/tendik)
// TIDAK difilter, selalu data keseluruhan sekolah.
func GetStatistik(idPengguna string, isSuperAdmin bool) (*Statistik, error) {

	var data Statistik

	contentFilter := ""
	if !isSuperAdmin {
		contentFilter = "WHERE id_pengguna = ?"
	}

	query := `
SELECT
	(
		SELECT COUNT(*)
		FROM pesertadidik
		WHERE (
			status_keluar IS NULL
			OR status_keluar = ''
		)
	) AS jumlah_murid,

	(
		SELECT COUNT(*)
		FROM tbl_rombel
		WHERE aktif = 1
	) AS jumlah_rombel,

	(
		SELECT COUNT(*)
		FROM pegawai
		WHERE jabatan = 'Guru'
		AND (
			status_keluar IS NULL
			OR status_keluar = ''
		)
	) AS jumlah_guru,

	(
		SELECT COUNT(*)
		FROM pegawai
		WHERE jabatan = 'Tatausaha'
		AND (
			status_keluar IS NULL
			OR status_keluar = ''
		)
	) AS jumlah_tendik,

	(SELECT COUNT(*) FROM berita ` + contentFilter + `) AS jumlah_berita,

	(SELECT COUNT(*) FROM pengumuman ` + contentFilter + `) AS jumlah_pengumuman,

	(SELECT COUNT(*) FROM fotoekskul ` + contentFilter + `) AS jumlah_galerifoto,

	(SELECT COUNT(*) FROM video ` + contentFilter + `) AS jumlah_video
`

	var err error
	if isSuperAdmin {
		err = database.DB.Get(&data, query)
	} else {
		err = database.DB.Get(&data, query, idPengguna, idPengguna, idPengguna, idPengguna)
	}

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func GetRekapRombel() ([]RekapRombel, error) {
	var data []RekapRombel

	query := `
		SELECT
			k.nama_kelas,
			r.nama_rombel,
			COUNT(*) AS total,
			SUM(CASE WHEN p.jenis_kelamin = 'L' THEN 1 ELSE 0 END) AS laki,
			SUM(CASE WHEN p.jenis_kelamin = 'P' THEN 1 ELSE 0 END) AS perempuan,
			SUM(CASE WHEN p.agama_siswa = 'Islam' THEN 1 ELSE 0 END) AS islam,
			SUM(CASE WHEN p.agama_siswa = 'Kristen' THEN 1 ELSE 0 END) AS kristen,
			SUM(CASE WHEN p.agama_siswa = 'Katholik' THEN 1 ELSE 0 END) AS katholik,
			SUM(CASE WHEN p.agama_siswa = 'Hindu' THEN 1 ELSE 0 END) AS hindu,
			SUM(CASE WHEN p.agama_siswa = 'Budha' THEN 1 ELSE 0 END) AS budha
		FROM pesertadidik p
		JOIN tbl_rombel r ON p.id_rombel = r.id_rombel
		JOIN tbl_kelas k ON r.id_kelas = k.id_kelas
		WHERE p.id_tahunpelajaran = (
			SELECT id_tahunpelajaran FROM tbl_tahunpelajaran WHERE aktif = 1 LIMIT 1
		)
		AND (p.status_keluar IS NULL OR p.status_keluar = '')
		GROUP BY k.id_kelas, k.nama_kelas, r.id_rombel, r.nama_rombel
		ORDER BY k.nama_kelas ASC, r.nama_rombel ASC
	`

	err := database.DB.Select(&data, query)
	if err != nil {
		return nil, err
	}
	return data, nil
}
