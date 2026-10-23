package database

import (
	"fmt"
	"net/url"
	"time"

	"sekolah/internal/config"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sqlx.DB

func Connect() error {

	loc := url.QueryEscape(config.App.DBTimezone)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=%s",
		config.App.DBUser,
		config.App.DBPass,
		config.App.DBHost,
		config.App.DBPort,
		config.App.DBName,
		loc,
	)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return err
	}

	DB = db

	return nil
}