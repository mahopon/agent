package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"agent/tool"
)

type LocalLLM struct {
	LLM
}

type SSEDecoder struct {
	reader *bufio.Reader
}

func NewLocalLLM(config *LLMConfig) *LocalLLM {
	return &LocalLLM{
		LLM: LLM{
			config: config,
		},
	}
}

func NewDecoder(r io.Reader) *SSEDecoder {
	return &SSEDecoder{
		reader: bufio.NewReader(r),
	}
}

func parseReply(respBody []byte) (content, reasoning, finishReason string, toolCalls []tool.ToolCallInfo, promptTokens, completionTokens, totalTokens int, err error) {
	var llmResp LLMResponse
	if err := json.Unmarshal(respBody, &llmResp); err != nil {
		return "", "", "", nil, 0, 0, 0, err
	}

	if len(llmResp.Choices) == 0 {
		return "", "", "", nil, 0, 0, 0, fmt.Errorf("no choices in LLM response")
	}

	msg := llmResp.Choices[0].Message
	content = msg.Content
	reasoning = msg.ReasoningContent
	finishReason = llmResp.Choices[0].FinishReason

	for _, tc := range msg.ToolCalls {
		toolCalls = append(toolCalls, tool.ToolCallInfo{
			Name:      tc.Function.Name,
			Arguments: tc.Function.Arguments,
			ID:        tc.ID,
		})
	}

	return content, reasoning, finishReason, toolCalls, llmResp.Usage.PromptTokens, llmResp.Usage.CompletionTokens, llmResp.Usage.TotalTokens, nil
}

func (d *SSEDecoder) Decode() (string, error) {
	line, err := d.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "data: ") {
		return "", nil
	}

	data := strings.TrimPrefix(line, "data: ")
	if data == "[DONE]" {
		return "", io.EOF
	}

	var chunk struct {
		Choices []struct {
			Delta struct {
				Content string `json:"content"`
			} `json:"delta"`
		} `json:"choices"`
	}

	if err := json.Unmarshal([]byte(data), &chunk); err != nil {
		return "", err
	}

	if len(chunk.Choices) > 0 {
		return chunk.Choices[0].Delta.Content, nil
	}
	return "", nil
}

func (llm *LLM) Call(body *LLMBody) (*ParsedResponse, error) {
	client := &http.Client{Timeout: 5 * time.Minute}
	llmReq := NewLLMRequest(llm.config.LLM_MODEL, body)
	jsonData, err := json.Marshal(llmReq)
	if err != nil {
		return nil, err
	}
	slog.Debug("LLM Request body", "details", body.ToLogBody(llm.config.LLM_MODEL))
	bodyReader := bytes.NewReader(jsonData)
	req, err := http.NewRequest(http.MethodPost, llm.config.LLM_URL, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+llm.config.LLM_KEY)

	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("HTTP error %d: failed to read response body", resp.StatusCode)
		}
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	inferenceTimeMs := time.Since(startTime).Milliseconds()
	logResp := ParseToLogResponse(respBody)
	logResp.InferenceTimeMs = inferenceTimeMs
	slog.Debug("LLM Response", "status", resp.Status, "details", logResp)

	content, reasoning, finishReason, toolCalls, promptTokens, completionTokens, totalTokens, err := parseReply(respBody)
	if err != nil {
		return nil, err
	}

	if content == "" && len(toolCalls) == 0 {
		slog.Warn("LLM returned empty response with no tool calls", "message_count", len(body.Msgs))
	} else if content == "" && len(toolCalls) > 0 {
		slog.Debug("LLM returned empty content with tool calls", "tool_calls", len(toolCalls), "message_count", len(body.Msgs))
	}

	return &ParsedResponse{
		Content:          content,
		Reasoning:        reasoning,
		FinishReason:     finishReason,
		ToolCalls:        toolCalls,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      totalTokens,
		InferenceTimeMs:  inferenceTimeMs,
		RequestBody:      string(jsonData),
	}, nil
}

func RetryWithResult[T any](
	maxRetries int,
	baseDelay time.Duration,
	fn func() (T, error),
) (T, error) {

	var zero T

	for attempt := range maxRetries {
		result, err := fn()
		if err == nil {
			return result, nil
		}

		delay := time.Duration(1<<attempt) * baseDelay
		time.Sleep(delay)
	}

	return zero, fmt.Errorf("max retries reached")
}

func (llm *LLM) CallWithRetry(body *LLMBody) (*ParsedResponse, error) {
	resp, err := llm.Call(body)
	if err != nil {
		return nil, err
	}

	if resp.FinishReason == "" {
		log.Printf("")
		slog.Warn("LLM returned empty finish reason")
		resp, err = RetryWithResult(3, 2*time.Second, func() (*ParsedResponse, error) { return llm.Call(body) })
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

func (llm *LLM) Stream(body *LLMBody, onChunk func(string)) error {
	client := &http.Client{Timeout: 5 * time.Minute}
	llmReq := NewLLMRequest(llm.config.LLM_MODEL, body)
	llmReq.Stream = true
	jsonData, err := json.Marshal(llmReq)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonData)
	req, err := http.NewRequest(http.MethodPost, llm.config.LLM_URL, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+llm.config.LLM_KEY)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := NewDecoder(resp.Body)
	for {
		token, err := decoder.Decode()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if token != "" {
			onChunk(token)
		}
	}
	return nil
}
