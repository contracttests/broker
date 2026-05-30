package record_deployment

const DeploymentRecorded string = "deployment recorded"

type RecordDeploymentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
