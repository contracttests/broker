package model

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func UuidFromStrings(parts ...string) string {
	data := strings.Join(parts, ";")
	hash := sha256.Sum256([]byte(data))
	b := hash[:16]

	b[6] = (b[6] & 0x0f) | 0x50
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4],
		b[4:6],
		b[6:8],
		b[8:10],
		b[10:16])
}
