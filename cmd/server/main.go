package main

import (
	"github.com/contracttests/broker/server/internal"
	"github.com/joho/godotenv"
)


func main() {
	godotenv.Load()
	components := internal.Run()
	components.Server.Listen(":3000")
}
