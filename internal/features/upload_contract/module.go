package upload_contract

import (
	"github.com/contracttests/broker/server/internal/components"
	"github.com/contracttests/broker/server/internal/dsl"
	"github.com/contracttests/broker/server/internal/features/upload_contract/wireout"
	"github.com/gofiber/fiber/v3"
)

func Register(components *components.Components) {
	components.Server.Post("/contracts", func(ctx fiber.Ctx) error {
		dslContract := &dsl.Contract{}
		if err := ctx.Bind().JSON(dslContract); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": string(wireout.ContractUploadFailed)})
		}

		_ = dslContract.ToContractModel()

		return ctx.Status(fiber.StatusOK).JSON(wireout.UploadResponse{
			Success: true,
			Message: string(wireout.ContractUploadSuccessful),
		})
	})
}
