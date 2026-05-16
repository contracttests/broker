package model

type ConsumerRestRequestArgs struct {
	Owner    string
	Provider string
	Endpoint string
	Method   string
}

func NewConsumerRestRequest(
	args ConsumerRestRequestArgs,
	schema Schema,
) ConsumerRequest {
	return ConsumerRequest{
		Owner:        args.Owner,
		Provider:     args.Provider,
		Schema:       schema,
		ResourceType: ConsumesRestRequest,
		RestRequest: RestRequest{
			Endpoint: args.Endpoint,
			Method:   args.Method,
		},
	}
}
