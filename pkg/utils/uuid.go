package utils

import (
	"errors"

	"github.com/google/uuid"
)

func ParseUUID(idStr string) (uuid.UUID, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid UUID")
	}
	return id, nil
}
