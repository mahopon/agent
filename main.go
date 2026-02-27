package main

import (
	"agent/agent"
	"agent/llm"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

func main() {
	godotenv.Load()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	llmConfig := llm.NewLLMConfig("", "")
	llm := llm.NewLocalLLM(llmConfig)
	titleAgent := agent.NewTitleAgent(llm)
	planningAgent := agent.NewPlanningAgent(llm)
	session := agent.NewSession()
	err := titleAgent.Run("Build flappy bird using HTML, CSS, and Vanilla JS", session)
	if err != nil {
		panic(err)
	}
	err = planningAgent.Run("Build flappy bird using HTML, CSS, and Vanilla JS", session)
	if err != nil {
		panic(err)
	}

	slog.Debug("Session details", "details", session)
}
