package main

import (
	"embed"
	"encoding/json"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

// SpecTask represents a single task in OpenSpec tasks.md
type SpecTask struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// OpenSpecService provides core APIs to the frontend
type OpenSpecService struct{}

func (o *OpenSpecService) GetConfig() LLMConfig {
	return LoadLLMConfig()
}

func (o *OpenSpecService) SaveConfig(apiKey, baseUrl, model string) error {
	return SaveLLMConfig(LLMConfig{
		APIKey:  apiKey,
		BaseURL: baseUrl,
		Model:   model,
	})
}

// RunPrompt takes a prompt and system instruction and returns the LLM response
func (o *OpenSpecService) RunPrompt(prompt string, system string) (string, error) {
	return SendPrompt(prompt, system)
}

// GetTasks mocks the task parsing from openspec/current/tasks.md
func (o *OpenSpecService) GetTasks() string {
	tasks := []SpecTask{
		{Title: "初始化 OpenSpec 目录结构", Completed: true},
		{Title: "解析 proposal.md 提取意图", Completed: true},
		{Title: "根据当前架构评估系统提示词", Completed: true},
		{Title: "生成针对 AI 上下文的规范差异", Completed: false},
		{Title: "将界面组件逻辑合并回主分支", Completed: false},
	}
	bytes, _ := json.Marshal(tasks)
	return string(bytes)
}

func main() {
	svc := &OpenSpecService{}
	StartLocalServer(svc)

	app := application.New(application.Options{
		Name:        "OpenSpec 可视化工具",
		Description: "OpenSpec SDD 工作流可视化仪表盘",
		Assets:      application.AssetOptions{Handler: application.AssetFileServerFS(assets)},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Services: []application.Service{
			application.NewService(svc),
		},
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "OpenSpec 可视化工具",
		Width:  1024,
		Height: 768,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
