package templates

type ErrorPrompt struct{}

func (e *ErrorPrompt) Name() string {
	return "error"
}

func (e *ErrorPrompt) Create(data map[string]any) (string, error) {
	return Render(errorTemplate, data)
}

const errorTemplate = `A tool execution failed:
Error: {{.Error}}

Respond to the user explaining what happened and whether you can recover from this error.`
