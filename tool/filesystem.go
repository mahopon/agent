package tool

import (
	"fmt"
	"os"
)

type FileSystemExecutor struct{}

func NewFileSystemExecutor() *FileSystemExecutor {
	return &FileSystemExecutor{}
}

func (f *FileSystemExecutor) Execute(name string, args map[string]any) (string, error) {
	switch name {
	case "read_file":
		path := args["path"].(string)
		return readFile(path)
	case "write_file":
		path := args["path"].(string)
		content := args["content"].(string)
		return "", writeFile(content, path)
	case "get_cwd":
		return getWd()
	case "create_folder":
		path := args["path"].(string)
		return "", createFolder(path)
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
		{
			Type: "function",
			Function: Function{
				Name:        "create_folder",
				Description: "Creates folder in specified path",
				Parameters: ToolParameters{
					Type: "object",
					Properties: map[string]ToolProperty{
						"path": {
							Type:        "string",
							Description: "Path to create folder in",
						},
					},
					Required:             []string{"path"},
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
	return os.WriteFile(path, []byte(content), 0644)
}

func getWd() (string, error) {
	return os.Getwd()
}

func createFolder(path string) error {
	pwd, err := getWd()
	if err != nil {
		return err
	}
	concPath := fmt.Sprintf("%s/%s", pwd, path)
	return os.Mkdir(concPath, 0644)
}
