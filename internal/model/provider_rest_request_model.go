package model

type ProviderRestRequestArgs struct {
	Owner    string
	Endpoint string
	Method   string
}

func NewProviderRestRequest(
	args ProviderRestRequestArgs,
	schema Schema,
) ProviderRequest {
	return ProviderRequest{
		Owner:        args.Owner,
		Schema:       schema,
		ResourceType: ProvidesRestRequest,
		RestRequest: RestRequest{
			Endpoint: args.Endpoint,
			Method:   args.Method,
		},
	}
}
