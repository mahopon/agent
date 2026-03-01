package agent

import (
	"agent/llm"
	"agent/prompt/templates"
	"agent/tool"
)

type AgentRunnable interface {
	Run(userQuery string, session *Session) (*llm.ParsedResponse, error)
}

type Agent struct {
	Name         string
	Tools        []tool.ToolExecutor
	LLM          llm.LLMCallable
	SystemPrompt *templates.SystemPrompt
}
