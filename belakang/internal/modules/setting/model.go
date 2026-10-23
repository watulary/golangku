package setting

import "time"

type Setting struct {
	IDSekolah  int    `db:"id_sekolah" json:"id_sekolah"`
	IDPengguna string `db:"id_pengguna" json:"-"`

	// Identitas
	NamaSekolah      string `db:"nama_sekolah" json:"nama_sekolah"`
	LogoSekolah      string `db:"logo_sekolah" json:"logo_sekolah"`
	LogoURL          string `json:"logo_url"`
	NPSN             string `db:"npsn" json:"npsn"`
	NSSS             string `db:"nsss" json:"nsss"`
	NSP              string `db:"nsp" json:"nsp"`
	DeskripsiSekolah string `db:"deskripsi_sekolah" json:"deskripsi_sekolah"`

	// Kontak
	AlamatSekolah string `db:"alamat_sekolah" json:"alamat_sekolah"`
	TelpSekolah   string `db:"telp_sekolah" json:"telp_sekolah"`
	EmailSekolah  string `db:"email_sekolah" json:"email_sekolah"`
	WebSekolah    string `db:"web_sekolah" json:"web_sekolah"`
	PetaSekolah   string `db:"peta_sekolah" json:"peta_sekolah"`

	// Visi Sekolah
	VisiSekolah string `db:"visi_sekolah" json:"visi_sekolah"`

	// Wilayah
	KelurahanSekolah string `db:"kelurahan_sekolah" json:"kelurahan_sekolah"`
	KecamatanSekolah string `db:"kecamatan_sekolah" json:"kecamatan_sekolah"`
	KabupatenSekolah string `db:"kabupaten_sekolah" json:"kabupaten_sekolah"`
	ProvinsiSekolah  string `db:"provinsi_sekolah" json:"provinsi_sekolah"`
	KodePosSekolah   string `db:"kodepos_sekolah" json:"kodepos_sekolah"`

	// Sosial Media
	Facebook  string `db:"facebook" json:"facebook"`
	Instagram string `db:"instagram" json:"instagram"`
	Twitter   string `db:"twitter" json:"twitter"`
	Youtube   string `db:"youtube" json:"youtube"`

	// Footer
	FooterSekolah string `db:"footer_sekolah" json:"footer_sekolah"`
	ProfilSekolah string `db:"profil_sekolah" json:"profil_sekolah"`

	// Heroarea
	HeroArea1 string `db:"heroarea1" json:"heroarea1"`
	HeroArea2 string `db:"heroarea2" json:"heroarea2"`
	HeroArea3 string `db:"heroarea3" json:"heroarea3"`
	HeroArea4 string `db:"heroarea4" json:"heroarea4"`

	// Audit
	UpdatedBy     string    `db:"updated_by" json:"updated_by"`
	UpdatedByNama string    `db:"updated_by_nama" json:"updated_by_nama"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
