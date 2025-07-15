package firestore_services

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func GenerateID(prefix string) (string, error) {
	// Pseudo-random generator seeded with the current nanosecond timestamp
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)

	// Generate a ULID using the current timestamp and the entropy source
	id, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err != nil {
		return "", err
	}

	// Add a prefix to the generated ID
	return prefix + "_" + id.String(), nil
}
