package model

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

const (
	Consumes     Direction    = "consumes"
	Provides     Direction    = "provides"
	RestRequest  ResourceKind = "rest_request"
	RestResponse ResourceKind = "rest_response"
)

type Direction string

func (direction *Direction) String() string {
	return string(*direction)
}

type ResourceKind string

func (resourceKind *ResourceKind) String() string {
	return string(*resourceKind)
}

type Resource struct {
	ID          int64
	Direction   Direction
	Kind        ResourceKind
	Provider    string
	Endpoint    string
	Method      string
	StatusCode  string
	Properties  map[string]Property
	Participant *Participant
}

func (resouce *Resource) AddParticipant(participant *Participant) {
	resouce.Participant = participant
}

func (resouce *Resource) ParticipantID() int64 {
	return resouce.Participant.ID
}

func (resouce *Resource) ProviderHash() string {
	providerName := resouce.Provider
	if resouce.Direction == Provides {
		providerName = resouce.ParticipantName()
	}

	parts := []string{providerName, resouce.Endpoint, resouce.Method}
	if resouce.Kind == RestResponse {
		parts = append(parts, resouce.StatusCode)
	}

	return hashParts(parts)
}

func (resouce *Resource) ConsumerHash() string {
	if resouce.Direction != Consumes {
		return ""
	}

	parts := []string{resouce.ParticipantName(), resouce.Provider, resouce.Endpoint, resouce.Method}
	if resouce.Kind == RestResponse {
		parts = append(parts, resouce.StatusCode)
	}

	return hashParts(parts)
}

func (resouce *Resource) ParticipantName() string {
	return resouce.Participant.Name
}

func (resouce *Resource) CanonicalKey() string {
	propertyKeys := make([]string, 0, len(resouce.Properties))

	for _, property := range resouce.Properties {
		propertyKeys = append(propertyKeys, property.CanonicalKey())
	}

	sort.Strings(propertyKeys)

	return strings.Join([]string{
		string(resouce.Direction),
		string(resouce.Kind),
		resouce.ParticipantName(),
		resouce.Provider,
		resouce.Endpoint,
		resouce.Method,
		resouce.StatusCode,
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
) *Resource {
	return &Resource{
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
) *Resource {
	return &Resource{
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
) *Resource {
	return &Resource{
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
) *Resource {
	return &Resource{
		Direction:  Provides,
		Kind:       RestResponse,
		Endpoint:   endpoint,
		Method:     method,
		StatusCode: statusCode,
		Properties: properties,
	}
}
