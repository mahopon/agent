package agent

import (
	"agent/llm"
	"agent/prompt/templates"
	"agent/tool"
	"agent/memory"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type OrchestratorAgent struct {
	Agent
}

func (a *OrchestratorAgent) Run(userQuery string, reasoning bool, session *Session) (string, error) {
	if	session.History.Size() == 0 {
		cwd, _ := os.Getwd()
		prompt, err := a.SystemPrompt.Create(map[string]any{
			"cwd": cwd,
		})
		if err != nil {
			slog.Error("failed to create system prompt", "agent", a.Name, "error", err)
			return "", fmt.Errorf("failed to create system prompt: %w", err)
		}
		session.History.Add(memory.NewSystemMessage(prompt))
		slog.Info("system prompt added to session", "agent", a.Name, "cwd", cwd)
	}

	if userQuery != "" {
		session.History.Add(memory.NewUserMessage(userQuery))
		slog.Debug("user query added to history", "agent", a.Name, "query", userQuery)
	}

	tools := tool.ToOpenAIScheme(a.Tools)
	body := a.LLM.Call(session.History.ToHistory(), tools)
	slog.Debug("preparing LLM request", "agent", a.Name, "tools_count", len(tools))

	response, err := a.LLM.CallWithRetry(body)
	if err != nil {
		slog.Error("LLM call failed", "agent", a.Name, "error", err)
		return nil, err
	}
	slog.Debug("LLM response received", "agent", a.Name, "has_tool_calls", response.HasToolCalls())

	for response.HasToolCalls() {
		slog.Info("processing tool calls", "agent", a.Name, "tool_calls", len(response.ToolCalls))
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
			slog.Debug("tool call added to history", "agent", a.Name, "tool", tc.Name, "call_id", tc.ID)

			executor, err := tool.FindExecutorForTool(a.Tools, tc.Name)
			if err != nil {
				slog.Error("executor not found for tool", "agent", a.Name, "tool", tc.Name, "error", err)
				session.History = append(session.History, map[string]any{
					"role":         "tool",
					"tool_call_id": tc.ID,
					"content":      fmt.Sprintf("error: %v", err),
				})
				continue
			}

			var args map[string]any
			if err := json.Unmarshal([]byte(tc.Arguments), &args); err != nil {
				slog.Error("failed to parse tool arguments", "agent", a.Name, "tool", tc.Name, "error", err)
				session.History = append(session.History, map[string]any{
					"role":         "tool",
					"tool_call_id": tc.ID,
					"content":      fmt.Sprintf("error: failed to parse arguments: %v", err),
				})
				continue
			}

			result, err := executor(tc.Name, args)
			if err != nil {
				slog.Error("tool execution failed", "agent", a.Name, "tool", tc.Name, "error", err)
				session.History = append(session.History, map[string]any{
					"role":         "tool",
					"tool_call_id": tc.ID,
					"content":      fmt.Sprintf("error: %v", err),
				})
				continue
			}

			slog.Info("tool execution successful", "agent", a.Name, "tool", tc.Name, "result", result)
			session.History = append(session.History, map[string]any{
				"role":         "tool",
				"tool_call_id": tc.ID,
				"content":      result,
			})
			time.Sleep(50 * time.Millisecond)
		}

		slog.Info("tool calls processed, making follow-up LLM call", "agent", a.Name)
		body = llm.NewLLMBody(session.History, reasoning, tools)
		response, err = a.LLM.CallWithRetry(body)
		if err != nil {
			slog.Error("follow-up LLM call failed", "agent", a.Name, "error", err)
			return nil, err
		}
		slog.Debug("follow-up response received", "agent", a.Name, "has_tool_calls", response.HasToolCalls())
	}

	if !response.HasToolCalls() {
		session.History = append(session.History, map[string]any{
			"role":    "assistant",
			"content": response.Content,
		})
	}

	slog.Debug("run completed", "agent", a.Name)
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
