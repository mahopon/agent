package tool

import (
	"errors"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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
	case "list_dir":
		path := args["path"].(string)
		return listDir(path)
	case "walk_dir":
		path := args["root"].(string)
		return walkDir(path)
	case "modify_file":
		path := args["path"].(string)
		modified := args["modified"].(string)
		return "", modifyFile(path, modified)
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
							Description: "Path to create folder in. Working directory is to be included before calling",
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
				Name:        "list_dir",
				Description: "Checks the current directory for files and directories and lists them. Output is returned with files prefixed with 'File:' and directoryes with 'Dir:'",
				Parameters: ToolParameters{
					Type: "object",
					Properties: map[string]ToolProperty{
						"path": {
							Type:        "string",
							Description: "Path to directory to be checked. Working directory is to be included before calling",
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
				Name:        "walk_dir",
				Description: "Recursively walks through all directories from the root directory. Output is returned with files prefixed with 'File:' and directories with 'Dir:'",
				Parameters: ToolParameters{
					Type: "object",
					Properties: map[string]ToolProperty{
						"root": {
							Type:        "string",
							Description: "Root directory to start checking from. Working directory is to be included before calling",
						},
					},
					Required:             []string{"root"},
					AdditionalProperties: false,
				},
				Strict: true,
			},
		},
		{
			Type: "function",
			Function: Function{
				Name:        "modify_file",
				Description: "Modifies a file safely using patch",
				Parameters: ToolParameters{
					Type: "object",
					Properties: map[string]ToolProperty{
						"path": {
							Type:        "string",
							Description: "Path of the file to be modfied. Working directory is to be included before calling",
						},
						"modified": {
							Type:        "string",
							Description: "Modified generated code based on the original file",
						},
					},
					Required:             []string{"root"},
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
	// pwd, err := getWd()
	// if err != nil {
	// 	return err
	// }
	// concPath := fmt.Sprintf("%s/%s", pwd, path)
	return os.MkdirAll(path, 0744)
}

func listDir(path string) (string, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}
	var builder strings.Builder
	for _, entry := range dir {
		var insertStr string
		if entry.IsDir() {
			insertStr = fmt.Sprintf("Dir: %s", entry.Name())
		} else {
			insertStr = fmt.Sprintf("File: %s", entry.Name())
		}
		builder.WriteString(fmt.Sprintf("%s\n", insertStr))
	}
	return builder.String(), nil
}

func walkDir(root string) (string, error) {
	var builder strings.Builder

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		var insertStr string
		if entry.IsDir() {
			insertStr = fmt.Sprintf("Dir: %s", path)
		} else {
			insertStr = fmt.Sprintf("File: %s", path)
		}

		builder.WriteString(fmt.Sprintf("%s\n", insertStr))
		return nil
	})

	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func modifyFile(path, modified string) error {
	original, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	dmp := diffmatchpatch.New()
	patches := dmp.PatchMake(string(original), modified)
	result, applied := dmp.PatchApply(patches, string(original))
	for _, ok := range applied {
		if !ok {
			return errors.New("patch failed to apply cleanly")
		}
	}
	absPath, _ := filepath.Abs(path)
	dir := filepath.Dir(absPath)
	tmp, err := os.CreateTemp(dir, ".tmp-*")
	_, err = tmp.WriteString(result)
	tmp.Sync()
	tmp.Close()
	err = os.Rename(tmp.Name(), path)
	if err != nil {
		return err
	}
	dirFd, _ := os.Open(dir)
	dirFd.Sync()
	dirFd.Close()
	return nil
}
