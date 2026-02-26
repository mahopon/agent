package agent

import (
	"agent/llm"
	"agent/tool"
	"log"
	"slices"
)

type AgentAction interface {
	Execute() (string, error)
}

type Agent struct {
	Name      string
	Tools     []tool.Tool
	SysPrompt string
}

func (a *Agent) Call(userQuery string, loop *Session, config *llm.LLMConfig) (string, error) {
	if len(loop.History) == 0 {
		loop.History = slices.Insert(loop.History, 0, map[string]string{
			"role":    "system",
			"content": a.SysPrompt,
		})
	}
	log.Printf("%v", loop.History)
	params := llm.LLMCallBody{
		Model: config.LLM_MODEL,
		Msgs:  loop.History,
		Reasoning: map[string]bool{
			"enabled": false,
		},
	}
	reply, err := llm.LLMCall(userQuery, params, *config)
	if err != nil {
		return "", err
	}
	loop.Summary = reply
	log.Printf("Summary: %v", loop.Summary)
	return reply, nil
}
