package git

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type LocalGit interface {
	CreateLocalTag(string) error
	GenerateCommitDiff(string, string, string) error
}

type localGit struct {
}

func NewGitLocal() LocalGit {
	return &localGit{}
}

func (l *localGit) CreateLocalTag(tag string) error {
	cmd := exec.Command("git", "checkout", "master")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	cmd = exec.Command("git", "pull", "origin", "master")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	cmd = exec.Command("git", "tag", tag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}

func (l *localGit) GenerateCommitDiff(previousTag, currentTag, destinationPath string) error {
	outFile, err := os.Create(destinationPath)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	cmd := exec.Command("git", "log", fmt.Sprintf("%s..%s", previousTag, currentTag))
	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr

	// Execute the command
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}
