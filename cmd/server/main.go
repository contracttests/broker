package main

import "github.com/contracttests/broker/server/internal"


func main() {
	components := internal.Run()
	components.Server.Listen(":3000")
}