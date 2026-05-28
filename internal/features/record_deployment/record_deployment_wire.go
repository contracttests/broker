package record_deployment

import (
	"github.com/contracttesting/broker/internal/components"
	"github.com/contracttesting/broker/internal/repository"
)

func Register(components *components.Components) {
	participantRepository := repository.NewParticipantRepository(components.Pool)
	contractRepository := repository.NewContractRepository(components.Pool)
	environmentRepository := repository.NewEnvironmentRepository(components.Pool)
	deploymentRepository := repository.NewDeploymentRepository(components.Pool)

	handler := NewRecordDeploymentHandler(
		participantRepository,
		contractRepository,
		environmentRepository,
		deploymentRepository,
	)

	components.Server.Post("/api/:participant/deployments", handler.Handle)
}
