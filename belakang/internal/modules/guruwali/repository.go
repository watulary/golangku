package guruwali

import (
	"database/sql"
	"fmt"

	"sekolah/internal/database" // sesuaikan path import ini dengan project Anda

	"github.com/jmoiron/sqlx"
)

var allowedSortColumns = map[string]string{
	"id":      "g.id_guruwali",
	"nama":    "g.nama_guruwali",
	"nuptk":   "g.nuptk",
	"pegawai": "p.nama_pegawai",
	"anggota": "jumlah_anggota",
	"aktif":   "g.aktif",
}

var allowedSortColumnsAnggota = map[string]string{
	"nis":    "pd.nis",
	"siswa":  "pd.nama_siswa",
	"kelas":  "k.nama_kelas",
	"rombel": "r.nama_rombel",
}

func condIf(cond bool, s string) string {
	if cond {
		return s
	}
	return ""
}

// ================= GURU WALI =================

func GetAllAdmin(search, aktif, sortBy, sortDir string, limit, offset int) ([]GuruWaliListItem, error) {
	col, ok := allowedSortColumns[sortBy]
	if !ok {
		col = "g.id_guruwali"
	}
	if sortDir != "desc" {
		sortDir = "asc"
	}

	query := fmt.Sprintf(`
		SELECT g.id_guruwali, g.nama_guruwali, g.nuptk, g.aktif,
			COALESCE(p.nama_pegawai,'') AS nama_pegawai,
			(SELECT COUNT(*) FROM tbl_anggota_guruwali a WHERE a.id_guruwali = g.id_guruwali) AS jumlah_anggota
		FROM tbl_guruwali g
		LEFT JOIN pegawai p ON p.nuptk = g.nuptk
		WHERE (g.nama_guruwali LIKE ? OR g.nuptk LIKE ? OR p.nama_pegawai LIKE ?)
			%s
		ORDER BY %s %s, g.nama_guruwali ASC, g.id_guruwali ASC
		LIMIT ? OFFSET ?
	`, condIf(aktif != "", "AND g.aktif = ?"), col, sortDir)

	kw := "%" + search + "%"
	args := []interface{}{kw, kw, kw}
	if aktif != "" {
		args = append(args, aktif)
	}
	args = append(args, limit, offset)

	var items []GuruWaliListItem
	err := database.DB.Select(&items, query, args...)
	return items, err
}

func CountAdmin(search, aktif string) (int, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) FROM tbl_guruwali g
		LEFT JOIN pegawai p ON p.nuptk = g.nuptk
		WHERE (g.nama_guruwali LIKE ? OR g.nuptk LIKE ? OR p.nama_pegawai LIKE ?)
			%s
	`, condIf(aktif != "", "AND g.aktif = ?"))

	kw := "%" + search + "%"
	args := []interface{}{kw, kw, kw}
	if aktif != "" {
		args = append(args, aktif)
	}

	var count int
	err := database.DB.Get(&count, query, args...)
	return count, err
}

func GetByID(id int) (*GuruWali, error) {
	var g GuruWali
	err := database.DB.Get(&g, `SELECT id_guruwali, nama_guruwali, nuptk, aktif FROM tbl_guruwali WHERE id_guruwali = ?`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &g, nil
}

// GetGuruWaliBySiswa mengambil guru wali yang membawahi seorang siswa (lewat tbl_anggota_guruwali),
// dipakai endpoint /guruwali/me sisi Peserta Didik. Kalau siswa belum terdaftar di guru wali manapun,
// mengembalikan slice kosong (bukan error) — mirror GetEkskulBySiswa ekskul.
func GetGuruWaliBySiswa(nis string) ([]GuruWaliListItem, error) {
	data := []GuruWaliListItem{}
	err := database.DB.Select(&data, `
		SELECT g.id_guruwali, g.nama_guruwali, g.nuptk, g.aktif,
			COALESCE(p.nama_pegawai, '') AS nama_pegawai
		FROM tbl_anggota_guruwali a
		JOIN tbl_guruwali g ON g.id_guruwali = a.id_guruwali
		LEFT JOIN pegawai p ON p.nuptk = g.nuptk
		WHERE a.nis = ?
		ORDER BY g.nama_guruwali ASC
	`, nis)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ExistsByNama mengecek apakah nama_guruwali sudah dipakai guru wali lain.
// excludeID dipakai saat Update supaya baris yang sedang diedit tidak menghitung dirinya sendiri
// (isi 0 saat Create, karena id_guruwali auto_increment tidak pernah bernilai 0).
func ExistsByNama(nama string, excludeID int) (bool, error) {
	var total int
	err := database.DB.Get(&total, `
		SELECT COUNT(*) FROM tbl_guruwali
		WHERE nama_guruwali = ? AND id_guruwali != ?
	`, nama, excludeID)
	if err != nil {
		return false, err
	}
	return total > 0, nil
}

// ExistsByNUPTK mengecek apakah seorang guru (NUPTK) sudah menjabat sebagai guru wali lain
// (satu guru cuma boleh jadi guru wali satu kali). excludeID sama pola dgn ExistsByNama.
func ExistsByNUPTK(nuptk string, excludeID int) (bool, error) {
	var total int
	err := database.DB.Get(&total, `
		SELECT COUNT(*) FROM tbl_guruwali
		WHERE nuptk = ? AND id_guruwali != ?
	`, nuptk, excludeID)
	if err != nil {
		return false, err
	}
	return total > 0, nil
}

func Create(data GuruWali) (int, error) {
	res, err := database.DB.Exec(`INSERT INTO tbl_guruwali (nama_guruwali, nuptk, aktif) VALUES (?, ?, ?)`,
		data.NamaGuruWali, data.NUPTK, data.Aktif)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func Update(id int, data GuruWali) error {
	_, err := database.DB.Exec(`UPDATE tbl_guruwali SET nama_guruwali=?, nuptk=?, aktif=? WHERE id_guruwali=?`,
		data.NamaGuruWali, data.NUPTK, data.Aktif, id)
	return err
}

func Delete(id int) error {
	// anggota ikut terhapus otomatis kalau FK id_guruwali di tbl_anggota_guruwali diberi ON DELETE CASCADE.
	_, err := database.DB.Exec(`DELETE FROM tbl_guruwali WHERE id_guruwali=?`, id)
	return err
}

func BulkDelete(ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	query, args, err := sqlx.In(`DELETE FROM tbl_guruwali WHERE id_guruwali IN (?)`, ids)
	if err != nil {
		return err
	}
	query = database.DB.Rebind(query)
	_, err = database.DB.Exec(query, args...)
	return err
}

func GetAllForExport(search string) ([]GuruWaliListItem, error) {
	kw := "%" + search + "%"
	var items []GuruWaliListItem
	query := `
		SELECT g.id_guruwali, g.nama_guruwali, g.nuptk, g.aktif,
			COALESCE(p.nama_pegawai,'') AS nama_pegawai,
			(SELECT COUNT(*) FROM tbl_anggota_guruwali a WHERE a.id_guruwali = g.id_guruwali) AS jumlah_anggota
		FROM tbl_guruwali g
		LEFT JOIN pegawai p ON p.nuptk = g.nuptk
		WHERE g.nama_guruwali LIKE ? OR g.nuptk LIKE ? OR p.nama_pegawai LIKE ?
		ORDER BY g.nama_guruwali ASC, g.id_guruwali ASC
	`
	err := database.DB.Select(&items, query, kw, kw, kw)
	return items, err
}

// SearchPegawai dipakai autocomplete pemilihan NUPTK di Form.vue.
// Dibatasi hanya pegawai dengan jabatan "Guru" DAN masih aktif (status_keluar kosong) —
// jabatan "Tatausaha" dan pegawai yang sudah status_keluar = "Keluar" tidak ikut muncul.
func SearchPegawai(query string) ([]PegawaiOption, error) {
	kw := "%" + query + "%"
	var items []PegawaiOption
	err := database.DB.Select(&items, `
		SELECT nuptk, nama_pegawai
		FROM pegawai
		WHERE (nuptk LIKE ? OR nama_pegawai LIKE ?)
			AND jabatan = 'Guru'
			AND status_keluar = ''
		ORDER BY nama_pegawai ASC
		LIMIT 20
	`, kw, kw)
	return items, err
}

// GetNamaPegawai dipakai GetByIDHandler supaya Form.vue mode Ubah bisa langsung menampilkan
// nama guru yang sudah terpilih (bukan cuma NUPTK mentah)
func GetNamaPegawai(nuptk string) (string, error) {
	var nama string
	err := database.DB.Get(&nama, `SELECT nama_pegawai FROM pegawai WHERE nuptk = ?`, nuptk)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return nama, nil
}

// ================= ANGGOTA (mirip GetAnggotaByEkskul ekskul) =================

func GetAnggotaByGuruWali(idGuruWali int, search, sortBy, sortDir string, idKelas, idRombel, limit, offset int) ([]AnggotaGuruWaliListItem, error) {
	col, ok := allowedSortColumnsAnggota[sortBy]
	if !ok {
		col = "r.nama_rombel"
	}
	if sortDir != "desc" {
		sortDir = "asc"
	}

	query := fmt.Sprintf(`
		SELECT a.id, a.id_guruwali, a.nis,
			pd.nama_siswa, COALESCE(k.nama_kelas,'') AS nama_kelas, COALESCE(r.nama_rombel,'') AS nama_rombel
		FROM tbl_anggota_guruwali a
		JOIN pesertadidik pd ON pd.nis = a.nis
		LEFT JOIN tbl_rombel r ON r.id_rombel = pd.id_rombel
		LEFT JOIN tbl_kelas k ON k.id_kelas = r.id_kelas
		WHERE a.id_guruwali = ?
			AND (pd.nis LIKE ? OR pd.nama_siswa LIKE ?)
			%s
			%s
		ORDER BY %s %s, pd.nama_siswa ASC, a.id ASC
		LIMIT ? OFFSET ?
	`,
		condIf(idKelas > 0, "AND k.id_kelas = ?"),
		condIf(idRombel > 0, "AND r.id_rombel = ?"),
		col, sortDir)

	kw := "%" + search + "%"
	args := []interface{}{idGuruWali, kw, kw}
	if idKelas > 0 {
		args = append(args, idKelas)
	}
	if idRombel > 0 {
		args = append(args, idRombel)
	}
	args = append(args, limit, offset)

	var items []AnggotaGuruWaliListItem
	err := database.DB.Select(&items, query, args...)
	return items, err
}

func CountAnggotaByGuruWali(idGuruWali int, search string, idKelas, idRombel int) (int, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM tbl_anggota_guruwali a
		JOIN pesertadidik pd ON pd.nis = a.nis
		LEFT JOIN tbl_rombel r ON r.id_rombel = pd.id_rombel
		LEFT JOIN tbl_kelas k ON k.id_kelas = r.id_kelas
		WHERE a.id_guruwali = ?
			AND (pd.nis LIKE ? OR pd.nama_siswa LIKE ?)
			%s
			%s
	`, condIf(idKelas > 0, "AND k.id_kelas = ?"), condIf(idRombel > 0, "AND r.id_rombel = ?"))

	kw := "%" + search + "%"
	args := []interface{}{idGuruWali, kw, kw}
	if idKelas > 0 {
		args = append(args, idKelas)
	}
	if idRombel > 0 {
		args = append(args, idRombel)
	}

	var count int
	err := database.DB.Get(&count, query, args...)
	return count, err
}

func GetAnggotaByGuruWaliForExport(idGuruWali int, search string, idKelas, idRombel int) ([]AnggotaGuruWaliListItem, error) {
	query := fmt.Sprintf(`
		SELECT a.id, a.id_guruwali, a.nis,
			pd.nama_siswa, COALESCE(k.nama_kelas,'') AS nama_kelas, COALESCE(r.nama_rombel,'') AS nama_rombel
		FROM tbl_anggota_guruwali a
		JOIN pesertadidik pd ON pd.nis = a.nis
		LEFT JOIN tbl_rombel r ON r.id_rombel = pd.id_rombel
		LEFT JOIN tbl_kelas k ON k.id_kelas = r.id_kelas
		WHERE a.id_guruwali = ? AND (pd.nis LIKE ? OR pd.nama_siswa LIKE ?)
			%s
			%s
		ORDER BY r.nama_rombel ASC, pd.nama_siswa ASC, a.id ASC
	`, condIf(idKelas > 0, "AND k.id_kelas = ?"), condIf(idRombel > 0, "AND r.id_rombel = ?"))

	kw := "%" + search + "%"
	args := []interface{}{idGuruWali, kw, kw}
	if idKelas > 0 {
		args = append(args, idKelas)
	}
	if idRombel > 0 {
		args = append(args, idRombel)
	}

	var items []AnggotaGuruWaliListItem
	err := database.DB.Select(&items, query, args...)
	return items, err
}

func ExistsPesertaDidik(nis string) (bool, error) {
	var count int
	err := database.DB.Get(&count, `SELECT COUNT(*) FROM pesertadidik WHERE nis = ?`, nis)
	return count > 0, err
}

func IsAnggotaExists(idGuruWali int, nis string) (bool, error) {
	var count int
	err := database.DB.Get(&count, `SELECT COUNT(*) FROM tbl_anggota_guruwali WHERE id_guruwali = ? AND nis = ?`, idGuruWali, nis)
	return count > 0, err
}

func AddAnggota(idGuruWali int, nis string) (int, error) {
	res, err := database.DB.Exec(`INSERT INTO tbl_anggota_guruwali (id_guruwali, nis) VALUES (?, ?)`, idGuruWali, nis)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func RemoveAnggota(id int) error {
	_, err := database.DB.Exec(`DELETE FROM tbl_anggota_guruwali WHERE id = ?`, id)
	return err
}

func BulkRemoveAnggota(ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	query, args, err := sqlx.In(`DELETE FROM tbl_anggota_guruwali WHERE id IN (?)`, ids)
	if err != nil {
		return err
	}
	query = database.DB.Rebind(query)
	_, err = database.DB.Exec(query, args...)
	return err
}

// SearchPesertaDidik dipakai autocomplete "Tambah Anggota" di Anggota.vue (pencarian global,
// duplikat terhadap guru wali tertentu tetap divalidasi di AddAnggotaHandler)
func SearchPesertaDidik(query string) ([]PesertaDidikOption, error) {
	kw := "%" + query + "%"
	var items []PesertaDidikOption
	err := database.DB.Select(&items, `
		SELECT pd.nis, pd.nama_siswa,
			COALESCE(k.nama_kelas,'') AS nama_kelas, COALESCE(r.nama_rombel,'') AS nama_rombel
		FROM pesertadidik pd
		LEFT JOIN tbl_rombel r ON r.id_rombel = pd.id_rombel
		LEFT JOIN tbl_kelas k ON k.id_kelas = r.id_kelas
		WHERE pd.nis LIKE ? OR pd.nama_siswa LIKE ?
		ORDER BY pd.nama_siswa ASC
		LIMIT 20
	`, kw, kw)
	return items, err
}

// GetAllAnggotaForExport mengambil semua anggota dari semua guru wali (tanpa paginasi),
// dipakai ExportAnggotaHandler (sheet per guru wali) — mirror GetAllAnggotaForExport ekskul.
func GetAllAnggotaForExport() ([]AnggotaGuruWaliListItem, error) {
	var items []AnggotaGuruWaliListItem
	err := database.DB.Select(&items, `
		SELECT a.id, a.id_guruwali, a.nis,
			pd.nama_siswa, COALESCE(k.nama_kelas,'') AS nama_kelas, COALESCE(r.nama_rombel,'') AS nama_rombel
		FROM tbl_anggota_guruwali a
		JOIN pesertadidik pd ON pd.nis = a.nis
		LEFT JOIN tbl_rombel r ON r.id_rombel = pd.id_rombel
		LEFT JOIN tbl_kelas k ON k.id_kelas = r.id_kelas
		ORDER BY a.id_guruwali ASC, pd.nama_siswa ASC
	`)
	return items, err
}
