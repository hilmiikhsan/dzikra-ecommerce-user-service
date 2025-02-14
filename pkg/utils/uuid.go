package utils

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/sonitx/uuidv7"
)

func GenerateUUIDv7String() (uuid.UUID, error) {
	// Generate UUID v7 as string
	uuidStr, err := uuidv7.GetUUIDv7String()
	if err != nil {
		log.Warn().Err(err).Msg("utils::GenerateUUIDv7String - Error while generating UUID V7")
		return uuid.UUID{}, fmt.Errorf("cannot generate UUID V7, error: %v", err)
	}

	// Parse string UUID to uuid.UUID
	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		log.Warn().Err(err).Msg("utils::GenerateUUIDv7String - Error while parsing UUID")
		return uuid.UUID{}, fmt.Errorf("error parsing UUID: %v", err)
	}

	return parsedUUID, nil
}
