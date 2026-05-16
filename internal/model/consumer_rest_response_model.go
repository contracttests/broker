package model

type ConsumerRestResponseArgs struct {
	Owner      string
	Provider   string
	Endpoint   string
	Method     string
	StatusCode string
}

func NewConsumerRestResponse(
	args ConsumerRestResponseArgs,
	schema Schema,
) ConsumerResponse {
	return ConsumerResponse{
		Owner:        args.Owner,
		Provider:     args.Provider,
		Schema:       schema,
		ResourceType: ConsumesRestResponse,
		RestResponse: RestResponse{
			Endpoint:   args.Endpoint,
			Method:     args.Method,
			StatusCode: args.StatusCode,
		},
	}
}
