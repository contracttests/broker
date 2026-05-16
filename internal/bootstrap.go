package internal

import (
	"github.com/contracttests/broker/server/internal/components"
	"github.com/contracttests/broker/server/internal/features/upload_contract"
)

func Run() *components.Components {
	components := components.New()
	upload_contract.Register(components)
	return components
}