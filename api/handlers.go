package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"openspec-visualizer/domain"
	"openspec-visualizer/service"
)

type Handlers struct {
	fsService      *service.FSService
	llmService     *service.LLMService
	reviewer       domain.Reviewer
	archiveService *service.ArchiveService
}

func NewHandlers(fs *service.FSService, llm *service.LLMService, rev domain.Reviewer, archive *service.ArchiveService) *Handlers {
	return &Handlers{
		fsService:      fs,
		llmService:     llm,
		reviewer:       rev,
		archiveService: archive,
	}
}

func SetupRouter(h *Handlers) *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	api := app.Group("/api")

	api.Post("/report", h.HandleReport)
	api.Get("/reports", h.GetReports)
	api.Get("/pending", h.GetPending)
	api.Post("/review", h.SubmitReview)
	
	// History Endpoints
	api.Get("/history", h.GetHistory)
	api.Get("/diff", h.GetDiff)
	api.Post("/rollback", h.Rollback)
	
	// Legacy endpoints mapping
	api.Get("/config", h.GetConfig)
	api.Post("/config", h.PostConfig)
	api.Post("/generate", h.Generate)
	api.Get("/list", h.List)
	api.Get("/read", h.Read)
	api.Post("/prompt", h.Prompt)

	return app
}

func (h *Handlers) HandleReport(c *fiber.Ctx) error {
	var req domain.ReportRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无法解析请求数据: " + err.Error(),
		})
	}

	response, err := h.reviewer.Review(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "评审过程出错: " + err.Error(),
		})
	}

	if !response.Approved {
		return c.Status(fiber.StatusBadRequest).JSON(response) // 400 with feedback/advice
	}
	
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *Handlers) GetReports(c *fiber.Ctx) error {
	if rev, ok := h.reviewer.(*service.InteractiveReviewer); ok {
		reports := rev.GetReports()
		return c.JSON(reports)
	}
	return c.JSON([]domain.ReportRequest{})
}

func (h *Handlers) GetPending(c *fiber.Ctx) error {
	if rev, ok := h.reviewer.(*service.InteractiveReviewer); ok {
		pending := rev.GetPending()
		if pending != nil {
			return c.JSON(pending)
		}
	}
	return c.JSON(fiber.Map{})
}

func (h *Handlers) SubmitReview(c *fiber.Ctx) error {
	var dec domain.ReviewDecision
	if err := c.BodyParser(&dec); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
	}

	if rev, ok := h.reviewer.(*service.InteractiveReviewer); ok {
		rev.SubmitDecision(dec)
		return c.SendStatus(fiber.StatusOK)
	}
	return c.SendStatus(fiber.StatusInternalServerError)
}

func (h *Handlers) GetConfig(c *fiber.Ctx) error {
	cfg := h.llmService.LoadLLMConfig()
	return c.JSON(cfg)
}

func (h *Handlers) PostConfig(c *fiber.Ctx) error {
	var p struct{ APIKey, BaseURL, Model string }
	if err := c.BodyParser(&p); err == nil {
		h.llmService.SaveLLMConfig(p.APIKey, p.BaseURL, p.Model)
		return c.SendStatus(fiber.StatusOK)
	}
	return c.SendStatus(fiber.StatusBadRequest)
}

func (h *Handlers) Generate(c *fiber.Ctx) error {
	var p struct{ FeatureName, Content string }
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	err := h.fsService.GenerateOpenSpecStructure(p.FeatureName, p.Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handlers) List(c *fiber.Ctx) error {
	node, err := h.fsService.ListOpenSpecFiles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(node)
}

func (h *Handlers) Read(c *fiber.Ctx) error {
	path := c.Query("path")
	data, err := h.fsService.ReadFileContent(path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString(data)
}

func (h *Handlers) ReadFile(c *fiber.Ctx) error {
	pathStr := c.Query("path")
	if pathStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing path")
	}

	content, err := h.fsService.ReadFileContent(pathStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString(content)
}

func (h *Handlers) GetHistory(c *fiber.Ctx) error {
	history, err := h.archiveService.GetHistory()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(history)
}

func (h *Handlers) GetDiff(c *fiber.Ctx) error {
	hash := c.Query("hash")
	if hash == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Hash required")
	}
	diff, err := h.archiveService.GetDiff(hash)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString(diff)
}

func (h *Handlers) Rollback(c *fiber.Ctx) error {
	hash := c.Query("hash")
	if hash == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Hash required")
	}
	err := h.archiveService.Rollback(hash)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handlers) Prompt(c *fiber.Ctx) error {
	var p struct{ Prompt, System string }
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	res, err := h.llmService.SendPrompt(p.Prompt, p.System)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(fiber.Map{"result": res})
}
