package nilai

// NilaiSiswa = satu baris layar input: siswa di 1 rombel + nilainya utk
// 1 mapel / 1 semester / 1 tahun pelajaran. Siswa yang belum punya nilai:
// IDNilai = 0, NilaiAkhir = nil, Capaian & KKTP kosong.
type NilaiSiswa struct {
	IDNilai    int      `db:"id_nilai" json:"id_nilai"`
	Nis        string   `db:"nis" json:"nis"`
	NamaSiswa  string   `db:"nama_siswa" json:"nama_siswa"`
	NamaKelas  string   `db:"nama_kelas" json:"nama_kelas"`
	NamaRombel string   `db:"nama_rombel" json:"nama_rombel"`
	NilaiAkhir *float64 `db:"nilai_akhir" json:"nilai_akhir"`
	Capaian    string   `db:"capaian_pembelajaran" json:"capaian_pembelajaran"`
	Kktp       string   `db:"kktp" json:"kktp"`
}

// BulkItem = satu siswa dalam payload simpan massal.
// nilai_akhir dikirim null kalau dikosongkan.
type BulkItem struct {
	Nis        string   `json:"nis"`
	NilaiAkhir *float64 `json:"nilai_akhir"`
	Capaian    string   `json:"capaian_pembelajaran"`
	Kktp       string   `json:"kktp"`
}

// BulkRequest = payload simpan nilai 1 rombel + 1 semester + 1 mapel.
type BulkRequest struct {
	IDTahunPelajaran int        `json:"id_tahunpelajaran"`
	Semester         string     `json:"semester"`
	IDRombel         int        `json:"id_rombel"`
	IDMapel          int        `json:"id_mapel"`
	IDKelasCatat     int        `json:"id_kelas_catat"`
	Items            []BulkItem `json:"items"`
}

// RombelInfo dipakai utk validasi rombel & menurunkan id_kelas (kelas siswa saat nilai diinput).
type RombelInfo struct {
	IDKelas    int    `db:"id_kelas"`
	NamaRombel string `db:"nama_rombel"`
	NamaKelas  string `db:"nama_kelas"`
}
