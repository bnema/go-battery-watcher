package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDB() (*sql.DB, error) {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "powerdata.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateTables(db *sql.DB) error {
	// Create the power_usage table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS power_usage (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		device_name TEXT UNIQUE,
		power_usage REAL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS power_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		device_id INTEGER,
		power_usage REAL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (device_id) REFERENCES power_usage(id)
	);
	`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return err
	}

	return nil
}
