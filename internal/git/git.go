package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// File represents a file in the git repository
type File struct {
	Path      string
	Status    string // M=modified, A=added, D=deleted, ??=untracked
	IsStaged  bool
	IsTracked bool
}

// GetStatus returns the current git status
func GetStatus() ([]File, error) {
	cmd := exec.Command("git", "status", "--porcelain", "-uall")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git status: %w", err)
	}

	var files []File
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		// Parse git status line
		// Format: XY filename
		// X = staged status, Y = unstaged status
		if len(line) < 3 {
			continue
		}
		
		stagedStatus := line[0]
		unstagedStatus := line[1]
		filename := strings.TrimSpace(line[3:])
		
		file := File{
			Path:      filename,
			IsStaged:  stagedStatus != ' ' && stagedStatus != '?',
			IsTracked: stagedStatus != '?' || unstagedStatus != '?',
		}
		
		// Determine status
		if stagedStatus == '?' && unstagedStatus == '?' {
			file.Status = "??"
		} else if stagedStatus == 'A' || unstagedStatus == 'A' {
			file.Status = "A"
		} else if stagedStatus == 'M' || unstagedStatus == 'M' {
			file.Status = "M"
		} else if stagedStatus == 'D' || unstagedStatus == 'D' {
			file.Status = "D"
		} else if stagedStatus == 'R' || unstagedStatus == 'R' {
			file.Status = "R"
		}
		
		files = append(files, file)
	}
	
	return files, nil
}

// StageFiles stages the specified files
func StageFiles(files []string) error {
	if len(files) == 0 {
		return nil
	}
	
	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stage files: %w\n%s", err, output)
	}
	
	return nil
}

// UnstageFiles unstages the specified files
func UnstageFiles(files []string) error {
	if len(files) == 0 {
		return nil
	}
	
	args := append([]string{"reset", "HEAD", "--"}, files...)
	cmd := exec.Command("git", args...)
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to unstage files: %w\n%s", err, output)
	}
	
	return nil
}

// GetStagedDiff returns the diff of staged changes
func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get staged diff: %w", err)
	}
	
	return string(output), nil
}

// GetStagedDiffForFile returns the diff of staged changes for a specific file
func GetStagedDiffForFile(file string) (string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--", file)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get staged diff for file: %w", err)
	}
	
	return string(output), nil
}

// GetStagedFiles returns a list of staged files
func GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var files []string
	for _, line := range lines {
		if line != "" {
			files = append(files, line)
		}
	}
	
	return files, nil
}

// Commit creates a commit with the given message
func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit: %w\n%s", err, stderr.String())
	}
	
	return nil
}

// GetLastCommitMessage returns the last commit message
func GetLastCommitMessage() (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=%B")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get last commit message: %w", err)
	}
	
	return strings.TrimSpace(string(output)), nil
} 