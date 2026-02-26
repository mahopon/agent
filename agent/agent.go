package agent

import (
	"agent/llm"
	"agent/tool"
	"log/slog"
	"slices"
)

type AgentAction interface {
	Execute() (string, error)
}

type Agent struct {
	Name            string
	Tools           []tool.Tool
	SysPrompt       string
	InternalHistory []map[string]string
}

func (a *Agent) Call(userQuery string, requiredContext []map[string]string, config *llm.LLMConfig) (string, error) {
	if len(a.InternalHistory) == 0 {
		var content string
		if requiredContext != nil {
			content = a.SysPrompt + requiredContext[0]["content"]
		} else {
			content = a.SysPrompt
		}
		a.InternalHistory = slices.Insert(a.InternalHistory, 0, map[string]string{
			"role":    "system",
			"content": content,
		})
	}
	params := llm.LLMCallBody{
		Model: config.LLM_MODEL,
		Msgs:  a.InternalHistory,
		Reasoning: map[string]bool{
			"enabled": false,
		},
	}
	reply, err := llm.LLMCall(userQuery, params, *config)
	if err != nil {
		return "", err
	}
	slog.Debug("Agent Call", "details", a)
	return reply, nil
}
