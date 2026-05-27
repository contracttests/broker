package can_i_deploy

import (
	"strings"

	"github.com/contracttesting/broker/server/internal/model"
	"github.com/contracttesting/broker/server/internal/repository"
	"github.com/gofiber/fiber/v3"
)

type CanIDeployHandler struct {
	participantRepository         *repository.ParticipantRepository
	compatibilityMatrixRepository *repository.CompatibilityMatrixRepository
}

func NewCanIDeployHandler(
	participantRepository *repository.ParticipantRepository,
	compatibilityMatrixRepository *repository.CompatibilityMatrixRepository,
) *CanIDeployHandler {
	return &CanIDeployHandler{
		participantRepository:         participantRepository,
		compatibilityMatrixRepository: compatibilityMatrixRepository,
	}
}

func (h *CanIDeployHandler) Handle(ctx fiber.Ctx) error {
	participantName := strings.TrimSpace(ctx.Params("participant"))
	version := strings.TrimSpace(ctx.Query("version"))
	if participantName == "" || version == "" {
		return h.respondInvalidInput(ctx)
	}

	asker, found := h.participantRepository.FindByName(ctx.Context(), participantName)
	if !found {
		return h.respondParticipantNotFound(ctx)
	}

	// Hardcoded verdict until the real compatibility check lands; proves persistence end-to-end.
	row := &model.CompatibilityMatrixRow{
		ParticipantID: asker.ID,
		Version:       version,
		Deployable:    true,
	}

	h.compatibilityMatrixRepository.Insert(ctx.Context(), row)

	return ctx.Status(fiber.StatusOK).JSON(CanIDeployResponse{
		Success:    true,
		Deployable: row.Deployable,
	})
}

func (h *CanIDeployHandler) respondInvalidInput(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(CanIDeployErrorResponse{
		Success: false,
		Message: CanIDeployInvalidInput,
	})
}

func (h *CanIDeployHandler) respondParticipantNotFound(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotFound).JSON(CanIDeployErrorResponse{
		Success: false,
		Message: ParticipantNotFound,
	})
}
