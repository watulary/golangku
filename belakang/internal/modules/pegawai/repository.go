package pegawai

import (
	"database/sql"
	"fmt"
	"sekolah/internal/config"
	"sekolah/internal/database"
	"strings"

	"github.com/jmoiron/sqlx"
)

func fotoDefaultFor(jenisKelamin string) string {
	if jenisKelamin == "P" {
		return "P.jpg"
	}
	return "L.jpg"
}

func buildFotoURL(foto, jenisKelamin string) string {
	if foto == "" {
		foto = fotoDefaultFor(jenisKelamin)
	}
	return config.App.AppURL + "/storage/fotopegawai/" + foto
}

var allowedSortColumns = map[string]string{
	"nama":           "nama_pegawai",
	"nuptk":          "nuptk",
	"nip_baru":       "nip_baru",
	"jabatan":        "jabatan",
	"agama":          "agama",
	"tempat_lahir":   "tempat_lahir",
	"tanggal_lahir":  "tanggal_lahir",
	"tanggal_keluar": "tanggal_keluar",
}

// GetAllAdmin — statusAktif: 1 = aktif saja, 0 = keluar saja, selain itu = semua.
func GetAllAdmin(search, kategori, pnsNon string, statusAktif int, sortBy, sortDir string, limit, offset int) ([]PegawaiListItem, error) {
	var data []PegawaiListItem

	query := `
    SELECT
        nuptk,
        nama_pegawai,
		COALESCE(nip_baru, '') AS nip_baru,
        jabatan,
        COALESCE(tempat_lahir, '') AS tempat_lahir,
        COALESCE(tanggal_lahir, '') AS tanggal_lahir,
        COALESCE(pns_non, '') AS pns_non,
        COALESCE(jenis_kelamin, '') AS jenis_kelamin,
        COALESCE(agama, '') AS agama,
        COALESCE(telp_hp, '') AS telp_hp,
        COALESCE(email, '') AS email,
        COALESCE(foto, '') AS foto,
		COALESCE(keterangan_keluar, '') AS keterangan_keluar,
        COALESCE(tanggal_keluar, '') AS tanggal_keluar,
        COALESCE(nomor_surat, '') AS nomor_surat,
        status_keluar
    FROM pegawai
    WHERE (
        nama_pegawai LIKE ?
        OR nuptk LIKE ?
        OR jabatan LIKE ?
        OR agama LIKE ?
        OR tempat_lahir LIKE ?
        OR tanggal_lahir LIKE ?
    )`
	keyword := "%" + search + "%"
	args := []interface{}{keyword, keyword, keyword, keyword, keyword, keyword}

	switch kategori {
	case "Guru":
		query += ` AND jabatan LIKE '%Guru%'`
	case "Tatausaha":
		query += ` AND jabatan LIKE '%Tatausaha%'`
	}

	if pnsNon != "" {
		query += ` AND pns_non = ?`
		args = append(args, pnsNon)
	}

	switch statusAktif {
	case 1:
		query += ` AND (status_keluar = '' OR status_keluar IS NULL)`
	case 0:
		query += ` AND status_keluar != '' AND status_keluar IS NOT NULL`
	}

	column, ok := allowedSortColumns[sortBy]
	if !ok {
		column = "nama_pegawai"
	}

	direction := "ASC"
	if strings.ToLower(sortDir) == "desc" {
		direction = "DESC"
	}

	query += ` ORDER BY ` + column + ` ` + direction + ` LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	err := database.DB.Select(&data, query, args...)
	if err != nil {
		return nil, err
	}

	for i := range data {
		data[i].FotoURL = buildFotoURL(data[i].Foto, data[i].JenisKelamin)
	}

	return data, nil
}

func CountAdmin(search, kategori, pnsNon string, statusAktif int) (int, error) {
	var total int

	query := `
    SELECT COUNT(*)
    FROM pegawai
    WHERE (
        nama_pegawai LIKE ?
        OR nuptk LIKE ?
        OR jabatan LIKE ?
        OR agama LIKE ?
        OR tempat_lahir LIKE ?
        OR tanggal_lahir LIKE ?
    )`
	keyword := "%" + search + "%"
	args := []interface{}{keyword, keyword, keyword, keyword, keyword, keyword}

	switch kategori {
	case "Guru":
		query += ` AND jabatan LIKE '%Guru%'`
	case "Tatausaha":
		query += ` AND jabatan LIKE '%Tatausaha%'`
	}

	if pnsNon != "" {
		query += ` AND pns_non = ?`
		args = append(args, pnsNon)
	}

	switch statusAktif {
	case 1:
		query += ` AND (status_keluar = '' OR status_keluar IS NULL)`
	case 0:
		query += ` AND status_keluar != '' AND status_keluar IS NOT NULL`
	}

	err := database.DB.Get(&total, query, args...)

	return total, err
}

func GetByID(nuptk string) (*Pegawai, error) {
	var data Pegawai

	err := database.DB.Get(
		&data,
		`SELECT
			nuptk,
			nama_pegawai,
			jabatan,
			COALESCE(strata, '') AS strata,
			COALESCE(pns_non, '') AS pns_non,
			COALESCE(nikki, '') AS nikki,
			COALESCE(simpkbid, '') AS simpkbid,
			COALESCE(dapodikid, '') AS dapodikid,
			COALESCE(nip_lama, '') AS nip_lama,
			COALESCE(nip_baru, '') AS nip_baru,
			COALESCE(nrk, '') AS nrk,
			COALESCE(nrg, '') AS nrg,
			COALESCE(tempat_lahir, '') AS tempat_lahir,
			COALESCE(tanggal_lahir, '') AS tanggal_lahir,
			COALESCE(jenis_kelamin, '') AS jenis_kelamin,
			COALESCE(agama, '') AS agama,
			COALESCE(nik, '') AS nik,
			COALESCE(no_kk, '') AS no_kk,
			COALESCE(alamat, '') AS alamat,
			COALESCE(rt, '') AS rt,
			COALESCE(rw, '') AS rw,
			COALESCE(kel, '') AS kel,
			COALESCE(kec, '') AS kec,
			COALESCE(kopos, '') AS kopos,
			COALESCE(telp_rmh, '') AS telp_rmh,
			COALESCE(telp_hp, '') AS telp_hp,
			COALESCE(email, '') AS email,
			COALESCE(sk_sekolah, '') AS sk_sekolah,
			COALESCE(tmt_sekolah, '') AS tmt_sekolah,
			COALESCE(sk_nonkki, '') AS sk_nonkki,
			COALESCE(tmt_nonkki, '') AS tmt_nonkki,
			COALESCE(sk_kki, '') AS sk_kki,
			COALESCE(tmt_kki, '') AS tmt_kki,
			COALESCE(sk_cpns, '') AS sk_cpns,
			COALESCE(tmt_cpns, '') AS tmt_cpns,
			COALESCE(sk_pns, '') AS sk_pns,
			COALESCE(tmt_pns, '') AS tmt_pns,
			COALESCE(tmt_kgb, '') AS tmt_kgb,
			COALESCE(pangkat_kgb, '') AS pangkat_kgb,
			COALESCE(gaji_kgb, '') AS gaji_kgb,
			COALESCE(sk_kgb, '') AS sk_kgb,
			COALESCE(tglsk_kgb, '') AS tglsk_kgb,
			COALESCE(karpeg, '') AS karpeg,
			COALESCE(karis, '') AS karis,
			COALESCE(nama_ibukdg, '') AS nama_ibukdg,
			COALESCE(nm_suami_istri, '') AS nm_suami_istri,
			COALESCE(pek_suami_istri, '') AS pek_suami_istri,
			COALESCE(npwp, '') AS npwp,
			status_keluar,
			COALESCE(keterangan_keluar, '') AS keterangan_keluar,
			COALESCE(tanggal_keluar, '') AS tanggal_keluar,
			COALESCE(nomor_surat, '') AS nomor_surat,
			COALESCE(foto, '') AS foto
		FROM pegawai
		WHERE nuptk = ?
		LIMIT 1`,
		nuptk,
	)

	if err != nil {
		return nil, err
	}

	data.FotoURL = buildFotoURL(data.Foto, data.JenisKelamin)

	return &data, nil
}

func ExistsByNuptk(nuptk string) (bool, error) {
	var total int

	err := database.DB.Get(
		&total,
		`SELECT COUNT(*) FROM pegawai WHERE nuptk = ?`,
		nuptk,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

// SearchOption mencari pegawai aktif untuk dropdown autocomplete (search-as-you-type).
func SearchOption(search string, limit int) ([]PegawaiOption, error) {
	var data []PegawaiOption

	err := database.DB.Select(
		&data,
		`SELECT
			nuptk,
			nama_pegawai,
			jabatan
		FROM pegawai
		WHERE (
			nama_pegawai LIKE ?
			OR nuptk LIKE ?
		)
		AND (status_keluar = '' OR status_keluar IS NULL)
		ORDER BY nama_pegawai ASC
		LIMIT ?`,
		"%"+search+"%",
		"%"+search+"%",
		limit,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func Create(data Pegawai) error {
	_, err := database.DB.Exec(
		`INSERT INTO pegawai (
			nuptk, nama_pegawai, jabatan, strata, pns_non, nikki, simpkbid, dapodikid,
			nip_lama, nip_baru, nrk, nrg, tempat_lahir, tanggal_lahir, jenis_kelamin, agama,
			nik, no_kk, alamat, rt, rw, kel, kec, kopos, telp_rmh, telp_hp, email,
			sk_sekolah, tmt_sekolah, sk_nonkki, tmt_nonkki, sk_kki, tmt_kki,
			sk_cpns, tmt_cpns, sk_pns, tmt_pns, tmt_kgb, pangkat_kgb, gaji_kgb, sk_kgb, tglsk_kgb,
			karpeg, karis, nama_ibukdg, nm_suami_istri, pek_suami_istri, npwp,
			status_keluar, keterangan_keluar, tanggal_keluar, nomor_surat, foto
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?
		)`,
		data.Nuptk, data.NamaPegawai, data.Jabatan, data.Strata, data.PnsNon, data.Nikki, data.Simpkbid, data.Dapodikid,
		data.NipLama, data.NipBaru, data.Nrk, data.Nrg, data.TempatLahir, data.TanggalLahir, data.JenisKelamin, data.Agama,
		data.Nik, data.NoKk, data.Alamat, data.Rt, data.Rw, data.Kel, data.Kec, data.Kopos, data.TelpRmh, data.TelpHp, data.Email,
		data.SkSekolah, data.TmtSekolah, data.SkNonkki, data.TmtNonkki, data.SkKki, data.TmtKki,
		data.SkCpns, data.TmtCpns, data.SkPns, data.TmtPns, data.TmtKgb, data.PangkatKgb, data.GajiKgb, data.SkKgb, data.TglskKgb,
		data.Karpeg, data.Karis, data.NamaIbukdg, data.NmSuamiIstri, data.PekSuamiIstri, data.Npwp,
		data.StatusKeluar, data.KetKeluar, data.TanggalKeluar, data.NomorSurat, data.Foto,
	)

	return err
}

func Update(nuptk string, data Pegawai) error {
	_, err := database.DB.Exec(
		`UPDATE pegawai SET
			nama_pegawai = ?, jabatan = ?, strata = ?, pns_non = ?, nikki = ?, simpkbid = ?, dapodikid = ?,
			nip_lama = ?, nip_baru = ?, nrk = ?, nrg = ?, tempat_lahir = ?, tanggal_lahir = ?, jenis_kelamin = ?, agama = ?,
			nik = ?, no_kk = ?, alamat = ?, rt = ?, rw = ?, kel = ?, kec = ?, kopos = ?, telp_rmh = ?, telp_hp = ?, email = ?,
			sk_sekolah = ?, tmt_sekolah = ?, sk_nonkki = ?, tmt_nonkki = ?, sk_kki = ?, tmt_kki = ?,
			sk_cpns = ?, tmt_cpns = ?, sk_pns = ?, tmt_pns = ?, tmt_kgb = ?, pangkat_kgb = ?, gaji_kgb = ?, sk_kgb = ?, tglsk_kgb = ?,
			karpeg = ?, karis = ?, nama_ibukdg = ?, nm_suami_istri = ?, pek_suami_istri = ?, npwp = ?,
			status_keluar = ?, keterangan_keluar = ?, tanggal_keluar = ?, nomor_surat = ?
		WHERE nuptk = ?`,
		data.NamaPegawai, data.Jabatan, data.Strata, data.PnsNon, data.Nikki, data.Simpkbid, data.Dapodikid,
		data.NipLama, data.NipBaru, data.Nrk, data.Nrg, data.TempatLahir, data.TanggalLahir, data.JenisKelamin, data.Agama,
		data.Nik, data.NoKk, data.Alamat, data.Rt, data.Rw, data.Kel, data.Kec, data.Kopos, data.TelpRmh, data.TelpHp, data.Email,
		data.SkSekolah, data.TmtSekolah, data.SkNonkki, data.TmtNonkki, data.SkKki, data.TmtKki,
		data.SkCpns, data.TmtCpns, data.SkPns, data.TmtPns, data.TmtKgb, data.PangkatKgb, data.GajiKgb, data.SkKgb, data.TglskKgb,
		data.Karpeg, data.Karis, data.NamaIbukdg, data.NmSuamiIstri, data.PekSuamiIstri, data.Npwp,
		data.StatusKeluar, data.KetKeluar, data.TanggalKeluar, data.NomorSurat,
		nuptk,
	)
	return err
}

func Delete(nuptk string) error {
	result, err := database.DB.Exec(
		`DELETE FROM pegawai WHERE nuptk = ?`,
		nuptk,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetAllForExport mengambil semua pegawai (tanpa paginasi) untuk export Excel.
// statusAktif: 1 = aktif saja, 0 = keluar saja, selain itu = semua.
func GetAllForExport(search, kategori, pnsNon string, statusAktif int) ([]Pegawai, error) {
	var data []Pegawai
	query := `
		SELECT
			nuptk,
			COALESCE(nama_pegawai, '')       AS nama_pegawai,
			COALESCE(jabatan, '')            AS jabatan,
			COALESCE(strata, '')             AS strata,
			COALESCE(pns_non, '')            AS pns_non,
			COALESCE(nikki, '')              AS nikki,
			COALESCE(simpkbid, '')           AS simpkbid,
			COALESCE(dapodikid, '')          AS dapodikid,
			COALESCE(nip_lama, '')           AS nip_lama,
			COALESCE(nip_baru, '')           AS nip_baru,
			COALESCE(nrk, '')                AS nrk,
			COALESCE(nrg, '')                AS nrg,
			COALESCE(tempat_lahir, '')       AS tempat_lahir,
			COALESCE(tanggal_lahir, '')      AS tanggal_lahir,
			COALESCE(jenis_kelamin, '')      AS jenis_kelamin,
			COALESCE(agama, '')              AS agama,
			COALESCE(nik, '')                AS nik,
			COALESCE(no_kk, '')              AS no_kk,
			COALESCE(alamat, '')             AS alamat,
			COALESCE(rt, '')                 AS rt,
			COALESCE(rw, '')                 AS rw,
			COALESCE(kel, '')                AS kel,
			COALESCE(kec, '')                AS kec,
			COALESCE(kopos, '')              AS kopos,
			COALESCE(telp_rmh, '')           AS telp_rmh,
			COALESCE(telp_hp, '')            AS telp_hp,
			COALESCE(email, '')              AS email,
			COALESCE(sk_sekolah, '')         AS sk_sekolah,
			COALESCE(tmt_sekolah, '')        AS tmt_sekolah,
			COALESCE(sk_nonkki, '')          AS sk_nonkki,
			COALESCE(tmt_nonkki, '')         AS tmt_nonkki,
			COALESCE(sk_kki, '')             AS sk_kki,
			COALESCE(tmt_kki, '')            AS tmt_kki,
			COALESCE(sk_cpns, '')            AS sk_cpns,
			COALESCE(tmt_cpns, '')           AS tmt_cpns,
			COALESCE(sk_pns, '')             AS sk_pns,
			COALESCE(tmt_pns, '')            AS tmt_pns,
			COALESCE(tmt_kgb, '')            AS tmt_kgb,
			COALESCE(pangkat_kgb, '')        AS pangkat_kgb,
			COALESCE(gaji_kgb, '')           AS gaji_kgb,
			COALESCE(sk_kgb, '')             AS sk_kgb,
			COALESCE(tglsk_kgb, '')          AS tglsk_kgb,
			COALESCE(karpeg, '')             AS karpeg,
			COALESCE(karis, '')              AS karis,
			COALESCE(nama_ibukdg, '')        AS nama_ibukdg,
			COALESCE(nm_suami_istri, '')     AS nm_suami_istri,
			COALESCE(pek_suami_istri, '')    AS pek_suami_istri,
			COALESCE(npwp, '')               AS npwp,
			COALESCE(status_keluar, '')      AS status_keluar,
			COALESCE(keterangan_keluar, '')  AS keterangan_keluar,
			COALESCE(tanggal_keluar, '')     AS tanggal_keluar,
			COALESCE(nomor_surat, '')        AS nomor_surat,
			COALESCE(foto, '')               AS foto
		FROM pegawai
		WHERE (
    nama_pegawai LIKE ?
    OR nuptk LIKE ?
    OR jabatan LIKE ?
    OR agama LIKE ?
    OR tempat_lahir LIKE ?
    OR tanggal_lahir LIKE ?
)`
	keyword := "%" + search + "%"
	args := []interface{}{keyword, keyword, keyword, keyword, keyword, keyword}

	switch kategori {
	case "Guru":
		query += ` AND jabatan LIKE '%Guru%'`
	case "Tatausaha":
		query += ` AND jabatan LIKE '%Tatausaha%'`
	}

	if pnsNon != "" {
		query += ` AND pns_non = ?`
		args = append(args, pnsNon)
	}

	switch statusAktif {
	case 1:
		query += ` AND (status_keluar = '' OR status_keluar IS NULL)`
	case 0:
		query += ` AND status_keluar != '' AND status_keluar IS NOT NULL`
	}

	query += ` ORDER BY nama_pegawai ASC`

	err := database.DB.Select(&data, query, args...)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func UpdateFoto(nuptk, filename string) error {
	_, err := database.DB.Exec(
		`UPDATE pegawai SET foto = ? WHERE nuptk = ?`,
		filename, nuptk,
	)
	return err
}

// DuplicateResult berisi info pemilik data yang sama (untuk pesan "Ganda dengan -> ...").
type DuplicateResult struct {
	Nuptk       string `db:"nuptk" json:"nuptk"`
	NamaPegawai string `db:"nama_pegawai" json:"nama_pegawai"`
}

var allowedDuplicateFields = map[string]bool{
	"nuptk":    true,
	"nip_baru": true,
	"nik":      true,
}

// FindDuplicate mencari pegawai lain (selain excludeNuptk) yang punya nilai sama di kolom field.
// field WAJIB divalidasi lewat allowedDuplicateFields sebelum dipanggil (mencegah SQL injection).
func FindDuplicate(field, value, excludeNuptk string) (*DuplicateResult, error) {
	if !allowedDuplicateFields[field] {
		return nil, fmt.Errorf("field tidak diizinkan: %s", field)
	}

	var data DuplicateResult

	query := `SELECT nuptk, nama_pegawai FROM pegawai WHERE ` + field + ` = ? AND nuptk != ? LIMIT 1`

	err := database.DB.Get(&data, query, value, excludeNuptk)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GetExistingConflicts mengecek apakah nuptk/nip_baru/nik yang mau diimport
// sudah ada di database. Slice kosong akan dilewati (tidak ikut query).
func GetExistingConflicts(nuptks, nipBarus, niks []string) ([]PegawaiConflictCheck, error) {
	var conditions []string
	var args []interface{}

	if len(nuptks) > 0 {
		conditions = append(conditions, "nuptk IN (?)")
		args = append(args, nuptks)
	}
	if len(nipBarus) > 0 {
		conditions = append(conditions, "(nip_baru IN (?) AND nip_baru != '')")
		args = append(args, nipBarus)
	}
	if len(niks) > 0 {
		conditions = append(conditions, "(nik IN (?) AND nik != '')")
		args = append(args, niks)
	}
	if len(conditions) == 0 {
		return nil, nil
	}

	query := `
		SELECT
			nuptk,
			COALESCE(nip_baru, '') AS nip_baru,
			COALESCE(nik, '') AS nik,
			COALESCE(nama_pegawai, '') AS nama_pegawai
		FROM pegawai
		WHERE ` + strings.Join(conditions, " OR ")

	query, args, err := sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	query = database.DB.Rebind(query)

	var result []PegawaiConflictCheck
	err = database.DB.Select(&result, query, args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ImportPegawai menyimpan seluruh baris dalam satu transaction (all-or-nothing)
func ImportPegawai(data []Pegawai) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO pegawai (
		nuptk, nama_pegawai, jabatan, strata, pns_non, nikki, simpkbid, dapodikid,
		nip_lama, nip_baru, nrk, nrg, tempat_lahir, tanggal_lahir, jenis_kelamin, agama,
		nik, no_kk, alamat, rt, rw, kel, kec, kopos, telp_rmh, telp_hp, email,
		sk_sekolah, tmt_sekolah, sk_nonkki, tmt_nonkki, sk_kki, tmt_kki,
		sk_cpns, tmt_cpns, sk_pns, tmt_pns, tmt_kgb, pangkat_kgb, gaji_kgb, sk_kgb, tglsk_kgb,
		karpeg, karis, nama_ibukdg, nm_suami_istri, pek_suami_istri, npwp,
		status_keluar, keterangan_keluar, tanggal_keluar, nomor_surat, foto
	) VALUES (
		:nuptk, :nama_pegawai, :jabatan, :strata, :pns_non, :nikki, :simpkbid, :dapodikid,
		:nip_lama, :nip_baru, :nrk, :nrg, :tempat_lahir, :tanggal_lahir, :jenis_kelamin, :agama,
		:nik, :no_kk, :alamat, :rt, :rw, :kel, :kec, :kopos, :telp_rmh, :telp_hp, :email,
		:sk_sekolah, :tmt_sekolah, :sk_nonkki, :tmt_nonkki, :sk_kki, :tmt_kki,
		:sk_cpns, :tmt_cpns, :sk_pns, :tmt_pns, :tmt_kgb, :pangkat_kgb, :gaji_kgb, :sk_kgb, :tglsk_kgb,
		:karpeg, :karis, :nama_ibukdg, :nm_suami_istri, :pek_suami_istri, :npwp,
		:status_keluar, :keterangan_keluar, :tanggal_keluar, :nomor_surat, :foto
	)`

	for _, row := range data {
		if _, err := tx.NamedExec(query, row); err != nil {
			return err
		}
	}

	return tx.Commit()
}
