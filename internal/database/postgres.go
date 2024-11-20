package database

import (
	"MusicLibrary/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func InitDB(cfg *config.Config) (*Postgres, error) {
	DSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Hostname,
		cfg.Postgres.Port,
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)
	db, _ := sql.Open("postgres", DSN)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Postgres{}, nil
}
