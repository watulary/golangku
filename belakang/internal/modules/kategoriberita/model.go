package kategoriberita

import "time"

// KategoriBerita merepresentasikan tabel kategori_berita.
// 3 field pertama (IDKategoriBerita, NamaKategoriBerita, SlugKategoriBerita)
// tetap dipakai apa adanya oleh query publik yang sudah berjalan.
type KategoriBerita struct {
	IDKategoriBerita   string `db:"id_kategoriberita" json:"id_kategoriberita"`
	NamaKategoriBerita string `db:"nama_kategoriberita" json:"nama_kategoriberita" form:"nama_kategoriberita"`
	SlugKategoriBerita string `db:"slug_kategoriberita" json:"slug_kategoriberita"`

	// Field tambahan untuk kebutuhan admin
	IDPengguna string    `db:"id_pengguna" json:"id_pengguna"`
	Urutan     int       `db:"urutan" json:"urutan" form:"urutan"`
	Hits       int       `db:"hits" json:"hits"`
	Tanggal    time.Time `db:"tanggal" json:"tanggal"`

	// Field turunan hasil JOIN, tanpa tag db
	NamaPengguna string `db:"nama_pengguna" json:"nama_pengguna"`
}

// KategoriBeritaRequest dipakai untuk binding body create/update dari form admin
type KategoriBeritaRequest struct {
	NamaKategoriBerita string `json:"nama_kategoriberita" form:"nama_kategoriberita" validate:"required,max=200"`
	Urutan             int    `json:"urutan" form:"urutan"`
}
