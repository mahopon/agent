package main

import (
	"agent/agent/v2"
	"agent/config"
	"agent/llm"
	"agent/tool"
	"log/slog"
	"os"

	templates "agent/prompt/templates"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config := config.NewConfig()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	llmConfig := llm.NewLLMConfig(config.LLM_URL, config.API_KEY, config.LLM_MODEL)
	llm := llm.NewLocalLLM(llmConfig)
	sysPrompt := templates.NewSystemPrompt()
	tools := make([]tool.ToolExecutor, 0)
	tools = append(tools, tool.NewFileSystemExecutor())
	orchestratorAgent := agent.NewOrchestratorAgent(sysPrompt, tools, llm)
	session := agent.NewSession()
	userQuery := "Dekete hello.txt and Create a text file in the current directory with the text 'Hello World!'"
	response, err := orchestratorAgent.Run(userQuery, session)
	if err != nil {
		panic(err)
	}

	for response.HasToolCalls() {
		response, err = orchestratorAgent.Run("", session)
		if err != nil {
			panic(err)
		}
	}

	slog.Debug("Session details", "details", session)
}
