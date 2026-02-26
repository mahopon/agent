package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"
)

type LLMCallBody struct {
	Model     string              `json:"model"`
	Msgs      []map[string]string `json:"messages"`
	Reasoning map[string]bool     `json:"reasoning"`
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

func LLMCall(content string, body LLMCallBody, config LLMConfig) (string, error) {
	client := &http.Client{Timeout: 2 * time.Minute}
	body.Msgs = append(body.Msgs, map[string]string{
		"role":    "user",
		"content": content,
	})
	log.Printf("Messages: %v", body.Msgs)
	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	bodyReader := bytes.NewReader(jsonData)
	req, err := http.NewRequest(http.MethodPost, config.LLM_URL, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.LLM_KEY)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	var llmResp LLMResponse
	err = json.Unmarshal(respBody, &llmResp)
	if err != nil {
		return "", err
	}
	slog.Debug("LLM Response", "status", resp.Status, "response", string(respBody))
	if len(llmResp.Choices) > 0 {
		return llmResp.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no choices in LLM response")
}
