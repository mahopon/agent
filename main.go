package main

import (
	"agent/agent"
	"agent/llm"
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	titleAgent := agent.NewTitleAgent()
	session := agent.NewSession()
	llmConfig := llm.NewLLMConfig()
	titleAgent.Call("Build flappy bird using HTML, CSS, and Vanilla JS", session, llmConfig)
}
