package upload_contract

import (
	"encoding/json"

	"github.com/contracttesting/broker/server/internal/dsl"
	"github.com/contracttesting/broker/server/internal/model"
	"github.com/contracttesting/broker/server/internal/repository"
	"github.com/gofiber/fiber/v3"
)

type UploadContractController struct {
	contractRepository    *repository.ContractRepository
	participantRepository *repository.ParticipantRepository
}

func NewUploadContractController(
	contractRepository *repository.ContractRepository,
	participantRepository *repository.ParticipantRepository,
) *UploadContractController {
	return &UploadContractController{
		contractRepository:    contractRepository,
		participantRepository: participantRepository,
	}
}

func (ctr *UploadContractController) Handle(ctx fiber.Ctx) error {
	input := &UploadContractInput{}

	if err := ctx.Bind().JSON(input); err != nil {
		return ctr.respondInvalidInput(ctx)
	}

	if len(input.Contract) == 0 || input.Participant == "" {
		return ctr.respondInvalidInput(ctx)
	}

	participant, ok := ctr.participantRepository.FindByName(ctx.Context(), input.Participant)
	if !ok {
		return ctr.respondParticipantNotFound(ctx)
	}

	dslContract := &dsl.Contract{}
	if err := json.Unmarshal(input.Contract, dslContract); err != nil {
		return ctr.respondInvalidInput(ctx)
	}

	contract := model.NewContract(participant, string(input.Contract))
	dslContract.HydrateContract(contract)
	ctr.upsert(ctx, contract)

	return ctr.respondSuccess(ctx)
}

func (ctr *UploadContractController) upsert(ctx fiber.Ctx, contract *model.Contract) {
	if ctr.contractRepository.HasContractsForParticipant(ctx.Context(), contract.ParticipantID()) {
		ctr.contractRepository.Update(ctx.Context(), contract)

		return
	}

	ctr.contractRepository.Create(ctx.Context(), contract)
}

func (ctr *UploadContractController) respondInvalidInput(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(UploadContractOutput{
		Success: false,
		Message: ContractInvalidInput,
	})
}

func (ctr *UploadContractController) respondParticipantNotFound(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(UploadContractOutput{
		Success: false,
		Message: ContractParticipantNotFound,
	})
}

func (ctr *UploadContractController) respondSuccess(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(UploadContractOutput{
		Success: true,
		Message: ContractUploadSuccessful,
	})
}
