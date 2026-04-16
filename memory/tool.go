package memory

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ToolReplyMessage struct {
	id uuid.UUID
	Role string `json:"role"`
	ToolCallId string `json:"tool_call_id"`
	Content map[string]any `json:"content"`
}

type ToolCallMessage struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments map[string]any `json:"arguments"`
	} `json:"function"`
}

func NewToolReplyMessage(toolCallId string) *ToolReplyMessage {
	id, _ := uuid.NewUUID()
	return &ToolReplyMessage{
		id: id,
		Role: "tool",
		ToolCallId: toolCallId,
		Content: nil,
	}
}


func NewToolCallMessage(id, func_type, name, args string) *ToolCallMessage {
	var content map[string]any
	_ = json.Unmarshal([]byte(args), &content)

	return &ToolCallMessage{
		ID: id,
		Type: func_type,
		Function: struct{Name string "json:\"name\""; Arguments map[string]any "json:\"arguments\""}{Name: name, Arguments: content},
	}
}
