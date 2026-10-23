package kelulusan

import "time"

// KelulusanWaktu - baris tunggal (id_waktu selalu 1). pengaturan pengumuman kelulusan (dipakai halaman
// Peserta Didik + countdown) dan isi konten surat SKL (kop*/isi*/ttd, dipakai
// generator PDF). Diupdate lewat 2 endpoint terpisah (pengumuman vs skl)
// persis seperti simpan()/simpanskl() di controller CI lama.
type KelulusanWaktu struct {
	IDWaktu        int       `db:"id_waktu" json:"id_waktu"`
	IDPengguna     string    `db:"id_pengguna" json:"id_pengguna"`
	AturWaktu      string    `db:"aturwaktu" json:"aturwaktu"`
	AturForm       string    `db:"aturform" json:"aturform"`
	HeroArea1      string    `db:"heroarea1" json:"heroarea1"`
	HeroArea2      string    `db:"heroarea2" json:"heroarea2"`
	Waktu          time.Time `db:"waktu" json:"waktu"`
	Judul          string    `db:"judul" json:"judul"`
	Nomor          string    `db:"nomor" json:"nomor"`
	Pengantar      string    `db:"pengantar" json:"pengantar"`
	Penutup        string    `db:"penutup" json:"penutup"`
	Akhir          string    `db:"akhir" json:"akhir"`
	Tanggal        string    `db:"tanggal" json:"tanggal"`
	Mengetahui     string    `db:"mengetahui" json:"mengetahui"`
	Kepsek         string    `db:"kepsek" json:"kepsek"`
	NipKepsek      string    `db:"nipkepsek" json:"nipkepsek"`
	Kop1SKL        string    `db:"kop1_skl" json:"kop1_skl"`
	Kop2SKL        string    `db:"kop2_skl" json:"kop2_skl"`
	Kop3SKL        string    `db:"kop3_skl" json:"kop3_skl"`
	Kop4SKL        string    `db:"kop4_skl" json:"kop4_skl"`
	Kop5SKL        string    `db:"kop5_skl" json:"kop5_skl"`
	Kop6SKL        string    `db:"kop6_skl" json:"kop6_skl"`
	Kop7SKL        string    `db:"kop7_skl" json:"kop7_skl"`
	Kop8SKL        string    `db:"kop8_skl" json:"kop8_skl"`
	Isi1SKL        string    `db:"isi1_skl" json:"isi1_skl"`
	Isi2SKL        string    `db:"isi2_skl" json:"isi2_skl"`
	Isi3SKL        string    `db:"isi3_skl" json:"isi3_skl"`
	Isi4SKL        string    `db:"isi4_skl" json:"isi4_skl"`
	Isi5SKL        string    `db:"isi5_skl" json:"isi5_skl"`
	TglSKL         string    `db:"tgl_skl" json:"tgl_skl"`
	SekolahSKL     string    `db:"sekolah_skl" json:"sekolah_skl"`
	KepsekSKL      string    `db:"kepsek_skl" json:"kepsek_skl"`
	NipKepsekSKL   string    `db:"nipkepsek_skl" json:"nipkepsek_skl"`
	TtdSKL         string    `db:"ttd_skl" json:"ttd_skl"`
	TahunPelajaran string    `db:"tahun_pelajaran" json:"tahun_pelajaran"`
	NamaPengguna   string    `db:"nama_pengguna" json:"nama_pengguna"` // JOIN pengguna, utk "Diperbarui oleh"
}

// Kelulusan - hasil kelulusan per siswa (PK nis). TIDAK mengubah data
// aktual di pesertadidik, sama seperti kenaikan.
type Kelulusan struct {
	IDTahunPelajaran int       `db:"id_tahunpelajaran" json:"id_tahunpelajaran"`
	IDKelas          int       `db:"id_kelas" json:"id_kelas"`
	IDRombel         int       `db:"id_rombel" json:"id_rombel"`
	NIS              string    `db:"nis" json:"nis"`
	Nopes            string    `db:"nopes" json:"nopes"`
	PaS              string    `db:"pa_s" json:"pa_s"`
	PpknS            string    `db:"ppkn_s" json:"ppkn_s"`
	BindS            string    `db:"bind_s" json:"bind_s"`
	BingS            string    `db:"bing_s" json:"bing_s"`
	MtkS             string    `db:"mtk_s" json:"mtk_s"`
	IpaS             string    `db:"ipa_s" json:"ipa_s"`
	IpsS             string    `db:"ips_s" json:"ips_s"`
	PenjasS          string    `db:"penjas_s" json:"penjas_s"`
	SenbudS          string    `db:"senbud_s" json:"senbud_s"`
	PrakaryaS        string    `db:"prakarya_s" json:"prakarya_s"`
	Jumlah           string    `db:"jumlah" json:"jumlah"`
	RataS            string    `db:"rata_s" json:"rata_s"`
	Ket              string    `db:"ket" json:"ket"`
	FotoLulus        string    `db:"fotolulus" json:"fotolulus"`
	Waktu            time.Time `db:"waktu" json:"waktu"`
}

// KelulusanAdmin - hasil JOIN utk list admin (mirip tampillulus di CI lama).
type KelulusanAdmin struct {
	Kelulusan
	NamaSiswa          string `db:"nama_siswa" json:"nama_siswa"`
	NamaKelas          string `db:"nama_kelas" json:"nama_kelas"`
	NamaRombel         string `db:"nama_rombel" json:"nama_rombel"`
	NamaTahunPelajaran string `db:"nama_tahunpelajaran" json:"nama_tahunpelajaran"`
	NISN               string `db:"nisn" json:"nisn"`
	TempatLahirSiswa   string `db:"tempat_lahirsiswa" json:"tempat_lahirsiswa"`
	TanggalLahirSiswa  string `db:"tanggal_lahirsiswa" json:"tanggal_lahirsiswa"`
	FotoURL            string `db:"-" json:"foto_url"` // computed di handler dari fotolulus, lihat konvensi field foto modul lain
}

// CetakData - data gabungan (pesertadidik+kelas+rombel+tahunpelajaran+
// setting+kelulusan) dipakai generator PDF SKL, mirip cetakpdf() model lama.
type CetakData struct {
	NamaSiswa          string `db:"nama_siswa"`
	TempatLahirSiswa   string `db:"tempat_lahirsiswa"`
	TanggalLahirSiswa  string `db:"tanggal_lahirsiswa"`
	NIS                string `db:"nis"`
	NISN               string `db:"nisn"`
	NamaSekolah        string `db:"nama_sekolah"`
	NPSN               string `db:"npsn"`
	NamaKelas          string `db:"nama_kelas"`
	NamaRombel         string `db:"nama_rombel"`
	NamaTahunPelajaran string `db:"nama_tahunpelajaran"`
	Nopes              string `db:"nopes"`
	Ket                string `db:"ket"`
	PaS                string `db:"pa_s"`
	PpknS              string `db:"ppkn_s"`
	BindS              string `db:"bind_s"`
	BingS              string `db:"bing_s"`
	MtkS               string `db:"mtk_s"`
	IpaS               string `db:"ipa_s"`
	IpsS               string `db:"ips_s"`
	PenjasS            string `db:"penjas_s"`
	SenbudS            string `db:"senbud_s"`
	PrakaryaS          string `db:"prakarya_s"`
	Jumlah             string `db:"jumlah"`
	RataS              string `db:"rata_s"`
	FotoLulus          string `db:"fotolulus"`
}

// FotoKelulusan - galeri foto (diisi via upload ZIP massal). fotolulus di
// tabel kelulusan menyimpan NAMA FILE yang dicocokkan manual ke galeri ini
// (pola lama dipertahankan, beda dari upload foto per-siswa di pesertadidik).
type FotoKelulusan struct {
	ID     int    `db:"id" json:"id"`
	Gambar string `db:"gambar" json:"gambar"`
}
