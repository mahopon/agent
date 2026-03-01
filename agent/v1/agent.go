package agent

import (
	"agent/llm"
	"agent/tool"
)

type AgentRunnable interface {
	Run(userQuery string, session *Session) error
}

type Agent struct {
	Name        string
	Tools       []tool.Tool
	SysPrompt   string
	InternalCoT []map[string]string
	LLMProvider llm.LLMCallable
}
