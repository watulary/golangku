package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppEnv  string
	AppPort string
	AppURL  string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPass     string
	DBTimezone string

	JWTSecret        string
	JWTExpireMinutes int

	DefaultPassword string

	UploadMaxSize int64
}

var App Config

func Load() error {
	if err := godotenv.Load(); err != nil {
		// Abaikan jika .env memang tidak ada (misalnya di production)
	}

	jwtExpireMinutes, err := strconv.Atoi(os.Getenv("JWT_EXPIRE_MINUTES"))
	if err != nil {
		jwtExpireMinutes = 1440 // default 1 hari
	}

	uploadMaxSize, err := strconv.ParseInt(os.Getenv("UPLOAD_MAX_SIZE"), 10, 64)
	if err != nil {
		uploadMaxSize = 2 * 1024 * 1024 // default 2 MB
	}

	App = Config{
		AppName: os.Getenv("APP_NAME"),
		AppEnv:  os.Getenv("APP_ENV"),
		AppPort: os.Getenv("APP_PORT"),
		AppURL:  os.Getenv("APP_URL"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPass:     os.Getenv("DB_PASS"),
		DBTimezone: os.Getenv("DB_TIMEZONE"),

		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTExpireMinutes: jwtExpireMinutes,

		DefaultPassword: os.Getenv("DEFAULT_PASSWORD"),

		UploadMaxSize: uploadMaxSize,
	}

	return nil
}
