package postgres

import (
	"database/sql"
	"log"
)

func Migrate(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS files (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        size BIGINT NOT NULL,
        link TEXT NOT NULL
    );`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	log.Println("Migration completed successfully")
	return nil
}
