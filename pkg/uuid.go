package pkg

import (
	"crypto/sha1"

	"github.com/google/uuid"
)

func ParseOrGenerateUUID(input string) string {
	id, err := uuid.Parse(input)
	if err == nil {
		return id.String()
	}

	hash := sha1.New()
	hash.Write([]byte(input))
	hashedBytes := hash.Sum(nil)

	var uuidBytes [16]byte
	copy(uuidBytes[:], hashedBytes)

	uuidBytes[6] = (uuidBytes[6] & 0x0f) | 0x50
	uuidBytes[8] = (uuidBytes[8] & 0x3f) | 0x80

	return uuid.UUID(uuidBytes).String()
}
