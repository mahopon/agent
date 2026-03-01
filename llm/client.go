package llm

type LLMRequest struct {
	Model string `json:"model"`
	LLMBody
}

type LLMBody struct {
	Msgs      []map[string]string `json:"messages"`
	Reasoning map[string]bool     `json:"reasoning,omitempty"`
}

type LLMChoice struct {
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
	Message      struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}

type LLMResponse struct {
	Choices []LLMChoice `json:"choices"`
	Model   string      `json:"model"`
}

type LLMCallable interface {
	Call(body *LLMBody) (string, error)
}

type LLM struct {
	config *LLMConfig
}

func NewLLMBody(messages []map[string]string, enableReasoning bool) *LLMBody {
	return &LLMBody{
		Msgs: messages,
		Reasoning: map[string]bool{
			"enabled": enableReasoning,
		},
	}
}

func NewLLMRequest(model string, body *LLMBody) *LLMRequest {
	return &LLMRequest{
		Model:   model,
		LLMBody: *body,
	}
}
