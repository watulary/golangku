package pengguna

import (
	"database/sql"
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"strings"

	"sekolah/internal/config"
	"sekolah/internal/database"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

// ========================================
// SearchInfo
// ========================================

type SearchInfo struct {
	Keyword string

	Level string

	Status string

	Password string

	LoginNull bool

	Day   int
	Month int
	Year  int
}

// ========================================
// Bulan Indonesia
// ========================================

var bulanIndonesia = map[string]int{
	"jan":       1,
	"januari":   1,
	"feb":       2,
	"februari":  2,
	"mar":       3,
	"maret":     3,
	"apr":       4,
	"april":     4,
	"mei":       5,
	"jun":       6,
	"juni":      6,
	"jul":       7,
	"juli":      7,
	"agu":       8,
	"agt":       8,
	"agustus":   8,
	"sep":       9,
	"sept":      9,
	"september": 9,
	"okt":       10,
	"oktober":   10,
	"nov":       11,
	"november":  11,
	"des":       12,
	"desember":  12,
}

// ========================================
// parseSearch
// ========================================

func parseSearch(search string) SearchInfo {

	info := SearchInfo{}

	search = strings.TrimSpace(search)

	info.Keyword = "%" + search + "%"

	lower := strings.ToLower(search)

	switch lower {

	// ==========================
	// Level
	// ==========================

	case "super", "super admin":

		info.Level = "Admin"

	case "admin":

		info.Level = "Staf"

	case "pegawai":

		info.Level = "Pegawai"

	case "peserta",
		"peserta didik":

		info.Level = "Pesertadidik"

	// ==========================
	// Status
	// ==========================

	case "aktif":

		info.Status = "1"

	case "nonaktif",
		"non aktif",
		"tidak aktif":

		info.Status = "0"

	// ==========================
	// Password
	// ==========================

	case "standar",
		"default",
		"belum diubah":

		info.Password = "0"

	case "diubah",
		"sudah diubah":

		info.Password = "1"

	// ==========================
	// Login
	// ==========================

	case "belum",
		"belum login",
		"belum masuk",
		"belum pernah",
		"belum pernah login",
		"belum pernah masuk":

		info.LoginNull = true
	}

	parseDateSearch(&info, lower)

	return info
}

// ========================================
// parseDateSearch
// ========================================

var (
	tanggalLengkapRegex = regexp.MustCompile(`^(\d{1,2})[-/](\d{1,2})[-/](\d{4})$`)

	tanggalBulanTahunRegex = regexp.MustCompile(`^(\d{1,2})[-/](\d{4})$`)
)

func parseDateSearch(
	info *SearchInfo,
	search string,
) {

	search = strings.TrimSpace(search)

	// ====================================
	// Contoh: 19-08-2026 / 19/08/2026
	// ====================================

	if m := tanggalLengkapRegex.FindStringSubmatch(search); m != nil {

		info.Day, _ = strconv.Atoi(m[1])
		info.Month, _ = strconv.Atoi(m[2])
		info.Year, _ = strconv.Atoi(m[3])

		return
	}

	// ====================================
	// Contoh: 08-2026 / 08/2026
	// ====================================

	if m := tanggalBulanTahunRegex.FindStringSubmatch(search); m != nil {

		info.Month, _ = strconv.Atoi(m[1])
		info.Year, _ = strconv.Atoi(m[2])

		return
	}

	fields := strings.Fields(search)

	switch len(fields) {

	// ====================================
	// Contoh:
	// 19
	// juli
	// 2026
	// ====================================

	case 1:

		token := fields[0]

		// angka
		if n, err := strconv.Atoi(token); err == nil {

			switch {

			case n >= 1 && n <= 31:
				info.Day = n

			case n >= 1900:
				info.Year = n
			}
		}

		// bulan
		if month, ok := bulanIndonesia[token]; ok {

			info.Month = month
		}

	// ====================================
	// Contoh:
	// 19 juli
	// ====================================

	case 2:

		if day, err := strconv.Atoi(fields[0]); err == nil {

			info.Day = day
		}

		if month, ok := bulanIndonesia[fields[1]]; ok {

			info.Month = month
		}

	// ====================================
	// Contoh:
	// 19 juli 2026
	// ====================================

	case 3:

		if day, err := strconv.Atoi(fields[0]); err == nil {

			info.Day = day
		}

		if month, ok := bulanIndonesia[fields[1]]; ok {

			info.Month = month
		}

		if year, err := strconv.Atoi(fields[2]); err == nil {

			info.Year = year
		}
	}
}

// ========================================
// applySearch
// ========================================

func applySearch(
	query string,
	args []interface{},
	search string,
) (string, []interface{}) {

	search = strings.TrimSpace(search)

	if search == "" {
		return query, args
	}

	info := parseSearch(search)

	query += `
	AND (
		nama_pengguna LIKE ?
		OR email LIKE ?
	`

	args = append(
		args,
		info.Keyword,
		info.Keyword,
	)

	// ========================================
	// Level
	// ========================================

	if info.Level != "" {
		query += `
		OR level = ?
		`
		args = append(args, info.Level)
	}

	// ========================================
	// Status
	// ========================================

	if info.Status != "" {
		query += `
		OR aktif = ?
		`
		args = append(args, info.Status)
	}

	// ========================================
	// Password
	// ========================================

	if info.Password != "" {
		query += `
		OR password_changed = ?
		`
		args = append(args, info.Password)
	}

	// ========================================
	// Belum Pernah Login
	// ========================================

	if info.LoginNull {
		query += `
		OR last_login IS NULL
		`
	}

	// ========================================
	// Login berdasarkan tanggal
	// (kombinasi menggunakan AND)
	// ========================================

	if info.Day > 0 || info.Month > 0 || info.Year > 0 {

		query += `
		OR (
			1=1
		`

		if info.Day > 0 {
			query += `
			AND DAY(last_login) = ?
			`
			args = append(args, info.Day)
		}

		if info.Month > 0 {
			query += `
			AND MONTH(last_login) = ?
			`
			args = append(args, info.Month)
		}

		if info.Year > 0 {
			query += `
			AND YEAR(last_login) = ?
			`
			args = append(args, info.Year)
		}

		query += `
		)
		`
	}

	query += `
	)
	`

	return query, args
}

// =========================================
// List
// =========================================

func List(
	page int,
	limit int,
	search string,
	level string,
	aktif string,
	sort string,
	order string,
) ([]Pengguna, error) {

	data := []Pengguna{}

	offset := (page - 1) * limit

	query := `
	SELECT
		id_pengguna,
		COALESCE(nuptk,'') AS nuptk,
		COALESCE(nis,'') AS nis,
		nama_pengguna,
		email,
		level,
		password_changed,
		last_login,
		aktif,
		foto,
		created_at,
		updated_at
	FROM pengguna
	WHERE 1=1
	`

	args := []interface{}{}

	// =========================================
	// Smart Search
	// =========================================

	query, args = applySearch(
		query,
		args,
		search,
	)

	// =========================================
	// Filter Level
	// =========================================

	if level != "" {

		query += `
		AND level = ?
		`

		args = append(args, level)
	}

	// =========================================
	// Filter Status
	// =========================================

	if aktif != "" {

		query += `
		AND aktif = ?
		`

		args = append(args, aktif)
	}

	// =========================================
	// Sorting
	// =========================================

	allowedSort := map[string]bool{
		"nama_pengguna":    true,
		"email":            true,
		"level":            true,
		"aktif":            true,
		"password_changed": true,
		"last_login":       true,
		"created_at":       true,
	}

	if !allowedSort[sort] {
		sort = "last_login"
	}

	order = strings.ToUpper(order)

	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	query += fmt.Sprintf(`
	ORDER BY %s %s, nama_pengguna ASC
	LIMIT ?
	OFFSET ?
	`, sort, order)

	args = append(
		args,
		limit,
		offset,
	)

	// =========================================
	// Execute Query
	// =========================================

	err := database.DB.Select(
		&data,
		query,
		args...,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

// =========================================
// Count List
// =========================================

func CountList(
	search string,
	level string,
	aktif string,
) (int, error) {

	var total int

	query := `
	SELECT COUNT(*)
	FROM pengguna
	WHERE 1=1
	`

	args := []interface{}{}

	// =========================================
	// Smart Search
	// =========================================

	query, args = applySearch(
		query,
		args,
		search,
	)

	// =========================================
	// Filter Level
	// =========================================

	if level != "" {

		query += `
		AND level = ?
		`

		args = append(args, level)
	}

	// =========================================
	// Filter Status
	// =========================================

	if aktif != "" {

		query += `
		AND aktif = ?
		`

		args = append(args, aktif)
	}

	// =========================================
	// Execute Query
	// =========================================

	err := database.DB.Get(
		&total,
		query,
		args...,
	)

	if err != nil {
		return 0, err
	}

	return total, nil
}

// =========================================
// Update Status
// =========================================

func updateStatus(
	ids []string,
	aktif int,
) error {

	query, args, err := sqlx.In(`
		UPDATE pengguna
		SET aktif = ?
		WHERE id_pengguna IN (?)
	`, aktif, ids)

	if err != nil {
		return err
	}

	query = database.DB.Rebind(query)

	_, err = database.DB.Exec(
		query,
		args...,
	)

	return err
}

// =========================================
// Update Password
// =========================================

func updatePassword(
	ids []string,
	password string,
	passwordChanged int,
) error {

	query, args, err := sqlx.In(`
		UPDATE pengguna
		SET
			password = ?,
			password_changed = ?,
			password_changed_at = NULL
		WHERE id_pengguna IN (?)
	`, password, passwordChanged, ids)

	if err != nil {
		return err
	}

	query = database.DB.Rebind(query)

	_, err = database.DB.Exec(
		query,
		args...,
	)

	return err
}

// =========================================
// Has Admin
// =========================================

func HasAdmin(ids []string) (bool, error) {

	var total int

	query, args, err := sqlx.In(`
		SELECT COUNT(*)
		FROM pengguna
		WHERE level = ?
		AND id_pengguna IN (?)
	`, "Admin", ids)

	if err != nil {
		return false, err
	}

	query = database.DB.Rebind(query)

	err = database.DB.Get(
		&total,
		query,
		args...,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

// =========================================
// Bulk Delete Pengguna
// =========================================

func BulkDelete(ids []string) error {

	query, args, err := sqlx.In(`
		DELETE FROM pengguna
		WHERE id_pengguna IN (?)
	`, ids)

	if err != nil {
		return err
	}

	query = database.DB.Rebind(query)

	_, err = database.DB.Exec(
		query,
		args...,
	)

	return err
}

// GetFotosByIDs mengambil daftar nama file foto dari sekumpulan pengguna,
// dipanggil SEBELUM BulkDelete supaya file fotonya masih bisa dihapus
// setelah baris pengguna itu sendiri sudah terhapus dari database.
func GetFotosByIDs(ids []string) ([]string, error) {

	var fotos []string

	query, args, err := sqlx.In(`
		SELECT COALESCE(foto,'') AS foto
		FROM pengguna
		WHERE id_pengguna IN (?)
	`, ids)

	if err != nil {
		return nil, err
	}

	query = database.DB.Rebind(query)

	err = database.DB.Select(&fotos, query, args...)

	if err != nil {
		return nil, err
	}

	return fotos, nil
}

// =========================================
// Bulk Activate
// =========================================

func BulkActivate(ids []string) error {

	return updateStatus(ids, 1)

}

// =========================================
// Bulk Deactivate
// =========================================

func BulkDeactivate(ids []string) error {

	return updateStatus(ids, 0)

}

// =========================================
// Bulk Reset Password
// =========================================

func BulkResetPassword(ids []string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(config.App.DefaultPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	return updatePassword(
		ids,
		string(hashedPassword),
		0,
	)

}

// =========================================
// Get By ID
// =========================================

func GetByID(id string) (*Pengguna, error) {

	var pengguna Pengguna

	err := database.DB.Get(
		&pengguna,
		`
		SELECT
			id_pengguna,
			nama_pengguna,
			email,
			level,
			aktif
		FROM pengguna
		WHERE id_pengguna = ?
		LIMIT 1
		`,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &pengguna, nil
}

// =========================================
// Get Pengguna By ID
// =========================================

func GetPenggunaByID(id string) (*Pengguna, error) {

	var pengguna Pengguna

	err := database.DB.Get(
		&pengguna,
		`
		SELECT
			id_pengguna,
			COALESCE(nuptk,'') AS nuptk,
			COALESCE(nis,'') AS nis,
			nama_pengguna,
			email,
			level,
			aktif,
			password_changed,
			last_login,
			foto,
			created_at,
			updated_at
		FROM pengguna
		WHERE id_pengguna = ?
		LIMIT 1
		`,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &pengguna, nil
}

// =========================================
// Get Password By ID
// =========================================

func GetPasswordByID(id string) (*UserPassword, error) {

	var user UserPassword

	err := database.DB.Get(
		&user,
		`
		SELECT
			id_pengguna,
			password
		FROM pengguna
		WHERE id_pengguna = ?
		LIMIT 1
		`,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

// =========================================
// Change Password
// =========================================

func ChangePassword(
	id string,
	hashedPassword string,
) error {

	_, err := database.DB.Exec(
		`
		UPDATE pengguna
		SET
			password = ?,
			password_changed = 1,
			password_changed_at = NOW()
		WHERE id_pengguna = ?
		`,
		hashedPassword,
		id,
	)

	return err

}

// =========================================
// Jenis Kelamin & Foto Default
// (L.jpg / P.jpg, dicari lewat NUPTK -> pegawai
// atau NIS -> pesertadidik)
// =========================================

func getJenisKelaminPengguna(nuptk string, nis string) string {

	var jenisKelamin string

	switch {

	case nuptk != "":

		_ = database.DB.Get(
			&jenisKelamin,
			`
			SELECT jenis_kelamin
			FROM pegawai
			WHERE nuptk = ?
			LIMIT 1
			`,
			nuptk,
		)

	case nis != "":

		_ = database.DB.Get(
			&jenisKelamin,
			`
			SELECT jenis_kelamin
			FROM pesertadidik
			WHERE nis = ?
			LIMIT 1
			`,
			nis,
		)
	}

	return jenisKelamin
}

func fotoDefaultForPengguna(jenisKelamin string) string {

	if jenisKelamin == "P" {
		return "P.jpg"
	}

	return "L.jpg"
}

func buildFotoURLPengguna(
	foto string,
	nuptk string,
	nis string,
) string {

	if foto == "" || foto == FotoDefaultPengguna {

		foto = fotoDefaultForPengguna(
			getJenisKelaminPengguna(nuptk, nis),
		)
	}

	return config.App.AppURL + "/storage/fotopengguna/" + foto
}

// =========================================
// Foto Pengguna (profil sendiri)
// =========================================

func GetFotoPengguna(id string) (string, error) {

	var foto string

	err := database.DB.Get(
		&foto,
		`
		SELECT COALESCE(foto,'') AS foto
		FROM pengguna
		WHERE id_pengguna = ?
		LIMIT 1
		`,
		id,
	)

	if err != nil {
		return "", err
	}

	return foto, nil
}

func GetFotoInfoPengguna(id string) (foto string, nuptk string, nis string, err error) {

	row := struct {
		Foto  string `db:"foto"`
		NUPTK string `db:"nuptk"`
		NIS   string `db:"nis"`
	}{}

	err = database.DB.Get(
		&row,
		`
		SELECT
			COALESCE(foto,'') AS foto,
			COALESCE(nuptk,'') AS nuptk,
			COALESCE(nis,'') AS nis
		FROM pengguna
		WHERE id_pengguna = ?
		LIMIT 1
		`,
		id,
	)

	if err != nil {
		return "", "", "", err
	}

	return row.Foto, row.NUPTK, row.NIS, nil
}

func UpdateFotoPengguna(id, namaFile string) error {

	_, err := database.DB.Exec(
		`
		UPDATE pengguna
		SET foto = ?
		WHERE id_pengguna = ?
		`,
		namaFile,
		id,
	)

	return err
}

// =========================================
// Export Excel
// =========================================

func ExportExcel() ([]ExportPengguna, error) {

	data := []ExportPengguna{}

	query := `
	SELECT
		nama_pengguna,
		email,
		level,
		aktif,
		password_changed,
		last_login
	FROM pengguna
	ORDER BY nama_pengguna ASC
	`

	err := database.DB.Select(
		&data,
		query,
	)

	if err != nil {
		return nil, err
	}

	return data, nil
}

// =========================================
// Read Upload Excel
// =========================================

func ReadUploadExcel(f *excelize.File) ([]UploadPengguna, error) {

	const sheet = "Data Pengguna"

	// =========================================
	// Pastikan Sheet Ada
	// =========================================

	rows, err := f.GetRows(sheet)

	if err != nil {
		return nil, fmt.Errorf(
			"sheet '%s' tidak ditemukan",
			sheet,
		)
	}

	lastRow := len(rows)

	if lastRow <= 1 {

		return nil, fmt.Errorf(
			"file Excel tidak memiliki data",
		)
	}

	// =========================================
	// Validasi Header
	// =========================================

	expected := []string{
		"Nama Pengguna",
		"Email",
		"Level",
		"NUPTK",
		"NIS",
	}

	for col, title := range expected {

		cell, _ := excelize.CoordinatesToCellName(col+1, 1)

		value, err := f.GetCellValue(sheet, cell)

		if err != nil {
			return nil, err
		}

		if strings.TrimSpace(value) != title {

			return nil, fmt.Errorf(
				"format template tidak sesuai, silakan gunakan template yang disediakan",
			)
		}
	}

	// =========================================
	// Membaca Data
	// =========================================

	data := make([]UploadPengguna, 0)

	for row := 2; row <= lastRow; row++ {

		item := UploadPengguna{}

		item.NamaPengguna, _ = f.GetCellValue(
			sheet,
			fmt.Sprintf("A%d", row),
		)

		item.Email, _ = f.GetCellValue(
			sheet,
			fmt.Sprintf("B%d", row),
		)

		item.Level, _ = f.GetCellValue(
			sheet,
			fmt.Sprintf("C%d", row),
		)

		item.NUPTK, _ = f.GetCellValue(
			sheet,
			fmt.Sprintf("D%d", row),
		)

		item.NIS, _ = f.GetCellValue(
			sheet,
			fmt.Sprintf("E%d", row),
		)

		item.NamaPengguna = strings.TrimSpace(item.NamaPengguna)
		item.Email = strings.TrimSpace(item.Email)
		item.Level = strings.TrimSpace(item.Level)
		item.NUPTK = strings.TrimSpace(item.NUPTK)
		item.NIS = strings.TrimSpace(item.NIS)

		// =====================================
		// Lewati Baris Kosong
		// =====================================

		if item.NamaPengguna == "" &&
			item.Email == "" &&
			item.Level == "" &&
			item.NUPTK == "" &&
			item.NIS == "" {

			continue
		}

		data = append(data, item)
	}

	if len(data) == 0 {

		return nil, fmt.Errorf(
			"file Excel tidak memiliki data",
		)
	}

	return data, nil
}

// =========================================
// Validate Upload
// =========================================

func ValidateUpload(
	data []UploadPengguna,
) ([]UploadError, error) {

	var errors []UploadError

	emailMap := make(map[string]int)
	nuptkMap := make(map[string]int)
	nisMap := make(map[string]int)

	validLevel := map[string]bool{
		"Admin":        true,
		"Staf":         true,
		"Pegawai":      true,
		"Pesertadidik": true,
	}

	for i, item := range data {

		row := i + 2

		// =====================================
		// Nama
		// =====================================

		if item.NamaPengguna == "" {

			errors = append(errors, UploadError{
				Row:     row,
				Message: "Nama Pengguna wajib diisi.",
			})
		}

		// =====================================
		// Email Wajib
		// =====================================

		if item.Email == "" {

			errors = append(errors, UploadError{
				Row:     row,
				Message: "Email wajib diisi.",
			})
		}

		// =====================================
		// Format Email
		// =====================================

		if item.Email != "" &&
			!IsValidEmail(item.Email) {

			errors = append(errors, UploadError{
				Row: row,
				Message: fmt.Sprintf(
					"Email %s tidak valid.",
					item.Email,
				),
			})
		}

		// =====================================
		// Level
		// =====================================

		if !validLevel[item.Level] {

			errors = append(errors, UploadError{
				Row:     row,
				Message: "Level tidak valid.",
			})
		}

		// =====================================
		// Pegawai
		// =====================================

		if item.Level == "Pegawai" && item.NUPTK == "" {

			errors = append(errors, UploadError{
				Row:     row,
				Message: "NUPTK wajib diisi.",
			})
		}

		// =====================================
		// Peserta Didik
		// =====================================

		if item.Level == "Pesertadidik" && item.NIS == "" {

			errors = append(errors, UploadError{
				Row:     row,
				Message: "NIS wajib diisi.",
			})
		}

		// =====================================
		// Duplikat Email di File
		// =====================================

		if item.Email != "" {

			if old, ok := emailMap[item.Email]; ok {

				errors = append(errors, UploadError{
					Row: row,
					Message: fmt.Sprintf(
						"Email duplikat dengan baris %d.",
						old,
					),
				})

			} else {

				emailMap[item.Email] = row
			}

			// =====================================
			// Email Database
			// =====================================

			exists, err := EmailExists(item.Email)

			if err != nil {
				return nil, err
			}

			if exists {

				errors = append(errors, UploadError{
					Row: row,
					Message: fmt.Sprintf(
						"Email %s sudah digunakan.",
						item.Email,
					),
				})
			}
		}

		// =====================================
		// Duplikat NUPTK di File
		// =====================================

		if item.NUPTK != "" {
			// =====================================
			// Format NUPTK
			// =====================================

			if !IsNumeric(item.NUPTK) ||
				len(item.NUPTK) != 16 {

				errors = append(errors, UploadError{
					Row:     row,
					Message: "NUPTK harus terdiri dari 16 digit angka.",
				})
			}

			if old, ok := nuptkMap[item.NUPTK]; ok {

				errors = append(errors, UploadError{
					Row: row,
					Message: fmt.Sprintf(
						"NUPTK duplikat dengan baris %d.",
						old,
					),
				})

			} else {

				nuptkMap[item.NUPTK] = row
			}

			// =====================================
			// NUPTK Database
			// =====================================

			exists, err := NUPTKExists(item.NUPTK)

			if err != nil {
				return nil, err
			}

			if exists {

				errors = append(errors, UploadError{
					Row: row,
					Message: fmt.Sprintf(
						"NUPTK %s sudah digunakan.",
						item.NUPTK,
					),
				})
			}
		}

		// =====================================
		// Duplikat NIS di File
		// =====================================

		if item.NIS != "" {
			// =====================================
			// Format NIS
			// =====================================

			if item.NIS != "" &&
				!IsNumeric(item.NIS) {

				errors = append(errors, UploadError{
					Row:     row,
					Message: "NIS hanya boleh berisi angka.",
				})
			}

			if old, ok := nisMap[item.NIS]; ok {

				errors = append(errors, UploadError{
					Row: row,
					Message: fmt.Sprintf(
						"NIS duplikat dengan baris %d.",
						old,
					),
				})

			} else {

				nisMap[item.NIS] = row
			}

			// =====================================
			// NIS Database
			// =====================================

			exists, err := NISExists(item.NIS)

			if err != nil {
				return nil, err
			}

			if exists {

				errors = append(errors, UploadError{
					Row: row,
					Message: fmt.Sprintf(
						"NIS %s sudah digunakan.",
						item.NIS,
					),
				})
			}
		}
	}

	return errors, nil
}

func EmailExists(email string) (bool, error) {

	var total int

	err := database.DB.Get(
		&total,
		`
		SELECT COUNT(*)
		FROM pengguna
		WHERE email = ?
		`,
		email,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

// =========================================
// Email Exists Except
// =========================================

func EmailExistsExcept(
	email string,
	id string,
) (bool, error) {

	var total int

	err := database.DB.Get(
		&total,
		`
		SELECT COUNT(*)
		FROM pengguna
		WHERE email = ?
		AND id_pengguna <> ?
		`,
		email,
		id,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

func NUPTKExists(nuptk string) (bool, error) {

	var total int

	err := database.DB.Get(
		&total,
		`
		SELECT COUNT(*)
		FROM pengguna
		WHERE nuptk = ?
		`,
		nuptk,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

// =========================================
// NUPTK Exists Except
// =========================================

func NUPTKExistsExcept(
	nuptk string,
	id string,
) (bool, error) {

	var total int

	err := database.DB.Get(
		&total,
		`
		SELECT COUNT(*)
		FROM pengguna
		WHERE nuptk = ?
		AND id_pengguna <> ?
		`,
		nuptk,
		id,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

// =========================================
// Check Duplicate NUPTK (realtime, Tambah/Ubah)
// =========================================

type DuplicateInfo struct {
	Exists       bool   `json:"exists"`
	NamaPengguna string `json:"nama_pengguna,omitempty"`
}

func CheckDuplicateNUPTK(
	nuptk string,
	exclude string,
) (DuplicateInfo, error) {

	info := DuplicateInfo{}

	query := `
	SELECT nama_pengguna
	FROM pengguna
	WHERE nuptk = ?
	`

	args := []interface{}{nuptk}

	if exclude != "" {

		query += `
		AND id_pengguna <> ?
		`

		args = append(args, exclude)
	}

	query += `
	LIMIT 1
	`

	var namaPengguna string

	err := database.DB.Get(&namaPengguna, query, args...)

	if err != nil {

		if err == sql.ErrNoRows {
			return info, nil
		}

		return info, err
	}

	info.Exists = true
	info.NamaPengguna = namaPengguna

	return info, nil
}

func NISExists(nis string) (bool, error) {

	var total int

	err := database.DB.Get(
		&total,
		`
		SELECT COUNT(*)
		FROM pengguna
		WHERE nis = ?
		`,
		nis,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

// =========================================
// NIS Exists Except
// =========================================

func NISExistsExcept(
	nis string,
	id string,
) (bool, error) {

	var total int

	err := database.DB.Get(
		&total,
		`
		SELECT COUNT(*)
		FROM pengguna
		WHERE nis = ?
		AND id_pengguna <> ?
		`,
		nis,
		id,
	)

	if err != nil {
		return false, err
	}

	return total > 0, nil
}

// =========================================
// Bulk Insert
// =========================================

func BulkInsert(data []UploadPengguna) error {

	tx, err := database.DB.Beginx()

	if err != nil {
		return err
	}

	defer func() {

		if err != nil {
			_ = tx.Rollback()
		}

	}()

	password, err := bcrypt.GenerateFromPassword(
		[]byte(config.App.DefaultPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	query := `
	INSERT INTO pengguna
(
	id_pengguna,
	nuptk,
	nis,
	nama_pengguna,
	email,
	password,
	password_changed,
	level,
	aktif,
	foto
)
VALUES
(
	?,
	?,
	?,
	?,
	?,
	?,
	0,
	?,
	1,
	?
)
	`

	for _, item := range data {

		idPengguna := uuid.NewString()

		var nuptk interface{}
		var nis interface{}

		if item.NUPTK != "" {
			nuptk = item.NUPTK
		}

		if item.NIS != "" {
			nis = item.NIS
		}

		_, err = tx.Exec(
			query,
			idPengguna,
			nuptk,
			nis,
			item.NamaPengguna,
			item.Email,
			string(password),
			item.Level,
			"1636812404_b39deec2f158cb789567.jpg",
		)

		if err != nil {
			return err
		}
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

// =========================================
// Create Pengguna
// =========================================

func CreatePengguna(req CreatePenggunaRequest) (string, error) {

	// =====================================
	// Nama
	// =====================================

	req.NamaPengguna = strings.TrimSpace(req.NamaPengguna)

	if req.NamaPengguna == "" {
		return "", fmt.Errorf("Nama Pengguna wajib diisi.")
	}

	// =====================================
	// Email
	// =====================================

	req.Email = strings.TrimSpace(req.Email)

	if req.Email == "" {
		return "", fmt.Errorf("Email wajib diisi.")
	}

	if !IsValidEmail(req.Email) {
		return "", fmt.Errorf("Format email tidak valid.")
	}

	exists, err := EmailExists(req.Email)

	if err != nil {
		return "", err
	}

	if exists {
		return "", fmt.Errorf("Email sudah digunakan.")
	}

	// =====================================
	// Level
	// =====================================

	req.Level = strings.TrimSpace(req.Level)

	validLevel := map[string]bool{
		"Admin":        true,
		"Staf":         true,
		"Pegawai":      true,
		"Pesertadidik": true,
	}

	if !validLevel[req.Level] {
		return "", fmt.Errorf("Level tidak valid.")
	}

	// =====================================
	// Pegawai
	// =====================================

	req.NUPTK = strings.TrimSpace(req.NUPTK)

	if req.Level == "Pegawai" {

		if req.NUPTK == "" {
			return "", fmt.Errorf("NUPTK wajib diisi.")
		}

		if !IsNumeric(req.NUPTK) || len(req.NUPTK) != 16 {
			return "", fmt.Errorf("NUPTK harus terdiri dari 16 digit angka.")
		}

		exists, err := NUPTKExists(req.NUPTK)

		if err != nil {
			return "", err
		}

		if exists {
			return "", fmt.Errorf("NUPTK sudah digunakan.")
		}
	}

	// =====================================
	// Peserta Didik
	// =====================================

	req.NIS = strings.TrimSpace(req.NIS)

	if req.Level == "Pesertadidik" {

		if req.NIS == "" {
			return "", fmt.Errorf("NIS wajib diisi.")
		}

		if !IsNumeric(req.NIS) {
			return "", fmt.Errorf("NIS hanya boleh berisi angka.")
		}

		exists, err := NISExists(req.NIS)

		if err != nil {
			return "", err
		}

		if exists {
			return "", fmt.Errorf("NIS sudah digunakan.")
		}
	}

	// =====================================
	// Password
	// =====================================

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(config.App.DefaultPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	// =====================================
	// NULL Value
	// =====================================

	var nuptk interface{}
	var nis interface{}

	if req.NUPTK != "" {
		nuptk = req.NUPTK
	}

	if req.NIS != "" {
		nis = req.NIS
	}

	// =====================================
	// Insert
	// =====================================

	idPengguna := uuid.NewString()

	_, err = database.DB.Exec(
		`
		INSERT INTO pengguna
		(
			id_pengguna,
			nuptk,
			nis,
			nama_pengguna,
			email,
			password,
			password_changed,
			level,
			aktif,
			foto
		)
		VALUES
		(
			?,
			?,
			?,
			?,
			?,
			?,
			0,
			?,
			1,
			?
		)
		`,
		idPengguna,
		nuptk,
		nis,
		req.NamaPengguna,
		req.Email,
		string(hashedPassword),
		req.Level,
		"1636812404_b39deec2f158cb789567.jpg",
	)

	if err != nil {
		return "", err
	}

	return idPengguna, nil
}

// =========================================
// Update Pengguna
// =========================================

func UpdatePengguna(req UpdatePenggunaRequest) (bool, error) {

	// =====================================
	// Ambil Data Lama
	// =====================================

	old, err := GetPenggunaByID(req.IDPengguna)

	if err != nil {
		return false, fmt.Errorf("Data pengguna tidak ditemukan.")
	}

	// =====================================
	// Nama
	// =====================================

	req.NamaPengguna = strings.TrimSpace(req.NamaPengguna)

	if req.NamaPengguna == "" {
		return false, fmt.Errorf("Nama Pengguna wajib diisi.")
	}

	// =====================================
	// Email
	// =====================================

	req.Email = strings.TrimSpace(req.Email)

	if req.Email == "" {
		return false, fmt.Errorf("Email wajib diisi.")
	}

	if !IsValidEmail(req.Email) {
		return false, fmt.Errorf("Format email tidak valid.")
	}

	exists, err := EmailExistsExcept(
		req.Email,
		req.IDPengguna,
	)

	if err != nil {
		return false, err
	}

	if exists {
		return false, fmt.Errorf("Email sudah digunakan.")
	}

	// =====================================
	// Level
	// =====================================

	req.Level = strings.TrimSpace(req.Level)

	validLevel := map[string]bool{
		"Admin":        true,
		"Staf":         true,
		"Pegawai":      true,
		"Pesertadidik": true,
	}

	if !validLevel[req.Level] {
		return false, fmt.Errorf("Level tidak valid.")
	}

	// =====================================
	// Trim
	// =====================================

	req.NUPTK = strings.TrimSpace(req.NUPTK)
	req.NIS = strings.TrimSpace(req.NIS)

	// =====================================
	// Tidak Ada Perubahan
	// =====================================

	if old.NamaPengguna == req.NamaPengguna &&
		old.Email == req.Email &&
		old.Level == req.Level &&
		old.NUPTK == req.NUPTK &&
		old.NIS == req.NIS {

		return false, nil
	}

	// =====================================
	// Validasi berdasarkan Level
	// =====================================

	switch req.Level {

	case "Pegawai":

		if req.NUPTK == "" {
			return false, fmt.Errorf("NUPTK wajib diisi.")
		}

		if !IsNumeric(req.NUPTK) || len(req.NUPTK) != 16 {
			return false, fmt.Errorf("NUPTK harus terdiri dari 16 digit angka.")
		}

		exists, err := NUPTKExistsExcept(
			req.NUPTK,
			req.IDPengguna,
		)

		if err != nil {
			return false, err
		}

		if exists {
			return false, fmt.Errorf("NUPTK sudah digunakan.")
		}

		req.NIS = ""

	case "Admin", "Staf":

		if req.NUPTK != "" {

			if !IsNumeric(req.NUPTK) || len(req.NUPTK) != 16 {
				return false, fmt.Errorf("NUPTK harus terdiri dari 16 digit angka.")
			}

			exists, err := NUPTKExistsExcept(
				req.NUPTK,
				req.IDPengguna,
			)

			if err != nil {
				return false, err
			}

			if exists {
				return false, fmt.Errorf("NUPTK sudah digunakan.")
			}
		}

		req.NIS = ""

	case "Pesertadidik":

		if req.NIS == "" {
			return false, fmt.Errorf("NIS wajib diisi.")
		}

		if !IsNumeric(req.NIS) {
			return false, fmt.Errorf("NIS hanya boleh berisi angka.")
		}

		exists, err := NISExistsExcept(
			req.NIS,
			req.IDPengguna,
		)

		if err != nil {
			return false, err
		}

		if exists {
			return false, fmt.Errorf("NIS sudah digunakan.")
		}

		req.NUPTK = ""
	}

	// =====================================
	// NULL Value
	// =====================================

	var nuptk interface{}
	var nis interface{}

	if req.NUPTK != "" {
		nuptk = req.NUPTK
	}

	if req.NIS != "" {
		nis = req.NIS
	}

	// =====================================
	// Update Database
	// =====================================

	_, err = database.DB.Exec(
		`
		UPDATE pengguna
		SET
			nuptk = ?,
			nis = ?,
			nama_pengguna = ?,
			email = ?,
			level = ?
		WHERE id_pengguna = ?
		`,
		nuptk,
		nis,
		req.NamaPengguna,
		req.Email,
		req.Level,
		req.IDPengguna,
	)

	if err != nil {
		return false, err
	}

	return true, nil
}

func IsValidEmail(email string) bool {

	_, err := mail.ParseAddress(email)

	return err == nil
}

func IsNumeric(value string) bool {

	for _, r := range value {

		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}
