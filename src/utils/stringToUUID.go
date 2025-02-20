package utils

import "github.com/google/uuid"

func StringToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func UuidToString(id uuid.UUID) (string, error) {
	return id.String(), nil
}