package evaluation

import (
	"io/ioutil"
	"os"
	"path"
)

func getTempDir() string {
	dir, err := ioutil.TempDir("/tmp", "jen_test_")
	if err != nil {
		panic(err)
	}
	return dir
}

func removeAll(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		panic(err)
	}
}

func getTempFile() string {
	file, err := ioutil.TempFile("/tmp", "jen_test_")
	if err != nil {
		panic(err)
	}
	return file.Name()
}

func writeTempFile(content string) string {
	file := getTempFile()
	err := ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		panic(err)
	}
	return file
}

func readFile(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func deleteFile(file string) {
	if err := os.Remove(file); err != nil {
		panic(err)
	}
}

func createEmptyFile(filePath string) {
	dir := path.Dir(filePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filePath, []byte(""), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
