package main

import (
	"log"

	"github.com/contracttesting/broker/server/internal"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	components := internal.Run()
	if err := components.Server.Listen(":9000"); err != nil {
		log.Fatalf("server: %v", err)
	}
}
