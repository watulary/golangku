package berita

import "time"

type Berita struct {
	IDBerita           string    `db:"id_berita" json:"id_berita" form:"id_berita"`
	JudulBerita        string    `db:"judul_berita" json:"judul_berita" form:"judul_berita"`
	SlugBerita         string    `db:"slug_berita" json:"slug_berita" form:"slug_berita"`
	IsiBerita          string    `db:"isi_berita" json:"isi_berita" form:"isi_berita"`
	Ringkasan          string    `db:"ringkasan" json:"ringkasan" form:"ringkasan"`
	Foto               string    `db:"foto" json:"foto" form:"foto"`
	FotoURL            string    `json:"foto_url"`
	Hits               int       `db:"hits" json:"hits"`
	TanggalPublish     time.Time `db:"tanggal_publish" json:"tanggal_publish"`
	TanggalUbah        time.Time `db:"tanggal" json:"tanggal_ubah"`
	NamaPengguna       string    `db:"nama_pengguna" json:"nama_pengguna"`
	IDKategoriBerita   string    `db:"id_kategoriberita" json:"id_kategoriberita" form:"id_kategoriberita"`
	NamaKategoriBerita string    `db:"nama_kategoriberita" json:"nama_kategoriberita"`
	SlugKategoriBerita string    `db:"slug_kategoriberita" json:"slug_kategoriberita"`
	StatusBerita       string    `db:"status_berita" json:"status_berita" form:"status_berita"`
	IDPengguna         string    `db:"id_pengguna" json:"id_pengguna"`
}
