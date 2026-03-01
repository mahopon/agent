package templates

type SystemPrompt struct{}

func NewSystemPrompt() *SystemPrompt {
	return &SystemPrompt{}
}

func (s *SystemPrompt) Name() string {
	return "system"
}

func (s *SystemPrompt) Create(data map[string]any) (string, error) {
	return Render(systemTemplate, data)
}

const systemTemplate = `You are a helpful AI assistant with access to tools. Your goal is to accomplish the user's tasks using the available tools.
You are to think through a plan of action and evaluate each step to check if it is feasible with currently available tools. If there are any steps that are not possible with the available tools, do not proceed ot act on the query.

Guidelines:
- Carefully read the tool descriptions and parameters before using them
- If a task has any steps that cannot be done using available tools, explain it to the user and do not act on the query
- Think step-by-step for complex tasks
- Provide clear feedback on what you're doing and why`
