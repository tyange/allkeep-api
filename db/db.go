package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		panic("Could not connect to database.")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT
	)
	`

	_, err := DB.Exec(createUsersTable)

	if err != nil {
		panic("Could not create users table.")
	}

	createWorksTable := `
	CREATE TABLE IF NOT EXISTS works (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		start_at TEXT NOT NULL,
		company_name TEXT NOT NULL,
		user_id INTEGER
	)
	`

	_, err = DB.Exec(createWorksTable)

	if err != nil {
		panic("Could not create works table.")
	}

	createCompaniesTable := `
	CREATE TABLE IF NOT EXISTS companies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		company_name TEXT NOT NULL,
		user_id INTEGER NOT NULL
	)
	`

	_, err = DB.Exec(createCompaniesTable)

	if err != nil {
		panic("Could not create companies table.")
	}
}
