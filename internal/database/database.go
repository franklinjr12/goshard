package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DbConnectionParams struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	Sslmode  string
}

var DbTestConnectionParams = DbConnectionParams{
	Host:     "localhost",
	Port:     5432,
	User:     "postgres",
	Password: "postgres",
	Dbname:   "testapplication",
	Sslmode:  "disable",
}

const DbTestConnectionString = "host=localhost port=5432 user=postgres password=postgres dbname=testapplication sslmode=disable"

func BuildConnectionString(params DbConnectionParams) string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s", params.Host, params.Port, params.User, params.Password)
	if params.Dbname != "" {
		dsn += fmt.Sprintf(" dbname=%s", params.Dbname)
	}
	if params.Sslmode != "" {
		dsn += fmt.Sprintf(" sslmode=%s", params.Sslmode)
	}
	return dsn
}

func Connect(dsn string) (*sql.DB, error) {
	if dsn == "" {
		dsn = DbTestConnectionString
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}

func Close(db *sql.DB) {
	db.Close()
}

func Query(db *sql.DB, query string) (*sql.Rows, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}
	return rows, nil
}

func CreateDatabaseFromSchema(db *sql.DB, schema string) error {
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create database from schema: %w", err)
	}
	return nil
}

func ReadSchemaFromFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read schema from file: %w", err)
	}
	return string(content), nil
}
