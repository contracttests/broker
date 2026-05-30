package rename_participant

import (
	"github.com/contracttesting/broker/internal/components"
	"github.com/contracttesting/broker/internal/middleware"
	"github.com/contracttesting/broker/internal/repository"
)

func Register(components *components.Components) {
	repo := repository.NewParticipantRepository(components.Pool)
	handler := NewRenameParticipantHandler(repo)
	components.Server.Post("/api/participants/rename",
		middleware.RequireParticipant(repo, middleware.FromBody("name")),
		handler.Handle,
	)
}
