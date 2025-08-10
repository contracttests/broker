package model

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func HashFromStrings(parts ...string) string {
	data := strings.Join(parts, ";")
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
