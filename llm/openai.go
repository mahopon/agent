package llm

import (
	"agent/tool"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
    Tools      []tool.Tool `json:"tools,omitempty"`
    ToolChoice string    `json:"tool_choice,omitempty"`
}

type ChatResponse struct {
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Message Message `json:"message"`
		FinishReason string `json:"finish_reason"` // "stop", "length", "tool_calls", etc.
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type LLMResponse struct {
	Content string
	ToolCalls []ToolCall
	FinishReason string
	TotalTokens int
}

func NewLLMResponse(content, finishReason string, toolCalls []ToolCall, tokensUsed int) *LLMResponse {
	return &LLMResponse{
		Content: content,
		FinishReason: finishReason,
		ToolCalls: toolCalls,
		TotalTokens: tokensUsed,
	}
}

type OpenAILLM struct {
	config LLMConfig
}

func NewOpenAILLM(url, key, model string) *OpenAILLM {
	return &OpenAILLM{
		config: NewLLMConfig(url, key, model),
	}
}

func (llm *OpenAILLM) CallLLM(messages []Message, tools []tool.ToolExecutor) (*LLMResponse, error) {
	reqBody := ChatRequest{
		Model:    llm.config.model,
		Messages: messages,
		Tools: tool.ToOpenAIScheme(tools),
		ToolChoice: "auto",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		llm.config.baseURL+"/v1/chat/completions",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if llm.config.apiKey != "" {
		req.Header.Set("Authorization", "Bearer " + llm.config.apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no response choices")
	}
	llmResp := NewLLMResponse(result.Choices[0].Message.Content, result.Choices[0].FinishReason, result.Choices[0].Message.ToolCalls, result.Usage.TotalTokens)
	return llmResp, nil
}
