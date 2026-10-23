package rombel

type Rombel struct {
	ID        int    `db:"id_rombel" json:"id_rombel"`
	IDKelas   int    `db:"id_kelas" json:"id_kelas"`
	Nama      string `db:"nama_rombel" json:"nama_rombel"`
	Aktif     int    `db:"aktif" json:"aktif"`
	NamaKelas string `db:"nama_kelas" json:"nama_kelas"`
}
