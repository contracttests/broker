package publish_contract

import (
	"github.com/contracttesting/broker/internal/components"
	"github.com/contracttesting/broker/internal/repository"
)

func Register(components *components.Components) {
	contractRepository := repository.NewContractRepository(components.Pool)
	participantRepository := repository.NewParticipantRepository(components.Pool)
	controller := NewPublishContractController(contractRepository, participantRepository)
	components.Server.Post("/api/:participant/contracts/:version", controller.Handle)
}
