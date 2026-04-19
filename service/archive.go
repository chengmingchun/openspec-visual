package service

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type ArchiveService struct{}

func NewArchiveService() *ArchiveService {
	return &ArchiveService{}
}

// CommitSpecs commits changes in the openspec folder to the local repository using go-git
func (s *ArchiveService) CommitSpecs(commitMsg string) error {
	r, err := git.PlainOpen(".")
	if err != nil {
		return fmt.Errorf("未找到 Git 仓库: %w", err)
	}

	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("无法获取 Worktree: %w", err)
	}

	// 自动跟踪 openspec 的变更
	_, err = w.Add("openspec")
	if err != nil {
		return fmt.Errorf("Git Add 失败: %w", err)
	}

	_, err = w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "OpenSpec Agent",
			Email: "agent@openspec.local",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("Git Commit 失败: %w", err)
	}

	return nil
}
