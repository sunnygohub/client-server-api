package internal

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitializeDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS exchange_rates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			rate TEXT,
			create_date DATETIME
		)
	`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InsertExchangeRate(ctx context.Context, db *sql.DB, rate ExchangeRate) error {
	query := "INSERT INTO exchange_rates (rate, create_date) VALUES (?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, err := stmt.Exec(rate.Bid, rate.CreateDate)
		return err
	}
}
