package migrator_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/contracttests/broker/server/pkg/migrator"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type MigratorSuite struct {
	suite.Suite
	container *postgres.PostgresContainer
	pool      *pgxpool.Pool
	dir       string
}

func TestMigratorSuite(t *testing.T) {
	suite.Run(t, new(MigratorSuite))
}

func (s *MigratorSuite) SetupTest() {
	ctx := context.Background()

	container, err := postgres.Run(
		ctx, "postgres:16.6-alpine",
		postgres.WithDatabase("contracttests"),
		postgres.WithUsername("contracttests"),
		postgres.WithPassword("s3cr3t"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		s.T().Fatalf("Failed to run postgres container: %v", err)
	}
	s.container = container

	connectionString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		s.T().Fatalf("Failed to get postgres connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		s.T().Fatalf("Failed to create database pool: %v", err)
	}
	s.pool = pool

	s.dir = s.T().TempDir()
}

func (s *MigratorSuite) TearDownTest() {
	if s.pool != nil {
		s.pool.Close()
	}
	if s.container != nil {
		_ = s.container.Terminate(context.Background())
	}
}

func (s *MigratorSuite) writeMigration(name, sql string) {
	path := filepath.Join(s.dir, name)
	if err := os.WriteFile(path, []byte(sql), 0644); err != nil {
		s.T().Fatalf("Failed to write migration file: %v", err)
	}
}

func (s *MigratorSuite) countMigrations() int {
	var count int
	if err := s.pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM public.schema_migrations").Scan(&count); err != nil {
		s.T().Fatalf("Failed to count migrations: %v", err)
	}
	return count
}

func (s *MigratorSuite) tableExists(name string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1)`
	if err := s.pool.QueryRow(context.Background(), query, name).Scan(&exists); err != nil {
		s.T().Fatalf("Failed to check table existence: %v", err)
	}
	return exists
}

func (s *MigratorSuite) TestAppliesAllPendingMigrations() {
	s.writeMigration("0001_create_foo.sql", "CREATE TABLE foo (id SERIAL PRIMARY KEY);")
	s.writeMigration("0002_create_bar.sql", "CREATE TABLE bar (id SERIAL PRIMARY KEY);")

	m := migrator.New(s.pool, s.dir, "public.schema_migrations")
	if err := m.Migrate(); err != nil {
		s.T().Fatalf("Migrate returned error: %v", err)
	}

	s.Equal(2, s.countMigrations())
	s.True(s.tableExists("foo"))
	s.True(s.tableExists("bar"))
}

func (s *MigratorSuite) TestIsIdempotent() {
	s.writeMigration("0001_create_foo.sql", "CREATE TABLE foo (id SERIAL PRIMARY KEY);")

	m := migrator.New(s.pool, s.dir, "public.schema_migrations")
	s.NoError(m.Migrate())
	s.NoError(m.Migrate())

	s.Equal(1, s.countMigrations())
}

func (s *MigratorSuite) TestPicksUpNewlyAddedMigrations() {
	s.writeMigration("0001_create_foo.sql", "CREATE TABLE foo (id SERIAL PRIMARY KEY);")

	m := migrator.New(s.pool, s.dir, "public.schema_migrations")
	s.NoError(m.Migrate())
	s.Equal(1, s.countMigrations())

	s.writeMigration("0002_create_bar.sql", "CREATE TABLE bar (id SERIAL PRIMARY KEY);")
	s.NoError(m.Migrate())

	s.Equal(2, s.countMigrations())
	s.True(s.tableExists("bar"))
}

func (s *MigratorSuite) TestEnsuresMigrationsTable() {
	m := migrator.New(s.pool, s.dir, "public.schema_migrations")
	s.NoError(m.Migrate())

	s.True(s.tableExists("schema_migrations"))
	s.Equal(0, s.countMigrations())
}

func (s *MigratorSuite) TestSkipsNonSqlAndDirectories() {
	s.writeMigration("0001_create_foo.sql", "CREATE TABLE foo (id SERIAL PRIMARY KEY);")
	if err := os.WriteFile(filepath.Join(s.dir, "README.md"), []byte("ignore me"), 0644); err != nil {
		s.T().Fatalf("Failed to write README: %v", err)
	}
	if err := os.Mkdir(filepath.Join(s.dir, "subdir"), 0755); err != nil {
		s.T().Fatalf("Failed to create subdir: %v", err)
	}

	m := migrator.New(s.pool, s.dir, "public.schema_migrations")
	s.NoError(m.Migrate())

	s.Equal(1, s.countMigrations())
}

func (s *MigratorSuite) TestReturnsErrorOnInvalidSQL() {
	s.writeMigration("0001_broken.sql", "THIS IS NOT VALID SQL;")

	m := migrator.New(s.pool, s.dir, "public.schema_migrations")
	err := m.Migrate()

	s.Error(err)
	s.Equal(0, s.countMigrations())
}
