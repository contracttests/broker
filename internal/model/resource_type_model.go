package model

type ResourceType string

const (
	ConsumesRestRequest  ResourceType = "consumes_rest_request"
	ProvidesRestRequest  ResourceType = "provides_rest_request"
	ConsumesRestResponse ResourceType = "consumes_rest_response"
	ProvidesRestResponse ResourceType = "provides_rest_response"
)
