package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// GenerateOpenSpecStructure creates the folder structure and writes mock/real data
func (o *OpenSpecService) GenerateOpenSpecStructure(featureName string, content string) error {
	if featureName == "" {
		featureName = "new-feature"
	}
	cwd, _ := os.Getwd()
	fmt.Printf("Generating OpenSpec structure in: %s\n", filepath.Join(cwd, "openspec"))
	
	baseDir := "openspec"
	
	// Create changes dir
	changesDir := filepath.Join(baseDir, "changes", featureName)
	specsDir := filepath.Join(changesDir, "specs")
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		return err
	}

	// Create root specs dir
	rootSpecsDir := filepath.Join(baseDir, "specs", "auth")
	if err := os.MkdirAll(rootSpecsDir, 0755); err != nil {
		return err
	}

	// Write proposal
	proposalContent := fmt.Sprintf("# 变更提案: %s\n\n## 目的\n\n此提案由于 AI 驱动生成。描述了 %s 的背景及目的。\n\n## 背景\n\n用户意图如下:\n%s\n", featureName, featureName, content)
	if err := ioutil.WriteFile(filepath.Join(changesDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		return err
	}

	tasksContent := fmt.Sprintf("# 实施任务 (tasks.md)\n\n- [ ] 1. 深入调研 %s 的架构可行性\n- [ ] 2. 在核心框架中引入支持包\n- [ ] 3. 开发 UI 组件\n- [ ] 4. 编写并运行自动化测试\n", featureName)
	if err := ioutil.WriteFile(filepath.Join(changesDir, "tasks.md"), []byte(tasksContent), 0644); err != nil {
		return err
	}

	specDeltasContent := "## ADDED Requirements (增量规格)\n\n### Requirement: User Profile Filters\n- 系统 MUST 允许用户根据角色(role)进行搜索\n- 过滤组件 SHOULD 响应式支持移动端展示\n\n"
	if err := ioutil.WriteFile(filepath.Join(specsDir, "spec.md"), []byte(specDeltasContent), 0644); err != nil {
		return err
	}

	// write root global spec if not present
	rootSpecPath := filepath.Join(rootSpecsDir, "spec.md")
	if _, err := os.Stat(rootSpecPath); os.IsNotExist(err) {
		ioutil.WriteFile(rootSpecPath, []byte("# 核心验证服务 (Auth Spec)\n\n当前真理源规范区块。\n"), 0644)
	}

	// write project.md
	projectPath := filepath.Join(baseDir, "project.md")
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		ioutil.WriteFile(projectPath, []byte("# 项目级别约定\n\n- 缩进: 4空格\n- API响应: RESTful JSON\n- AI契约引擎: OpenSpec v1\n"), 0644)
	}

	return nil
}

type FileNode struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	IsDir    bool        `json:"isDir"`
	Children []*FileNode `json:"children"`
}

func (o *OpenSpecService) ListOpenSpecFiles() (*FileNode, error) {
	baseDir := "openspec"
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return nil, nil
	}
	return buildTree(baseDir)
}

func buildTree(path string) (*FileNode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	node := &FileNode{
		Name:  info.Name(),
		Path:  filepath.ToSlash(path),
		IsDir: info.IsDir(),
	}

	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			childPath := filepath.Join(path, entry.Name())
			childNode, err := buildTree(childPath)
			if err == nil {
				node.Children = append(node.Children, childNode)
			}
		}
	}
	return node, nil
}

func (o *OpenSpecService) ReadFileContent(path string) (string, error) {
	if !strings.HasPrefix(filepath.ToSlash(path), "openspec") {
		return "", fmt.Errorf("非 openspec 核心目录，禁止跨目录读取")
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
