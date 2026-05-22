package upload_contract

import (
	"github.com/contracttesting/broker/server/internal/dsl"
	"github.com/contracttesting/broker/server/internal/features/upload_contract/wireout"
	"github.com/contracttesting/broker/server/internal/model"
	"github.com/contracttesting/broker/server/internal/repository"
	"github.com/gofiber/fiber/v3"
)

type UploadContractController struct {
	contractRepository   *repository.ContractRepository
	compatibilityChecker *CompatibilityChecker
}

func NewUploadContractController(
	repo *repository.ContractRepository,
	checker *CompatibilityChecker,
) *UploadContractController {
	return &UploadContractController{contractRepository: repo, compatibilityChecker: checker}
}

func (ctr *UploadContractController) Handle(ctx fiber.Ctx) error {
	dslContract := &dsl.Contract{}

	if err := ctx.Bind().JSON(dslContract); err != nil {
		return ctr.respondInvalidInput(ctx)
	}

	contract := dslContract.ToContractModel()
	contract.RawPayload = string(ctx.Body())

	if contract.Name == "" {
		return ctr.respondInvalidInput(ctx)
	}

	report := ctr.compatibilityChecker.Run(ctx.Context(), &contract)

	if report.HasBreaks() {
		return ctr.respondIncompatible(ctx, report)
	}

	ctr.persist(ctx, &contract)

	return ctr.respondSuccess(ctx)
}

func (ctr *UploadContractController) persist(ctx fiber.Ctx, contract *model.Contract) {
	if ctr.contractRepository.ExistsByName(ctx.Context(), contract.Name) {
		ctr.contractRepository.Update(ctx.Context(), contract)

		return
	}

	ctr.contractRepository.Save(ctx.Context(), contract)
}

func (ctr *UploadContractController) respondInvalidInput(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(wireout.UploadResponse{
		Success: false,
		Message: string(wireout.ContractInvalidInput),
	})
}

func (ctr *UploadContractController) respondIncompatible(
	ctx fiber.Ctx,
	report *model.CompatibilityReport,
) error {
	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(wireout.NewBreakingChangesResponse(report))
}

func (ctr *UploadContractController) respondSuccess(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(wireout.UploadResponse{
		Success: true,
		Message: string(wireout.ContractUploadSuccessful),
	})
}
