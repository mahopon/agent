package llm

import "agent/tool"

type LLMRequest struct {
	Model  string `json:"model"`
	Stream bool   `json:"stream,omitempty"`
	LLMBody
}

type LLMBody struct {
	Msgs      []map[string]any `json:"messages"`
	Reasoning map[string]bool  `json:"reasoning,omitempty"`
	Tools     []tool.Tool      `json:"tools,omitempty"`
}

type LLMChoice struct {
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
	Message      struct {
		Role             string     `json:"role"`
		Content          string     `json:"content"`
		ReasoningContent string     `json:"reasoning_content"`
		ToolCalls        []ToolCall `json:"tool_calls"`
	} `json:"message"`
}

type LLMResponse struct {
	Choices []LLMChoice `json:"choices"`
	Model   string      `json:"model"`
}

type LLMCallable interface {
	Call(body *LLMBody) (*ParsedResponse, error)
	Stream(body *LLMBody, onChunk func(string)) error
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type ParsedResponse struct {
	Content   string
	Reasoning string
	ToolCalls []tool.ToolCallInfo
}

func (p *ParsedResponse) HasToolCalls() bool {
	return len(p.ToolCalls) > 0
}

func (p *ParsedResponse) NumToolCalls() int {
	return len(p.ToolCalls)
}

type LLM struct {
	config *LLMConfig
}

func NewLLMBody(messages []map[string]any, enableReasoning bool, tools []tool.Tool) *LLMBody {
	return &LLMBody{
		Msgs: messages,
		Reasoning: map[string]bool{
			"enabled": enableReasoning,
		},
		Tools: tools,
	}
}

func NewLLMRequest(model string, body *LLMBody) *LLMRequest {
	return &LLMRequest{
		Model:   model,
		LLMBody: *body,
	}
}
