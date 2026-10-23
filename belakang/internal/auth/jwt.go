package auth

import (
	"time"

	"sekolah/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	IDPengguna string `json:"id_pengguna"`
	Email      string `json:"email"`
	Level      string `json:"level"`
	NUPTK      string `json:"nuptk,omitempty"`
	NIS        string `json:"nis,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken membuat JWT untuk user yang berhasil login.
func GenerateToken(user *User) (string, error) {

	now := time.Now()

	expire := now.Add(
		time.Duration(config.App.JWTExpireMinutes) * time.Minute,
	)

	claims := Claims{
		IDPengguna: user.IDPengguna,
		Email:      user.Email,
		Level:      user.Level,
		NUPTK:      user.NUPTK,
		NIS:        user.NIS,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.App.AppName,
			Subject:   user.IDPengguna,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expire),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.App.JWTSecret))
}

// ParseToken membaca dan memvalidasi JWT.
func ParseToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		getJWTKey,
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

// getJWTKey memvalidasi algoritma JWT dan mengembalikan secret key.
func getJWTKey(token *jwt.Token) (interface{}, error) {

	// Pastikan hanya menerima HMAC (HS256)
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, jwt.ErrTokenSignatureInvalid
	}

	return []byte(config.App.JWTSecret), nil
}
