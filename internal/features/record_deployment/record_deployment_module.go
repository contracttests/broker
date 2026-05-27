package record_deployment

const (
	DeploymentRecorded     string = "deployment recorded"
	DeploymentInvalidInput string = "deployment invalid input"
	ParticipantNotFound    string = "participant not found"
	VersionNotPublished    string = "version not published"
	EnvironmentNotFound    string = "environment not found"
)

type RecordDeploymentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
