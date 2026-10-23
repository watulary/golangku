package pesertadidik

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"sekolah/internal/database"

	"github.com/jmoiron/sqlx"
)

// pesertaDidikColumns urutan kolom asli tabel pesertadidik (dipakai untuk build query insert/update
// dan untuk INSERT ... SELECT saat mutasi/arsip)
var pesertaDidikColumns = []string{
	"id_tahunpelajaran", "id_kelas", "id_rombel", "nama_siswa", "nis", "nisn",
	"no_ujiansd", "no_simdik", "jenis_kelamin", "tempat_lahirsiswa", "tanggal_lahirsiswa",
	"no_akte", "akte_kota", "agama_siswa", "warga_siswa", "anak_ke", "saudara_kandung",
	"saudara_tiri", "saudara_angkat", "yatim_piatu", "bahasa_siswa", "bantuan", "no_kk",
	"nik_siswa", "alamat_siswa", "rt_siswa", "rw_siswa", "kopos_siswa", "kel_siswa",
	"kec_siswa", "telprmh_siswa", "hp_siswa", "tinggal_siswa", "jarak_sekolah", "gol_darah",
	"pernah_sakit", "status_jasmani", "kebutuhan_khusus", "tinggi_badan", "berat_badan",
	"lingkar_kepala", "moda_transportasi", "lulus_sd", "tanggal_ijazah", "kota_ijazah",
	"nomor_ijazah", "tahun_ijazah", "nomor_skhun", "tahun_skhun", "nilai_ijazah", "nilai_nem",
	"lama_belajar", "jalur", "semester", "asal_smp", "tanggal_masuk", "email",
	"nama_ayah", "tempat_lahirayah", "tanggal_lahirayah", "agama_ayah", "warga_ayah",
	"pend_ayah", "pek_ayah", "gaji_ayah", "gaji_ayahktg", "nik_ayah", "alamat_ayah",
	"rt_ayah", "rw_ayah", "kopos_ayah", "kel_ayah", "kec_ayah", "telp_ayah", "status_ayah",
	"nama_ibu", "tempat_lahiribu", "tanggal_lahiribu", "agama_ibu", "warga_ibu",
	"pend_ibu", "pek_ibu", "gaji_ibu", "gaji_ibuktg", "nik_ibu", "alamat_ibu",
	"rt_ibu", "rw_ibu", "kopos_ibu", "kel_ibu", "kec_ibu", "telp_ibu", "status_ibu",
	"nama_wali", "tempat_lahirwali", "tanggal_lahirwali", "agama_wali", "warga_wali",
	"pend_wali", "pek_wali", "gaji_wali", "gaji_waliktg", "alamat_wali", "rt_wali",
	"rw_wali", "kopos_wali", "kel_wali", "kec_wali", "telp_wali",
	"hobi_seni", "hobi_orkes", "hobi_organisasi", "hobi_lain",
	"status_keluar", "alasan_keluar", "tanggal_keluar", "nomor_surat", "tujuan", "foto",
}

var notNullColumns = map[string]bool{
	"id_tahunpelajaran": true,
	"id_kelas":          true,
	"id_rombel":         true,
	"nis":               true,
}

// nullSafeColumns membungkus kolom nullable dengan COALESCE(kolom, ”) supaya aman di-scan ke string Go
func nullSafeColumns(alias string) []string {
	cols := make([]string, len(pesertaDidikColumns))
	for i, c := range pesertaDidikColumns {
		if notNullColumns[c] {
			cols[i] = alias + "." + c
		} else {
			cols[i] = fmt.Sprintf("COALESCE(%s.%s, '') AS %s", alias, c, c)
		}
	}
	return cols
}

var allowedSortColumns = map[string]string{
	"nis":                 "p.nis",
	"nama_siswa":          "p.nama_siswa",
	"nisn":                "p.nisn",
	"jenis_kelamin":       "p.jenis_kelamin",
	"agama_siswa":         "p.agama_siswa",
	"tempat_lahirsiswa":   "p.tempat_lahirsiswa",
	"tanggal_lahirsiswa":  "p.tanggal_lahirsiswa",
	"nama_kelas":          "k.nama_kelas",
	"nama_rombel":         "r.nama_rombel",
	"nama_tahunpelajaran": "t.nama_tahunpelajaran",
	"status_keluar":       "p.status_keluar",
}

const joinClause = `
	FROM pesertadidik p
	LEFT JOIN tbl_kelas k ON p.id_kelas = k.id_kelas
	LEFT JOIN tbl_rombel r ON p.id_rombel = r.id_rombel
	LEFT JOIN tbl_tahunpelajaran t ON p.id_tahunpelajaran = t.id_tahunpelajaran
`

var duplicateCheckableFields = map[string]bool{
	"nis":        true,
	"nik_siswa":  true,
	"nik_ayah":   true,
	"nik_ibu":    true,
	"nisn":       true,
	"no_ujiansd": true,
}

// CheckDuplicate mengecek apakah value sudah dipakai siswa lain (kecuali excludeNis, untuk mode edit)
func CheckDuplicate(field, value, excludeNis string) (bool, string, string, error) {
	if !duplicateCheckableFields[field] {
		return false, "", "", errors.New("field tidak valid untuk pengecekan duplikat")
	}
	if value == "" {
		return false, "", "", nil
	}

	query := fmt.Sprintf(`SELECT nis, nama_siswa FROM pesertadidik WHERE %s = ? AND nis != ? LIMIT 1`, field)
	var nis, nama string
	err := database.DB.QueryRow(query, value, excludeNis).Scan(&nis, &nama)
	if err == sql.ErrNoRows {
		return false, "", "", nil
	}
	if err != nil {
		return false, "", "", err
	}
	return true, nis, nama, nil
}

// GetAllAdmin mengembalikan daftar peserta didik aktif untuk halaman admin (pagination+search+sort+filter)
func GetAllAdmin(search string, idTahunPelajaran, idKelas, idRombel int, sortBy, sortDir string, limit, offset int) ([]PesertaDidikListItem, error) {
	query := `
		SELECT p.nis, COALESCE(p.nama_siswa,'') AS nama_siswa, COALESCE(p.nisn,'') AS nisn,
			COALESCE(p.jenis_kelamin,'') AS jenis_kelamin,
			COALESCE(p.agama_siswa,'') AS agama_siswa,
			COALESCE(p.tempat_lahirsiswa,'') AS tempat_lahirsiswa,
			COALESCE(p.tanggal_lahirsiswa,'') AS tanggal_lahirsiswa,
			p.id_tahunpelajaran, p.id_kelas, p.id_rombel,
			COALESCE(t.nama_tahunpelajaran,'') AS nama_tahunpelajaran,
			COALESCE(k.nama_kelas,'') AS nama_kelas,
			COALESCE(r.nama_rombel,'') AS nama_rombel,
			COALESCE(p.status_keluar,'') AS status_keluar,
			COALESCE(p.foto,'') AS foto
	` + joinClause + ` WHERE 1=1`

	args := []interface{}{}

	if search != "" {
		query += ` AND (
		p.nis LIKE ? OR p.nama_siswa LIKE ? OR p.nisn LIKE ?
		OR p.agama_siswa LIKE ? OR p.tempat_lahirsiswa LIKE ?
		OR p.tanggal_lahirsiswa LIKE ?
		OR DATE_FORMAT(STR_TO_DATE(p.tanggal_lahirsiswa, '%Y-%m-%d'), '%d-%m-%Y') LIKE ?
	)`
		s := "%" + search + "%"
		args = append(args, s, s, s, s, s, s, s)
	}
	if idTahunPelajaran > 0 {
		query += ` AND p.id_tahunpelajaran = ?`
		args = append(args, idTahunPelajaran)
	}
	if idKelas > 0 {
		query += ` AND p.id_kelas = ?`
		args = append(args, idKelas)
	}
	if idRombel > 0 {
		query += ` AND p.id_rombel = ?`
		args = append(args, idRombel)
	}

	dir := "ASC"
	if strings.ToUpper(sortDir) == "DESC" {
		dir = "DESC"
	}

	// Default (sortBy kosong/nama_siswa): urut kelas -> rombel -> nama, biar rapi per rombongan belajar.
	// Kalau user klik sort kolom lain secara eksplisit, pakai kolom itu saja.
	if sortBy == "" || sortBy == "nama_siswa" {
		query += fmt.Sprintf(
			" ORDER BY k.nama_kelas %s, r.nama_rombel %s, p.nama_siswa %s",
			dir, dir, dir,
		)
	} else {
		col, ok := allowedSortColumns[sortBy]
		if !ok {
			col = "p.nama_siswa"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", col, dir)
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var data []PesertaDidikListItem
	err := database.DB.Select(&data, query, args...)
	return data, err
}

// CountAdmin menghitung total baris untuk filter yang sama dengan GetAllAdmin (tanpa limit/offset)
func CountAdmin(search string, idTahunPelajaran, idKelas, idRombel int) (int, error) {
	query := `SELECT COUNT(*) ` + joinClause + ` WHERE 1=1`
	args := []interface{}{}

	if search != "" {
		query += ` AND (
		p.nis LIKE ? OR p.nama_siswa LIKE ? OR p.nisn LIKE ?
		OR p.agama_siswa LIKE ? OR p.tempat_lahirsiswa LIKE ?
		OR p.tanggal_lahirsiswa LIKE ?
		OR DATE_FORMAT(STR_TO_DATE(p.tanggal_lahirsiswa, '%Y-%m-%d'), '%d-%m-%Y') LIKE ?
	)`
		s := "%" + search + "%"
		args = append(args, s, s, s, s, s, s, s)
	}
	if idTahunPelajaran > 0 {
		query += ` AND p.id_tahunpelajaran = ?`
		args = append(args, idTahunPelajaran)
	}
	if idKelas > 0 {
		query += ` AND p.id_kelas = ?`
		args = append(args, idKelas)
	}
	if idRombel > 0 {
		query += ` AND p.id_rombel = ?`
		args = append(args, idRombel)
	}

	var total int
	err := database.DB.Get(&total, query, args...)
	return total, err
}

// GetByID mengembalikan data lengkap satu peserta didik aktif (dengan nama kelas/rombel/tahun pelajaran)
func GetByID(nis string) (*PesertaDidik, error) {
	query := fmt.Sprintf(`
		SELECT %s, t.nama_tahunpelajaran, k.nama_kelas, r.nama_rombel
		FROM pesertadidik p
		LEFT JOIN tbl_kelas k ON p.id_kelas = k.id_kelas
		LEFT JOIN tbl_rombel r ON p.id_rombel = r.id_rombel
		LEFT JOIN tbl_tahunpelajaran t ON p.id_tahunpelajaran = t.id_tahunpelajaran
		WHERE p.nis = ?
	`, strings.Join(nullSafeColumns("p"), ", "))

	var data PesertaDidik
	err := database.DB.Get(&data, query, nis)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// Create menambahkan peserta didik baru ke tabel aktif
func Create(data *PesertaDidik) error {
	cols := strings.Join(pesertaDidikColumns, ", ")
	placeholders := make([]string, len(pesertaDidikColumns))
	for i, c := range pesertaDidikColumns {
		placeholders[i] = ":" + c
	}
	query := fmt.Sprintf("INSERT INTO pesertadidik (%s) VALUES (%s)", cols, strings.Join(placeholders, ", "))

	_, err := database.DB.NamedExec(query, data)
	return err
}

// Update memperbarui data peserta didik aktif berdasarkan nis (PK tidak ikut di-update)
func Update(nis string, data *PesertaDidik) error {
	sets := make([]string, 0, len(pesertaDidikColumns)-1)
	for _, c := range pesertaDidikColumns {
		if c == "nis" {
			continue
		}
		sets = append(sets, fmt.Sprintf("%s = :%s", c, c))
	}
	query := fmt.Sprintf("UPDATE pesertadidik SET %s WHERE nis = :nis", strings.Join(sets, ", "))

	data.NIS = nis
	_, err := database.DB.NamedExec(query, data)
	return err
}

// UpdateFoto memperbarui hanya kolom foto
func UpdateFoto(nis, filename string) error {
	_, err := database.DB.Exec(`UPDATE pesertadidik SET foto = ? WHERE nis = ?`, filename, nis)
	return err
}

// Delete menghapus permanen data peserta didik aktif (koreksi data oleh admin, bukan alur mutasi/arsip)
func Delete(nis string) error {
	_, err := database.DB.Exec(`DELETE FROM pesertadidik WHERE nis = ?`, nis)
	return err
}

// MutasiSiswa memindahkan siswa aktif ke pesertadidik_mutasi (transaction: update status_keluar dkk → insert → delete)
func MutasiSiswa(nis, statusKeluar, alasanKeluar, tanggalKeluar, tujuan, nomorSurat string) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE pesertadidik
		SET status_keluar = ?, alasan_keluar = ?, tanggal_keluar = ?, tujuan = ?, nomor_surat = ?
		WHERE nis = ?
	`, statusKeluar, alasanKeluar, tanggalKeluar, tujuan, nomorSurat, nis)
	if err != nil {
		return err
	}

	cols := strings.Join(pesertaDidikColumns, ", ")
	insertQuery := fmt.Sprintf(`INSERT INTO pesertadidik_mutasi (%s) SELECT %s FROM pesertadidik WHERE nis = ?`, cols, cols)
	res, err := tx.Exec(insertQuery, nis)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("data peserta didik tidak ditemukan")
	}

	if _, err = tx.Exec(`DELETE FROM pesertadidik WHERE nis = ?`, nis); err != nil {
		return err
	}

	return tx.Commit()
}

// ArsipkanSiswa menyimpan snapshot siswa ke pesertadidik_arsip TANPA menghapusnya
// dari tabel aktif — siswa tetap tampil di Tampil.vue supaya admin tinggal ubah
// tahun pelajaran/kelas/rombel untuk proses naik kelas. Dibatalkan (error) kalau
// siswa ini sudah pernah diarsipkan untuk tahun pelajaran yang sama (cegah duplikat).
// Setelah snapshot DB berhasil, foto aktif ikut disalin fisik ke storage/fotosiswaarsip
// (lihat copyFotoToArsip) supaya foto arsip independen dari data aktif.
func ArsipkanSiswa(nis string) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var aktif struct {
		IDTahunPelajaran int    `db:"id_tahunpelajaran"`
		Foto             string `db:"foto"`
	}
	if err := tx.Get(&aktif, `SELECT id_tahunpelajaran, COALESCE(foto,'') AS foto FROM pesertadidik WHERE nis = ?`, nis); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data peserta didik tidak ditemukan")
		}
		return err
	}

	var count int
	if err := tx.Get(&count, `SELECT COUNT(*) FROM pesertadidik_arsip WHERE nis = ? AND id_tahunpelajaran = ?`, nis, aktif.IDTahunPelajaran); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Sudah pernah diarsipkan.")
	}

	cols := strings.Join(pesertaDidikColumns, ", ")
	insertQuery := fmt.Sprintf(`INSERT INTO pesertadidik_arsip (%s) SELECT %s FROM pesertadidik WHERE nis = ?`, cols, cols)
	res, err := tx.Exec(insertQuery, nis)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("data peserta didik tidak ditemukan")
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if aktif.Foto != "" {
		go copyFotoToArsip(aktif.Foto)
	}

	return nil
}

// ==================== MUTASI (read-only list) ====================

func GetAllMutasi(search string, limit, offset int) ([]PesertaDidikListItem, error) {
	query := `
	SELECT p.nis, COALESCE(p.nama_siswa,'') AS nama_siswa, COALESCE(p.nisn,'') AS nisn,
		COALESCE(p.jenis_kelamin,'') AS jenis_kelamin,
		p.id_tahunpelajaran, p.id_kelas, p.id_rombel,
		COALESCE(t.nama_tahunpelajaran,'') AS nama_tahunpelajaran,
		COALESCE(k.nama_kelas,'') AS nama_kelas,
		COALESCE(r.nama_rombel,'') AS nama_rombel,
		COALESCE(p.status_keluar,'') AS status_keluar,
		COALESCE(p.tanggal_keluar,'') AS tanggal_keluar,
		COALESCE(p.alasan_keluar,'') AS alasan_keluar,
		COALESCE(p.tujuan,'') AS tujuan,
		COALESCE(p.nomor_surat,'') AS nomor_surat,
		COALESCE(p.foto,'') AS foto
	FROM pesertadidik_mutasi p
	LEFT JOIN tbl_kelas k ON p.id_kelas = k.id_kelas
	LEFT JOIN tbl_rombel r ON p.id_rombel = r.id_rombel
	LEFT JOIN tbl_tahunpelajaran t ON p.id_tahunpelajaran = t.id_tahunpelajaran
	WHERE 1=1
`
	args := []interface{}{}
	if search != "" {
		query += ` AND (
		p.nis LIKE ? OR p.nama_siswa LIKE ? OR r.nama_rombel LIKE ?
		OR p.tanggal_keluar LIKE ? OR p.nomor_surat LIKE ?
		OR p.tujuan LIKE ? OR p.alasan_keluar LIKE ?
		OR DATE_FORMAT(STR_TO_DATE(p.tanggal_keluar, '%Y-%m-%d'), '%d-%m-%Y') LIKE ?
	)`
		s := "%" + search + "%"
		args = append(args, s, s, s, s, s, s, s, s)
	}
	query += ` ORDER BY p.tanggal_keluar DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	var data []PesertaDidikListItem
	err := database.DB.Select(&data, query, args...)
	return data, err
}

func CountMutasi(search string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM pesertadidik_mutasi p
		LEFT JOIN tbl_rombel r ON p.id_rombel = r.id_rombel
		WHERE 1=1
	`
	args := []interface{}{}
	if search != "" {
		query += ` AND (
			p.nis LIKE ? OR p.nama_siswa LIKE ? OR r.nama_rombel LIKE ?
			OR p.tanggal_keluar LIKE ? OR p.nomor_surat LIKE ?
			OR p.tujuan LIKE ? OR p.alasan_keluar LIKE ?
			OR DATE_FORMAT(STR_TO_DATE(p.tanggal_keluar, '%Y-%m-%d'), '%d-%m-%Y') LIKE ?
		)`
		s := "%" + search + "%"
		args = append(args, s, s, s, s, s, s, s, s)
	}
	var total int
	err := database.DB.Get(&total, query, args...)
	return total, err
}

// ==================== ARSIP (read-only list, filter wajib per tahun pelajaran) ====================

// prefixedColumns membangun daftar kolom dengan alias tabel, dipakai untuk SELECT p.col1, p.col2, ...
func prefixedColumns(alias string) string {
	cols := make([]string, len(pesertaDidikColumns))
	for i, c := range pesertaDidikColumns {
		cols[i] = alias + "." + c
	}
	return strings.Join(cols, ", ")
}

// GetAllArsip diurutkan per kelas -> rombel -> nama, supaya tampil terkelompok
func GetAllArsip(search string, idTahunPelajaran, idKelas, idRombel int, limit, offset int) ([]PesertaDidikListItem, error) {
	query := `
	SELECT p.kode_siswa AS id, p.nis, COALESCE(p.nama_siswa,'') AS nama_siswa, COALESCE(p.nisn,'') AS nisn,
		COALESCE(p.jenis_kelamin,'') AS jenis_kelamin,
		p.id_tahunpelajaran, p.id_kelas, p.id_rombel,
		COALESCE(t.nama_tahunpelajaran,'') AS nama_tahunpelajaran,
		COALESCE(k.nama_kelas,'') AS nama_kelas,
		COALESCE(r.nama_rombel,'') AS nama_rombel,
		COALESCE(p.status_keluar,'') AS status_keluar,
		COALESCE(p.foto,'') AS foto
	FROM pesertadidik_arsip p
	LEFT JOIN tbl_kelas k ON p.id_kelas = k.id_kelas
	LEFT JOIN tbl_rombel r ON p.id_rombel = r.id_rombel
	LEFT JOIN tbl_tahunpelajaran t ON p.id_tahunpelajaran = t.id_tahunpelajaran
	WHERE 1=1
`
	args := []interface{}{}
	if idTahunPelajaran > 0 {
		query += ` AND p.id_tahunpelajaran = ?`
		args = append(args, idTahunPelajaran)
	}
	if idKelas > 0 {
		query += ` AND p.id_kelas = ?`
		args = append(args, idKelas)
	}
	if idRombel > 0 {
		query += ` AND p.id_rombel = ?`
		args = append(args, idRombel)
	}
	if search != "" {
		query += ` AND (p.nis LIKE ? OR p.nama_siswa LIKE ?)`
		s := "%" + search + "%"
		args = append(args, s, s)
	}
	// diurutkan per kelas -> rombel -> nama siswa, supaya tampil terkelompok di frontend
	query += ` ORDER BY k.nama_kelas ASC, r.nama_rombel ASC, p.nama_siswa ASC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	var data []PesertaDidikListItem
	err := database.DB.Select(&data, query, args...)
	return data, err
}

// GetArsipByID mengambil detail lengkap satu baris arsip berdasarkan kode_siswa (bukan nis, karena bisa dobel)
func GetArsipByID(id int) (*PesertaDidik, error) {
	query := fmt.Sprintf(`
		SELECT %s, t.nama_tahunpelajaran, k.nama_kelas, r.nama_rombel
		FROM pesertadidik_arsip p
		LEFT JOIN tbl_kelas k ON p.id_kelas = k.id_kelas
		LEFT JOIN tbl_rombel r ON p.id_rombel = r.id_rombel
		LEFT JOIN tbl_tahunpelajaran t ON p.id_tahunpelajaran = t.id_tahunpelajaran
		WHERE p.kode_siswa = ?
	`, strings.Join(nullSafeColumns("p"), ", "))

	var data PesertaDidik
	err := database.DB.Get(&data, query, id)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func CountArsip(search string, idTahunPelajaran, idKelas, idRombel int) (int, error) {
	query := `SELECT COUNT(*) FROM pesertadidik_arsip WHERE 1=1`
	args := []interface{}{}
	if idTahunPelajaran > 0 {
		query += ` AND id_tahunpelajaran = ?`
		args = append(args, idTahunPelajaran)
	}
	if idKelas > 0 {
		query += ` AND id_kelas = ?`
		args = append(args, idKelas)
	}
	if idRombel > 0 {
		query += ` AND id_rombel = ?`
		args = append(args, idRombel)
	}
	if search != "" {
		query += ` AND (nis LIKE ? OR nama_siswa LIKE ?)`
		s := "%" + search + "%"
		args = append(args, s, s)
	}
	var total int
	err := database.DB.Get(&total, query, args...)
	return total, err
}

// DeleteFoto mengosongkan kolom foto (default ditentukan saat tampil, berdasarkan jenis_kelamin)
func DeleteFoto(nis string) error {
	_, err := database.DB.Exec(`UPDATE pesertadidik SET foto = '' WHERE nis = ?`, nis)
	return err
}

// FotoMasihDipakai mengecek apakah suatu nama file foto masih direferensikan oleh
// snapshot pesertadidik_mutasi atau pesertadidik_arsip. Dibutuhkan karena Mutasi/Arsip
// menyalin NAMA FILE foto dari data aktif (bukan file fisiknya) saat snapshot dibuat,
// dan Arsipkan TIDAK menghapus baris aktif — jadi kalau foto di baris aktif yang sama
// diganti/dihapus belakangan, file fisik lama tidak boleh langsung dihapus selama masih
// dipakai snapshot mutasi/arsip, supaya foto di halaman Mutasi/Arsip tidak ikut pecah (404).
func FotoMasihDipakai(filename string) (bool, error) {
	if filename == "" {
		return false, nil
	}
	var count int
	err := database.DB.Get(&count, `
		SELECT
			(SELECT COUNT(*) FROM pesertadidik_mutasi WHERE foto = ?) +
			(SELECT COUNT(*) FROM pesertadidik_arsip WHERE foto = ?)
	`, filename, filename)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAllForExport mengambil semua peserta didik aktif (tanpa paginasi) untuk export Excel,
// mengikuti filter yang sama dengan GetAllAdmin (search + tahun pelajaran + kelas + rombel).
func GetAllForExport(search string, idTahunPelajaran, idKelas, idRombel int) ([]PesertaDidik, error) {
	query := fmt.Sprintf(`
		SELECT %s, t.nama_tahunpelajaran, k.nama_kelas, r.nama_rombel
	`, strings.Join(nullSafeColumns("p"), ", ")) + joinClause + ` WHERE 1=1`

	args := []interface{}{}

	if search != "" {
		query += ` AND (
			p.nis LIKE ? OR p.nama_siswa LIKE ? OR p.nisn LIKE ?
			OR p.agama_siswa LIKE ? OR p.tempat_lahirsiswa LIKE ?
			OR p.tanggal_lahirsiswa LIKE ?
			OR DATE_FORMAT(STR_TO_DATE(p.tanggal_lahirsiswa, '%Y-%m-%d'), '%d-%m-%Y') LIKE ?
		)`
		s := "%" + search + "%"
		args = append(args, s, s, s, s, s, s, s)
	}
	if idTahunPelajaran > 0 {
		query += ` AND p.id_tahunpelajaran = ?`
		args = append(args, idTahunPelajaran)
	}
	if idKelas > 0 {
		query += ` AND p.id_kelas = ?`
		args = append(args, idKelas)
	}
	if idRombel > 0 {
		query += ` AND p.id_rombel = ?`
		args = append(args, idRombel)
	}

	query += ` ORDER BY k.nama_kelas ASC, r.nama_rombel ASC, p.nama_siswa ASC`

	var data []PesertaDidik
	err := database.DB.Select(&data, query, args...)
	return data, err
}

// GetAllMutasiForExport mengambil SEMUA data mutasi (tanpa filter/paginasi) untuk export Excel.
func GetAllMutasiForExport() ([]PesertaDidik, error) {
	query := fmt.Sprintf(`
		SELECT %s, t.nama_tahunpelajaran, k.nama_kelas, r.nama_rombel
		FROM pesertadidik_mutasi p
		LEFT JOIN tbl_kelas k ON p.id_kelas = k.id_kelas
		LEFT JOIN tbl_rombel r ON p.id_rombel = r.id_rombel
		LEFT JOIN tbl_tahunpelajaran t ON p.id_tahunpelajaran = t.id_tahunpelajaran
		ORDER BY p.tanggal_keluar DESC
	`, strings.Join(nullSafeColumns("p"), ", "))

	var data []PesertaDidik
	err := database.DB.Select(&data, query)
	return data, err
}

// GetAllArsipForExport mengambil data arsip (tanpa paginasi) untuk export Excel,
// mengikuti filter yang sama dengan GetAllArsip (search + tahun pelajaran + kelas + rombel).
func GetAllArsipForExport(search string, idTahunPelajaran, idKelas, idRombel int) ([]PesertaDidik, error) {
	query := fmt.Sprintf(`
		SELECT %s, t.nama_tahunpelajaran, k.nama_kelas, r.nama_rombel
		FROM pesertadidik_arsip p
		LEFT JOIN tbl_kelas k ON p.id_kelas = k.id_kelas
		LEFT JOIN tbl_rombel r ON p.id_rombel = r.id_rombel
		LEFT JOIN tbl_tahunpelajaran t ON p.id_tahunpelajaran = t.id_tahunpelajaran
		WHERE 1=1
	`, strings.Join(nullSafeColumns("p"), ", "))

	args := []interface{}{}

	if idTahunPelajaran > 0 {
		query += ` AND p.id_tahunpelajaran = ?`
		args = append(args, idTahunPelajaran)
	}
	if idKelas > 0 {
		query += ` AND p.id_kelas = ?`
		args = append(args, idKelas)
	}
	if idRombel > 0 {
		query += ` AND p.id_rombel = ?`
		args = append(args, idRombel)
	}
	if search != "" {
		query += ` AND (p.nis LIKE ? OR p.nama_siswa LIKE ?)`
		s := "%" + search + "%"
		args = append(args, s, s)
	}

	query += ` ORDER BY k.nama_kelas ASC, r.nama_rombel ASC, p.nama_siswa ASC`

	var data []PesertaDidik
	err := database.DB.Select(&data, query, args...)
	return data, err
}

// GetMutasiByID mengambil detail lengkap satu baris mutasi berdasarkan nis
func GetMutasiByID(nis string) (*PesertaDidik, error) {
	query := fmt.Sprintf(`
		SELECT %s, t.nama_tahunpelajaran, k.nama_kelas, r.nama_rombel
		FROM pesertadidik_mutasi p
		LEFT JOIN tbl_kelas k ON p.id_kelas = k.id_kelas
		LEFT JOIN tbl_rombel r ON p.id_rombel = r.id_rombel
		LEFT JOIN tbl_tahunpelajaran t ON p.id_tahunpelajaran = t.id_tahunpelajaran
		WHERE p.nis = ?
	`, strings.Join(nullSafeColumns("p"), ", "))

	var data PesertaDidik
	err := database.DB.Get(&data, query, nis)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// GetExistingConflicts mengecek apakah nis/nisn/no_ujiansd/nik_siswa/nik_ayah/nik_ibu
// yang mau diimport sudah ada di database. Slice kosong dilewati (tidak ikut query).
func GetExistingConflicts(nisList, nisnList, noUjianList, nikSiswaList, nikAyahList, nikIbuList []string) ([]PesertaDidikConflictCheck, error) {
	var conditions []string
	var args []interface{}

	addCond := func(col string, vals []string) {
		if len(vals) > 0 {
			conditions = append(conditions, fmt.Sprintf("(%s IN (?) AND %s != '')", col, col))
			args = append(args, vals)
		}
	}
	if len(nisList) > 0 {
		conditions = append(conditions, "nis IN (?)")
		args = append(args, nisList)
	}
	addCond("nisn", nisnList)
	addCond("no_ujiansd", noUjianList)
	addCond("nik_siswa", nikSiswaList)
	addCond("nik_ayah", nikAyahList)
	addCond("nik_ibu", nikIbuList)

	if len(conditions) == 0 {
		return nil, nil
	}

	query := `
		SELECT
			nis,
			COALESCE(nisn, '') AS nisn,
			COALESCE(no_ujiansd, '') AS no_ujiansd,
			COALESCE(nik_siswa, '') AS nik_siswa,
			COALESCE(nik_ayah, '') AS nik_ayah,
			COALESCE(nik_ibu, '') AS nik_ibu,
			COALESCE(nama_siswa, '') AS nama_siswa
		FROM pesertadidik
		WHERE ` + strings.Join(conditions, " OR ")

	query, args, err := sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	query = database.DB.Rebind(query)

	var result []PesertaDidikConflictCheck
	err = database.DB.Select(&result, query, args...)
	return result, err
}

// ImportPesertaDidik menyimpan seluruh baris dalam satu transaction (all-or-nothing)
func ImportPesertaDidik(data []PesertaDidik) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	cols := strings.Join(pesertaDidikColumns, ", ")
	placeholders := make([]string, len(pesertaDidikColumns))
	for i, c := range pesertaDidikColumns {
		placeholders[i] = ":" + c
	}
	query := fmt.Sprintf("INSERT INTO pesertadidik (%s) VALUES (%s)", cols, strings.Join(placeholders, ", "))

	for _, row := range data {
		if _, err := tx.NamedExec(query, row); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// BulkDelete menghapus banyak peserta didik aktif sekaligus
func BulkDelete(nisList []string) error {
	if len(nisList) == 0 {
		return nil
	}
	query, args, err := sqlx.In(`DELETE FROM pesertadidik WHERE nis IN (?)`, nisList)
	if err != nil {
		return err
	}
	query = database.DB.Rebind(query)
	_, err = database.DB.Exec(query, args...)
	return err
}

// BulkResult menampung ringkasan sukses/gagal per nis untuk operasi bulk mutasi/arsip
type BulkResult struct {
	Success []string          `json:"success"`
	Failed  map[string]string `json:"failed"` // nis -> pesan error
}

// BulkArsipkan mengarsipkan banyak siswa sekaligus, tiap nis diproses independen
// (satu gagal tidak menggagalkan yang lain, karena tiap nis punya transaction sendiri)
func BulkArsipkan(nisList []string) BulkResult {
	result := BulkResult{Failed: map[string]string{}}
	for _, nis := range nisList {
		if err := ArsipkanSiswa(nis); err != nil {
			result.Failed[nis] = err.Error()
			continue
		}
		result.Success = append(result.Success, nis)
	}
	return result
}

// BulkMutasi memutasi banyak siswa sekaligus dengan status/alasan/tanggal yang sama,
// tiap nis diproses independen (satu gagal tidak menggagalkan yang lain)
func BulkMutasi(nisList []string, statusKeluar, alasanKeluar, tanggalKeluar, tujuan, nomorSurat string) BulkResult {
	result := BulkResult{Failed: map[string]string{}}
	for _, nis := range nisList {
		if err := MutasiSiswa(nis, statusKeluar, alasanKeluar, tanggalKeluar, tujuan, nomorSurat); err != nil {
			result.Failed[nis] = err.Error()
			continue
		}
		result.Success = append(result.Success, nis)
	}
	return result
}
