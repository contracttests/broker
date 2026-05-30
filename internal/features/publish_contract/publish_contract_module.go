package publish_contract

import (
	"github.com/contracttesting/broker/internal/components"
	"github.com/contracttesting/broker/internal/middleware"
	"github.com/contracttesting/broker/internal/repository"
)

func Register(components *components.Components) {
	contractRepository := repository.NewContractRepository(components.Pool)
	participantRepository := repository.NewParticipantRepository(components.Pool)
	controller := NewPublishContractController(contractRepository)
	components.Server.Post("/api/contracts",
		middleware.RequireParticipant(participantRepository, middleware.FromBody("name")),
		controller.Handle,
	)
}
