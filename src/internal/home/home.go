package home

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Samasource/jen/src/internal/helpers"
	"github.com/Samasource/jen/src/internal/logging"
	"github.com/Samasource/jen/src/internal/shell"
	"github.com/mitchellh/go-homedir"
)

// getJenRepo reads the repoitory value from the environment and returns an error if it is not set
func getJenRepo() (string, error) {
	jenRepo, ok := os.LookupEnv("JEN_REPO")
	if !ok {
		return "", fmt.Errorf("please specify a JEN_REPO env var pointing to your jen templates git repo")
	}
	return jenRepo, nil
}

// getJenHomeDir returns the path to the jen home folder defaulting to ~/.jen if not provided
func getJenHomeDir() (jenHomeDir string, err error) {
	defer func() {
		if err == nil {
			logging.Log("Using jen home dir: %s", jenHomeDir)
		}
	}()

	jenHomeDir, ok := os.LookupEnv("JEN_HOME")
	if ok && jenHomeDir != "" {
		return
	}

	home, err := homedir.Dir()
	if err != nil {
		err = fmt.Errorf("failed to detect home directory: %v", err)
		return
	}
	jenHomeDir = filepath.Join(home, ".jen")
	return
}

// CloneJenRepo will clone the jenRepo if it does not exist, and return the path to where it was cloned
func CloneJenRepo() (string, error) {
	jenHome, err := getJenHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to locate jen home: %w", err)
	}

	if helpers.PathExists(jenHome) {
		if helpers.PathExists(filepath.Join(jenHome, ".git")) {
			return jenHome, nil
		}

		infos, err := ioutil.ReadDir(jenHome)
		if err != nil {
			return jenHome, fmt.Errorf("listing content of jen dir %q to ensure it's empty before cloning into it: %w", jenHome, err)
		}

		if len(infos) > 0 {
			return jenHome, fmt.Errorf("jen dir %q already exists, is not a valid git working copy and already contains files so we cannot clone into it (please delete or empty it)", jenHome)
		}
	}

	jenRepo, err := getJenRepo()
	if err != nil {
		return jenHome, fmt.Errorf("failed to detect jen repo: %w", err)
	}

	logging.Log("Cloning jen templates repo %q into jen dir %q", jenRepo, jenHome)
	return jenHome, shell.Execute(nil, "", nil, fmt.Sprintf("git clone %s %s", jenRepo, jenHome))
}
