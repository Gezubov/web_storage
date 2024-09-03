package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"web_storage/internal/config"
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
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL database successfully")

	if err := Migrate(db); err != nil {
		log.Fatal(err)
	}
	log.Println("Migrated  successfully")
	return db, nil
}
