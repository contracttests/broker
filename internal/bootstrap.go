package internal

import (
	"github.com/contracttesting/broker/server/internal/components"
	"github.com/contracttesting/broker/server/internal/features/create_environment"
	"github.com/contracttesting/broker/server/internal/features/create_participant"
	"github.com/contracttesting/broker/server/internal/features/publish_contract"
)

func Run() *components.Components {
	components := components.New()

	create_participant.Register(components)
	create_environment.Register(components)
	publish_contract.Register(components)

	return components
}
