package agent

import (
	"agent/llm"
	"agent/tool"
)

const (
	PLANNING_PROMPT string = "You are a planning agent. Your job is to plan out the required steps to accomplish what the user has queried. You are to return the steps in a numbered list form."
)

type PlanningAgent struct {
	Agent
}

func NewPlanningAgent(llm llm.LLMCallable) *PlanningAgent {
	return &PlanningAgent{
		Agent: Agent{
			Name:        "Planner",
			Tools:       []tool.Tool{},
			SysPrompt:   PLANNING_PROMPT,
			InternalCoT: make([]map[string]string, 0),
			LLMProvider: llm,
		},
	}
}

func (a *PlanningAgent) Run(userQuery string, session *Session) error {
	content := []map[string]string{
		{
			"role":    "system",
			"content": a.SysPrompt,
		},
	}
	content = append(content, session.History...)
	body := llm.NewLLMBody(content, true)
	reply, err := a.LLMProvider.Call(body)
	if err != nil {
		return err
	}
	session.History = append(session.History, map[string]string{
		"role":    "assistant",
		"content": a.Name + ":\n" + reply,
	})
	return nil
}
