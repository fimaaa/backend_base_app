package util

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateID
// TODO replace this with uuid or gonanoid later
func GenerateID() string {
	// simple id by time
	return time.Now().Format("060102150405")
}

func GenerateIDFourChar() string {
	// simple id by time
	return time.Now().Format("0405")
}

func GenerateUUID() string {
	id := uuid.New()
	return (id.String())
}

func GenerateUuidWithoutDash() string {
	id := uuid.New().String()
	uuidWithoutHyphens := strings.Replace(id, "-", "", -1)
	return (uuidWithoutHyphens)
}
