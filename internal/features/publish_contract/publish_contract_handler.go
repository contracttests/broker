package publish_contract

import (
	"encoding/json"
	"strings"

	"github.com/contracttesting/broker/internal/dsl"
	"github.com/contracttesting/broker/internal/middleware"
	"github.com/contracttesting/broker/internal/model"
	"github.com/contracttesting/broker/internal/repository"
	"github.com/gofiber/fiber/v3"
)

type PublishContractController struct {
	contractRepository *repository.ContractRepository
}

func NewPublishContractController(
	contractRepository *repository.ContractRepository,
) *PublishContractController {
	return &PublishContractController{
		contractRepository: contractRepository,
	}
}

type publishContractRequest struct {
	Version  string          `json:"version"`
	Contract json.RawMessage `json:"contract"`
}

func (ctr *PublishContractController) Handle(ctx fiber.Ctx) error {
	request := &publishContractRequest{}
	if err := json.Unmarshal(ctx.Body(), request); err != nil {
		return ctr.respondInvalidInput(ctx)
	}

	version := strings.TrimSpace(request.Version)
	if version == "" || len(request.Contract) == 0 {
		return ctr.respondInvalidInput(ctx)
	}

	dslContract := &dsl.Contract{}
	if err := json.Unmarshal(request.Contract, dslContract); err != nil {
		return ctr.respondInvalidInput(ctx)
	}

	participant := middleware.ParticipantFrom(ctx)

	contract := model.NewContract(participant, version, string(request.Contract))
	dslContract.HydrateContract(contract)

	if existing, found := ctr.contractRepository.LoadChecksumForVersion(ctx.Context(), contract.ParticipantID(), version); found {
		if existing == contract.Checksum() {
			return ctr.respondSuccess(ctx)
		}
		return ctr.respondVersionConflict(ctx)
	}

	ctr.upsert(ctx, contract)

	return ctr.respondSuccess(ctx)
}

func (ctr *PublishContractController) upsert(ctx fiber.Ctx, contract *model.Contract) {
	if ctr.contractRepository.HasContractsForParticipant(ctx.Context(), contract.ParticipantID()) {
		ctr.contractRepository.Update(ctx.Context(), contract)

		return
	}

	ctr.contractRepository.Create(ctx.Context(), contract)
}

func (ctr *PublishContractController) respondInvalidInput(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(PublishContractOutput{
		Success: false,
		Message: ContractInvalidInput,
	})
}

func (ctr *PublishContractController) respondVersionConflict(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusConflict).JSON(PublishContractOutput{
		Success: false,
		Message: ContractVersionConflict,
	})
}

func (ctr *PublishContractController) respondSuccess(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(PublishContractOutput{
		Success: true,
		Message: ContractPublishSuccessful,
	})
}
