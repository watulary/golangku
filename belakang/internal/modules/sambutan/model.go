package sambutan

import "time"

type Sambutan struct {
	IDSambutan   int    `db:"id_sambutan" json:"id_sambutan"`
	IDPengguna   string `db:"id_pengguna" json:"id_pengguna"`
	NamaPengguna string `db:"nama_pengguna" json:"nama_pengguna"`

	JudulSambutan string `db:"judul_sambutan" json:"judul_sambutan"`
	IsiSambutan   string `db:"isi_sambutan" json:"isi_sambutan"`

	FotoKepsek string `db:"foto_kepsek" json:"foto_kepsek"`
	FotoURL    string `json:"foto_url"`

	Hits int `db:"hits" json:"hits"`

	Waktu time.Time `db:"waktu" json:"waktu"`
}
