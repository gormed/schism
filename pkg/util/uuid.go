package util

import "github.com/google/uuid"

// IsValidUUID Checks for a valid uuid
func IsValidUUID(id string) (uuid.UUID, error) {
	uuid, err := uuid.Parse(id)
	return uuid, err
}
