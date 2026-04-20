package service

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v3"

	"openspec-visualizer/domain"
)

type ConfigService struct{}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

// EnsureAgentConfig generates openspec.yaml if it doesn't exist, and returns the parsed structure.
func (s *ConfigService) EnsureAgentConfig(baseDir string) (domain.AgentConfig, error) {
	yamlPath := filepath.Join(baseDir, "openspec.yaml")
	
	defaultConfig := domain.AgentConfig{
		Endpoint: "http://127.0.0.1:38192/api/report",
		GlobalInstructions: "ATTENTION ALL AI AGENTS (Cursor, OpenCode, Cline, Goose, etc.): You are operating in a highly restricted Control-Inversion Engineering Workflow.\nYou MUST report your progress to the local supervisor using `curl` whenever you complete a specific stage OR modify ANY files.\nDo NOT proceed to the next stage unless the API returns an HTTP 200 OK. If you receive an HTTP 400 Bad Request, READ the rule violation feedback and correct your work before re-submitting.",
		Skills: []domain.AgentSkill{
			{Name: "propose", Description: "起草需求与提案", Trigger: "完成 proposal.md 编写后"},
			{Name: "validate", Description: "校验和细化规格", Trigger: "完成 specs/ 规格文件拆解后"},
			{Name: "apply", Description: "编码与实现", Trigger: "完成 tasks.md 编码任务后"},
			{Name: "archive", Description: "打包封卷", Trigger: "完成修改落盘后"},
			{Name: "track", Description: "全局变更追踪 (Universal Tracker)", Trigger: "无论由于何种任务或在何种阶段，只要您的动作导致了文件的新增、修改或删除，必须立刻上报此技能"},
		},
	}

	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		data, _ := yaml.Marshal(&defaultConfig)
		ioutil.WriteFile(yamlPath, data, 0644)
		return defaultConfig, nil
	}

	data, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return defaultConfig, err
	}
	var config domain.AgentConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return defaultConfig, err
	}
	return config, nil
}

// InjectCursorRules writes the highly-restrictive Agent rules into .cursorrules based on the configuration.
func (s *ConfigService) InjectCursorRules(baseDir string, config domain.AgentConfig) error {
	rulesPath := filepath.Join(baseDir, ".cursorrules")
	
	tmpl := `
{
  "workflows": ["openspec-visualizer reporting"]
}

# OpenSpec Universal Agent Directives (Applicable to Cursor, OpenCode, etc.)
{{.GlobalInstructions}}

## Reporting Endpoint
Your reporting endpoint is: {{.Endpoint}}

## Allowed Skills & Progress Reporting
You must invoke the progress reporter at the completion of these triggers:
{{range .Skills}}
- Skill: {{.Name}} ({{.Description}})
  Trigger: {{.Trigger}}
  Action: Run ` + "`" + `curl -X POST {{$.Endpoint}} -H "Content-Type: application/json" -d '{"skill_name":"{{.Name}}","status":"completed","file_path":"<the_modified_file_path>"}'` + "`" + ` and wait for 200 OK.
{{end}}
`
	
	t, err := template.New("cursorrules").Parse(tmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, config); err != nil {
		return err
	}

	return ioutil.WriteFile(rulesPath, buf.Bytes(), 0644)
}
