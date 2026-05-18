package wireout

type ContractMessage string

const (
    ContractUploadSuccessful ContractMessage = "contract upload successful"
    ContractInvalidInput     ContractMessage = "contract invalid input"
)

type UploadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}