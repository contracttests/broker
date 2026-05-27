package can_i_deploy

const (
	CanIDeployInvalidInput = "can-i-deploy invalid input"
	ParticipantNotFound    = "participant not found"
)

type CanIDeployResponse struct {
	Success    bool `json:"success"`
	Deployable bool `json:"deployable"`
}

type CanIDeployErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
