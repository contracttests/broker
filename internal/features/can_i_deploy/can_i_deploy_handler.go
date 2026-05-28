package can_i_deploy

import (
	"strings"

	"github.com/contracttesting/broker/internal/compatibility_checker"
	"github.com/contracttesting/broker/internal/model"
	"github.com/contracttesting/broker/internal/repository"
	"github.com/gofiber/fiber/v3"
)

type CanIDeployHandler struct {
	participantRepository         *repository.ParticipantRepository
	contractRepository            *repository.ContractRepository
	compatibilityMatrixRepository *repository.CompatibilityMatrixRepository
	environmentRepository         *repository.EnvironmentRepository
	compatibilityChecker          *compatibility_checker.CompatibilityChecker
}

func NewCanIDeployHandler(
	participantRepository *repository.ParticipantRepository,
	contractRepository *repository.ContractRepository,
	compatibilityMatrixRepository *repository.CompatibilityMatrixRepository,
	environmentRepository *repository.EnvironmentRepository,
	compatibilityChecker *compatibility_checker.CompatibilityChecker,
) *CanIDeployHandler {
	return &CanIDeployHandler{
		participantRepository:         participantRepository,
		contractRepository:            contractRepository,
		compatibilityMatrixRepository: compatibilityMatrixRepository,
		environmentRepository:         environmentRepository,
		compatibilityChecker:          compatibilityChecker,
	}
}

func (h *CanIDeployHandler) Handle(ctx fiber.Ctx) error {
	participantName := strings.TrimSpace(ctx.Params("participant"))
	environmentName := strings.TrimSpace(ctx.Query("environment"))
	version := strings.TrimSpace(ctx.Query("version"))

	if participantName == "" || version == "" || environmentName == "" {
		return h.respondInvalidInput(ctx)
	}

	deployableParticipant, found := h.participantRepository.FindByName(ctx.Context(), participantName)
	if !found {
		return h.respondParticipantNotFound(ctx)
	}

	environment, found := h.environmentRepository.FindByName(ctx.Context(), environmentName)
	if !found {
		return h.respondEnvironmentNotFound(ctx)
	}

	matrixItem := &model.CompatibilityMatrix{
		ParticipantID: deployableParticipant.ID,
		Version:       version,
		Deployable:    true,
	}

	contract, found := h.contractRepository.LoadContractByNameAndVersion(ctx.Context(), deployableParticipant.Name, version)
	if !found {
		return h.respondContractNotFound(ctx)
	}

	report := h.compatibilityChecker.Check(ctx.Context(), contract, environment)

	if len(report.Breaks) > 0 {
		matrixItem.Deployable = false
	}

	if matrixItem.Deployable {
		h.compatibilityMatrixRepository.Insert(ctx.Context(), matrixItem)
	}

	return ctx.Status(fiber.StatusOK).JSON(CanIDeployResponse{
		Success:    true,
		Deployable: matrixItem.Deployable,
		Breaks:     report.Breaks,
	})
}

func (h *CanIDeployHandler) respondEnvironmentNotFound(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(CanIDeployErrorResponse{
		Success: false,
		Message: EnvironmentNotFound,
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

func (h *CanIDeployHandler) respondContractNotFound(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotFound).JSON(CanIDeployErrorResponse{
		Success: false,
		Message: ContractNotFound,
	})
}
