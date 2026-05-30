package can_i_deploy

import (
	"github.com/contracttesting/broker/internal/compatibility_checker"
	"github.com/contracttesting/broker/internal/middleware"
	"github.com/contracttesting/broker/internal/model"
	"github.com/contracttesting/broker/internal/repository"
	"github.com/gofiber/fiber/v3"
)

type CanIDeployHandler struct {
	contractRepository            *repository.ContractRepository
	compatibilityMatrixRepository *repository.CompatibilityMatrixRepository
	compatibilityChecker          *compatibility_checker.CompatibilityChecker
}

func NewCanIDeployHandler(
	contractRepository *repository.ContractRepository,
	compatibilityMatrixRepository *repository.CompatibilityMatrixRepository,
	compatibilityChecker *compatibility_checker.CompatibilityChecker,
) *CanIDeployHandler {
	return &CanIDeployHandler{
		contractRepository:            contractRepository,
		compatibilityMatrixRepository: compatibilityMatrixRepository,
		compatibilityChecker:          compatibilityChecker,
	}
}

func (h *CanIDeployHandler) Handle(ctx fiber.Ctx) error {
	deployableParticipant := middleware.ParticipantFrom(ctx)
	environment := middleware.EnvironmentFrom(ctx)
	version := middleware.VersionFrom(ctx)

	contract, found := h.contractRepository.LoadContractByNameAndVersion(ctx.Context(), deployableParticipant.Name, version)
	if !found {
		return h.respondContractNotFound(ctx)
	}

	report := h.compatibilityChecker.Check(ctx.Context(), contract, environment)

	for _, result := range report.Results {
		h.compatibilityMatrixRepository.Insert(ctx.Context(), &model.CompatibilityMatrix{
			ParticipantID:            contract.ParticipantID(),
			Version:                  contract.Version,
			CounterpartParticipantID: result.CounterpartParticipantID,
			CounterpartVersion:       result.CounterpartVersion,
			Deployable:               result.Deployable,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(CanIDeployResponse{
		Success:    true,
		Deployable: len(report.Breaks) == 0,
		Breaks:     report.Breaks,
	})
}

func (h *CanIDeployHandler) respondContractNotFound(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotFound).JSON(CanIDeployErrorResponse{
		Success: false,
		Message: ContractNotFound,
	})
}
