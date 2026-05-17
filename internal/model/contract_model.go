package model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Contract struct {
	Name      string              `json:"name"`
	Owner     string              `json:"owner"`
	Resources map[string]Resource `json:"resources,omitzero"`
}

func (c *Contract) AddResource(r Resource) {
	if c.Resources == nil {
		c.Resources = make(map[string]Resource)
	}
	c.Resources[r.Key()] = r
}

func (c *Contract) Checksum() string {
	payload, _ := json.Marshal(c)
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}
