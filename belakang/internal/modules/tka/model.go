package tka

type Tka struct {
	IDTahunPelajaran int    `db:"id_tahunpelajaran" json:"id_tahunpelajaran"`
	IDKelas          int    `db:"id_kelas" json:"id_kelas"`
	IDRombel         int    `db:"id_rombel" json:"id_rombel"`
	NIS              string `db:"nis" json:"nis"`
	Nopes            string `db:"nopes" json:"nopes"`
	MTK              string `db:"mtk" json:"mtk"`
	KetMTK           string `db:"ket_mtk" json:"ket_mtk"`
	Bind             string `db:"bind" json:"bind"`
	KetBind          string `db:"ket_bind" json:"ket_bind"`
	Waktu            string `db:"waktu" json:"waktu"`
}

// TkaAdmin = baris list admin, join referensi + nama siswa
type TkaAdmin struct {
	Tka
	NamaTahunPelajaran string `db:"nama_tahunpelajaran" json:"nama_tahunpelajaran"`
	NamaKelas          string `db:"nama_kelas" json:"nama_kelas"`
	NamaRombel         string `db:"nama_rombel" json:"nama_rombel"`
	NamaSiswa          string `db:"nama_siswa" json:"nama_siswa"`
}

// TkaExport = baris export Excel, tambahan kolom pribadi siswa
type TkaExport struct {
	TkaAdmin
	NISN              string `db:"nisn" json:"nisn"`
	JenisKelamin      string `db:"jenis_kelamin" json:"jenis_kelamin"`
	TempatLahirSiswa  string `db:"tempat_lahirsiswa" json:"tempat_lahirsiswa"`
	TanggalLahirSiswa string `db:"tanggal_lahirsiswa" json:"tanggal_lahirsiswa"`
}

type TkaWaktu struct {
	IDWaktu        int    `db:"id_waktu" json:"id_waktu"`
	IDPengguna     string `db:"id_pengguna" json:"id_pengguna"`
	NamaPengguna   string `db:"nama_pengguna" json:"nama_pengguna"`
	AturWaktu      string `db:"aturwaktu" json:"aturwaktu"`
	AturForm       string `db:"aturform" json:"aturform"`
	Waktu          string `db:"waktu" json:"waktu"`
	Judul          string `db:"judul" json:"judul"`
	Nomor          string `db:"nomor" json:"nomor"`
	Pengantar      string `db:"pengantar" json:"pengantar"`
	Penutup        string `db:"penutup" json:"penutup"`
	Akhir          string `db:"akhir" json:"akhir"`
	Tanggal        string `db:"tanggal" json:"tanggal"`
	Mengetahui     string `db:"mengetahui" json:"mengetahui"`
	Kepsek         string `db:"kepsek" json:"kepsek"`
	NIPKepsek      string `db:"nipkepsek" json:"nipkepsek"`
	TtdSKL         string `db:"ttd_skl" json:"ttd_skl"`
	TahunPelajaran string `db:"tahun_pelajaran" json:"tahun_pelajaran"`
}
