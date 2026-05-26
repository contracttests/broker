package create_environment

const (
	EnvironmentCreated       string = "environment created"
	EnvironmentInvalidInput  string = "environment invalid input"
	EnvironmentAlreadyExists string = "environment already exists"
)

type CreateEnvironmentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
