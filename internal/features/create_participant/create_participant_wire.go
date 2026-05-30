package create_participant

const (
	ParticipantCreated       string = "participant created"
	ParticipantInvalidInput  string = "participant invalid input"
	ParticipantAlreadyExists string = "participant already exists"
)

type CreateParticipantResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
