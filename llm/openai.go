package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type LocalLLM struct {
	LLM
}

func NewLocalLLM(config *LLMConfig) *LocalLLM {
	return &LocalLLM{
		LLM: LLM{
			config: config,
		},
	}
}

func (llm *LLM) Call(body *LLMBody) (string, error) {
	client := &http.Client{Timeout: 5 * time.Minute}
	llmReq := NewLLMRequest("z.ai/glm-4.7-flash", body)
	jsonData, err := json.Marshal(llmReq)
	if err != nil {
		return "", err
	}
	bodyReader := bytes.NewReader(jsonData)
	req, err := http.NewRequest(http.MethodPost, llm.config.LLM_URL, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+llm.config.LLM_KEY)
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
