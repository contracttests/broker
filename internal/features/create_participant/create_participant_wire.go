package create_participant

type ParticipantMessage string

const (
	ParticipantCreated       ParticipantMessage = "participant created"
	ParticipantInvalidInput  ParticipantMessage = "participant invalid input"
	ParticipantAlreadyExists ParticipantMessage = "participant already exists"
)

type CreateParticipantResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
