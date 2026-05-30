package rename_participant

const (
	ParticipantRenamed       string = "participant renamed"
	ParticipantInvalidInput  string = "participant invalid input"
	ParticipantAlreadyExists string = "participant already exists"
)

type RenameParticipantResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
