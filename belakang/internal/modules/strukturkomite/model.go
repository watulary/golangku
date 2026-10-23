package strukturkomite

import "time"

type StrukturKomite struct {
	ID           int       `db:"id" json:"id"`
	IDPengguna   string    `db:"id_pengguna" json:"id_pengguna"`
	NamaPengguna string    `db:"nama_pengguna" json:"nama_pengguna"`
	Judul        string    `db:"judul" json:"judul"`
	Gambar       string    `db:"gambar" json:"gambar"`
	FotoURL      string    `json:"foto_url"`
	Hits         int       `db:"hits" json:"hits"`
	Waktu        time.Time `db:"waktu" json:"waktu"`
}
