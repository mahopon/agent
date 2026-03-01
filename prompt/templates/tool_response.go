package templates

type ToolResponse struct{}

func (t *ToolResponse) Name() string {
	return "tool_response"
}

func (t *ToolResponse) Create(data map[string]any) (string, error) {
	return Render(toolResponseTemplate, data)
}

const toolResponseTemplate = `You just executed a tool. Here's the result:
{{.Result}}

What does this result tell you? Do you need to:
- Use another tool to continue?
- Process this information to form your final answer?
- Ask the user for clarification?

Provide your next action or your final response.`
