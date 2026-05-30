package rename_participant

import (
	"github.com/contracttesting/broker/internal/middleware"
	"github.com/contracttesting/broker/internal/repository"
	"github.com/gofiber/fiber/v3"
)

type renameParticipantRequest struct {
	NewName string `json:"newName"`
}

type RenameParticipantHandler struct {
	participantRepository *repository.ParticipantRepository
}

func NewRenameParticipantHandler(repo *repository.ParticipantRepository) *RenameParticipantHandler {
	return &RenameParticipantHandler{participantRepository: repo}
}

func (h *RenameParticipantHandler) Handle(ctx fiber.Ctx) error {
	oldName := middleware.ParticipantFrom(ctx).Name

	request := &renameParticipantRequest{}
	if err := ctx.Bind().JSON(request); err != nil {
		return h.respondInvalidInput(ctx)
	}

	if request.NewName == "" {
		return h.respondInvalidInput(ctx)
	}

	if _, conflict := h.participantRepository.Rename(ctx.Context(), oldName, request.NewName); conflict {
		return h.respondAlreadyExists(ctx)
	}

	return ctx.Status(fiber.StatusOK).JSON(RenameParticipantResponse{
		Success: true,
		Message: ParticipantRenamed,
	})
}

func (h *RenameParticipantHandler) respondInvalidInput(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(RenameParticipantResponse{
		Success: false,
		Message: ParticipantInvalidInput,
	})
}

func (h *RenameParticipantHandler) respondAlreadyExists(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(RenameParticipantResponse{
		Success: false,
		Message: ParticipantAlreadyExists,
	})
}
