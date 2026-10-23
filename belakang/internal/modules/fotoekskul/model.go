package fotoekskul

import "time"

type FotoEkskul struct {
	ID           string    `db:"id" json:"id" form:"id"`
	Nama         string    `db:"nama" json:"nama" form:"nama"`
	Kegiatan     string    `db:"kegiatan" json:"kegiatan" form:"kegiatan"`
	Foto         string    `db:"foto" json:"foto" form:"foto"`
	FotoURL      string    `json:"foto_url"`
	Waktu        time.Time `db:"waktu" json:"waktu"`
	IDPengguna   string    `db:"id_pengguna" json:"id_pengguna"`
	NamaPengguna string    `db:"nama_pengguna" json:"nama_pengguna"`
}
