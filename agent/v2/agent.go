package agent

import (
	"agent/llm"
	"agent/prompt/templates"
	"agent/tool"
)

type AgentRunner interface {
	Run(userQuery string, session *Session) error
}

type Agent struct {
	Name         string
	Tools        []tool.ToolExecutor
	LLM          llm.LLMCallable
	SystemPrompt *templates.SystemPrompt
}

type AgentRunResponse struct {
	Continue bool
	Tools    []AgentToolCall
}

type AgentToolCall struct {
	Function   string
	Parameters []any
}

func NewAgentRunResponse(con bool, tools []AgentToolCall) *AgentRunResponse {
	return &AgentRunResponse{
		Continue: con,
		Tools:    tools,
	}
}

func NewAgentToolCall(functionName string, parameters []any) *AgentToolCall {
	return &AgentToolCall{
		Function:   functionName,
		Parameters: parameters,
	}
}
