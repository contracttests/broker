package can_i_deploy

import "github.com/contracttesting/broker/internal/compatibility_checker"

const ContractNotFound = "contract not found"

type CanIDeployResponse struct {
	Success    bool                                              `json:"success"`
	Deployable bool                                              `json:"deployable"`
	Breaks     map[string][]compatibility_checker.BreakingChange `json:"breaks,omitempty"`
}

type CanIDeployErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
