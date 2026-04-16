package agent

import (
	"agent/prompt/templates"
	"agent/tool"
)

type LLMCallable interface {
	Call()
}
type AgentRunnable interface {
	Run(userQuery string, reasoning bool, session *Session) (string, error)
}

type Agent struct {
	Name         string
	Tools        []tool.ToolExecutor
	LLM          LLMCallable
	SystemPrompt *templates.SystemPrompt
}
