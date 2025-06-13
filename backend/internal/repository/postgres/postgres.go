package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"

	"github.com/qwaq-dev/macan-ai/internal/config"
)

func InitDatabase(cfg config.Database, log *slog.Logger) (*sql.DB, error) {
	const op = "repository.posgres.InitDatabase"
	log.With("op", op)

	connect := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.DBusername, cfg.DBname, cfg.DBpassword, cfg.SSLmode)

	db, err := sql.Open("postgres", connect)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Info("Database connecting successfully")
	return db, nil
}
