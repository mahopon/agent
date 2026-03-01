package templates

import (
	"strings"
	"text/template"
)

type Template interface {
	Name() string
	Create(data map[string]any) (string, error)
}

func Render(tpl string, data map[string]any) (string, error) {
	tmpl, err := template.New("").Parse(tpl)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", err
	}

	return result.String(), nil
}
