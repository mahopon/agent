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
	planningAgent := agent.NewPlanningAgent()
	session := agent.NewSession()
	llmConfig := llm.NewLLMConfig()
	reply, err := titleAgent.Call("Build flappy bird using HTML, CSS, and Vanilla JS", nil, llmConfig)
	if err != nil {
		panic(err)
	}
	session.Summary = reply
	reply2, err := planningAgent.Call("Build flappy bird using HTML, CSS, and Vanilla JS", nil, llmConfig)
	if err != nil {
		panic(err)
	}
	session.Instructions = reply2

	slog.Debug("Session details", "details", session)
}
