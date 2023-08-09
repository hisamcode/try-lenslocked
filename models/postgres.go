package models

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// Open will open a SQL connection with the provided postgres databases.
// Caller of Open need to ensure that connection is eventually closed via the
// db.Close() method.
func Open(config PostgreConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.String())
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	return db, nil
}

func DefaultPostgresConfig() PostgreConfig {
	return PostgreConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "admin",
		Password: "password",
		Database: "lenslocked",
		SSLMode:  "disable",
	}
}

type PostgreConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgreConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}
