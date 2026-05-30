package publish_contract

type ContractMessage string

const (
	ContractPublishSuccessful ContractMessage = "contract publish successful"
	ContractInvalidInput      ContractMessage = "contract invalid input"
	ContractVersionConflict   ContractMessage = "contract version already exists with different content"
)

type PublishContractOutput struct {
	Success bool            `json:"success"`
	Message ContractMessage `json:"message"`
}
