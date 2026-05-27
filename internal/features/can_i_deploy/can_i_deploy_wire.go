package can_i_deploy

import (
	"github.com/contracttesting/broker/server/internal/components"
	"github.com/contracttesting/broker/server/internal/repository"
)

func Register(components *components.Components) {
	participantRepository := repository.NewParticipantRepository(components.Pool)
	compatibilityMatrixRepository := repository.NewCompatibilityMatrixRepository(components.Pool)

	handler := NewCanIDeployHandler(participantRepository, compatibilityMatrixRepository)

	components.Server.Get("/api/:participant/can-i-deploy", handler.Handle)
}
