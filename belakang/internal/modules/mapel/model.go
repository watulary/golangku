package mapel

// Mapel = satu baris tbl_mapel.
// IDInduk NULL  -> mapel utama (mis. "Bahasa Indonesia", "Seni Budaya & Prakarya")
// IDInduk terisi -> sub-mapel dari mapel utama tsb (mis. "Seni Rupa" di bawah Seni Budaya & Prakarya)
// NamaInduk hanya hasil JOIN (read-only), tidak ikut disimpan.
type Mapel struct {
	ID        int     `db:"id_mapel" json:"id_mapel"`
	IDInduk   *int    `db:"id_induk" json:"id_induk"`
	Nama      string  `db:"nama_mapel" json:"nama_mapel"`
	Urutan    int     `db:"urutan" json:"urutan"`
	Aktif     int     `db:"aktif" json:"aktif"`
	NamaInduk *string `db:"nama_induk" json:"nama_induk"`
}
