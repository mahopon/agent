package tool

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type FileSystemExecutor struct{}

type ReadFileInput struct {
	Path string `json:"path" yaml:"path"`
}

type WriteFileInput struct {
	Path    string `json:"path" yaml:"path"`
	Content string `json:"content" yaml:"content"`
}

type CreateFolderInput struct {
	Path string `json:"path" yaml:"path"`
}

type ListDirInput struct {
	Path string `json:"path" yaml:"path"`
}

type WalkDirInput struct {
	Root string `json:"root" yaml:"root"`
}

type ModifyFileInput struct {
	Path     string `json:"path" yaml:"path"`
	Modified string `json:"modified" yaml:"modified"`
}

var errMissingInput = errors.New("missing required input")

func NewFileSystemExecutor() *FileSystemExecutor {
	return &FileSystemExecutor{}
}

func (r ReadFileInput) Validate() error {
	if strings.TrimSpace(r.Path) == "" {
		return fmt.Errorf("%w: path", errMissingInput)
	}
	return nil
}

func (w WriteFileInput) Validate() error {
	if strings.TrimSpace(w.Path) == "" {
		return fmt.Errorf("%w: path", errMissingInput)
	}
	if strings.TrimSpace(w.Content) == "" {
		return fmt.Errorf("%w: content", errMissingInput)
	}
	return nil
}

func (c CreateFolderInput) Validate() error {
	if strings.TrimSpace(c.Path) == "" {
		return fmt.Errorf("%w: path", errMissingInput)
	}
	return nil
}

func (l ListDirInput) Validate() error {
	if strings.TrimSpace(l.Path) == "" {
		return fmt.Errorf("%w: path", errMissingInput)
	}
	return nil
}

func (w WalkDirInput) Validate() error {
	if strings.TrimSpace(w.Root) == "" {
		return fmt.Errorf("%w: root", errMissingInput)
	}
	return nil
}

func (m ModifyFileInput) Validate() error {
	if strings.TrimSpace(m.Path) == "" {
		return fmt.Errorf("%w: path", errMissingInput)
	}
	if strings.TrimSpace(m.Modified) == "" {
		return fmt.Errorf("%w: modified", errMissingInput)
	}
	return nil
}

func (f *FileSystemExecutor) Execute(name string, args map[string]any) (string, error) {
	switch name {
	case "read_file":
		path, ok := args["path"].(string)
		if !ok {
			return "", fmt.Errorf("%w: path", errMissingInput)
		}
		input := ReadFileInput{Path: path}
		if err := input.Validate(); err != nil {
			return "", err
		}
		return readFile(input.Path)

	case "write_file":
		path, ok := args["path"].(string)
		if !ok {
			return "", fmt.Errorf("%w: path", errMissingInput)
		}
		content, ok := args["content"].(string)
		if !ok {
			return "", fmt.Errorf("%w: content", errMissingInput)
		}
		input := WriteFileInput{Path: path, Content: content}
		if err := input.Validate(); err != nil {
			return "", err
		}
		return "", writeFile(input.Content, input.Path)

	case "get_cwd":
		return getWd()

	case "create_folder":
		path, ok := args["path"].(string)
		if !ok {
			return "", fmt.Errorf("%w: path", errMissingInput)
		}
		input := CreateFolderInput{Path: path}
		if err := input.Validate(); err != nil {
			return "", err
		}
		return "", createFolder(input.Path)

	case "list_dir":
		path, ok := args["path"].(string)
		if !ok {
			return "", fmt.Errorf("%w: path", errMissingInput)
		}
		input := ListDirInput{Path: path}
		if err := input.Validate(); err != nil {
			return "", err
		}
		return listDir(input.Path)

	case "walk_dir":
		root, ok := args["root"].(string)
		if !ok {
			return "", fmt.Errorf("%w: root", errMissingInput)
		}
		input := WalkDirInput{Root: root}
		if err := input.Validate(); err != nil {
			return "", err
		}
		return walkDir(input.Root)

	case "modify_file":
		path, ok := args["path"].(string)
		if !ok {
			return "", fmt.Errorf("%w: path", errMissingInput)
		}
		modified, ok := args["modified"].(string)
		if !ok {
			return "", fmt.Errorf("%w: modified", errMissingInput)
		}
		input := ModifyFileInput{Path: path, Modified: modified}
		if err := input.Validate(); err != nil {
			return "", err
		}
		return "", modifyFile(input.Path, input.Modified)
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
					Required:             []string{"path", "modified"},
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
	patches, err := dmp.PatchFromText(modified)
	if err != nil {
		return fmt.Errorf("invalid patch format: %w", err)
	}

	if len(patches) == 0 {
		return errors.New("empty patch")
	}

	result, applied := dmp.PatchApply(patches, string(original))
	for i, ok := range applied {
		if !ok {
			return fmt.Errorf("patch failed to apply at index %d", i)
		}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	dir := filepath.Dir(absPath)

	tmp, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if _, err = tmp.WriteString(result); err != nil {
		tmp.Close()
		return err
	}
	if err = tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}

	if err = os.Rename(tmp.Name(), path); err != nil {
		return err
	}

	if dirFd, err := os.Open(dir); err == nil {
		dirFd.Sync()
		dirFd.Close()
	}

	return nil
}
