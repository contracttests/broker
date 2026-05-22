package main

import (
	"github.com/contracttesting/broker/server/internal"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	components := internal.Run()
	components.Server.Listen(":3000")
}
