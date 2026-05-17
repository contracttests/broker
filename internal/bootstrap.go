package internal

import (
	"fmt"
	"os"

	"github.com/contracttests/broker/server/internal/components"
	"github.com/contracttests/broker/server/internal/features/upload_contract"
	"github.com/contracttests/broker/server/pkg/migrator"
)

func Run() *components.Components {
	components := components.New()
	runMigrations(components)
	upload_contract.Register(components)
	return components
}

func runMigrations(c *components.Components) {
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}

	m := migrator.New(c.Pool, migrationsDir, "public.schema_migrations")
	if err := m.Migrate(); err != nil {
		panic(fmt.Errorf("Failed to run migrations: %v", err))
	}
}
