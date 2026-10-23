package walikelas

type WaliKelas struct {
	ID          int    `db:"id_walikelas" json:"id_walikelas"`
	IDKelas     int    `db:"id_kelas" json:"id_kelas"`
	IDRombel    int    `db:"id_rombel" json:"id_rombel"`
	NUPTK       string `db:"nuptk" json:"nuptk"`
	NamaKelas   string `db:"nama_kelas" json:"nama_kelas"`
	NamaRombel  string `db:"nama_rombel" json:"nama_rombel"`
	NamaPegawai string `db:"nama_pegawai" json:"nama_pegawai"`
}

type PegawaiOption struct {
	NUPTK       string `db:"nuptk" json:"nuptk"`
	NamaPegawai string `db:"nama_pegawai" json:"nama_pegawai"`
	Jabatan     string `db:"jabatan" json:"jabatan"`
}
