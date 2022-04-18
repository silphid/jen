package home

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/silphid/jen/cmd/jen/internal/constant"
	"github.com/silphid/jen/cmd/jen/internal/helpers"
	"github.com/silphid/jen/cmd/jen/internal/logging"
	"github.com/silphid/jen/cmd/jen/internal/shell"
)

const (
	// jenCloneVar is the name of env var specifying the local path where to clone the jen git repo containing user scripts
	// and templates.
	jenCloneVar = "JEN_CLONE"

	// jenRepoVar is the name of env var specifying the URL of jen git repo containing user scripts and templates.
	jenRepoVar = "JEN_REPO"

	// jenRepoSubDirVar is the name of env var specifying the sub-directory within jen git repo, where user shared scripts
	// "bin" and "templates" directories are located.
	jenRepoSubDirVar = "JEN_SUBDIR"
)

// GetOrCloneRepo clones the jen git repo if it does not exist and returns the path to where it was cloned
func GetOrCloneRepo() (string, error) {
	cloneDir, err := getCloneDir()
	if err != nil {
		return "", fmt.Errorf("failed to locate clone dir: %w", err)
	}

	if helpers.PathExists(cloneDir) {
		if helpers.PathExists(filepath.Join(cloneDir, ".git")) {
			// Valid git repo
			logging.Log("Using jen templates clone at %q", cloneDir)
			return cloneDir, nil
		}

		// Not a valid git repo, therefore must be empty, so we can clone into it
		infos, err := ioutil.ReadDir(cloneDir)
		if err != nil {
			return cloneDir, fmt.Errorf("listing content of target dir %q to ensure it's empty before cloning into it: %w", cloneDir, err)
		}
		if len(infos) > 0 {
			return cloneDir, fmt.Errorf("dir %q already exists, but is not a valid git working copy and already contains files so we cannot clone into it (please delete or empty it)", cloneDir)
		}
	}

	jenRepoURL, err := getJenRepoURL()
	if err != nil {
		return cloneDir, fmt.Errorf("failed to detect jen repo: %w", err)
	}

	logging.Log("Cloning jen templates repo %q into dir %q", jenRepoURL, cloneDir)
	return cloneDir, shell.Execute(nil, "", fmt.Sprintf("git clone %s %s", jenRepoURL, cloneDir))
}

// getJenRepoURL returns URL of templates git repo to clone, as specified by required JEN_REPO env var.
func getJenRepoURL() (string, error) {
	jenRepo, ok := os.LookupEnv(jenRepoVar)
	if !ok {
		return "", fmt.Errorf("please specify a JEN_REPO env var pointing to your jen templates git repo")
	}
	return jenRepo, nil
}

// getCloneDir returns the path where jen will clone the templates git repo, as specified by JEN_CLONE
// env var, defaulting to "~/.jen/repo".
func getCloneDir() (cloneDir string, err error) {
	defer func() {
		if err == nil {
			logging.Log("Using clone dir: %s", cloneDir)
		}
	}()

	cloneDir, ok := os.LookupEnv(jenCloneVar)
	if ok && cloneDir != "" {
		return
	}

	home, err := homedir.Dir()
	if err != nil {
		err = fmt.Errorf("failed to detect home directory: %w", err)
		return
	}
	cloneDir = filepath.Join(home, constant.DefaultCloneDir)
	return
}

// GetCloneSubDir returns the path within cloned git repo where to look
// for "bin" and "templates" directories.
func GetCloneSubDir() (string, error) {
	repoCloneDir, err := getCloneDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(repoCloneDir, os.Getenv(jenRepoSubDirVar)), nil
}

// GetTemplatesDir returns the path within cloned templates git repo where
// all templates all located.
func GetTemplatesDir() (string, error) {
	repoSubDir, err := GetCloneSubDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(repoSubDir, constant.TemplatesDirName), nil
}
