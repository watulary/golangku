package ekskul

// Ekskul merepresentasikan data ekstrakurikuler,
// digabung dengan nama pembina dari tabel pegawai.
type Ekskul struct {
	IDEkskul      int    `db:"id_ekskul" json:"id_ekskul"`
	NamaEkskul    string `db:"nama_ekskul" json:"nama_ekskul"`
	NUPTK         string `db:"nuptk" json:"nuptk"`
	Aktif         bool   `db:"aktif" json:"aktif"`
	NamaPegawai   string `db:"nama_pegawai" json:"nama_pegawai"`
	JumlahAnggota int    `db:"jumlah_anggota" json:"jumlah_anggota"`
}

// AnggotaEkskul merepresentasikan relasi siswa (pesertadidik)
// terhadap satu ekskul, digabung dengan nama siswa dan nama rombel.
type AnggotaEkskul struct {
	ID         int    `db:"id" json:"id"`
	IDEkskul   int    `db:"id_ekskul" json:"id_ekskul"`
	NIS        string `db:"nis" json:"nis"`
	NamaSiswa  string `db:"nama_siswa" json:"nama_siswa"`
	NamaRombel string `db:"nama_rombel" json:"nama_rombel"`
	NamaEkskul string `db:"nama_ekskul" json:"nama_ekskul"`
}

// PegawaiOption dipakai untuk dropdown pemilihan pembina (jabatan = Guru) saat create/edit ekskul.
type PegawaiOption struct {
	NUPTK       string `db:"nuptk" json:"nuptk"`
	NamaPegawai string `db:"nama_pegawai" json:"nama_pegawai"`
	Jabatan     string `db:"jabatan" json:"jabatan"`
}

// PesertaDidikOption dipakai untuk dropdown/pencarian siswa aktif tahun berjalan
// saat menambah anggota ekskul.
type PesertaDidikOption struct {
	NIS        string `db:"nis" json:"nis"`
	NamaSiswa  string `db:"nama_siswa" json:"nama_siswa"`
	NamaRombel string `db:"nama_rombel" json:"nama_rombel"`
}
