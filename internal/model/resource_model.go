package model

type Direction string

const (
	Consumes Direction = "consumes"
	Provides Direction = "provides"
)

type ResourceKind string

const (
	RestRequest  ResourceKind = "rest_request"
	RestResponse ResourceKind = "rest_response"
)

type Resource struct {
	Direction  Direction           `json:"direction"`
	Kind       ResourceKind        `json:"kind"`
	Provider   string              `json:"provider,omitzero"`
	Endpoint   string              `json:"endpoint"`
	Method     string              `json:"method"`
	StatusCode string              `json:"statusCode,omitzero"`
	Properties map[string]Property `json:"properties,omitzero"`
}

func (r Resource) Key() string {
	return string(r.Direction) + ";;" + string(r.Kind) + ";;" + r.Provider + ";;" + r.Endpoint + ";;" + r.Method + ";;" + r.StatusCode
}

func NewConsumedRestRequest(provider, endpoint, method string, properties map[string]Property) Resource {
	return Resource{
		Direction:  Consumes,
		Kind:       RestRequest,
		Provider:   provider,
		Endpoint:   endpoint,
		Method:     method,
		Properties: properties,
	}
}

func NewProvidedRestRequest(endpoint, method string, properties map[string]Property) Resource {
	return Resource{
		Direction:  Provides,
		Kind:       RestRequest,
		Endpoint:   endpoint,
		Method:     method,
		Properties: properties,
	}
}

func NewConsumedRestResponse(provider, endpoint, method, statusCode string, properties map[string]Property) Resource {
	return Resource{
		Direction:  Consumes,
		Kind:       RestResponse,
		Provider:   provider,
		Endpoint:   endpoint,
		Method:     method,
		StatusCode: statusCode,
		Properties: properties,
	}
}

func NewProvidedRestResponse(endpoint, method, statusCode string, properties map[string]Property) Resource {
	return Resource{
		Direction:  Provides,
		Kind:       RestResponse,
		Endpoint:   endpoint,
		Method:     method,
		StatusCode: statusCode,
		Properties: properties,
	}
}
