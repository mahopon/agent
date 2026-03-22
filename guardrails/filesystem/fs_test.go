package filesystem

import (
	"testing"
)

func TestExcludedPath(t *testing.T) {
	for _, path := range EXCLUDED_PATHS {
		path := "test/" + path + "/test2"
		err := CheckPathSafe(path)
		if err == nil {
			t.Errorf("CheckPathSafe(%s) = %v. want %v", path, err, ErrForbiddenAccess.Error())
		}
	}
}
