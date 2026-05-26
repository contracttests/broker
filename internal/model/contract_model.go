package model

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

type Contract struct {
	ID          int64
	Version     string
	RawContract string
	Resources   map[string]Resource
	Participant *Participant
}

func NewContract(participant *Participant, version string, rawContract string) *Contract {
	return &Contract{
		Participant: participant,
		Version:     version,
		RawContract: rawContract,
	}
}

func (contract *Contract) ParticipantID() int64 {
	return contract.Participant.ID
}

func (contract *Contract) AddResource(resource *Resource) {
	if contract.Resources == nil {
		contract.Resources = make(map[string]Resource)
	}

	resource.AddParticipant(contract.Participant)

	contract.Resources[resource.PrimaryHash()] = *resource
}

func (resouce Resource) PrimaryHash() string {
	if resouce.Direction == Provides {
		return resouce.ProviderHash()
	}

	return resouce.ConsumerHash()
}

func (contract *Contract) CanonicalKey() string {
	resourceKeys := make([]string, 0, len(contract.Resources))

	for _, resource := range contract.Resources {
		resourceKeys = append(resourceKeys, resource.CanonicalKey())
	}

	sort.Strings(resourceKeys)

	return strings.Join([]string{
		contract.Participant.Name,
		strings.Join(resourceKeys, ";;"),
	}, ";;")
}

func (contract *Contract) Checksum() string {
	sum := sha256.Sum256([]byte(contract.CanonicalKey()))
	return hex.EncodeToString(sum[:])
}
