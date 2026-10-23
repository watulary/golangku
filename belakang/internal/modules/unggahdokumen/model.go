package unggahdokumen

import "time"

type Unggahdokumen struct {
	ID           string    `db:"id" json:"id"`
	IDPengguna   string    `db:"id_pengguna" json:"id_pengguna"`
	NamaPengguna string    `db:"nama_pengguna" json:"nama_pengguna"`
	Nama         string    `db:"nama" json:"nama"`
	Dokumen      string    `db:"dokumen" json:"dokumen"`
	DokumenURL   string    `json:"dokumen_url"`
	Waktu        time.Time `db:"waktu" json:"waktu"`
}
