package guruwali

// GuruWali merepresentasikan tbl_guruwali
type GuruWali struct {
	IDGuruWali   int    `db:"id_guruwali" json:"id_guruwali"`
	NamaGuruWali string `db:"nama_guruwali" json:"nama_guruwali"`
	NUPTK        string `db:"nuptk" json:"nuptk"`
	Aktif        bool   `db:"aktif" json:"aktif"`
}

// GuruWaliListItem dipakai di GetAllAdmin / GetAllForExport (ada kolom computed)
type GuruWaliListItem struct {
	GuruWali
	NamaPegawai   string `db:"nama_pegawai" json:"nama_pegawai"`
	JumlahAnggota int    `db:"jumlah_anggota" json:"jumlah_anggota"`
}

// AnggotaGuruWali merepresentasikan tbl_anggota_guruwali
type AnggotaGuruWali struct {
	ID         int    `db:"id" json:"id"`
	IDGuruWali int    `db:"id_guruwali" json:"id_guruwali"`
	NIS        string `db:"nis" json:"nis"`
}

// AnggotaGuruWaliListItem dipakai di GetAnggotaByGuruWali (JOIN ke pesertadidik/rombel/kelas)
type AnggotaGuruWaliListItem struct {
	AnggotaGuruWali
	NamaSiswa  string `db:"nama_siswa" json:"nama_siswa"`
	NamaKelas  string `db:"nama_kelas" json:"nama_kelas"`
	NamaRombel string `db:"nama_rombel" json:"nama_rombel"`
}

// PegawaiOption dipakai di SearchPegawaiHandler (autocomplete pilih NUPTK di Form.vue)
type PegawaiOption struct {
	NUPTK       string `db:"nuptk" json:"nuptk"`
	NamaPegawai string `db:"nama_pegawai" json:"nama_pegawai"`
}

// PesertaDidikOption dipakai di SearchPesertaDidikHandler (autocomplete Tambah Anggota di Anggota.vue)
type PesertaDidikOption struct {
	NIS        string `db:"nis" json:"nis"`
	NamaSiswa  string `db:"nama_siswa" json:"nama_siswa"`
	NamaKelas  string `db:"nama_kelas" json:"nama_kelas"`
	NamaRombel string `db:"nama_rombel" json:"nama_rombel"`
}
