package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bugrakocabay/konsume/pkg/config"

	_ "github.com/lib/pq"
)

type PostgresPlugin struct {
	db *sql.DB
}

// Connect establishes a connection to the PostgreSQL database
func (p *PostgresPlugin) Connect(connectionString, dbName string) error {
	slog.Info("Connecting to PostgreSQL database")
	var err error
	p.db, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	if err = p.db.Ping(); err != nil {
		return err
	}
	slog.Info("Connected to the PostgreSQL database")
	return nil
}

// Insert stores data into the PostgreSQL database
func (p *PostgresPlugin) Insert(data map[string]interface{}, dbRouteConfig config.DatabaseRouteConfig) error {
	columns := []string{}
	placeholders := []string{}
	values := []interface{}{}

	i := 1
	for key, value := range data {
		dbColumn, ok := dbRouteConfig.Mapping[key]
		if !ok {
			slog.Warn("No mapping found for", "key", key)
			continue
		}
		columns = append(columns, dbColumn)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		values = append(values, value)
		i++
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		dbRouteConfig.Table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	stmt, err := p.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		return fmt.Errorf("error executing insert statement: %w", err)
	}

	slog.Info("Inserted data into the database", "table", dbRouteConfig.Table, "values", values)
	return nil
}

// Close terminates the database connection
func (p *PostgresPlugin) Close() error {
	if p.db != nil {
		slog.Info("Closing the PostgreSQL database connection")
		return p.db.Close()
	}
	return nil
}

var Plugin PostgresPlugin
