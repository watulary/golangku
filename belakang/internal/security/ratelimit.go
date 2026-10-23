package security

import (
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	MaxLoginAttempts = 5
	BlockDuration    = 10 * time.Minute
)

type LoginAttempt struct {
	Count        int
	BlockedUntil time.Time
	LastAttempt  time.Time
}

var (
	loginAttempts = make(map[string]*LoginAttempt)
	mu            sync.RWMutex
)

// key = IP + Email
func getLoginKey(c *fiber.Ctx, email string) string {
	ip := c.IP()
	email = strings.TrimSpace(strings.ToLower(email))
	return ip + ":" + email
}

// CheckLoginLimit mengecek apakah login sedang diblok.
// Return:
//
//	blocked    = true jika masih diblok
//	retryAfter = sisa waktu blokir (detik)
func CheckLoginLimit(c *fiber.Ctx, email string) (blocked bool, retryAfter int) {

	key := getLoginKey(c, email)

	mu.RLock()
	attempt, exists := loginAttempts[key]
	mu.RUnlock()

	if !exists {
		return false, 0
	}

	now := time.Now()

	// Blokir sudah habis
	if !attempt.BlockedUntil.IsZero() && now.After(attempt.BlockedUntil) {

		mu.Lock()
		delete(loginAttempts, key)
		mu.Unlock()

		return false, 0
	}

	// Masih diblok
	if !attempt.BlockedUntil.IsZero() && now.Before(attempt.BlockedUntil) {

		return true, int(time.Until(attempt.BlockedUntil).Seconds())
	}

	return false, 0
}

// Dipanggil ketika login gagal
func IncrementLoginFail(c *fiber.Ctx, email string) {

	key := getLoginKey(c, email)

	mu.Lock()
	defer mu.Unlock()

	now := time.Now()

	attempt, exists := loginAttempts[key]

	if !exists {

		loginAttempts[key] = &LoginAttempt{
			Count:       1,
			LastAttempt: now,
		}

		return
	}

	// Kalau sudah diblok, jangan tambah counter lagi
	if !attempt.BlockedUntil.IsZero() && now.Before(attempt.BlockedUntil) {
		return
	}

	attempt.Count++
	attempt.LastAttempt = now

	if attempt.Count >= MaxLoginAttempts {
		attempt.BlockedUntil = now.Add(BlockDuration)
	}
}

// Dipanggil ketika login berhasil
func ResetLoginFail(c *fiber.Ctx, email string) {

	key := getLoginKey(c, email)

	mu.Lock()
	delete(loginAttempts, key)
	mu.Unlock()
}
