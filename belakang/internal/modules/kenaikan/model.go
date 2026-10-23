package kenaikan

import "time"

// Kenaikan merepresentasikan tabel `kenaikan` (PK: nis)
// Tabel ini murni pencatatan/pengumuman kenaikan kelas — tidak mengubah
// data aktual di tabel pesertadidik.
type Kenaikan struct {
	IDTahunPelajaran int    `db:"id_tahunpelajaran" json:"id_tahunpelajaran"`
	IDKelas          int    `db:"id_kelas" json:"id_kelas"`
	IDRombel         int    `db:"id_rombel" json:"id_rombel"`
	NIS              string `db:"nis" json:"nis"`
	Ket              string `db:"ket" json:"ket"`
}

// KenaikanListItem dipakai untuk list admin (List.vue), gabungan JOIN ke
// tbl_tahunpelajaran, tbl_kelas, tbl_rombel, dan pesertadidik.
type KenaikanListItem struct {
	IDTahunPelajaran int    `db:"id_tahunpelajaran" json:"id_tahunpelajaran"`
	IDKelas          int    `db:"id_kelas" json:"id_kelas"`
	IDRombel         int    `db:"id_rombel" json:"id_rombel"`
	NIS              string `db:"nis" json:"nis"`
	Ket              string `db:"ket" json:"ket"`

	// Hasil JOIN nama referensi
	NamaTahunPelajaran string `db:"nama_tahunpelajaran" json:"nama_tahunpelajaran"`
	NamaKelas          string `db:"nama_kelas" json:"nama_kelas"`
	NamaRombel         string `db:"nama_rombel" json:"nama_rombel"`

	// Hasil JOIN ke pesertadidik (data siswa, bukan kolom tabel kenaikan)
	NamaSiswa         string `db:"nama_siswa" json:"nama_siswa"`
	NISN              string `db:"nisn" json:"nisn"`
	NoUjianSD         string `db:"no_ujiansd" json:"no_ujiansd"`
	TempatLahirSiswa  string `db:"tempat_lahirsiswa" json:"tempat_lahirsiswa"`
	TanggalLahirSiswa string `db:"tanggal_lahirsiswa" json:"tanggal_lahirsiswa"`
	JenisKelamin      string `db:"jenis_kelamin" json:"jenis_kelamin"`

	// Foto diambil dari pesertadidik via JOIN nis, bukan kolom kenaikan
	Foto    string `db:"foto" json:"-"`
	FotoURL string `db:"-" json:"foto_url,omitempty"`
}

// KenaikanWaktu merepresentasikan tabel `kenaikan_waktu` — pengaturan
// pengumuman kenaikan kelas, single-row (id_waktu selalu = 1).
// Kolom kop*_skl, isi*_skl, tgl_skl, sekolah_skl, kepsek_skl, nipkepsek_skl
// sengaja TIDAK dimasukkan — sisa modul lain (SKL), tidak dipakai kenaikan.
type KenaikanWaktu struct {
	IDWaktu        int       `db:"id_waktu" json:"id_waktu"`
	IDPengguna     string    `db:"id_pengguna" json:"id_pengguna"`
	Aturwaktu      string    `db:"aturwaktu" json:"aturwaktu"`
	Waktu          time.Time `db:"waktu" json:"waktu"`
	Judul          string    `db:"judul" json:"judul"`
	Nomor          string    `db:"nomor" json:"nomor"`
	Pengantar      string    `db:"pengantar" json:"pengantar"`
	Penutup        string    `db:"penutup" json:"penutup"`
	Akhir          string    `db:"akhir" json:"akhir"`
	Tanggal        string    `db:"tanggal" json:"tanggal"`
	Mengetahui     string    `db:"mengetahui" json:"mengetahui"`
	Kepsek         string    `db:"kepsek" json:"kepsek"`
	NIPKepsek      string    `db:"nipkepsek" json:"nipkepsek"`
	TTDSkl         string    `db:"ttd_skl" json:"ttd_skl"`
	TahunPelajaran string    `db:"tahun_pelajaran" json:"tahun_pelajaran"`

	// Hasil JOIN ke pengguna, untuk footer "diubah oleh siapa & kapan"
	NamaPengguna string `db:"nama_pengguna" json:"nama_pengguna,omitempty"`
}