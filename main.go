package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const (
	API_KEY string = ""
	LLM_URL string = "https://openrouter.ai/api/v1/chat/completions"
)

func main() {
	params := LLMCallBody{
		Model: "z-ai/glm-4.7-flash",
		Msgs: []map[string]string{
			{
				"role":    "system",
				"content": "You are an AI Assistant",
			},
		},
		Reasoning: map[string]bool{
			"enabled": false,
		},
	}
	_, err := LLMCall("How are you?", params)
	check(err)
}

func check(e error) {
	if e != nil {
		slog.Error(e.Error())
	}
}

type LLMCallBody struct {
	Model     string              `json:"model"`
	Msgs      []map[string]string `json:"messages"`
	Reasoning map[string]bool     `json:"reasoning"`
}

func LLMCall(content string, body LLMCallBody) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	body.Msgs = append(body.Msgs, map[string]string{
		"role":    "user",
		"content": content,
	})
	jsonData, err := json.Marshal(body)
	check(err)
	bodyReader := bytes.NewReader(jsonData)
	req, err := http.NewRequest(http.MethodPost, LLM_URL, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+API_KEY)
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response:", string(respBody))
	return string(respBody), nil
}
