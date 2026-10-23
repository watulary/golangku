package nilai

// EkskulOpsi = pilihan ekskul (aktif) di layar input nilai ekskul.
type EkskulOpsi struct {
	IDEkskul   int    `db:"id_ekskul" json:"id_ekskul"`
	NamaEkskul string `db:"nama_ekskul" json:"nama_ekskul"`
}

// NilaiEkskulSiswa = satu baris layar input: siswa anggota ekskul di 1 rombel + nilainya
// untuk 1 tahun pelajaran / 1 semester. Belum ada nilai: IDNilai = 0, Nilai & Keterangan kosong.
type NilaiEkskulSiswa struct {
	IDNilai    int    `db:"id_nilai_ekskul" json:"id_nilai_ekskul"`
	Nis        string `db:"nis" json:"nis"`
	NamaSiswa  string `db:"nama_siswa" json:"nama_siswa"`
	NamaRombel string `db:"nama_rombel" json:"nama_rombel"`
	Nilai      string `db:"nilai" json:"nilai"`
	Keterangan string `db:"keterangan" json:"keterangan"`
}

type EkskulBulkItem struct {
	Nis        string `json:"nis"`
	Nilai      string `json:"nilai"`
	Keterangan string `json:"keterangan"`
}

type EkskulBulkRequest struct {
	IDTahunPelajaran int              `json:"id_tahunpelajaran"`
	Semester         string           `json:"semester"`
	IDRombel         int              `json:"id_rombel"`
	IDEkskul         int              `json:"id_ekskul"`
	IDKelasCatat     int              `json:"id_kelas_catat"`
	Items            []EkskulBulkItem `json:"items"`
}
