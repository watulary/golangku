package auth

import (
	"database/sql"
	"errors"
	"time"

	"sekolah/internal/config"
	"sekolah/internal/database"
)

// Sentinel error
var ErrUserNotFound = errors.New("user not found")

// =====================================================
// QUERY
// =====================================================

// GetByID mengambil data pengguna berdasarkan ID.
func GetByID(idPengguna string) (*User, error) {

	var user User

	query := `
	SELECT
    id_pengguna,
    COALESCE(nuptk, '') AS nuptk,
    COALESCE(nis, '') AS nis,
    nama_pengguna,
    email,
    password,
    level,
    aktif,
    password_changed,
    last_login,
    foto
FROM pengguna
WHERE id_pengguna = ?
LIMIT 1
	`

	err := database.DB.Get(&user, query, idPengguna)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	buildFotoURL(&user)

	return &user, nil
}

// GetByEmail mengambil data pengguna berdasarkan email.
func GetByEmail(email string) (*User, error) {

	var user User

	query := `
	SELECT
    id_pengguna,
    COALESCE(nuptk, '') AS nuptk,
    COALESCE(nis, '') AS nis,
    nama_pengguna,
    email,
    password,
    level,
    aktif,
    password_changed,
    last_login,
    foto
FROM pengguna
WHERE email = ?
LIMIT 1
	`

	err := database.DB.Get(&user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	buildFotoURL(&user)

	return &user, nil
}

// =====================================================
// COMMAND
// =====================================================

// UpdateLastLogin memperbarui waktu login terakhir pengguna.
func UpdateLastLogin(idPengguna string, loginTime time.Time) error {

	query := `
	UPDATE pengguna
	SET last_login = ?
	WHERE id_pengguna = ?
	`

	_, err := database.DB.Exec(query, loginTime, idPengguna)
	return err
}

// fotoDefaultPengguna: nama file default LAMA (sebelum ada L.jpg/P.jpg
// berbasis jenis_kelamin). Masih dicek utk kompatibilitas data lama yg
// tersimpan foto = fotoDefaultPengguna di DB.
const fotoDefaultPengguna = "1636812404_b39deec2f158cb789567.jpg"

// getJenisKelaminUser mencari jenis_kelamin lewat NUPTK -> pegawai
// atau NIS -> pesertadidik. Duplikasi kecil dari package pengguna
// (tidak bisa diimpor langsung krn package pengguna sudah impor package
// auth utk auth.Claims -> circular import kalau dibalik).
func getJenisKelaminUser(nuptk string, nis string) string {

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

func fotoDefaultForUser(jenisKelamin string) string {

	if jenisKelamin == "P" {
		return "P.jpg"
	}

	return "L.jpg"
}

func buildFotoURL(user *User) {

	if user == nil {
		return
	}

	foto := user.Foto

	if foto == "" || foto == fotoDefaultPengguna {
		foto = fotoDefaultForUser(getJenisKelaminUser(user.NUPTK, user.NIS))
	}

	user.FotoURL = config.App.AppURL +
		"/storage/fotopengguna/" +
		foto
}
