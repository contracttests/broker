package can_i_deploy

import (
	"github.com/contracttesting/broker/server/internal/components"
	"github.com/contracttesting/broker/server/internal/repository"
)

func Register(components *components.Components) {
	participantRepository := repository.NewParticipantRepository(components.Pool)
	contractRepository := repository.NewContractRepository(components.Pool)
	environmentRepository := repository.NewEnvironmentRepository(components.Pool)
	deploymentRepository := repository.NewDeploymentRepository(components.Pool)
	compatibilityMatrixRepository := repository.NewCompatibilityMatrixRepository(components.Pool)
	clockRepository := repository.NewClockRepository(components.Pool)

	handler := NewCanIDeployHandler(
		participantRepository,
		contractRepository,
		environmentRepository,
		deploymentRepository,
		compatibilityMatrixRepository,
		clockRepository,
	)

	components.Server.Get("/api/:participant/can-i-deploy", handler.Handle)
}
