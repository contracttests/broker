package wireout

import (
	"github.com/contracttesting/broker/server/internal/model"
)

type ContractMessage string

const (
	ContractUploadSuccessful ContractMessage = "contract upload successful"
	ContractInvalidInput     ContractMessage = "contract invalid input"
	ContractIncompatible     ContractMessage = "contract incompatible with stored counterparts"
)

type UploadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type BreakingChangesResponse struct {
	Success         bool                 `json:"success"`
	Message         string               `json:"message"`
	BreakingChanges []BreakingChangeItem `json:"breakingChanges"`
}

type BreakingChangeItem struct {
	ContractName  string         `json:"contractName"`
	ContractOwner string         `json:"contractOwner"`
	Resource      BrokenResource `json:"resource"`
	Property      string         `json:"property,omitempty"`
	Reason        string         `json:"reason"`
	ExpectedType  string         `json:"expectedType,omitempty"`
	ActualType    string         `json:"actualType,omitempty"`
}

type BrokenResource struct {
	Direction  string `json:"direction"`
	Kind       string `json:"kind"`
	Provider   string `json:"provider"`
	Endpoint   string `json:"endpoint"`
	Method     string `json:"method"`
	StatusCode string `json:"statusCode"`
}

func NewBreakingChangesResponse(report *model.CompatibilityReport) BreakingChangesResponse {
	return BreakingChangesResponse{
		Success:         false,
		Message:         string(ContractIncompatible),
		BreakingChanges: breakingChangeItemsFrom(report),
	}
}

func breakingChangeItemsFrom(report *model.CompatibilityReport) []BreakingChangeItem {
	breaks := report.Canonical()
	items := make([]BreakingChangeItem, 0, len(breaks))
	for _, breakingChange := range breaks {
		items = append(items, newBreakingChangeItem(breakingChange))
	}
	return items
}

func newBreakingChangeItem(breakChange model.BreakingChange) BreakingChangeItem {
	return BreakingChangeItem{
		ContractName:  breakChange.ContractInfo.Name,
		ContractOwner: breakChange.ContractInfo.Owner,
		Resource: BrokenResource{
			Direction:  string(breakChange.Resource.Direction),
			Kind:       string(breakChange.Resource.Kind),
			Provider:   breakChange.Resource.Provider,
			Endpoint:   breakChange.Resource.Endpoint,
			Method:     breakChange.Resource.Method,
			StatusCode: breakChange.Resource.StatusCode,
		},
		Property:     breakChange.Property,
		Reason:       string(breakChange.Reason),
		ExpectedType: breakChange.ExpectedType,
		ActualType:   breakChange.ActualType,
	}
}
