package integration_test

import (
	"context"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/contracttests/broker/server/internal"
	"github.com/contracttests/broker/server/internal/components"
	"github.com/contracttests/broker/server/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Suite struct {
	suite.Suite
	Pool              *pgxpool.Pool
	Repo              *repository.ContractRepository
	Components        *components.Components
	PostgresContainer *postgres.PostgresContainer
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (suite *Suite) StartPostgressContainer() *postgres.PostgresContainer {
	ctx := context.Background()

	postgresContainer, err := postgres.Run(
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
		panic(fmt.Errorf("Failed to run postgres container: %v", err))
	}

	connectionString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(fmt.Errorf("Failed to get postgres connection string: %v", err))
	}

	os.Setenv("DATABASE_URL", connectionString)

	return postgresContainer
}

func (suite *Suite) SetupTest() {
	godotenv.Load()
	os.Setenv("MIGRATIONS_DIR", "../../migrations")
	suite.PostgresContainer = suite.StartPostgressContainer()
	suite.Components = internal.Run()
	suite.Pool = suite.Components.Pool
	suite.Repo = repository.NewContractRepository(suite.Pool)
}

func (suite *Suite) TearDownTest() {
	if suite.PostgresContainer != nil {
		if err := suite.PostgresContainer.Terminate(context.Background()); err != nil {
			suite.T().Fatalf("Failed to terminate postgres container: %v", err)
		}
	}
}

type Request struct {
	Method  string
	Path    string
	Body    string
	Headers map[string]string
}

type Response struct {
	StatusCode int
	Body       string
}

func (suite *Suite) Request(args Request) (*Response, error) {
	request := httptest.NewRequest(args.Method, args.Path, strings.NewReader(args.Body))

	for key, value := range args.Headers {
		request.Header.Set(key, value)
	}

	response, err := suite.Components.Server.Test(request)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: response.StatusCode,
		Body:       string(body),
	}, nil
}
