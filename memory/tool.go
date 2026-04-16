package memory

import "github.com/google/uuid"

type ToolMessage struct {
	id uuid.UUID
	Role string `json:"role"`
	ToolCallId string `json:"tool_call_id"`
	Content map[string]any `json:"content"`
}

func NewToolMessage(toolCallId string) ToolMessage {
	id, _ := uuid.NewUUID()
	return ToolMessage{
		id: id,
		Role: "assistant",
		ToolCallId: toolCallId,
		Content: nil,
	}
}

func NewToolCallResponse(toolCallId string, content map[string]any) ToolMessage {
	id, _ := uuid.NewUUID()
	return ToolMessage{
		id: id,
		Role: "tool",
		ToolCallId: toolCallId,
		Content: content,
	}
}
