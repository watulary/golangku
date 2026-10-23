package auth

import "time"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	IDPengguna string `db:"id_pengguna" json:"id_pengguna"`

	NUPTK string `db:"nuptk" json:"nuptk,omitempty"`
	NIS   string `db:"nis" json:"nis,omitempty"`

	NamaPengguna string `db:"nama_pengguna" json:"nama_pengguna"`
	Email        string `db:"email" json:"email"`

	Password string `db:"password" json:"-"`

	Level string `db:"level" json:"level"`

	Aktif int `db:"aktif" json:"aktif"`

	PasswordChanged bool       `db:"password_changed" json:"password_changed"`
	LastLogin       *time.Time `db:"last_login" json:"last_login,omitempty"`

	Foto    string `db:"foto" json:"foto"`
	FotoURL string `json:"foto_url"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	User    *User  `json:"user,omitempty"`
}
