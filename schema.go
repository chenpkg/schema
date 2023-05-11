package schema

import (
	"context"
	"database/sql"
	"errors"
)

type Config struct {
	DB           *sql.DB // database handle
	Database     string  // current database
	Prefix       string  // database prefix
	Engine       string  // engin, default InnoDB
	Charset      string  // charset, default utf8mb4
	Collation    string  // collation, default utf8mb4_unicode_ci
	StringLength int     // string length, default 255
}

type Schema struct {
	ctx    context.Context
	config *Config
}

// NewSchema new schema
func NewSchema(ctx context.Context, config *Config) *Schema {
	return &Schema{
		ctx:    ctx,
		config: config,
	}
}

// Table Modify a table on the schema.
func (s *Schema) Table(table string, callback func(table *Blueprint)) error {
	return s.build(NewBlueprint(s, table, callback))
}

// Create a new table on the schema.
func (s *Schema) Create(table string, callback func(table *Blueprint)) error {
	return s.build(tap(NewBlueprint(s, table), func(table *Blueprint) {
		table.create()
		callback(table)
	}))
}

// Drop a table from the schema.
func (s *Schema) Drop(table string) error {
	return s.build(tap(NewBlueprint(s, table), func(table *Blueprint) {
		table.drop()
	}))
}

// DropIfExists Drop a table from the schema if it exists.
func (s *Schema) DropIfExists(table string) error {
	return s.build(tap(NewBlueprint(s, table), func(table *Blueprint) {
		table.dropIfExists()
	}))
}

// DropColumns Drop columns from a table schema.
// columns can string or []string
func (s *Schema) DropColumns(table string, columns ...string) error {
	return s.Table(table, func(table *Blueprint) {
		table.DropColumn(columns...)
	})
}

// HasTable check table exists
func (s *Schema) HasTable(table string) (bool, error) {
	if s.config.Database == "" {
		return false, errors.New("schema err: config.Database is empty")
	}

	rows, err := s.config.DB.Query(localGrammar.CompileTableExists(), s.config.Database, s.config.Prefix+table)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		i++
	}

	return i > 0, nil
}

// Rename table
func (s *Schema) Rename(from string, to string) error {
	return s.build(tap(NewBlueprint(s, from), func(table *Blueprint) {
		table.Rename(to)
	}))
}

func (s *Schema) build(blueprint *Blueprint) error {
	return blueprint.build(localGrammar)
}
