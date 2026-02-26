package tool

import (
// "os"
)

type FileSystemExecutor struct{}

func (f *FileSystemExecutor) Execute(name string, args map[string]any) (string, error) {
	return "", nil
}
