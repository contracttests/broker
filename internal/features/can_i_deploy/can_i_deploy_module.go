package can_i_deploy

import (
	"github.com/contracttesting/broker/internal/compatibility_checker"
	"github.com/contracttesting/broker/internal/components"
	"github.com/contracttesting/broker/internal/repository"
)

func Register(components *components.Components) {
	contractRepository := repository.NewContractRepository(components.Pool)
	compatibilityChecker := compatibility_checker.NewCompatibilityChecker(contractRepository)
	participantRepository := repository.NewParticipantRepository(components.Pool)
	compatibilityMatrixRepository := repository.NewCompatibilityMatrixRepository(components.Pool)
	environmentRepository := repository.NewEnvironmentRepository(components.Pool)

	handler := NewCanIDeployHandler(
		participantRepository,
		contractRepository,
		compatibilityMatrixRepository,
		environmentRepository,
		compatibilityChecker,
	)

	components.Server.Get("/api/:participant/can-i-deploy", handler.Handle)
}
