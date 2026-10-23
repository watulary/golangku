package pesertadidik

// PesertaDidik merepresentasikan data lengkap satu peserta didik (118 kolom asli + field JOIN)
type PesertaDidik struct {
	// Relasi
	IDTahunPelajaran int `db:"id_tahunpelajaran" json:"id_tahunpelajaran"`
	IDKelas          int `db:"id_kelas" json:"id_kelas"`
	IDRombel         int `db:"id_rombel" json:"id_rombel"`

	// Data Siswa
	NamaSiswa         string `db:"nama_siswa" json:"nama_siswa"`
	NIS               string `db:"nis" json:"nis"`
	NISN              string `db:"nisn" json:"nisn"`
	NoUjianSD         string `db:"no_ujiansd" json:"no_ujiansd"`
	NoSimdik          string `db:"no_simdik" json:"no_simdik"`
	JenisKelamin      string `db:"jenis_kelamin" json:"jenis_kelamin"`
	TempatLahirSiswa  string `db:"tempat_lahirsiswa" json:"tempat_lahirsiswa"`
	TanggalLahirSiswa string `db:"tanggal_lahirsiswa" json:"tanggal_lahirsiswa"`
	NoAkte            string `db:"no_akte" json:"no_akte"`
	AkteKota          string `db:"akte_kota" json:"akte_kota"`
	AgamaSiswa        string `db:"agama_siswa" json:"agama_siswa"`
	WargaSiswa        string `db:"warga_siswa" json:"warga_siswa"`
	AnakKe            string `db:"anak_ke" json:"anak_ke"`
	SaudaraKandung    string `db:"saudara_kandung" json:"saudara_kandung"`
	SaudaraTiri       string `db:"saudara_tiri" json:"saudara_tiri"`
	SaudaraAngkat     string `db:"saudara_angkat" json:"saudara_angkat"`
	YatimPiatu        string `db:"yatim_piatu" json:"yatim_piatu"`
	BahasaSiswa       string `db:"bahasa_siswa" json:"bahasa_siswa"`
	Bantuan           string `db:"bantuan" json:"bantuan"`
	NoKK              string `db:"no_kk" json:"no_kk"`
	NIKSiswa          string `db:"nik_siswa" json:"nik_siswa"`
	AlamatSiswa       string `db:"alamat_siswa" json:"alamat_siswa"`
	RTSiswa           string `db:"rt_siswa" json:"rt_siswa"`
	RWSiswa           string `db:"rw_siswa" json:"rw_siswa"`
	KoposSiswa        string `db:"kopos_siswa" json:"kopos_siswa"`
	KelSiswa          string `db:"kel_siswa" json:"kel_siswa"`
	KecSiswa          string `db:"kec_siswa" json:"kec_siswa"`
	TelpRmhSiswa      string `db:"telprmh_siswa" json:"telprmh_siswa"`
	HPSiswa           string `db:"hp_siswa" json:"hp_siswa"`
	TinggalSiswa      string `db:"tinggal_siswa" json:"tinggal_siswa"`
	JarakSekolah      string `db:"jarak_sekolah" json:"jarak_sekolah"`
	GolDarah          string `db:"gol_darah" json:"gol_darah"`
	PernahSakit       string `db:"pernah_sakit" json:"pernah_sakit"`
	StatusJasmani     string `db:"status_jasmani" json:"status_jasmani"`
	KebutuhanKhusus   string `db:"kebutuhan_khusus" json:"kebutuhan_khusus"`
	TinggiBadan       string `db:"tinggi_badan" json:"tinggi_badan"`
	BeratBadan        string `db:"berat_badan" json:"berat_badan"`
	LingkarKepala     string `db:"lingkar_kepala" json:"lingkar_kepala"`
	ModaTransportasi  string `db:"moda_transportasi" json:"moda_transportasi"`
	LulusSD           string `db:"lulus_sd" json:"lulus_sd"`
	TanggalIjazah     string `db:"tanggal_ijazah" json:"tanggal_ijazah"`
	KotaIjazah        string `db:"kota_ijazah" json:"kota_ijazah"`
	NomorIjazah       string `db:"nomor_ijazah" json:"nomor_ijazah"`
	TahunIjazah       string `db:"tahun_ijazah" json:"tahun_ijazah"`
	NomorSKHUN        string `db:"nomor_skhun" json:"nomor_skhun"`
	TahunSKHUN        string `db:"tahun_skhun" json:"tahun_skhun"`
	NilaiIjazah       string `db:"nilai_ijazah" json:"nilai_ijazah"`
	NilaiNEM          string `db:"nilai_nem" json:"nilai_nem"`
	LamaBelajar       string `db:"lama_belajar" json:"lama_belajar"`
	Jalur             string `db:"jalur" json:"jalur"`
	Semester          string `db:"semester" json:"semester"`
	AsalSMP           string `db:"asal_smp" json:"asal_smp"`
	TanggalMasuk      string `db:"tanggal_masuk" json:"tanggal_masuk"`
	Email             string `db:"email" json:"email"`

	// Data Ayah
	NamaAyah         string `db:"nama_ayah" json:"nama_ayah"`
	TempatLahirAyah  string `db:"tempat_lahirayah" json:"tempat_lahirayah"`
	TanggalLahirAyah string `db:"tanggal_lahirayah" json:"tanggal_lahirayah"`
	AgamaAyah        string `db:"agama_ayah" json:"agama_ayah"`
	WargaAyah        string `db:"warga_ayah" json:"warga_ayah"`
	PendAyah         string `db:"pend_ayah" json:"pend_ayah"`
	PekAyah          string `db:"pek_ayah" json:"pek_ayah"`
	GajiAyah         string `db:"gaji_ayah" json:"gaji_ayah"`
	GajiAyahKtg      string `db:"gaji_ayahktg" json:"gaji_ayahktg"`
	NIKAyah          string `db:"nik_ayah" json:"nik_ayah"`
	AlamatAyah       string `db:"alamat_ayah" json:"alamat_ayah"`
	RTAyah           string `db:"rt_ayah" json:"rt_ayah"`
	RWAyah           string `db:"rw_ayah" json:"rw_ayah"`
	KoposAyah        string `db:"kopos_ayah" json:"kopos_ayah"`
	KelAyah          string `db:"kel_ayah" json:"kel_ayah"`
	KecAyah          string `db:"kec_ayah" json:"kec_ayah"`
	TelpAyah         string `db:"telp_ayah" json:"telp_ayah"`
	StatusAyah       string `db:"status_ayah" json:"status_ayah"`

	// Data Ibu
	NamaIbu         string `db:"nama_ibu" json:"nama_ibu"`
	TempatLahirIbu  string `db:"tempat_lahiribu" json:"tempat_lahiribu"`
	TanggalLahirIbu string `db:"tanggal_lahiribu" json:"tanggal_lahiribu"`
	AgamaIbu        string `db:"agama_ibu" json:"agama_ibu"`
	WargaIbu        string `db:"warga_ibu" json:"warga_ibu"`
	PendIbu         string `db:"pend_ibu" json:"pend_ibu"`
	PekIbu          string `db:"pek_ibu" json:"pek_ibu"`
	GajiIbu         string `db:"gaji_ibu" json:"gaji_ibu"`
	GajiIbuKtg      string `db:"gaji_ibuktg" json:"gaji_ibuktg"`
	NIKIbu          string `db:"nik_ibu" json:"nik_ibu"`
	AlamatIbu       string `db:"alamat_ibu" json:"alamat_ibu"`
	RTIbu           string `db:"rt_ibu" json:"rt_ibu"`
	RWIbu           string `db:"rw_ibu" json:"rw_ibu"`
	KoposIbu        string `db:"kopos_ibu" json:"kopos_ibu"`
	KelIbu          string `db:"kel_ibu" json:"kel_ibu"`
	KecIbu          string `db:"kec_ibu" json:"kec_ibu"`
	TelpIbu         string `db:"telp_ibu" json:"telp_ibu"`
	StatusIbu       string `db:"status_ibu" json:"status_ibu"`

	// Data Wali
	NamaWali         string `db:"nama_wali" json:"nama_wali"`
	TempatLahirWali  string `db:"tempat_lahirwali" json:"tempat_lahirwali"`
	TanggalLahirWali string `db:"tanggal_lahirwali" json:"tanggal_lahirwali"`
	AgamaWali        string `db:"agama_wali" json:"agama_wali"`
	WargaWali        string `db:"warga_wali" json:"warga_wali"`
	PendWali         string `db:"pend_wali" json:"pend_wali"`
	PekWali          string `db:"pek_wali" json:"pek_wali"`
	GajiWali         string `db:"gaji_wali" json:"gaji_wali"`
	GajiWaliKtg      string `db:"gaji_waliktg" json:"gaji_waliktg"`
	AlamatWali       string `db:"alamat_wali" json:"alamat_wali"`
	RTWali           string `db:"rt_wali" json:"rt_wali"`
	RWWali           string `db:"rw_wali" json:"rw_wali"`
	KoposWali        string `db:"kopos_wali" json:"kopos_wali"`
	KelWali          string `db:"kel_wali" json:"kel_wali"`
	KecWali          string `db:"kec_wali" json:"kec_wali"`
	TelpWali         string `db:"telp_wali" json:"telp_wali"`

	// Hobi
	HobiSeni       string `db:"hobi_seni" json:"hobi_seni"`
	HobiOrkes      string `db:"hobi_orkes" json:"hobi_orkes"`
	HobiOrganisasi string `db:"hobi_organisasi" json:"hobi_organisasi"`
	HobiLain       string `db:"hobi_lain" json:"hobi_lain"`

	// Status Keluar
	StatusKeluar  string `db:"status_keluar" json:"status_keluar"`
	AlasanKeluar  string `db:"alasan_keluar" json:"alasan_keluar"`
	TanggalKeluar string `db:"tanggal_keluar" json:"tanggal_keluar"`
	NomorSurat    string `db:"nomor_surat" json:"nomor_surat"`
	Tujuan        string `db:"tujuan" json:"tujuan"`

	Foto string `db:"foto" json:"foto"`

	// Computed non-DB, hasil JOIN ke tbl_kelas/tbl_rombel/tbl_tahunpelajaran
	NamaKelas          string `db:"nama_kelas" json:"nama_kelas"`
	NamaRombel         string `db:"nama_rombel" json:"nama_rombel"`
	NamaTahunPelajaran string `db:"nama_tahunpelajaran" json:"nama_tahunpelajaran"`

	// Computed non-DB, di-generate manual di handler
	FotoURL string `db:"-" json:"foto_url,omitempty"`
}

// PesertaDidikListItem versi ringkas dipakai GetAllAdmin/list (kolom yang tampil di tabel)
type PesertaDidikListItem struct {
	ID                 int    `db:"id" json:"id,omitempty"` // hanya terisi untuk arsip (kode_siswa); kosong untuk pesertadidik & mutasi
	NIS                string `db:"nis" json:"nis"`
	NamaSiswa          string `db:"nama_siswa" json:"nama_siswa"`
	NISN               string `db:"nisn" json:"nisn"`
	JenisKelamin       string `db:"jenis_kelamin" json:"jenis_kelamin"`
	AgamaSiswa         string `db:"agama_siswa" json:"agama_siswa"`
	TempatLahirSiswa   string `db:"tempat_lahirsiswa" json:"tempat_lahirsiswa"`
	TanggalLahirSiswa  string `db:"tanggal_lahirsiswa" json:"tanggal_lahirsiswa"`
	IDTahunPelajaran   int    `db:"id_tahunpelajaran" json:"id_tahunpelajaran"`
	IDKelas            int    `db:"id_kelas" json:"id_kelas"`
	IDRombel           int    `db:"id_rombel" json:"id_rombel"`
	NamaTahunPelajaran string `db:"nama_tahunpelajaran" json:"nama_tahunpelajaran"`
	NamaKelas          string `db:"nama_kelas" json:"nama_kelas"`
	NamaRombel         string `db:"nama_rombel" json:"nama_rombel"`
	StatusKeluar       string `db:"status_keluar" json:"status_keluar"`
	TanggalKeluar      string `db:"tanggal_keluar" json:"tanggal_keluar,omitempty"`
	NomorSurat         string `db:"nomor_surat" json:"nomor_surat,omitempty"`
	AlasanKeluar       string `db:"alasan_keluar" json:"alasan_keluar,omitempty"`
	Tujuan             string `db:"tujuan" json:"tujuan,omitempty"`
	Foto               string `db:"foto" json:"foto"`
	FotoURL            string `db:"-" json:"foto_url,omitempty"`
}

// ImportRowError merepresentasikan satu baris bermasalah saat import
type ImportRowError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}

// PesertaDidikConflictCheck dipakai untuk cek bentrok data unik dgn database saat import
type PesertaDidikConflictCheck struct {
	NIS       string `db:"nis" json:"nis"`
	NISN      string `db:"nisn" json:"nisn"`
	NoUjianSD string `db:"no_ujiansd" json:"no_ujiansd"`
	NIKSiswa  string `db:"nik_siswa" json:"nik_siswa"`
	NIKAyah   string `db:"nik_ayah" json:"nik_ayah"`
	NIKIbu    string `db:"nik_ibu" json:"nik_ibu"`
	NamaSiswa string `db:"nama_siswa" json:"nama_siswa"`
}
