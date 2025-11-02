package utils

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
)

// GenerateULID generates a new ULID as string
func GenerateULID() string {
	entropy := ulid.DefaultEntropy()
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}

// GenerateULIDObject generates a new ULID object
func GenerateULIDObject() ulid.ULID {
	entropy := ulid.DefaultEntropy()
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
}

// ParseULID parses a string into a ULID object
func ParseULID(ulidStr string) (ulid.ULID, error) {
	if ulidStr == "" {
		return ulid.ULID{}, errors.New("ULID string cannot be empty")
	}
	
	parsedULID, err := ulid.Parse(ulidStr)
	if err != nil {
		return ulid.ULID{}, errors.New("invalid ULID format")
	}
	
	return parsedULID, nil
}

// ULIDToString converts a ULID object to string
func ULIDToString(id ulid.ULID) string {
	return id.String()
}

// IsValidULID checks if a string is a valid ULID
func IsValidULID(ulidStr string) bool {
	_, err := ulid.Parse(ulidStr)
	return err == nil
}

// IsEmptyULID checks if a ULID is empty/zero
func IsEmptyULID(id ulid.ULID) bool {
	return id == ulid.ULID{}
}

// SetTimestamps sets CreatedAt and UpdatedAt to current time
func SetTimestamps() time.Time {
	return time.Now()
}