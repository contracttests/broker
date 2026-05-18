package upload_contract

import (
	"github.com/contracttests/broker/server/internal/dsl"
	"github.com/contracttests/broker/server/internal/features/upload_contract/wireout"
	"github.com/contracttests/broker/server/internal/repository"
	"github.com/gofiber/fiber/v3"
)

type UploadContractController struct {
	repo *repository.ContractRepository
}

func NewUploadContractController(repo *repository.ContractRepository) *UploadContractController {
	return &UploadContractController{repo: repo}
}

func (ctr *UploadContractController) Handle(ctx fiber.Ctx) error {
	dslContract := &dsl.Contract{}
	if err := ctx.Bind().JSON(dslContract); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(wireout.UploadResponse{
			Success: false,
			Message: string(wireout.ContractInvalidInput),
		})
	}

	contract := dslContract.ToContractModel()
	if contract.Name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(wireout.UploadResponse{
			Success: false,
			Message: string(wireout.ContractInvalidInput),
		})
	}

	if existing := ctr.repo.FindByName(ctx.Context(), contract.Name); existing != nil {
		ctr.repo.Update(ctx.Context(), &contract, ctx.Body())
	} else {
		ctr.repo.Save(ctx.Context(), &contract, ctx.Body())
	}

	return ctx.Status(fiber.StatusOK).JSON(wireout.UploadResponse{
		Success: true,
		Message: string(wireout.ContractUploadSuccessful),
	})
}
