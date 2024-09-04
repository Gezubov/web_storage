package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"web_storage/internal/config"
	"web_storage/pkg/logger"
)

func Connect(cfg *config.Config) (*sql.DB, error) {

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.DBName,
		cfg.DB.Host,
		cfg.DB.Port,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Logger.Error("Failed to open database connection", slog.String("error", err.Error()))
		return nil, err
	}

	if err := db.Ping(); err != nil {
		logger.Logger.Error("Failed to ping database", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Logger.Info("Connected to PostgreSQL database successfully")

	if err := Migrate(db); err != nil {
		logger.Logger.Error("Database migration failed", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Logger.Info("Database migration completed successfully")
	return db, nil
}
