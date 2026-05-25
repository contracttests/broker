package model

type BreakingReason string

const (
	ReasonProviderNotSpecified                 BreakingReason = "provider_not_specified"
	ReasonProviderNotFound                     BreakingReason = "provider_not_found"
	ReasonProviderResourceNotFound             BreakingReason = "provider_resource_not_found"
	ReasonMissingInProvider                    BreakingReason = "missing_in_provider"
	ReasonMissingInConsumer                    BreakingReason = "missing_in_consumer"
	ReasonTypeMismatch                         BreakingReason = "type_mismatch"
	ReasonOptionalInProviderRequiredInConsumer BreakingReason = "optional_in_provider_required_in_consumer"
	ReasonOptionalInConsumerRequiredInProvider BreakingReason = "optional_in_consumer_required_in_provider"
)

type UploaderRole string

const (
	UploaderProvider UploaderRole = "provider"
	UploaderConsumer UploaderRole = "consumer"
)

type CounterpartInfo struct {
	Role UploaderRole
	Name string
}

type BrokenResource struct {
	Kind       ResourceKind
	Provider   string
	Endpoint   string
	Method     string
	StatusCode string
}

func NewBrokenResource(resource Resource) BrokenResource {
	return BrokenResource{
		Kind:       resource.Kind,
		Provider:   resource.Provider,
		Endpoint:   resource.Endpoint,
		Method:     resource.Method,
		StatusCode: resource.StatusCode,
	}
}

type BreakingChange struct {
	UploaderRole UploaderRole
	Counterpart  *CounterpartInfo
	Resource     BrokenResource
	Reason       BreakingReason
	Property     string
	ProviderType string
	ConsumerType string
}
