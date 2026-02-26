package agent

import (
	"agent/tool"
)

const (
	PLANNING_PROMPT string = "You are a planning agent. Your job is to plan out the required steps to accomplish what the user has queried. You are to return the steps in a numbered list form."
)

func NewPlanningAgent(session *Session) *Agent {
	return &Agent{
		Name:      "Planner",
		Tools:     []tool.Tool{},
		SysPrompt: PLANNING_PROMPT,
		Session:   session,
	}
}
