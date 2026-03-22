package filesystem

import (
	"strings"
)

var (
	EXCLUDED_PATHS = []string{".venv", "venv", "node_modules", ".git", ".gitlab", "__pycache__"}
)

type errorConst string

func (e errorConst) Error() string {
	return string(e)
}

const (
	ErrForbiddenAccess errorConst = "forbidden access"
)

func CheckPathSafe(path string) error {
	for _, p := range EXCLUDED_PATHS {
		if strings.Contains(path, p) {
			return ErrForbiddenAccess
		}
	}
	return nil
}
