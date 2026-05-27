package can_i_deploy

import (
	"strings"

	"github.com/contracttesting/broker/server/internal/repository"
	"github.com/gofiber/fiber/v3"
)

type CanIDeployHandler struct {
	participantRepository         *repository.ParticipantRepository
	contractRepository            *repository.ContractRepository
	environmentRepository         *repository.EnvironmentRepository
	deploymentRepository          *repository.DeploymentRepository
	compatibilityMatrixRepository *repository.CompatibilityMatrixRepository
	clockRepository               *repository.ClockRepository
}

func NewCanIDeployHandler(
	participantRepository *repository.ParticipantRepository,
	contractRepository *repository.ContractRepository,
	environmentRepository *repository.EnvironmentRepository,
	deploymentRepository *repository.DeploymentRepository,
	compatibilityMatrixRepository *repository.CompatibilityMatrixRepository,
	clockRepository *repository.ClockRepository,
) *CanIDeployHandler {
	return &CanIDeployHandler{
		participantRepository:         participantRepository,
		contractRepository:            contractRepository,
		environmentRepository:         environmentRepository,
		deploymentRepository:          deploymentRepository,
		compatibilityMatrixRepository: compatibilityMatrixRepository,
		clockRepository:               clockRepository,
	}
}

type canIDeployRequest struct {
	participantName string
	version         string
	environmentName string
}

func parseRequest(ctx fiber.Ctx) canIDeployRequest {
	return canIDeployRequest{
		participantName: strings.TrimSpace(ctx.Params("participant")),
		version:         strings.TrimSpace(ctx.Query("version")),
		environmentName: strings.TrimSpace(ctx.Query("environment")),
	}
}

func (r canIDeployRequest) hasMissingField() bool {
	return r.participantName == "" || r.version == "" || r.environmentName == ""
}

func (h *CanIDeployHandler) Handle(ctx fiber.Ctx) error {
	request := parseRequest(ctx)
	if request.hasMissingField() {
		return h.respondInvalidInput(ctx)
	}

	asker, found := h.participantRepository.FindByName(ctx.Context(), request.participantName)
	if !found {
		return h.respondParticipantNotFound(ctx)
	}

	if !h.contractRepository.HasContractForVersion(ctx.Context(), asker.ID, request.version) {
		return h.respondVersionNotPublished(ctx)
	}

	environment, found := h.environmentRepository.FindByName(ctx.Context(), request.environmentName)
	if !found {
		return h.respondEnvironmentNotFound(ctx)
	}

	deployable := h.evaluate(ctx, asker.ID, request.version, environment.ID, asker.Name)
	return ctx.Status(fiber.StatusOK).JSON(CanIDeployResponse{
		Success:    true,
		Deployable: deployable,
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

func (h *CanIDeployHandler) respondVersionNotPublished(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(CanIDeployErrorResponse{
		Success: false,
		Message: VersionNotPublished,
	})
}

func (h *CanIDeployHandler) respondEnvironmentNotFound(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(CanIDeployErrorResponse{
		Success: false,
		Message: EnvironmentNotFound,
	})
}
