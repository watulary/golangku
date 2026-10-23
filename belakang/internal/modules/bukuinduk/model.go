package bukuinduk

// SumberData menandakan asal tabel data siswa yang dicetak.
type SumberData string

const (
	SumberAktif  SumberData = "aktif"  // tabel pesertadidik
	SumberMutasi SumberData = "mutasi" // tabel pesertadidik_mutasi
	SumberArsip  SumberData = "arsip"  // tabel pesertadidik_arsip (PK kode_siswa)
)

// StatusKeluar menandakan cabang Bagian D (Meninggalkan Sekolah) yang berlaku.
type StatusKeluar string

const (
	StatusKeluarNone   StatusKeluar = ""       // masih aktif, Bagian D dikosongkan
	StatusKeluarLulus  StatusKeluar = "lulus"  // dari cek tabel kelulusan
	StatusKeluarPindah StatusKeluar = "pindah" // dari pesertadidik_mutasi.status_keluar = "Pindah"
	StatusKeluarKeluar StatusKeluar = "keluar" // dari pesertadidik_mutasi.status_keluar = "Keluar"/"Meninggal"
)

// BukuIndukData adalah struct unifikasi Lembar IV, diisi oleh 3 fungsi
// converter terpisah (mapAktifToBukuInduk, mapMutasiToBukuInduk,
// mapArsipToBukuInduk) di repository.go, ditambah enrichment dari
// tabel kelulusan. Generator PDF hanya bekerja dari struct ini,
// tidak peduli asal tabelnya.
type BukuIndukData struct {
	Sumber SumberData

	// Identitas
	NIS              string
	NISN             string
	Kelas            string
	NIK              string // nik_siswa
	NoKK             string // no_kk
	WaliKelasNama    string
	WaliKelasNIPNRK  string
	AnakKe           string // anak_ke
	SaudaraTiri      string // saudara_tiri
	SaudaraAngkat    string // saudara_angkat
	YatimPiatu       string // yatim_piatu
	PernahSakit      string // pernah_sakit
	KebutuhanKhusus  string // kebutuhan_khusus
	LamaBelajar      string // lama_belajar
	TempatLahirAyah  string
	TanggalLahirAyah string
	AgamaAyah        string
	WargaAyah        string
	GajiAyah         string
	AlamatAyah       string
	TelpAyah         string
	StatusAyah       string

	TempatLahirIbu  string
	TanggalLahirIbu string
	AgamaIbu        string
	WargaIbu        string
	GajiIbu         string
	AlamatIbu       string
	TelpIbu         string
	StatusIbu       string

	TempatLahirWali  string
	TanggalLahirWali string
	AgamaWali        string
	WargaWali        string
	GajiWali         string
	AlamatWali       string
	TelpWali         string

	// A. Keterangan Siswa
	NamaLengkap          string // nama_siswa
	NamaPanggilan        string // computed: kata pertama dari nama_siswa
	JenisKelamin         string // jenis_kelamin ("L"/"P")
	TempatLahir          string // tempat_lahirsiswa
	TanggalLahir         string // tanggal_lahirsiswa
	Agama                string // agama_siswa
	Kewarganegaraan      string // warga_siswa
	JumlahSaudara        string // saudara_kandung
	BahasaSehariHari     string // bahasa_siswa
	GolonganDarah        string // gol_darah
	Alamat               string // computed: gabungan alamat_siswa+rt+rw+kel+kec+kopos
	NoTelpHp             string // computed: telprmh_siswa / hp_siswa
	BertempatTinggalPada string // tinggal_siswa
	JarakSekolah         string // jarak_sekolah

	// B. Keterangan Orangtua/Wali
	NamaAyah       string // nama_ayah
	PendidikanAyah string // pend_ayah
	PekerjaanAyah  string // pek_ayah
	NamaIbu        string // nama_ibu
	PendidikanIbu  string // pend_ibu
	PekerjaanIbu   string // pek_ibu
	NamaWali       string // nama_wali
	HubunganWali   string // selalu kosong (titik-titik, isi manual)
	PendidikanWali string // pend_wali
	PekerjaanWali  string // pek_wali

	// C. Perkembangan Siswa (poin 16)
	AsalSekolahMasuk        string // lulus_sd (asal & nama SD/MI, dipakai utk 2 baris form)
	TanggalIjazahMasuk      string // tanggal_ijazah
	NomorIjazahMasuk        string // nomor_ijazah
	AsalSekolahPindahan     string // asal_smp (kosong -> titik-titik)
	DariTingkatPindahan     string // selalu kosong (tidak ada kolom)
	DiterimaTanggalPindahan string // tanggal_masuk (ASUMSI, tolong dikoreksi kalau keliru)

	// D. Meninggalkan Sekolah
	StatusKeluar StatusKeluar

	// isi hanya salah satu blok berikut, tergantung StatusKeluar
	LulusTahun         string // dari tbl_tahunpelajaran.nama_tahunpelajaran (JOIN kelulusan.id_tahunpelajaran)
	LulusNomorIjazah   string // selalu kosong (tidak ada kolomnya di tabel kelulusan)
	LulusMelanjutkanKe string // selalu kosong

	PindahTingkatKelas   string // selalu kosong (tidak ada kolom)
	PindahKeSekolah      string // tujuan
	PindahKeTingkatKelas string // selalu kosong

	KeluarAlasan  string // alasan_keluar, atau "Meninggal" kalau status_keluar="Meninggal"
	KeluarTanggal string // tanggal_keluar

	// E. Fisik & Kesehatan (hanya nilai snapshot terkini, sisa kolom form dikosongkan saat render PDF)
	TinggiBadan string // tinggi_badan
	BeratBadan  string // berat_badan

	HobiSeni       string
	HobiOrkes      string
	HobiOrganisasi string
	HobiLain       string
	EkskulList     []string

	FotoURL string
}

// SekolahProfil adalah data profil sekolah dari tabel setting, dipakai
// utk halaman Cover Buku Induk.
type SekolahProfil struct {
	NamaSekolah      string `db:"nama_sekolah"`
	AlamatSekolah    string `db:"alamat_sekolah"`
	KelurahanSekolah string `db:"kelurahan_sekolah"`
	KecamatanSekolah string `db:"kecamatan_sekolah"`
	KabupatenSekolah string `db:"kabupaten_sekolah"`
	ProvinsiSekolah  string `db:"provinsi_sekolah"`
	NSS              string `db:"nss"`
	NPSN             string `db:"npsn"`
	NSP              string `db:"nsp"`
}

// KepalaSekolahInfo dari tabel kelulusan_waktu, dipakai utk blok tanda
// tangan "Mengetahui, Kepala Sekolah" di beberapa halaman statis.
type KepalaSekolahInfo struct {
	NamaKepsek string `db:"nama_kepsek"`
	NIPKepsek  string `db:"nip_kepsek"`
}

// NomorUrutItem adalah satu baris di Bagian II (Nomor Urut Peserta
// Didik Dalam Buku Induk Register) — daftar ringkas seluruh siswa.
type NomorUrutItem struct {
	No           int
	NIS          string
	NISN         string
	Nama         string
	JenisKelamin string
	Agama        string
	TempatLahir  string
	TanggalLahir string
	NIK          string
	NoKK         string
	Kelas        string
}
