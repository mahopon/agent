package tool

import (
	"errors"
	"os"
)

type FileSystemExecutor struct{}

func (f *FileSystemExecutor) Execute(name string, args map[string]any) (string, error) {
	path := args["path"].(string)
	switch name {
	case "read_file":
		return readFile(path)
	case "write_file":
		content := args["content"].(string)
		return "", writeFile(content, path)
	case "get_cwd":
		return os.Getwd()
	}
	return "", nil
}

func (f *FileSystemExecutor) Schema() []Tool {
	return []Tool{
		{
			Type: "function",
			Function: Function{
				Name:        "read_file",
				Description: "Read the contents of a file",
				Parameters: ToolParameters{
					Type: "object",
					Properties: map[string]ToolProperty{
						"path": {
							Type:        "string",
							Description: "The path to the file to read",
						},
					},
					Required:             []string{"path"},
					AdditionalProperties: false,
				},
				Strict: true,
			},
		},
		{
			Type: "function",
			Function: Function{
				Name:        "write_file",
				Description: "Write content to a file",
				Parameters: ToolParameters{
					Type: "object",
					Properties: map[string]ToolProperty{
						"path": {
							Type:        "string",
							Description: "The path to the file to write",
						},
						"content": {
							Type:        "string",
							Description: "The content to write to the file",
						},
					},
					Required:             []string{"path", "content"},
					AdditionalProperties: false,
				},
				Strict: true,
			},
		},
		{
			Type: "function",
			Function: Function{
				Name:        "get_cwd",
				Description: "Get the current working directory",
				Parameters: ToolParameters{
					Type:                 "object",
					Properties:           map[string]ToolProperty{},
					Required:             []string{},
					AdditionalProperties: false,
				},
				Strict: true,
			},
		},
	}
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func writeFile(content, path string) error {
	fileExists, err := isFileExists(path)
	var file *os.File
	if !fileExists {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		file, err = os.Create(path)
		if err != nil {
			return err
		}
	} else {
		file, err = os.OpenFile(path, os.O_RDWR, 0644)
		if err != nil {
			return err
		}
	}
	if _, err := file.Write([]byte(content)); err != nil {
		return err
	}
	return nil
}

func isFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
