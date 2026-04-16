package memory

import (
	"github.com/google/uuid"
)

type Message struct {
	id uuid.UUID
	Role string `json:"Role"`
	Content string `json:"Content"`
}

type MessageContent = map[string]string

func NewUserMessage(Content string) Message {
	uuid, _ := uuid.NewUUID()
	return Message{
		id: uuid,
		Role: "user",
		Content: Content,
	}
}

func NewSystemMessage(Content string) Message {
	uuid, _ := uuid.NewUUID()
	return Message{
		id: uuid,
		Role: "system",
		Content: Content,
	}
}

func NewAssistantMessage(Content string) Message {
	uuid, _ := uuid.NewUUID()
	return Message{
		id: uuid,
		Role: "assistant",
		Content: Content,
	}
}


func (m Message) Id() uuid.UUID {
	return m.id
}

