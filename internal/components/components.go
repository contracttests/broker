package components

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Components struct {
	Server *fiber.App
	Pool *pgxpool.Pool
}

func createDatabasePool() *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(fmt.Errorf("Failed to create database pool: %v", err))
	}

	if err := pool.Ping(context.Background()); err != nil {
		panic(fmt.Errorf("Failed to ping database: %v", err))
	}

	return pool
}

func createHttpServer() *fiber.App {
	server := fiber.New()
	return server
}

func New() *Components {
	pool := createDatabasePool()
	server := createHttpServer()

	return &Components{
		Server: server,
		Pool: pool,
	}
}
