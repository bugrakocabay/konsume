package e2e

import (
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
)

// connectToPostgres establishes a connection to PostgreSQL and returns the connection
func connectToPostgres(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// createTable creates a table in the PostgreSQL database
func createTable(db *sql.DB, tableName string, columns map[string]string) error {
	query := "CREATE TABLE IF NOT EXISTS " + tableName + " ("
	for column, dataType := range columns {
		query += column + " " + dataType + ", "
	}
	// Trim the trailing comma and space
	query = strings.TrimSuffix(query, ", ")
	query += ")"
	_, err := db.Exec(query)
	return err
}

// queryTable queries the table in the PostgreSQL database
func queryTable(db *sql.DB, tableName string) (*sql.Rows, error) {
	query := "SELECT * FROM " + tableName
	return db.Query(query)
}

// dropTable drops the table in the PostgreSQL database
func dropTable(db *sql.DB, tableName string) error {
	query := "DROP TABLE " + tableName
	_, err := db.Exec(query)
	return err
}
