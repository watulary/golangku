package fotobanner

import "time"

type FotoBanner struct {
	ID           string    `db:"id" json:"id" form:"id"`
	Nama         string    `db:"nama" json:"nama" form:"nama"`
	Kegiatan     string    `db:"kegiatan" json:"kegiatan" form:"kegiatan"`
	Gambar       string    `db:"gambar" json:"gambar" form:"gambar"`
	GambarURL    string    `json:"gambar_url"`
	Waktu        time.Time `db:"waktu" json:"waktu"`
	IDPengguna   string    `db:"id_pengguna" json:"id_pengguna"`
	NamaPengguna string    `db:"nama_pengguna" json:"nama_pengguna"`
}
