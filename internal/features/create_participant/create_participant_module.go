package create_participant

import (
	"github.com/contracttesting/broker/server/internal/components"
	"github.com/contracttesting/broker/server/internal/repository"
)

func Register(components *components.Components) {
	repo := repository.NewParticipantRepository(components.Pool)
	handler := NewCreateParticipantHandler(repo)
	components.Server.Post("/api/participants", handler.Handle)
}
