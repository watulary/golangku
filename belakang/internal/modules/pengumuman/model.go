package pengumuman

import "time"

type Pengumuman struct {
	IDPengumuman     string    `db:"id_pengumuman" json:"id_pengumuman" form:"id_pengumuman"`
	JudulPengumuman  string    `db:"judul_pengumuman" json:"judul_pengumuman" form:"judul_pengumuman"`
	SlugPengumuman   string    `db:"slug_pengumuman" json:"slug_pengumuman" form:"slug_pengumuman"`
	IsiPengumuman    string    `db:"isi_pengumuman" json:"isi_pengumuman" form:"isi_pengumuman"`
	Ringkasan        string    `db:"ringkasan" json:"ringkasan" form:"ringkasan"`
	Foto             string    `db:"foto" json:"foto" form:"foto"`
	FotoURL          string    `json:"foto_url"`
	Hits             int       `db:"hits" json:"hits"`
	TanggalPublish   time.Time `db:"tanggal_publish" json:"tanggal_publish"`
	TanggalUbah      time.Time `db:"tanggal" json:"tanggal_ubah"`
	NamaPengguna     string    `db:"nama_pengguna" json:"nama_pengguna"`
	StatusPengumuman string    `db:"status_pengumuman" json:"status_pengumuman" form:"status_pengumuman"`
	IDPengguna       string    `db:"id_pengguna" json:"id_pengguna"`
}
