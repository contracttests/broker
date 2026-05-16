package model

type ProviderRestResponseArgs struct {
	Owner      string
	Endpoint   string
	Method     string
	StatusCode string
}

func NewProviderRestResponse(
	args ProviderRestResponseArgs,
	schema Schema,
) ProviderResponse {
	return ProviderResponse{
		Owner:        args.Owner,
		Schema:       schema,
		ResourceType: ProvidesRestResponse,
		RestResponse: RestResponse{
			Endpoint:   args.Endpoint,
			Method:     args.Method,
			StatusCode: args.StatusCode,
		},
	}
}
