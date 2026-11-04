package util

import (
	"github.com/google/uuid"
	"strings"
)

// NextID uses UUIDv7, avoiding per-instance node configuration and startup ordering.
func NextID() string {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return strings.ReplaceAll(id.String(), "-", "")
}
