package components

import "github.com/gofiber/fiber/v3"

type Components struct {
	Server *fiber.App
}

func New() *Components {
	server := fiber.New()

	return &Components{
		Server: server,
	}
}