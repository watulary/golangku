package kelas

type Kelas struct {
	ID    int    `db:"id_kelas" json:"id_kelas"`
	Nama  string `db:"nama_kelas" json:"nama_kelas"`
	Aktif int    `db:"aktif" json:"aktif"`
}
