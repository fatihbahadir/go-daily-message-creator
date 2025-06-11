package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gdmc/internal/config"
)

type CommitFetcher struct {
	config *config.Config
}

func NewCommitFetcher(cfg *config.Config) *CommitFetcher {
	return &CommitFetcher{config: cfg}
}

func (cf *CommitFetcher) FetchCommits(author, intervalKey string) ([]string, error) {
	interval, exists := cf.config.Intervals[intervalKey]
	if !exists {
		return nil, fmt.Errorf("unknown interval: %s", intervalKey)
	}

	// Check if current directory is a git repository
	if !cf.isGitRepository() {
		return nil, fmt.Errorf("current directory is not a git repository. Please run gdmc from within a git repository")
	}

	// Get current repository info
	repoInfo, err := cf.getRepositoryInfo()
	if err != nil {
		fmt.Printf("Warning: Could not get repository info: %v\n", err)
	} else {
		fmt.Printf("Repository: %s\n", repoInfo)
	}

	args := cf.buildGitArgs(author, interval)

	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git command failed: %w. Make sure you have commits from author '%s'", err, author)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}, nil
	}

	commitCount := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "commit ") {
			commitCount++
		}
	}

	fmt.Printf("Found %d commits from %s to %s period\n", commitCount, interval.Since, interval.Until)

	return lines, nil
}

func (cf *CommitFetcher) isGitRepository() bool {
	_, err := os.Stat(".git")
	if err == nil {
		return true
	}

	// Check if we're in a subdirectory of a git repository
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err = cmd.Run()
	return err == nil
}

func (cf *CommitFetcher) getRepositoryInfo() (string, error) {
	// Try to get remote origin URL
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		// Fallback to current directory name
		wd, wdErr := os.Getwd()
		if wdErr != nil {
			return "", fmt.Errorf("could not determine repository info")
		}
		return fmt.Sprintf("Local repository: %s", wd), nil
	}

	return strings.TrimSpace(string(output)), nil
}

func (cf *CommitFetcher) buildGitArgs(author string, interval config.Interval) []string {
	args := []string{"log"}

	args = append(args, fmt.Sprintf("--since=%s", interval.Since))
	args = append(args, fmt.Sprintf("--until=%s", interval.Until))

	// Add branches
	for _, branch := range cf.config.GitSettings.Branches {
		args = append(args, branch)
	}

	// Add merge settings
	if !cf.config.GitSettings.IncludeMerges {
		args = append(args, "--no-merges")
	}

	args = append(args, fmt.Sprintf("--author=%s", author))

	// Add exclude paths
	for _, path := range cf.config.GitSettings.ExcludePaths {
		args = append(args, fmt.Sprintf("-- ':!%s'", path))
	}

	return args
}

func (cf *CommitFetcher) GetAvailableIntervals() []string {
	intervals := make([]string, 0, len(cf.config.Intervals))
	for key := range cf.config.Intervals {
		intervals = append(intervals, key)
	}
	return intervals
}
