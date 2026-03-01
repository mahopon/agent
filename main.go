package main

import (
	"agent/agent"
	"agent/config"
	"agent/llm"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config := config.NewConfig()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	llmConfig := llm.NewLLMConfig(config.LLM_URL, config.API_KEY, config.LLM_MODEL)
	llm := llm.NewLocalLLM(llmConfig)
	titleAgent := agent.NewTitleAgent(llm)
	planningAgent := agent.NewPlanningAgent(llm)
	codingAgent := agent.NewCodingAgent(llm)
	session := agent.NewSession()
	userQuery := "Run flappy bird using HTML, CSS, and Vanilla JS"
	session.History = append(session.History, map[string]string{
		"role":    "user",
		"content": userQuery,
	})
	err := titleAgent.Run(userQuery, session)
	if err != nil {
		panic(err)
	}
	err = planningAgent.Run(userQuery, session)
	if err != nil {
		panic(err)
	}

	err = codingAgent.Run(userQuery, session)
	if err != nil {
		panic(err)
	}

	slog.Debug("Session details", "details", session)
}
