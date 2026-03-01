package tool

type AgentExecutor struct{}

func NewAgentExecutor() *AgentExecutor {
	return &AgentExecutor{}
}

func (a *AgentExecutor) Execute(name string, args map[string]any) (string, error) {
	return "", nil
}
