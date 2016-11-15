package util

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/CloudCoreo/cli/cmd/content"
)

// CheckGitInstall check if user has git installed
func CheckGitInstall() error {
	cmd := exec.Command("git", "version")

	if _, err := cmd.Output(); err != nil {
		println("Unable to execute 'git': " + err.Error())
		return err
	}

	return nil
}

// CreateGitSubmodule create composite git submodule
func CreateGitSubmodule(directory, gitRepoURL string) error {

	fmt.Printf(content.INFO_CREATING_GITSUBMODULE, gitRepoURL, directory)
	if directory == "" {
		return fmt.Errorf(content.ERROR_DIRECTORY_PATH_NOT_PROVIDED)
	}

	if gitRepoURL == "" {
		return fmt.Errorf(content.ERROR_GIT_URL_NOT_PROVIDED)
	}

	err := os.Chdir(directory)
	if err != nil {
		return fmt.Errorf(content.ERROR_INVALID_DIRECTORY, directory)
	}

	_, err = exec.Command("git", "rev-parse", "--is-inside-work-tree").CombinedOutput()
	if err != nil {
		return fmt.Errorf(content.ERROR_GIT_INIT_NOT_RAN)
	}

	output, err := exec.Command("git", "submodule", "add", "-f", gitRepoURL, "extends").CombinedOutput()
	if err != nil {
		return fmt.Errorf(content.ERROR_GIT_SUBMODULE_FAILED, output)
	}

	return nil
}
