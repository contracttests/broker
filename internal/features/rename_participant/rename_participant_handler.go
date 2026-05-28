package rename_participant

import (
	"strings"

	"github.com/contracttesting/broker/internal/repository"
	"github.com/gofiber/fiber/v3"
)

type renameParticipantRequest struct {
	Name string `json:"name"`
}

type RenameParticipantHandler struct {
	participantRepository *repository.ParticipantRepository
}

func NewRenameParticipantHandler(repo *repository.ParticipantRepository) *RenameParticipantHandler {
	return &RenameParticipantHandler{participantRepository: repo}
}

func (h *RenameParticipantHandler) Handle(ctx fiber.Ctx) error {
	oldName := strings.TrimSpace(ctx.Params("participant"))
	if oldName == "" {
		return h.respondInvalidInput(ctx)
	}

	request := &renameParticipantRequest{}
	if err := ctx.Bind().JSON(request); err != nil {
		return h.respondInvalidInput(ctx)
	}

	if request.Name == "" {
		return h.respondInvalidInput(ctx)
	}

	found, conflict := h.participantRepository.Rename(ctx.Context(), oldName, request.Name)
	if conflict {
		return h.respondAlreadyExists(ctx)
	}
	if !found {
		return h.respondNotFound(ctx)
	}

	return ctx.Status(fiber.StatusOK).JSON(RenameParticipantResponse{
		Success: true,
		Message: string(ParticipantRenamed),
	})
}

func (h *RenameParticipantHandler) respondInvalidInput(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(RenameParticipantResponse{
		Success: false,
		Message: string(ParticipantInvalidInput),
	})
}

func (h *RenameParticipantHandler) respondAlreadyExists(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(RenameParticipantResponse{
		Success: false,
		Message: string(ParticipantAlreadyExists),
	})
}

func (h *RenameParticipantHandler) respondNotFound(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotFound).JSON(RenameParticipantResponse{
		Success: false,
		Message: string(ParticipantNotFound),
	})
}
