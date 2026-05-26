package create_environment

import (
	"github.com/contracttesting/broker/server/internal/model"
	"github.com/contracttesting/broker/server/internal/repository"
	"github.com/gofiber/fiber/v3"
)

type createEnvironmentRequest struct {
	Name string `json:"name"`
}

type CreateEnvironmentHandler struct {
	environmentRepository *repository.EnvironmentRepository
}

func NewCreateEnvironmentHandler(repo *repository.EnvironmentRepository) *CreateEnvironmentHandler {
	return &CreateEnvironmentHandler{environmentRepository: repo}
}

func (ctr *CreateEnvironmentHandler) Handle(ctx fiber.Ctx) error {
	request := &createEnvironmentRequest{}

	if err := ctx.Bind().JSON(request); err != nil {
		return ctr.respondInvalidInput(ctx)
	}

	if request.Name == "" {
		return ctr.respondInvalidInput(ctx)
	}

	if ctr.environmentRepository.ExistsByName(ctx.Context(), request.Name) {
		return ctr.respondAlreadyExists(ctx)
	}

	ctr.environmentRepository.Create(ctx.Context(), model.NewEnvironment(request.Name))

	return ctx.Status(fiber.StatusOK).JSON(CreateEnvironmentResponse{
		Success: true,
		Message: EnvironmentCreated,
	})
}

func (ctr *CreateEnvironmentHandler) respondInvalidInput(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(CreateEnvironmentResponse{
		Success: false,
		Message: EnvironmentInvalidInput,
	})
}

func (ctr *CreateEnvironmentHandler) respondAlreadyExists(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(CreateEnvironmentResponse{
		Success: false,
		Message: EnvironmentAlreadyExists,
	})
}
