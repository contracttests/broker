package internal

import (
	"github.com/contracttesting/broker/server/internal/components"
	"github.com/contracttesting/broker/server/internal/features/create_participant"
	"github.com/contracttesting/broker/server/internal/features/upload_contract"
)

func Run() *components.Components {
	components := components.New()

	create_participant.Register(components)
	upload_contract.Register(components)

	return components
}
