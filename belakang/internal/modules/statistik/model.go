package statistik

type Statistik struct {
	JumlahMurid      int `db:"jumlah_murid" json:"jumlah_murid"`
	JumlahRombel     int `db:"jumlah_rombel" json:"jumlah_rombel"`
	JumlahGuru       int `db:"jumlah_guru" json:"jumlah_guru"`
	JumlahTendik     int `db:"jumlah_tendik" json:"jumlah_tendik"`
	JumlahBerita     int `db:"jumlah_berita" json:"jumlah_berita"`
	JumlahPengumuman int `db:"jumlah_pengumuman" json:"jumlah_pengumuman"`
	JumlahGaleriFoto int `db:"jumlah_galerifoto" json:"jumlah_galerifoto"`
	JumlahVideo      int `db:"jumlah_video" json:"jumlah_video"`
}

type RekapRombel struct {
	NamaKelas  string `db:"nama_kelas" json:"nama_kelas"`
	NamaRombel string `db:"nama_rombel" json:"nama_rombel"`
	Total      int    `db:"total" json:"total"`
	Laki       int    `db:"laki" json:"laki"`
	Perempuan  int    `db:"perempuan" json:"perempuan"`
	Islam      int    `db:"islam" json:"islam"`
	Kristen    int    `db:"kristen" json:"kristen"`
	Katholik   int    `db:"katholik" json:"katholik"`
	Hindu      int    `db:"hindu" json:"hindu"`
	Budha      int    `db:"budha" json:"budha"`
}
