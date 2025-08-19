package dsl

import (
	"fmt"
	"strings"

	"github.com/contracttests/broker/internal/model"
)

type ResourcePath string

func NewResourcePath(resourcePath string) ResourcePath {
	return ResourcePath(resourcePath)
}

func (f *ResourcePath) Append(parts ...string) ResourcePath {
	separator := ";"

	if string(*f) == "" {
		return ResourcePath(strings.Join(parts, separator))
	}

	return ResourcePath(strings.Join([]string{string(*f), strings.Join(parts, separator)}, separator))
}

func (f *ResourcePath) String() string {
	return string(*f)
}

func (f *ResourcePath) Chunks() []string {
	return strings.Split(f.String(), ";")
}

func (f *ResourcePath) IsConsumer() bool {
	return strings.Contains(f.String(), "consumes")
}

func (f *ResourcePath) IsProvider() bool {
	return strings.Contains(f.String(), "provides")
}

func (f *ResourcePath) IsRequest() bool {
	return strings.Contains(f.String(), "request")
}

func (f *ResourcePath) IsResponse() bool {
	return strings.Contains(f.String(), "response")
}

func (f *ResourcePath) IsRequestConsumer() bool {
	return f.IsConsumer() && f.IsRequest()
}

func (f *ResourcePath) IsResponseConsumer() bool {
	return f.IsConsumer() && f.IsResponse()
}

func (f *ResourcePath) IsRequestProvider() bool {
	return f.IsProvider() && f.IsRequest()
}

func (f *ResourcePath) IsResponseProvider() bool {
	return f.IsProvider() && f.IsResponse()
}

func NewResource(resourcePath ResourcePath, schema model.Schema) model.Resource {
	chunks := resourcePath.Chunks()

	if resourcePath.IsRequestConsumer() {
		consumerName := chunks[0]
		providerName := chunks[2]
		endpoint := chunks[4]
		method := chunks[5]

		return model.NewConsumerRequestBody(consumerName, providerName, endpoint, method, schema)
	}

	if resourcePath.IsResponseConsumer() {
		consumerName := chunks[0]
		providerName := chunks[2]
		endpoint := chunks[4]
		method := chunks[5]
		statusCode := chunks[7]

		return model.NewConsumerResponse(consumerName, providerName, endpoint, method, statusCode, schema)
	}

	if resourcePath.IsRequestProvider() {
		providerName := chunks[0]
		endpoint := chunks[3]
		method := chunks[4]

		return model.NewProviderRequestBody(providerName, endpoint, method, schema)
	}

	if resourcePath.IsResponseProvider() {
		providerName := chunks[0]
		endpoint := chunks[3]
		method := chunks[4]
		statusCode := chunks[6]

		return model.NewProviderResponse(providerName, endpoint, method, statusCode, schema)
	}

	panic(fmt.Sprintf("Invalid resource path: %s", resourcePath))
}
