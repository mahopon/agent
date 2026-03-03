package agent

import (
	"agent/llm"
	"agent/prompt/templates"
	"agent/tool"
	"encoding/json"
	"fmt"
	"os"
)

type OrchestratorAgent struct {
	Agent
}

func (a *OrchestratorAgent) Run(userQuery string, reasoning bool, session *Session) (*llm.ParsedResponse, error) {
	if len(session.History) == 0 {
		cwd, _ := os.Getwd()
		prompt, err := a.SystemPrompt.Create(map[string]any{
			"cwd": cwd,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create system prompt: %w", err)
		}

		session.History = append(session.History, map[string]any{
			"role":    "system",
			"content": prompt,
		})
	}

	if userQuery != "" {
		session.History = append(session.History, map[string]any{
			"role":    "user",
			"content": userQuery,
		})
	}

	tools := tool.ToOpenAIScheme(a.Tools)
	body := llm.NewLLMBody(session.History, reasoning, tools)

	response, err := a.LLM.CallWithRetry(body)
	if err != nil {
		return nil, err
	}

	session.History = append(session.History, map[string]any{
		"role":    "assistant",
		"content": response.Content,
	})

	if response.HasToolCalls() {
		for _, tc := range response.ToolCalls {
			session.History = append(session.History, map[string]any{
				"role": "assistant",
				"tool_calls": []map[string]any{
					{
						"id":   tc.ID,
						"type": "function",
						"function": map[string]any{
							"name":      tc.Name,
							"arguments": tc.Arguments,
						},
					},
				},
			})

			executor, err := tool.FindExecutorForTool(a.Tools, tc.Name)
			if err != nil {
				session.History = append(session.History, map[string]any{
					"role":         "tool",
					"tool_call_id": tc.ID,
					"content":      fmt.Sprintf("error: %v", err),
				})
				continue
			}

			var args map[string]any
			if err := json.Unmarshal([]byte(tc.Arguments), &args); err != nil {
				session.History = append(session.History, map[string]any{
					"role":         "tool",
					"tool_call_id": tc.ID,
					"content":      fmt.Sprintf("error: failed to parse arguments: %v", err),
				})
				continue
			}

			result, err := executor(tc.Name, args)
			if err != nil {
				session.History = append(session.History, map[string]any{
					"role":         "tool",
					"tool_call_id": tc.ID,
					"content":      fmt.Sprintf("error: %v", err),
				})
				continue
			}

			session.History = append(session.History, map[string]any{
				"role":         "tool",
				"tool_call_id": tc.ID,
				"content":      result,
			})
		}
	}

	return response, nil
}

func NewOrchestratorAgent(sysPrompt *templates.SystemPrompt, tools []tool.ToolExecutor, llm llm.LLMCallable) *OrchestratorAgent {
	return &OrchestratorAgent{
		Agent{
			Name:         "Orchestrator",
			LLM:          llm,
			Tools:        tools,
			SystemPrompt: sysPrompt,
		},
	}
}
