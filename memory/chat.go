package memory

import "github.com/google/uuid"

type ChatMessage interface {
	Id() uuid.UUID
}
