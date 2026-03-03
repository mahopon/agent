package main

import (
	"agent/agent/v2"
	"agent/config"
	"agent/llm"
	"agent/tool"
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"agent/prompt/templates"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config := config.NewConfig()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
	llmConfig := llm.NewLLMConfig(config.LLM_URL, config.API_KEY, config.LLM_MODEL)
	llm := llm.NewLocalLLM(llmConfig)
	sysPrompt := templates.NewSystemPrompt()
	tools := make([]tool.ToolExecutor, 0)
	tools = append(tools, tool.NewFileSystemExecutor())
	orchestratorAgent := agent.NewOrchestratorAgent(sysPrompt, tools, llm)
	session := agent.NewSession(nil)

	for {
		fmt.Print("Enter a query: ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			break
		}
		userQuery := scanner.Text()
		if userQuery == "STOP" {
			break
		}
		response, err := orchestratorAgent.Run(userQuery, session)
		if err != nil {
			panic(err)
		}

		if response.Content != "" {
			fmt.Printf("Assistant: %s\n", response.Content)
		}

		for response.HasToolCalls() {
			response, err = orchestratorAgent.Run("", session)
			if response.Content != "" {
				fmt.Printf("Assistant: %s\n", response.Content)
			}
			if err != nil {
				panic(err)
			}
		}
	}

	slog.Debug("Session details", "details", session)
}
