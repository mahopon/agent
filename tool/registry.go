package tool

type ToolExecutor interface {
	Execute(name string, args map[string]any) (string, error)
}

func ToOpenAIScheme(executors []ToolExecutor) []Tool {
	return nil
}
