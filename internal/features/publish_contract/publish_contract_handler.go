package publish_contract

import (
	"encoding/json"
	"strings"

	"github.com/contracttesting/broker/server/internal/dsl"
	"github.com/contracttesting/broker/server/internal/model"
	"github.com/contracttesting/broker/server/internal/repository"
	"github.com/gofiber/fiber/v3"
)

type PublishContractController struct {
	contractRepository    *repository.ContractRepository
	participantRepository *repository.ParticipantRepository
}

func NewPublishContractController(
	contractRepository *repository.ContractRepository,
	participantRepository *repository.ParticipantRepository,
) *PublishContractController {
	return &PublishContractController{
		contractRepository:    contractRepository,
		participantRepository: participantRepository,
	}
}

func (ctr *PublishContractController) Handle(ctx fiber.Ctx) error {
	participantName := strings.TrimSpace(ctx.Params("participant"))
	version := strings.TrimSpace(ctx.Params("version"))

	if participantName == "" || version == "" {
		return ctr.respondInvalidInput(ctx)
	}

	body := ctx.Body()
	if len(body) == 0 {
		return ctr.respondInvalidInput(ctx)
	}

	dslContract := &dsl.Contract{}
	if err := json.Unmarshal(body, dslContract); err != nil {
		return ctr.respondInvalidInput(ctx)
	}

	participant, ok := ctr.participantRepository.FindByName(ctx.Context(), participantName)
	if !ok {
		return ctr.respondParticipantNotFound(ctx)
	}

	contract := model.NewContract(participant, version, string(body))
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

func (ctr *PublishContractController) respondParticipantNotFound(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(PublishContractOutput{
		Success: false,
		Message: ContractParticipantNotFound,
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
