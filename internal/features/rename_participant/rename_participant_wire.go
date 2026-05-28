package rename_participant

type ParticipantMessage string

const (
	ParticipantRenamed       ParticipantMessage = "participant renamed"
	ParticipantInvalidInput  ParticipantMessage = "participant invalid input"
	ParticipantAlreadyExists ParticipantMessage = "participant already exists"
	ParticipantNotFound      ParticipantMessage = "participant not found"
)

type RenameParticipantResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
