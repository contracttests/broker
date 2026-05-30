package middleware

import (
	"github.com/contracttesting/broker/internal/repository"
	"github.com/gofiber/fiber/v3"
)

// RequireParticipant resolves the participant from the request, rejecting with
// 400 when it is missing or unknown, and stashes it for the handler.
func RequireParticipant(repo *repository.ParticipantRepository, extract Extractor) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		name := extract(ctx)
		if name == "" {
			return reject(ctx, ParticipantRequired)
		}

		participant, found := repo.FindByName(ctx.Context(), name)
		if !found {
			return reject(ctx, ParticipantNotFound)
		}

		ctx.Locals(participantKey, participant)
		return ctx.Next()
	}
}
