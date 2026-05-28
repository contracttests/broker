package rename_participant

import (
	"github.com/contracttesting/broker/internal/components"
	"github.com/contracttesting/broker/internal/repository"
)

func Register(components *components.Components) {
	repo := repository.NewParticipantRepository(components.Pool)
	handler := NewRenameParticipantHandler(repo)
	components.Server.Post("/api/:participant/rename", handler.Handle)
}
