package model

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"

	"github.com/google/uuid"
)

type ContractInfo struct {
	Name  string
	Owner string
}

type Contract struct {
	ID         int64
	UUID       uuid.UUID
	Name       string
	Owner      string
	RawPayload string
	Resources  map[string]Resource
}

func NewContract(name string, owner string, rawPayload string) *Contract {
	return &Contract{
		Name:       name,
		Owner:      owner,
		RawPayload: rawPayload,
	}
}

func (c *Contract) AddResource(r Resource) string {
	if c.Resources == nil {
		c.Resources = make(map[string]Resource)
	}
	r.ContractInfo = &ContractInfo{Name: c.Name, Owner: c.Owner}
	key := r.PrimaryHash()
	c.Resources[key] = r
	return key
}

func (r Resource) PrimaryHash() string {
	if r.Direction == Provides {
		return r.ProviderHash()
	}

	return r.ConsumerHash()
}

func (c *Contract) CanonicalKey() string {
	resourceKeys := make([]string, 0, len(c.Resources))

	for _, resource := range c.Resources {
		resourceKeys = append(resourceKeys, resource.CanonicalKey())
	}

	sort.Strings(resourceKeys)

	return strings.Join([]string{
		c.Name,
		c.Owner,
		strings.Join(resourceKeys, ";;"),
	}, ";;")
}

func (c *Contract) Checksum() string {
	sum := sha256.Sum256([]byte(c.CanonicalKey()))
	return hex.EncodeToString(sum[:])
}
