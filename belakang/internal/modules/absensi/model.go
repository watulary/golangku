package absensi

// Absensi merepresentasikan satu baris record di tabel absensi
type Absensi struct {
	ID         int    `db:"id" json:"id"`
	Nis        string `db:"nis" json:"nis"`
	Tanggal    string `db:"tanggal" json:"tanggal"`
	Status     string `db:"status" json:"status"`
	Keterangan string `db:"keterangan" json:"keterangan"`
}

// AbsensiSiswa dipakai untuk daftar tampilan per rombel pada tanggal tertentu.
// Hasil LEFT JOIN absensi -> pesertadidik, urut alfabetis by nama_siswa.
// Kalau siswa belum punya record absensi utk tanggal ini: ID=0, Status default "Hadir".
type AbsensiSiswa struct {
	ID         int    `db:"id" json:"id"`
	Nis        string `db:"nis" json:"nis"`
	NamaSiswa  string `db:"nama_siswa" json:"nama_siswa"`
	NamaKelas  string `db:"nama_kelas" json:"nama_kelas"`
	NamaRombel string `db:"nama_rombel" json:"nama_rombel"`
	Status     string `db:"status" json:"status"`
	Keterangan string `db:"keterangan" json:"keterangan"`
}

// AbsensiBulkItem adalah satu baris item dalam request body submit bulk
type AbsensiBulkItem struct {
	Nis        string `json:"nis"`
	Status     string `json:"status"`
	Keterangan string `json:"keterangan"`
}

// AbsensiBulkRequest adalah payload submit absensi harian per rombel
type AbsensiBulkRequest struct {
	Tanggal  string            `json:"tanggal"`
	IDKelas  int               `json:"id_kelas"`
	IDRombel int               `json:"id_rombel"`
	Items    []AbsensiBulkItem `json:"items"`
}

// RekapBulananSiswa = 1 baris matriks rekap bulanan (1 siswa, status per tanggal)
type RekapBulananSiswa struct {
	Nis        string   `json:"nis"`
	NamaSiswa  string   `json:"nama_siswa"`
	NamaKelas  string   `json:"nama_kelas"`
	NamaRombel string   `json:"nama_rombel"`
	Harian     []string `json:"harian"` // index 0 = tanggal 1, dst. Isi: "H"/"S"/"I"/"A"/"" (kosong = belum ada record)
}

// AbsensiSayaHarian dipakai utk tampilan absensi milik siswa sendiri (portal Peserta Didik)
type AbsensiSayaHarian struct {
	Tanggal    string `db:"tanggal" json:"tanggal"`
	Status     string `db:"status" json:"status"`
	Keterangan string `db:"keterangan" json:"keterangan"`
}
