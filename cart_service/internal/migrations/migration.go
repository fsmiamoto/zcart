package migrations

import (
	"database/sql"
	_ "embed"
)

//go:embed migration.sql
var migrations string

func Apply(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(migrations)
	if err != nil {
		defer tx.Rollback()
		return err
	}

	return tx.Commit()
}
