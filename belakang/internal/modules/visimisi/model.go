package visimisi

import "time"

type VisiMisi struct {
	ID           int       `db:"id" json:"id"`
	IDPengguna   string    `db:"id_pengguna" json:"id_pengguna"`
	NamaPengguna string    `db:"nama_pengguna" json:"nama_pengguna"`
	Judul        string    `db:"judul" json:"judul"`
	Keterangan   string    `db:"keterangan" json:"keterangan"`
	Hits         int       `db:"hits" json:"hits"`
	Waktu        time.Time `db:"waktu" json:"waktu"`
}
