package pegawai

// Pegawai merepresentasikan data lengkap kepegawaian (tabel legacy, PK: nuptk)
type Pegawai struct {
	Nuptk         string `db:"nuptk" json:"nuptk"`
	NamaPegawai   string `db:"nama_pegawai" json:"nama_pegawai"`
	Jabatan       string `db:"jabatan" json:"jabatan"`
	Strata        string `db:"strata" json:"strata"`
	PnsNon        string `db:"pns_non" json:"pns_non"`
	Nikki         string `db:"nikki" json:"nikki"`
	Simpkbid      string `db:"simpkbid" json:"simpkbid"`
	Dapodikid     string `db:"dapodikid" json:"dapodikid"`
	NipLama       string `db:"nip_lama" json:"nip_lama"`
	NipBaru       string `db:"nip_baru" json:"nip_baru"`
	Nrk           string `db:"nrk" json:"nrk"`
	Nrg           string `db:"nrg" json:"nrg"`
	TempatLahir   string `db:"tempat_lahir" json:"tempat_lahir"`
	TanggalLahir  string `db:"tanggal_lahir" json:"tanggal_lahir"`
	JenisKelamin  string `db:"jenis_kelamin" json:"jenis_kelamin"`
	Agama         string `db:"agama" json:"agama"`
	Nik           string `db:"nik" json:"nik"`
	NoKk          string `db:"no_kk" json:"no_kk"`
	Alamat        string `db:"alamat" json:"alamat"`
	Rt            string `db:"rt" json:"rt"`
	Rw            string `db:"rw" json:"rw"`
	Kel           string `db:"kel" json:"kel"`
	Kec           string `db:"kec" json:"kec"`
	Kopos         string `db:"kopos" json:"kopos"`
	TelpRmh       string `db:"telp_rmh" json:"telp_rmh"`
	TelpHp        string `db:"telp_hp" json:"telp_hp"`
	Email         string `db:"email" json:"email"`
	SkSekolah     string `db:"sk_sekolah" json:"sk_sekolah"`
	TmtSekolah    string `db:"tmt_sekolah" json:"tmt_sekolah"`
	SkNonkki      string `db:"sk_nonkki" json:"sk_nonkki"`
	TmtNonkki     string `db:"tmt_nonkki" json:"tmt_nonkki"`
	SkKki         string `db:"sk_kki" json:"sk_kki"`
	TmtKki        string `db:"tmt_kki" json:"tmt_kki"`
	SkCpns        string `db:"sk_cpns" json:"sk_cpns"`
	TmtCpns       string `db:"tmt_cpns" json:"tmt_cpns"`
	SkPns         string `db:"sk_pns" json:"sk_pns"`
	TmtPns        string `db:"tmt_pns" json:"tmt_pns"`
	TmtKgb        string `db:"tmt_kgb" json:"tmt_kgb"`
	PangkatKgb    string `db:"pangkat_kgb" json:"pangkat_kgb"`
	GajiKgb       string `db:"gaji_kgb" json:"gaji_kgb"`
	SkKgb         string `db:"sk_kgb" json:"sk_kgb"`
	TglskKgb      string `db:"tglsk_kgb" json:"tglsk_kgb"`
	Karpeg        string `db:"karpeg" json:"karpeg"`
	Karis         string `db:"karis" json:"karis"`
	NamaIbukdg    string `db:"nama_ibukdg" json:"nama_ibukdg"`
	NmSuamiIstri  string `db:"nm_suami_istri" json:"nm_suami_istri"`
	PekSuamiIstri string `db:"pek_suami_istri" json:"pek_suami_istri"`
	Npwp          string `db:"npwp" json:"npwp"`
	StatusKeluar  string `db:"status_keluar" json:"status_keluar"`
	KetKeluar     string `db:"keterangan_keluar" json:"keterangan_keluar"`
	TanggalKeluar string `db:"tanggal_keluar" json:"tanggal_keluar"`
	NomorSurat    string `db:"nomor_surat" json:"nomor_surat"`
	Foto          string `db:"foto" json:"foto"`
	FotoURL       string `db:"-" json:"foto_url"`
}

// PegawaiListItem — versi ringkas untuk GetAllAdmin (tabel listing)
type PegawaiListItem struct {
	Nuptk         string `db:"nuptk" json:"nuptk"`
	NamaPegawai   string `db:"nama_pegawai" json:"nama_pegawai"`
	NipBaru       string `db:"nip_baru" json:"nip_baru"`
	Jabatan       string `db:"jabatan" json:"jabatan"`
	TempatLahir   string `db:"tempat_lahir" json:"tempat_lahir"`
	TanggalLahir  string `db:"tanggal_lahir" json:"tanggal_lahir"`
	Agama         string `db:"agama" json:"agama"`
	PnsNon        string `db:"pns_non" json:"pns_non"`
	JenisKelamin  string `db:"jenis_kelamin" json:"jenis_kelamin"`
	TelpHp        string `db:"telp_hp" json:"telp_hp"`
	Email         string `db:"email" json:"email"`
	Foto          string `db:"foto" json:"foto"`
	FotoURL       string `db:"-" json:"foto_url"`
	StatusKeluar  string `db:"status_keluar" json:"status_keluar"`
	KetKeluar     string `db:"keterangan_keluar" json:"keterangan_keluar"`
	TanggalKeluar string `db:"tanggal_keluar" json:"tanggal_keluar"`
	NomorSurat    string `db:"nomor_surat" json:"nomor_surat"`
}

// PegawaiCreateRequest — payload create; Nuptk wajib diisi manual (PK, bukan auto/UUID)
type PegawaiCreateRequest struct {
	Nuptk         string `json:"nuptk"`
	NamaPegawai   string `json:"nama_pegawai"`
	Jabatan       string `json:"jabatan"`
	Strata        string `json:"strata"`
	PnsNon        string `json:"pns_non"`
	Nikki         string `json:"nikki"`
	Simpkbid      string `json:"simpkbid"`
	Dapodikid     string `json:"dapodikid"`
	NipLama       string `json:"nip_lama"`
	NipBaru       string `json:"nip_baru"`
	Nrk           string `json:"nrk"`
	Nrg           string `json:"nrg"`
	TempatLahir   string `json:"tempat_lahir"`
	TanggalLahir  string `json:"tanggal_lahir"`
	JenisKelamin  string `json:"jenis_kelamin"`
	Agama         string `json:"agama"`
	Nik           string `json:"nik"`
	NoKk          string `json:"no_kk"`
	Alamat        string `json:"alamat"`
	Rt            string `json:"rt"`
	Rw            string `json:"rw"`
	Kel           string `json:"kel"`
	Kec           string `json:"kec"`
	Kopos         string `json:"kopos"`
	TelpRmh       string `json:"telp_rmh"`
	TelpHp        string `json:"telp_hp"`
	Email         string `json:"email"`
	SkSekolah     string `json:"sk_sekolah"`
	TmtSekolah    string `json:"tmt_sekolah"`
	SkNonkki      string `json:"sk_nonkki"`
	TmtNonkki     string `json:"tmt_nonkki"`
	SkKki         string `json:"sk_kki"`
	TmtKki        string `json:"tmt_kki"`
	SkCpns        string `json:"sk_cpns"`
	TmtCpns       string `json:"tmt_cpns"`
	SkPns         string `json:"sk_pns"`
	TmtPns        string `json:"tmt_pns"`
	TmtKgb        string `json:"tmt_kgb"`
	PangkatKgb    string `json:"pangkat_kgb"`
	GajiKgb       string `json:"gaji_kgb"`
	SkKgb         string `json:"sk_kgb"`
	TglskKgb      string `json:"tglsk_kgb"`
	Karpeg        string `json:"karpeg"`
	Karis         string `json:"karis"`
	NamaIbukdg    string `json:"nama_ibukdg"`
	NmSuamiIstri  string `json:"nm_suami_istri"`
	PekSuamiIstri string `json:"pek_suami_istri"`
	Npwp          string `json:"npwp"`
	StatusKeluar  string `json:"status_keluar"`
	KetKeluar     string `json:"keterangan_keluar"`
	TanggalKeluar string `json:"tanggal_keluar"`
	NomorSurat    string `json:"nomor_surat"`
	Foto          string `json:"foto"`
}

// PegawaiUpdateRequest — sama seperti Create, tanpa Nuptk (diambil dari param URL :nuptk)
type PegawaiUpdateRequest struct {
	NamaPegawai   string `json:"nama_pegawai"`
	Jabatan       string `json:"jabatan"`
	Strata        string `json:"strata"`
	PnsNon        string `json:"pns_non"`
	Nikki         string `json:"nikki"`
	Simpkbid      string `json:"simpkbid"`
	Dapodikid     string `json:"dapodikid"`
	NipLama       string `json:"nip_lama"`
	NipBaru       string `json:"nip_baru"`
	Nrk           string `json:"nrk"`
	Nrg           string `json:"nrg"`
	TempatLahir   string `json:"tempat_lahir"`
	TanggalLahir  string `json:"tanggal_lahir"`
	JenisKelamin  string `json:"jenis_kelamin"`
	Agama         string `json:"agama"`
	Nik           string `json:"nik"`
	NoKk          string `json:"no_kk"`
	Alamat        string `json:"alamat"`
	Rt            string `json:"rt"`
	Rw            string `json:"rw"`
	Kel           string `json:"kel"`
	Kec           string `json:"kec"`
	Kopos         string `json:"kopos"`
	TelpRmh       string `json:"telp_rmh"`
	TelpHp        string `json:"telp_hp"`
	Email         string `json:"email"`
	SkSekolah     string `json:"sk_sekolah"`
	TmtSekolah    string `json:"tmt_sekolah"`
	SkNonkki      string `json:"sk_nonkki"`
	TmtNonkki     string `json:"tmt_nonkki"`
	SkKki         string `json:"sk_kki"`
	TmtKki        string `json:"tmt_kki"`
	SkCpns        string `json:"sk_cpns"`
	TmtCpns       string `json:"tmt_cpns"`
	SkPns         string `json:"sk_pns"`
	TmtPns        string `json:"tmt_pns"`
	TmtKgb        string `json:"tmt_kgb"`
	PangkatKgb    string `json:"pangkat_kgb"`
	GajiKgb       string `json:"gaji_kgb"`
	SkKgb         string `json:"sk_kgb"`
	TglskKgb      string `json:"tglsk_kgb"`
	Karpeg        string `json:"karpeg"`
	Karis         string `json:"karis"`
	NamaIbukdg    string `json:"nama_ibukdg"`
	NmSuamiIstri  string `json:"nm_suami_istri"`
	PekSuamiIstri string `json:"pek_suami_istri"`
	Npwp          string `json:"npwp"`
	StatusKeluar  string `json:"status_keluar"`
	KetKeluar     string `json:"keterangan_keluar"`
	TanggalKeluar string `json:"tanggal_keluar"`
	NomorSurat    string `json:"nomor_surat"`
	Foto          string `json:"foto"`
}

// PegawaiOption dipakai untuk dropdown autocomplete (search-as-you-type)
type PegawaiOption struct {
	Nuptk       string `db:"nuptk" json:"nuptk"`
	NamaPegawai string `db:"nama_pegawai" json:"nama_pegawai"`
	Jabatan     string `db:"jabatan" json:"jabatan"`
}

// PegawaiConflictCheck dipakai untuk cek bentrok NUPTK/NIP Baru/NIK dgn data existing di DB
type PegawaiConflictCheck struct {
	Nuptk       string `db:"nuptk" json:"nuptk"`
	NipBaru     string `db:"nip_baru" json:"nip_baru"`
	Nik         string `db:"nik" json:"nik"`
	NamaPegawai string `db:"nama_pegawai" json:"nama_pegawai"`
}

// ImportRowError merepresentasikan satu baris bermasalah saat import
type ImportRowError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}
