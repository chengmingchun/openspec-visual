package service

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	
	"openspec-visualizer/domain"
)

type ArchiveService struct{}

func NewArchiveService() *ArchiveService {
	return &ArchiveService{}
}

// CommitSpecs commits changes in the openspec folder to the local repository using go-git
func (s *ArchiveService) CommitSpecs(commitMsg string) error {
	// Workaround: Use exec git add to ensure deletions are captured perfectly
	exec.Command("git", "add", "openspec").Run()

	r, err := git.PlainOpen(".")
	if err != nil {
		return fmt.Errorf("未找到 Git 仓库: %w", err)
	}

	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("无法获取 Worktree: %w", err)
	}

	_, err = w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "OpenSpec Agent",
			Email: "agent@openspec.local",
			When:  time.Now(),
		},
	})
	if err != nil {
		// Suppress error if working tree is clean
		return nil
	}

	return nil
}

// GetHistory fetches git commit history mapped to openspec directory
func (s *ArchiveService) GetHistory() ([]domain.CommitLog, error) {
	cmd := exec.Command("git", "log", "--pretty=format:%H|%s|%an|%cI", "--", "openspec")
	out, err := cmd.Output()
	if err != nil {
		// Possibly no commits yet, return empty
		return []domain.CommitLog{}, nil
	}

	lines := strings.Split(string(out), "\n")
	logs := make([]domain.CommitLog, 0)
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			continue
		}
		parts := strings.SplitN(l, "|", 4)
		if len(parts) == 4 {
			logs = append(logs, domain.CommitLog{
				Hash:    parts[0],
				Message: parts[1],
				Author:  parts[2],
				Date:    parts[3],
			})
		}
	}
	return logs, nil
}

// GetDiff fetches the diff of a target commit against its parent
func (s *ArchiveService) GetDiff(hash string) (string, error) {
	cmd := exec.Command("git", "show", "--color=never", hash, "--", "openspec")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git show 失败: %v", err)
	}
	return string(out), nil
}

// Rollback checks out the openspec working tree directory to a previous state and commits it
func (s *ArchiveService) Rollback(hash string) error {
	cmd := exec.Command("git", "checkout", hash, "--", "openspec")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("回退状态失败: %v %s", err, string(out))
	}
	
	return s.CommitSpecs(fmt.Sprintf("Agent Timeline Rollback to %s", hash[:7]))
}
