package kelulusan

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"sekolah/internal/database"
)

var allowedSortColumns = map[string]string{
	"nama_siswa":          "pesertadidik.nama_siswa",
	"nis":                 "kelulusan.nis",
	"nama_kelas":          "tbl_kelas.nama_kelas",
	"nama_rombel":         "tbl_rombel.nama_rombel",
	"nama_tahunpelajaran": "tbl_tahunpelajaran.nama_tahunpelajaran",
	"ket":                 "kelulusan.ket",
	"waktu":               "kelulusan.waktu",
	"nopes":               "kelulusan.nopes",
	// rata_s kolom varchar -- CAST ke DECIMAL biar sort NUMERIK, bukan
	// alfabet (kalau alfabet, "9.5" dianggap lebih besar dari "10.2").
	"rata_s": "CAST(kelulusan.rata_s AS DECIMAL(6,2))",
}

const listSelectSQL = `SELECT kelulusan.*,
	COALESCE(pesertadidik.nama_siswa,'')            AS nama_siswa,
	COALESCE(tbl_kelas.nama_kelas,'')                AS nama_kelas,
	COALESCE(tbl_rombel.nama_rombel,'')              AS nama_rombel,
	COALESCE(tbl_tahunpelajaran.nama_tahunpelajaran,'') AS nama_tahunpelajaran,
	COALESCE(pesertadidik.nisn,'')                   AS nisn,
	COALESCE(pesertadidik.tempat_lahirsiswa,'')      AS tempat_lahirsiswa,
	COALESCE(pesertadidik.tanggal_lahirsiswa,'')     AS tanggal_lahirsiswa`

const listJoinSQL = `FROM kelulusan
	LEFT JOIN pesertadidik ON pesertadidik.nis = kelulusan.nis
	LEFT JOIN tbl_kelas ON tbl_kelas.id_kelas = kelulusan.id_kelas
	LEFT JOIN tbl_rombel ON tbl_rombel.id_rombel = kelulusan.id_rombel
	LEFT JOIN tbl_tahunpelajaran ON tbl_tahunpelajaran.id_tahunpelajaran = kelulusan.id_tahunpelajaran`

const searchWhereSQL = ` WHERE pesertadidik.nama_siswa LIKE ? OR kelulusan.nis LIKE ? OR kelulusan.nopes LIKE ?`

func GetAllAdmin(search, sortBy, sortDir string, limit, offset int) ([]KelulusanAdmin, error) {
	col, ok := allowedSortColumns[sortBy]
	if !ok {
		col = "pesertadidik.nama_siswa"
	}
	if sortDir != "desc" {
		sortDir = "asc"
	}

	query := listSelectSQL + " " + listJoinSQL
	args := []interface{}{}
	if search != "" {
		query += searchWhereSQL
		s := "%" + search + "%"
		args = append(args, s, s, s)
	}
	// Tie-breaker permanen Kelas -> Nama Siswa, apa pun kolom sort yg dipilih
	// user, supaya urutan selalu stabil dan enak dibaca per rombongan kelas.
	query += fmt.Sprintf(
		" ORDER BY %s %s, tbl_kelas.nama_kelas ASC, pesertadidik.nama_siswa ASC LIMIT ? OFFSET ?",
		col, sortDir,
	)
	args = append(args, limit, offset)

	var rows []KelulusanAdmin
	if err := database.DB.Select(&rows, query, args...); err != nil {
		return nil, err
	}
	return rows, nil
}

func CountAdmin(search string) (int, error) {
	query := "SELECT COUNT(*) " + listJoinSQL
	args := []interface{}{}
	if search != "" {
		query += searchWhereSQL
		s := "%" + search + "%"
		args = append(args, s, s, s)
	}
	var count int
	if err := database.DB.Get(&count, query, args...); err != nil {
		return 0, err
	}
	return count, nil
}

func GetAllForExport(search string) ([]KelulusanAdmin, error) {
	query := listSelectSQL + " " + listJoinSQL
	args := []interface{}{}
	if search != "" {
		query += searchWhereSQL
		s := "%" + search + "%"
		args = append(args, s, s, s)
	}
	query += " ORDER BY tbl_rombel.nama_rombel, pesertadidik.nama_siswa ASC"

	var rows []KelulusanAdmin
	if err := database.DB.Select(&rows, query, args...); err != nil {
		return nil, err
	}
	return rows, nil
}

func GetByNIS(nis string) (*Kelulusan, error) {
	var k Kelulusan
	if err := database.DB.Get(&k, `SELECT * FROM kelulusan WHERE nis = ?`, nis); err != nil {
		return nil, err
	}
	return &k, nil
}

// GetByNISAdmin - data kelulusan siswa tunggal berikut nama_siswa/kelas/
// rombel/tahunpelajaran, dipakai baik oleh Form.vue edit admin maupun
// endpoint /kelulusan/me milik Peserta Didik (reuse fungsi ini, filter
// pakai claims.NIS -- lihat pola RESMI di modul Kenaikan).
func GetByNISAdmin(nis string) (*KelulusanAdmin, error) {
	query := listSelectSQL + " " + listJoinSQL + " WHERE kelulusan.nis = ?"
	var row KelulusanAdmin
	if err := database.DB.Get(&row, query, nis); err != nil {
		return nil, err
	}
	return &row, nil
}

// ExistsPesertaDidik - validasi NIS terdaftar di pesertadidik sebelum
// insert/update/import, sama pola dengan kenaikan.
func ExistsPesertaDidik(nis string) (bool, error) {
	var count int
	err := database.DB.Get(&count, `SELECT COUNT(*) FROM pesertadidik WHERE nis = ?`, nis)
	return count > 0, err
}

func ExistsNIS(nis string) (bool, error) {
	var count int
	err := database.DB.Get(&count, `SELECT COUNT(*) FROM kelulusan WHERE nis = ?`, nis)
	return count > 0, err
}

func ExistsNopes(nopes string) (bool, error) {
	var count int
	err := database.DB.Get(&count, `SELECT COUNT(*) FROM kelulusan WHERE nopes = ?`, nopes)
	return count > 0, err
}

func ExistsFotoLulus(fotolulus string) (bool, error) {
	var count int
	err := database.DB.Get(&count, `SELECT COUNT(*) FROM kelulusan WHERE fotolulus = ?`, fotolulus)
	return count > 0, err
}

func Create(data map[string]interface{}) error {
	_, err := database.DB.NamedExec(`INSERT INTO kelulusan
		(id_tahunpelajaran, id_kelas, id_rombel, nis, nopes, pa_s, ppkn_s, bind_s, bing_s,
		 mtk_s, ipa_s, ips_s, penjas_s, senbud_s, prakarya_s, jumlah, rata_s, ket, fotolulus)
		VALUES
		(:id_tahunpelajaran, :id_kelas, :id_rombel, :nis, :nopes, :pa_s, :ppkn_s, :bind_s, :bing_s,
		 :mtk_s, :ipa_s, :ips_s, :penjas_s, :senbud_s, :prakarya_s, :jumlah, :rata_s, :ket, :fotolulus)`, data)
	return err
}

func Update(nis string, data map[string]interface{}) error {
	data["nis"] = nis
	_, err := database.DB.NamedExec(`UPDATE kelulusan SET
		id_tahunpelajaran=:id_tahunpelajaran, id_kelas=:id_kelas, id_rombel=:id_rombel, nopes=:nopes,
		pa_s=:pa_s, ppkn_s=:ppkn_s, bind_s=:bind_s, bing_s=:bing_s, mtk_s=:mtk_s, ipa_s=:ipa_s,
		ips_s=:ips_s, penjas_s=:penjas_s, senbud_s=:senbud_s, prakarya_s=:prakarya_s,
		jumlah=:jumlah, rata_s=:rata_s, ket=:ket, fotolulus=:fotolulus
		WHERE nis=:nis`, data)
	return err
}

func Delete(nis string) error {
	_, err := database.DB.Exec(`DELETE FROM kelulusan WHERE nis = ?`, nis)
	return err
}

func BulkDelete(nisList []string) error {
	if len(nisList) == 0 {
		return nil
	}
	query, args, err := sqlx.In(`DELETE FROM kelulusan WHERE nis IN (?)`, nisList)
	if err != nil {
		return err
	}
	query = database.DB.Rebind(query)
	_, err = database.DB.Exec(query, args...)
	return err
}

// --- kelulusan_waktu (baris tunggal id_waktu = 1) ---

func GetWaktu() (*KelulusanWaktu, error) {
	var w KelulusanWaktu
	// Kolom nullable (lihat DESCRIBE kelulusan_waktu) WAJIB dibungkus
	// COALESCE(...,'') -- kalau tidak, Scan error "converting NULL to string"
	// begitu ada baris yang belum pernah diisi salah satu field opsional.
	query := `SELECT
		kelulusan_waktu.id_waktu,
		kelulusan_waktu.id_pengguna,
		DATE_FORMAT(kelulusan_waktu.aturwaktu, '%Y-%m-%d %H:%i:%s') AS aturwaktu,
		kelulusan_waktu.aturform,
		kelulusan_waktu.heroarea1,
		kelulusan_waktu.heroarea2,
		kelulusan_waktu.waktu,
		COALESCE(kelulusan_waktu.judul,'')          AS judul,
		COALESCE(kelulusan_waktu.nomor,'')          AS nomor,
		COALESCE(kelulusan_waktu.pengantar,'')      AS pengantar,
		COALESCE(kelulusan_waktu.penutup,'')        AS penutup,
		COALESCE(kelulusan_waktu.akhir,'')          AS akhir,
		COALESCE(kelulusan_waktu.tanggal,'')        AS tanggal,
		COALESCE(kelulusan_waktu.mengetahui,'')     AS mengetahui,
		COALESCE(kelulusan_waktu.kepsek,'')         AS kepsek,
		COALESCE(kelulusan_waktu.nipkepsek,'')      AS nipkepsek,
		COALESCE(kelulusan_waktu.kop1_skl,'')       AS kop1_skl,
		COALESCE(kelulusan_waktu.kop2_skl,'')       AS kop2_skl,
		COALESCE(kelulusan_waktu.kop3_skl,'')       AS kop3_skl,
		COALESCE(kelulusan_waktu.kop4_skl,'')       AS kop4_skl,
		COALESCE(kelulusan_waktu.kop5_skl,'')       AS kop5_skl,
		COALESCE(kelulusan_waktu.kop6_skl,'')       AS kop6_skl,
		COALESCE(kelulusan_waktu.kop7_skl,'')       AS kop7_skl,
		COALESCE(kelulusan_waktu.kop8_skl,'')       AS kop8_skl,
		COALESCE(kelulusan_waktu.isi1_skl,'')       AS isi1_skl,
		COALESCE(kelulusan_waktu.isi2_skl,'')       AS isi2_skl,
		COALESCE(kelulusan_waktu.isi3_skl,'')       AS isi3_skl,
		COALESCE(kelulusan_waktu.isi4_skl,'')       AS isi4_skl,
		COALESCE(kelulusan_waktu.isi5_skl,'')       AS isi5_skl,
		kelulusan_waktu.tgl_skl,
		COALESCE(kelulusan_waktu.sekolah_skl,'')    AS sekolah_skl,
		COALESCE(kelulusan_waktu.kepsek_skl,'')     AS kepsek_skl,
		COALESCE(kelulusan_waktu.nipkepsek_skl,'')  AS nipkepsek_skl,
		COALESCE(kelulusan_waktu.ttd_skl,'')        AS ttd_skl,
		COALESCE(kelulusan_waktu.tahun_pelajaran,'') AS tahun_pelajaran,
		COALESCE(pengguna.nama_pengguna,'') AS nama_pengguna
		FROM kelulusan_waktu
		LEFT JOIN pengguna ON pengguna.id_pengguna = kelulusan_waktu.id_pengguna
		WHERE kelulusan_waktu.id_waktu = 1`
	if err := database.DB.Get(&w, query); err != nil {
		return nil, err
	}
	return &w, nil
}

func UpdatePengumuman(data map[string]interface{}) error {
	_, err := database.DB.NamedExec(`UPDATE kelulusan_waktu SET
		id_pengguna=:id_pengguna, aturwaktu=:aturwaktu, aturform=:aturform, heroarea1=:heroarea1,
		heroarea2=:heroarea2, judul=:judul, nomor=:nomor, pengantar=:pengantar, penutup=:penutup,
		akhir=:akhir, tanggal=:tanggal, mengetahui=:mengetahui, kepsek=:kepsek, nipkepsek=:nipkepsek,
		tahun_pelajaran=:tahun_pelajaran
		WHERE id_waktu = 1`, data)
	return err
}

func UpdateSKL(data map[string]interface{}) error {
	_, err := database.DB.NamedExec(`UPDATE kelulusan_waktu SET
		id_pengguna=:id_pengguna,
		kop1_skl=:kop1_skl, kop2_skl=:kop2_skl, kop3_skl=:kop3_skl, kop4_skl=:kop4_skl, kop5_skl=:kop5_skl,
		kop6_skl=:kop6_skl, kop7_skl=:kop7_skl, kop8_skl=:kop8_skl, isi1_skl=:isi1_skl, isi2_skl=:isi2_skl,
		isi3_skl=:isi3_skl, isi4_skl=:isi4_skl, isi5_skl=:isi5_skl, tgl_skl=:tgl_skl, sekolah_skl=:sekolah_skl,
		kepsek_skl=:kepsek_skl, nipkepsek_skl=:nipkepsek_skl, ttd_skl=:ttd_skl
		WHERE id_waktu = 1`, data)
	return err
}

func TruncateKelulusan() error {
	_, err := database.DB.Exec(`TRUNCATE TABLE kelulusan`)
	return err
}

// --- data gabungan utk cetak PDF ---

func GetCetakData(nis string) (*CetakData, error) {
	var c CetakData
	query := `SELECT
		COALESCE(pesertadidik.nama_siswa,'')         AS nama_siswa,
		COALESCE(pesertadidik.tempat_lahirsiswa,'')  AS tempat_lahirsiswa,
		COALESCE(pesertadidik.tanggal_lahirsiswa,'') AS tanggal_lahirsiswa,
		COALESCE(pesertadidik.nisn,'')                AS nisn,
		COALESCE(setting.nama_sekolah,'')             AS nama_sekolah,
		COALESCE(setting.npsn,'')                     AS npsn,
		COALESCE(tbl_kelas.nama_kelas,'')              AS nama_kelas,
		COALESCE(tbl_rombel.nama_rombel,'')            AS nama_rombel,
		COALESCE(tbl_tahunpelajaran.nama_tahunpelajaran,'') AS nama_tahunpelajaran,
		kelulusan.nis, kelulusan.nopes, kelulusan.ket, kelulusan.pa_s, kelulusan.ppkn_s, kelulusan.bind_s,
		kelulusan.bing_s, kelulusan.mtk_s, kelulusan.ipa_s, kelulusan.ips_s, kelulusan.penjas_s,
		kelulusan.senbud_s, kelulusan.prakarya_s, kelulusan.jumlah, kelulusan.rata_s, kelulusan.fotolulus
		FROM kelulusan
		LEFT JOIN pesertadidik ON pesertadidik.nis = kelulusan.nis
		LEFT JOIN tbl_kelas ON tbl_kelas.id_kelas = kelulusan.id_kelas
		LEFT JOIN tbl_rombel ON tbl_rombel.id_rombel = kelulusan.id_rombel
		LEFT JOIN tbl_tahunpelajaran ON tbl_tahunpelajaran.id_tahunpelajaran = kelulusan.id_tahunpelajaran
		LEFT JOIN setting ON 1=1
		WHERE kelulusan.nis = ?`
	if err := database.DB.Get(&c, query, nis); err != nil {
		return nil, err
	}
	return &c, nil
}

// --- galeri foto (fotokelulusan) ---

func ListFoto() ([]FotoKelulusan, error) {
	var rows []FotoKelulusan
	err := database.DB.Select(&rows, `SELECT * FROM fotokelulusan ORDER BY gambar ASC`)
	return rows, err
}

func InsertBatchFoto(names []string) error {
	if len(names) == 0 {
		return nil
	}
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	for _, name := range names {
		if _, err := tx.Exec(`INSERT INTO fotokelulusan (gambar) VALUES (?)`, name); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func GetFotoByIDs(ids []int) ([]FotoKelulusan, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query, args, err := sqlx.In(`SELECT * FROM fotokelulusan WHERE id IN (?)`, ids)
	if err != nil {
		return nil, err
	}
	query = database.DB.Rebind(query)
	var rows []FotoKelulusan
	err = database.DB.Select(&rows, query, args...)
	return rows, err
}

func DeleteFotoRows(ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	query, args, err := sqlx.In(`DELETE FROM fotokelulusan WHERE id IN (?)`, ids)
	if err != nil {
		return err
	}
	query = database.DB.Rebind(query)
	_, err = database.DB.Exec(query, args...)
	return err
}

// --- referensi (dipakai TemplateHandler utk isi sheet Petunjuk) ---

// RefItem - pasangan id+nama generik utk daftar referensi di sheet Petunjuk.
type RefItem struct {
	ID   int    `db:"id"`
	Nama string `db:"nama"`
}

func GetActiveTahunPelajaranRef() ([]RefItem, error) {
	var rows []RefItem
	err := database.DB.Select(&rows, `SELECT id_tahunpelajaran AS id, nama_tahunpelajaran AS nama
		FROM tbl_tahunpelajaran WHERE aktif = 1 ORDER BY id_tahunpelajaran`)
	return rows, err
}

func GetActiveKelasRef() ([]RefItem, error) {
	var rows []RefItem
	err := database.DB.Select(&rows, `SELECT id_kelas AS id, nama_kelas AS nama
		FROM tbl_kelas WHERE aktif = 1 ORDER BY id_kelas`)
	return rows, err
}

func GetActiveRombelRef() ([]RefItem, error) {
	var rows []RefItem
	err := database.DB.Select(&rows, `SELECT id_rombel AS id, nama_rombel AS nama
		FROM tbl_rombel WHERE aktif = 1 ORDER BY id_rombel`)
	return rows, err
}
