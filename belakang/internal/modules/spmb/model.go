package spmb

type SPMB struct {
	NoPeserta   string `db:"no_peserta" json:"no_peserta"`
	NamaPeserta string `db:"nama_peserta" json:"nama_peserta"`
	Kelamin     string `db:"kelamin" json:"kelamin"`
	AsalSekolah string `db:"asal_sekolah" json:"asal_sekolah"`
	Jalur       int    `db:"jalur" json:"jalur"`
	LaporDiri   int    `db:"lapor_diri" json:"lapor_diri"`
}
