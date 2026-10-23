package pengguna

import (
	"time"
)

type Pengguna struct {
	IDPengguna string `db:"id_pengguna" json:"id_pengguna"`

	NUPTK string `db:"nuptk" json:"nuptk,omitempty"`
	NIS   string `db:"nis" json:"nis,omitempty"`

	NamaPengguna string `db:"nama_pengguna" json:"nama_pengguna"`
	Email        string `db:"email" json:"email"`

	Level string `db:"level" json:"level"`

	Aktif int `db:"aktif" json:"aktif"`

	PasswordChanged bool       `db:"password_changed" json:"password_changed"`
	LastLogin       *time.Time `db:"last_login" json:"last_login,omitempty"`

	Foto    string `db:"foto" json:"foto"`
	FotoURL string `json:"foto_url"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// =========================================
// Upload Excel
// =========================================

type UploadPengguna struct {
	NamaPengguna string
	Email        string
	Level        string
	NUPTK        string
	NIS          string
}

type UploadError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}

// =========================================
// Tambah Pengguna
// =========================================

type CreatePenggunaRequest struct {
	NamaPengguna string `json:"nama_pengguna"`
	Email        string `json:"email"`
	Level        string `json:"level"`

	NUPTK string `json:"nuptk"`
	NIS   string `json:"nis"`
}

// =========================================
// Update Pengguna
// =========================================

type UpdatePenggunaRequest struct {
	IDPengguna string `json:"id_pengguna"`

	NamaPengguna string `json:"nama_pengguna"`
	Email        string `json:"email"`
	Level        string `json:"level"`

	NUPTK string `json:"nuptk"`
	NIS   string `json:"nis"`
}

// =========================================
// Export Excel
// =========================================

type ExportPengguna struct {
	NamaPengguna string     `db:"nama_pengguna"`
	Email        string     `db:"email"`
	Level        string     `db:"level"`
	Aktif        int        `db:"aktif"`
	Password     bool       `db:"password_changed"`
	LastLogin    *time.Time `db:"last_login"`
}

type Pagination struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

type ListResponse struct {
	Success    bool       `json:"success"`
	Message    string     `json:"message"`
	Data       []Pengguna `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type BulkStatusRequest struct {
	IDs []string `json:"ids"`
}

// =========================================
// Change Password
// =========================================

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

// =========================================
// Password User
// =========================================

type UserPassword struct {
	IDPengguna string `db:"id_pengguna"`

	Password string `db:"password"`
}
