package middleware

import (
	"github.com/contracttesting/broker/internal/repository"
	"github.com/gofiber/fiber/v3"
)

func RequireEnvironment(repo *repository.EnvironmentRepository, extract Extractor) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		name := extract(ctx)
		if name == "" {
			return reject(ctx, EnvironmentRequired)
		}

		environment, found := repo.FindByName(ctx.Context(), name)
		if !found {
			return reject(ctx, EnvironmentNotFound)
		}

		ctx.Locals(environmentKey, environment)
		return ctx.Next()
	}
}
