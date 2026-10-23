package suratmasuk

import "time"

type SuratMasuk struct {
	ID            string    `db:"id" json:"id" form:"id"`
	IDPengguna    string    `db:"id_pengguna" json:"id_pengguna"`
	NamaPengguna  string    `db:"nama_pengguna" json:"nama_pengguna"`
	Nama          string    `db:"nama" json:"nama" form:"nama"`
	Nomor         string    `db:"nomor" json:"nomor" form:"nomor"`
	Tanggal       string    `db:"tanggal" json:"tanggal" form:"tanggal"`
	TanggalTerima string    `db:"tanggal_terima" json:"tanggal_terima" form:"tanggal_terima"`
	Perihal       string    `db:"perihal" json:"perihal" form:"perihal"`
	Asal          string    `db:"asal" json:"asal" form:"asal"`
	Tujuan        string    `db:"tujuan" json:"tujuan" form:"tujuan"`
	Dokumen       string    `db:"dokumen" json:"dokumen" form:"dokumen"`
	DokumenURL    string    `json:"dokumen_url"`
	Waktu         time.Time `db:"waktu" json:"waktu"`
}
