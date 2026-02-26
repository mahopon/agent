package tool

type ToolExecutor interface {
	Execute(name string, args map[string]any) (string, error)
	Schema() []Tool
}

func ToOpenAIScheme(executors []ToolExecutor) []Tool {
	tools := make([]Tool, 0)
	for _, executor := range executors {
		tools = append(tools, executor.Schema()...)
	}
	return tools
}
