package migrator

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const ensureMigrationsTableSQL = `
CREATE SCHEMA IF NOT EXISTS %s;
CREATE TABLE IF NOT EXISTS %s (
    id SERIAL PRIMARY KEY,
    migration VARCHAR(255) NOT NULL UNIQUE,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const selectMigrationSQL = `SELECT COUNT(*) FROM %s WHERE migration = $1;`

const insertMigrationSQL = `INSERT INTO %s (migration) VALUES ($1);`

type Migrator struct {
	pool            *pgxpool.Pool
	migrationsDir   string
	migrationsTable string
}

func New(pool *pgxpool.Pool, migrationsDir, migrationsTable string) *Migrator {
	return &Migrator{
		pool:            pool,
		migrationsDir:   migrationsDir,
		migrationsTable: migrationsTable,
	}
}

func (m *Migrator) Migrate() error {
	if err := m.ensureMigrationsTable(); err != nil {
		return fmt.Errorf("error ensuring migrations table: %w", err)
	}

	migrationFiles, err := m.getMigrationFiles()
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %w", err)
	}

	for _, migrationFile := range migrationFiles {
		applied, err := m.isMigrationApplied(migrationFile)
		if err != nil {
			return fmt.Errorf("error checking if migration %s is applied: %w", migrationFile, err)
		}

		if applied {
			continue
		}

		if err := m.applyMigration(migrationFile); err != nil {
			return err
		}

		log.Printf("migrated %s", migrationFile)
	}

	return nil
}

func (m *Migrator) ensureMigrationsTable() error {
	schema := strings.SplitN(m.migrationsTable, ".", 2)[0]
	sql := fmt.Sprintf(ensureMigrationsTableSQL, schema, m.migrationsTable)

	if _, err := m.pool.Exec(context.Background(), sql); err != nil {
		return err
	}

	return nil
}

func (m *Migrator) getMigrationFiles() ([]string, error) {
	entries, err := os.ReadDir(m.migrationsDir)
	if err != nil {
		return nil, err
	}

	var migrationFiles []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		migrationFiles = append(migrationFiles, filepath.Join(m.migrationsDir, entry.Name()))
	}

	return migrationFiles, nil
}

func (m *Migrator) isMigrationApplied(migrationFile string) (bool, error) {
	sql := fmt.Sprintf(selectMigrationSQL, m.migrationsTable)

	var count int
	if err := m.pool.QueryRow(context.Background(), sql, migrationFile).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m *Migrator) applyMigration(migrationFile string) error {
	contents, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("error reading migration file %s: %w", migrationFile, err)
	}

	if _, err := m.pool.Exec(context.Background(), string(contents)); err != nil {
		return fmt.Errorf("error executing migration %s: %w", migrationFile, err)
	}

	sql := fmt.Sprintf(insertMigrationSQL, m.migrationsTable)
	if _, err := m.pool.Exec(context.Background(), sql, migrationFile); err != nil {
		return fmt.Errorf("error recording migration %s: %w", migrationFile, err)
	}

	return nil
}
