package model

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

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
	ID           int64
	Direction    Direction
	Kind         ResourceKind
	Provider     string
	Endpoint     string
	Method       string
	StatusCode   string
	Properties   map[string]Property
	ContractInfo *ContractInfo
}

func (r Resource) ProviderHash() string {
	providerName := r.Provider
	if r.Direction == Provides {
		providerName = r.ContractName()
	}

	parts := []string{providerName, r.Endpoint, r.Method}
	if r.Kind == RestResponse {
		parts = append(parts, r.StatusCode)
	}

	return hashParts(parts)
}

func (r Resource) ConsumerHash() string {
	if r.Direction != Consumes {
		return ""
	}
	parts := []string{r.ContractName(), r.Provider, r.Endpoint, r.Method}
	if r.Kind == RestResponse {
		parts = append(parts, r.StatusCode)
	}
	return hashParts(parts)
}

func (r Resource) ContractName() string {
	if r.ContractInfo == nil {
		return ""
	}

	return r.ContractInfo.Name
}

func (r Resource) CanonicalKey() string {
	propertyKeys := make([]string, 0, len(r.Properties))

	for _, property := range r.Properties {
		propertyKeys = append(propertyKeys, property.CanonicalKey())
	}

	sort.Strings(propertyKeys)

	return strings.Join([]string{
		string(r.Direction),
		string(r.Kind),
		r.Provider,
		r.Endpoint,
		r.Method,
		r.StatusCode,
		strings.Join(propertyKeys, ";;"),
	}, ";;")
}

func hashParts(parts []string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, ";;")))
	return hex.EncodeToString(sum[:])
}

func NewConsumedRestRequest(
	provider, endpoint, method string,
	properties map[string]Property,
) Resource {
	return Resource{
		Direction:  Consumes,
		Kind:       RestRequest,
		Provider:   provider,
		Endpoint:   endpoint,
		Method:     method,
		Properties: properties,
	}
}

func NewProvidedRestRequest(
	endpoint, method string,
	properties map[string]Property,
) Resource {
	return Resource{
		Direction:  Provides,
		Kind:       RestRequest,
		Endpoint:   endpoint,
		Method:     method,
		Properties: properties,
	}
}

func NewConsumedRestResponse(
	provider, endpoint, method, statusCode string,
	properties map[string]Property,
) Resource {
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

func NewProvidedRestResponse(
	endpoint, method, statusCode string,
	properties map[string]Property,
) Resource {
	return Resource{
		Direction:  Provides,
		Kind:       RestResponse,
		Endpoint:   endpoint,
		Method:     method,
		StatusCode: statusCode,
		Properties: properties,
	}
}
