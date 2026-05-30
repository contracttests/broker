package record_deployment

import (
	"github.com/contracttesting/broker/internal/middleware"
	"github.com/contracttesting/broker/internal/model"
	"github.com/contracttesting/broker/internal/repository"
	"github.com/gofiber/fiber/v3"
)

type RecordDeploymentHandler struct {
	deploymentRepository *repository.DeploymentRepository
}

func NewRecordDeploymentHandler(
	deploymentRepository *repository.DeploymentRepository,
) *RecordDeploymentHandler {
	return &RecordDeploymentHandler{
		deploymentRepository: deploymentRepository,
	}
}

func (h *RecordDeploymentHandler) Handle(ctx fiber.Ctx) error {
	participant := middleware.ParticipantFrom(ctx)
	version := middleware.VersionFrom(ctx)
	environment := middleware.EnvironmentFrom(ctx)

	h.deploymentRepository.Insert(ctx.Context(), model.NewDeployment(participant, version, environment))

	return ctx.Status(fiber.StatusOK).JSON(RecordDeploymentResponse{
		Success: true,
		Message: DeploymentRecorded,
	})
}
