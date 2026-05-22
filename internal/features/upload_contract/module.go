package upload_contract

import (
	"github.com/contracttesting/broker/server/internal/components"
	"github.com/contracttesting/broker/server/internal/repository"
)

func Register(components *components.Components) {
	repo := repository.NewContractRepository(components.Pool)
	checker := NewCompatibilityChecker(repo)
	controller := NewUploadContractController(repo, checker)
	components.Server.Post("/contracts", controller.Handle)
}
