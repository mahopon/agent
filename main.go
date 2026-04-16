package main

import (
	"agent/agent/v2"
	"agent/config"
	"agent/llm"
	"agent/tool"
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"agent/prompt/templates"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	debugFlag := flag.Bool("debug", false, "Enable debugging")
	shortDebugFlag := flag.Bool("d", false, "Enables debugging")

	flag.Parse()

	var debugging bool

	visitedFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		visitedFlags[f.Name] = true
	})

	if visitedFlags["debug"] {
		debugging = *debugFlag
	} else if visitedFlags["d"] {
		debugging = *shortDebugFlag
	}

	config := config.NewConfig()
	if debugging == true {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	} else {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})))
	}
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
		response, err := orchestratorAgent.Run(userQuery, true, session)
		if err != nil {
			panic(err)
		}

		if response.Content != "" {
			fmt.Printf("Assistant: %s\n", response.Content)
		}

		const maxIterations = 50
		iterations := 0

		for response.HasToolCalls() {
			iterations++
			response, err = orchestratorAgent.Run("", true, session)
			if err != nil {
				slog.Error("Tool call iteration failed", "error", err, "iteration", iterations)
				fmt.Println("Assistant: Encountered an error during tool execution.")
				break
			}
			if response == nil {
				slog.Error("Run returned nil response with no error")
				slog.Error("Response error:", "response", response, "history", session.History)
				fmt.Println("Assistant: Received empty response.")
				continue
			}
			if response.Content != "" {
				fmt.Printf("Assistant: %s\n", response.Content)
			}
		}

		if iterations >= maxIterations {
			fmt.Println("Assistant: Reached maximum iterations limit. Task may not be complete.")
		} else if response.Content == "" {
			fmt.Println("Assistant: (No response content generated)")
		}
	}

	slog.Debug("Session details", "details", session)
}
