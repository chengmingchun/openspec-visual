package cmd

import (
	"embed"
	"encoding/json"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"

	"openspec-visualizer/api"
	"openspec-visualizer/domain"
	"openspec-visualizer/service"
)

type OpenSpecWailsService struct {
	llmService *service.LLMService
}

func (o *OpenSpecWailsService) GetConfig() domain.LLMConfig {
	return o.llmService.LoadLLMConfig()
}

func (o *OpenSpecWailsService) SaveConfig(apiKey, baseUrl, model string) error {
	return o.llmService.SaveLLMConfig(apiKey, baseUrl, model)
}

func (o *OpenSpecWailsService) RunPrompt(prompt string, system string) (string, error) {
	return o.llmService.SendPrompt(prompt, system)
}

func (o *OpenSpecWailsService) GetTasks() string {
	tasks := []domain.SpecTask{
		{Title: "初始化 OpenSpec 目录结构", Completed: true},
		{Title: "解析 proposal.md 提取意图", Completed: true},
		{Title: "根据当前架构评估系统提示词", Completed: true},
		{Title: "生成针对 AI 上下文的规范差异", Completed: false},
		{Title: "将界面组件逻辑合并回主分支", Completed: false},
	}
	bytes, _ := json.Marshal(tasks)
	return string(bytes)
}

func Run(assets embed.FS) {
	fsSvc := service.NewFSService()
	llmSvc := service.NewLLMService()
	archiveSvc := service.NewArchiveService()
	mockReviewer := service.NewMockReviewer()

	apiHandlers := api.NewHandlers(fsSvc, llmSvc, mockReviewer, archiveSvc)
	fiberApp := api.SetupRouter(apiHandlers)

	go func() {
		log.Println("Starting Fiber Agent server on 127.0.0.1:38192")
		if err := fiberApp.Listen("127.0.0.1:38192"); err != nil {
			log.Fatalf("Fiber Error: %v", err)
		}
	}()

	wailsSvc := &OpenSpecWailsService{
		llmService: llmSvc,
	}

	app := application.New(application.Options{
		Name:        "OpenSpec 可视化工具",
		Description: "OpenSpec SDD 工作流可视化仪表盘",
		Assets:      application.AssetOptions{Handler: application.AssetFileServerFS(assets)},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Services: []application.Service{
			application.NewService(wailsSvc),
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
