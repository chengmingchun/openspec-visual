package service

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"

	"openspec-visualizer/domain"
)

type CheckerService struct {
	config domain.TDDConfig
}

func NewCheckerService(baseDir string) *CheckerService {
	svc := &CheckerService{}
	yamlPath := filepath.Join(baseDir, "tdd_rules.yaml")

	// Default TDD rules
	defaultConfig := domain.TDDConfig{
		Rules: []domain.TDDRule{
			{Name: "Given/When/Then", Description: "核心逻辑需求必须包含 Given/When/Then 关键字", Regex: "(?i)(given.*when.*then)"},
			{Name: "Markdown Format", Description: "产物必须是标准的 Markdown 结构", Regex: "^#\\s+"},
		},
	}

	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		data, _ := yaml.Marshal(&defaultConfig)
		ioutil.WriteFile(yamlPath, data, 0644)
		svc.config = defaultConfig
	} else {
		data, err := ioutil.ReadFile(yamlPath)
		if err == nil {
			yaml.Unmarshal(data, &svc.config)
		} else {
			svc.config = defaultConfig
		}
	}
	return svc
}

// Evaluate runs all rules against the content of the file
func (s *CheckerService) Evaluate(baseDir string, filePath string) []domain.CheckerResult {
	results := make([]domain.CheckerResult, 0)

	if filePath == "" || filePath == "auto-inferred" {
		for _, r := range s.config.Rules {
			results = append(results, domain.CheckerResult{RuleName: r.Name, Passed: true, Message: "SKIP: No file context provided."})
		}
		return results
	}

	fullPath := filepath.Join(baseDir, filePath)
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		for _, r := range s.config.Rules {
			results = append(results, domain.CheckerResult{RuleName: r.Name, Passed: false, Message: "FAIL: File missing or access denied."})
		}
		return results
	}

	for _, rule := range s.config.Rules {
		res := domain.CheckerResult{RuleName: rule.Name, Passed: false}
		if rule.Regex != "" {
			matched, _ := regexp.Match(rule.Regex, content)
			if matched {
				res.Passed = true
				res.Message = "PASS"
			} else {
				res.Message = "FAIL: Missing required pattern in content."
			}
		} else {
			res.Passed = true
			res.Message = "PASS (No logic)"
		}
		results = append(results, res)
	}

	return results
}
