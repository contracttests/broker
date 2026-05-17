package upload_contract

import (
	"github.com/contracttests/broker/server/internal/components"
	"github.com/contracttests/broker/server/internal/repository"
)

func Register(components *components.Components) {
	repo := repository.NewContractRepository(components.Pool)
	controller := NewUploadContractController(repo)
	components.Server.Post("/contracts", controller.Handle)
}
