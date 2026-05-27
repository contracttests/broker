package internal

import (
	"github.com/contracttesting/broker/server/internal/components"
	"github.com/contracttesting/broker/server/internal/features/can_i_deploy"
	"github.com/contracttesting/broker/server/internal/features/create_environment"
	"github.com/contracttesting/broker/server/internal/features/create_participant"
	"github.com/contracttesting/broker/server/internal/features/publish_contract"
	"github.com/contracttesting/broker/server/internal/features/record_deployment"
)

func Run() *components.Components {
	components := components.New()

	create_participant.Register(components)
	create_environment.Register(components)
	publish_contract.Register(components)
	record_deployment.Register(components)
	can_i_deploy.Register(components)

	return components
}
