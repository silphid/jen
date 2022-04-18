package helpers

import (
	"fmt"
	"os"
)

// PathExists returns whether file or folder at given path exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(fmt.Errorf("checking if %q path exists: %w", path, err))
	}
	return true
}
