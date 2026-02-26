package agent

import (
	"agent/tool"
)

type Agent struct {
	Name  string
	Tools []tool.Tool
}

func (a *Agent) ReadFile(path string) string {
	return ""
}

func (a *Agent) WriteFile(path, content string) error {
	return nil
}

func (a *Agent) ListFiles(dir string) []string {
	return []string{}
}
