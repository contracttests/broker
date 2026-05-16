package wireout

type ContractMessage string

const (
    ContractUploadSuccessful ContractMessage = "contract upload successful"
    ContractUploadFailed     ContractMessage = "contract upload failed"
    ContractAlreadyUploaded  ContractMessage = "contract already uploaded"
    ContractInvalidInput     ContractMessage = "contract invalid input"
)

type UploadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}