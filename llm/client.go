package llm

import (
	"agent/tool"
	"encoding/json"
)

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

type LLMLogBody struct {
	Model            string `json:"model"`
	MessageCount     int    `json:"message_count"`
	HasTools         bool   `json:"has_tools"`
	ReasoningEnabled bool   `json:"reasoning_enabled"`
}

func (b *LLMBody) ToLogBody(model string) LLMLogBody {
	hasTools := len(b.Tools) > 0
	reasoningEnabled := false
	if b.Reasoning != nil {
		reasoningEnabled = b.Reasoning["enabled"]
	}
	return LLMLogBody{
		Model:            model,
		MessageCount:     len(b.Msgs),
		HasTools:         hasTools,
		ReasoningEnabled: reasoningEnabled,
	}
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
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type LLMLogResponse struct {
	Model            string              `json:"model"`
	Content          string              `json:"content"`
	FinishReason     string              `json:"finish_reason"`
	ToolCalls        []map[string]string `json:"tool_calls"`
	PromptTokens     int                 `json:"prompt_tokens"`
	CompletionTokens int                 `json:"completion_tokens"`
	TotalTokens      int                 `json:"total_tokens"`
	InferenceTimeMs  int64               `json:"inference_time_ms"`
}

func ParseToLogResponse(respBody []byte) LLMLogResponse {
	var llmResp LLMResponse
	if err := json.Unmarshal(respBody, &llmResp); err != nil {
		return LLMLogResponse{Model: "error", Content: err.Error()}
	}

	if len(llmResp.Choices) == 0 {
		return LLMLogResponse{Model: llmResp.Model, Content: ""}
	}

	msg := llmResp.Choices[0].Message
	toolCallInfo := make([]map[string]string, 0, len(msg.ToolCalls))
	for _, tc := range msg.ToolCalls {
		toolCallInfo = append(toolCallInfo, map[string]string{
			"name":      tc.Function.Name,
			"arguments": tc.Function.Arguments,
		})
	}

	return LLMLogResponse{
		Model:            llmResp.Model,
		Content:          msg.Content,
		FinishReason:     llmResp.Choices[0].FinishReason,
		ToolCalls:        toolCallInfo,
		PromptTokens:     llmResp.Usage.PromptTokens,
		CompletionTokens: llmResp.Usage.CompletionTokens,
		TotalTokens:      llmResp.Usage.TotalTokens,
	}
}

type LLMCallable interface {
	Call(body *LLMBody) (*ParsedResponse, error)
	CallWithRetry(body *LLMBody) (*ParsedResponse, error)
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
	Content          string
	Reasoning        string
	FinishReason     string
	ToolCalls        []tool.ToolCallInfo
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	InferenceTimeMs  int64
	RequestBody      string
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
