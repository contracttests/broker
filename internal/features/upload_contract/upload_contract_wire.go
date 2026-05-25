package upload_contract

import "encoding/json"

type ContractMessage string

const (
	ContractUploadSuccessful    ContractMessage = "contract upload successful"
	ContractInvalidInput        ContractMessage = "contract invalid input"
	ContractParticipantNotFound ContractMessage = "participant not found"
)

type UploadContractInput struct {
	Version     string          `json:"version,omitempty"`
	Environment string          `json:"environment,omitempty"`
	Participant string          `json:"participant"`
	Contract    json.RawMessage `json:"contract"`
}

type UploadContractOutput struct {
	Success bool            `json:"success"`
	Message ContractMessage `json:"message"`
}
