package middleware

import (
	"github.com/contracttesting/broker/internal/repository"
	"github.com/gofiber/fiber/v3"
)

// RequireVersion ensures a contract exists for the (participant, version) pair,
// rejecting with 400 when the version is missing or unpublished. It must run
// after RequireParticipant, as it reads the resolved participant.
func RequireVersion(repo *repository.ContractRepository, extract Extractor) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		version := extract(ctx)
		if version == "" {
			return reject(ctx, VersionRequired)
		}

		participant := ParticipantFrom(ctx)
		if participant == nil || !repo.HasContractForVersion(ctx.Context(), participant.ID, version) {
			return reject(ctx, VersionNotFound)
		}

		ctx.Locals(versionKey, version)
		return ctx.Next()
	}
}
