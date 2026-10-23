package tahunpelajaran

type TahunPelajaran struct {
	ID       int    `db:"id_tahunpelajaran" json:"id_tahunpelajaran"`
	Nama     string `db:"nama_tahunpelajaran" json:"nama_tahunpelajaran"`
	Semester string `db:"semester" json:"semester"`
	Aktif    int    `db:"aktif" json:"aktif"`
}
