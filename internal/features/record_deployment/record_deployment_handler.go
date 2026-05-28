package record_deployment

import (
	"strings"

	"github.com/contracttesting/broker/internal/model"
	"github.com/contracttesting/broker/internal/repository"
	"github.com/gofiber/fiber/v3"
)

type recordDeploymentRequest struct {
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

type RecordDeploymentHandler struct {
	participantRepository *repository.ParticipantRepository
	contractRepository    *repository.ContractRepository
	environmentRepository *repository.EnvironmentRepository
	deploymentRepository  *repository.DeploymentRepository
}

func NewRecordDeploymentHandler(
	participantRepository *repository.ParticipantRepository,
	contractRepository *repository.ContractRepository,
	environmentRepository *repository.EnvironmentRepository,
	deploymentRepository *repository.DeploymentRepository,
) *RecordDeploymentHandler {
	return &RecordDeploymentHandler{
		participantRepository: participantRepository,
		contractRepository:    contractRepository,
		environmentRepository: environmentRepository,
		deploymentRepository:  deploymentRepository,
	}
}

func (h *RecordDeploymentHandler) Handle(ctx fiber.Ctx) error {
	participantName := strings.TrimSpace(ctx.Params("participant"))
	if participantName == "" {
		return h.respondInvalidInput(ctx)
	}

	request := &recordDeploymentRequest{}
	if err := ctx.Bind().JSON(request); err != nil {
		return h.respondInvalidInput(ctx)
	}

	if request.Version == "" || request.Environment == "" {
		return h.respondInvalidInput(ctx)
	}

	participant, ok := h.participantRepository.FindByName(ctx.Context(), participantName)
	if !ok {
		return h.respondParticipantNotFound(ctx)
	}

	if !h.contractRepository.HasContractForVersion(ctx.Context(), participant.ID, request.Version) {
		return h.respondVersionNotPublished(ctx)
	}

	environment, ok := h.environmentRepository.FindByName(ctx.Context(), request.Environment)
	if !ok {
		return h.respondEnvironmentNotFound(ctx)
	}

	h.deploymentRepository.Insert(ctx.Context(), model.NewDeployment(participant, request.Version, environment))

	return ctx.Status(fiber.StatusOK).JSON(RecordDeploymentResponse{
		Success: true,
		Message: DeploymentRecorded,
	})
}

func (h *RecordDeploymentHandler) respondInvalidInput(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(RecordDeploymentResponse{
		Success: false,
		Message: DeploymentInvalidInput,
	})
}

func (h *RecordDeploymentHandler) respondParticipantNotFound(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotFound).JSON(RecordDeploymentResponse{
		Success: false,
		Message: ParticipantNotFound,
	})
}

func (h *RecordDeploymentHandler) respondVersionNotPublished(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(RecordDeploymentResponse{
		Success: false,
		Message: VersionNotPublished,
	})
}

func (h *RecordDeploymentHandler) respondEnvironmentNotFound(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(RecordDeploymentResponse{
		Success: false,
		Message: EnvironmentNotFound,
	})
}
