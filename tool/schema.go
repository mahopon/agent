package tool

import (
	"os"
	"path/filepath"
)

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  ToolParameters `json:"parameters"`
	Strict      bool           `json:"strict"`
}

type ToolParameters struct {
	Type                 string                  `json:"type"`
	Properties           map[string]ToolProperty `json:"properties"`
	Required             []string                `json:"required"`
	AdditionalProperties bool                    `json:"additionalProperties"`
}

type ToolProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}


func getDataPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir := filepath.Join(cwd, "data")
	return dir
}
