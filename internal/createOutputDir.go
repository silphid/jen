package internal

import "os"

func createOutputDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}
