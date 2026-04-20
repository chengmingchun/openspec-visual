package service

import (
	"log"
	"os"
	"sync"

	"openspec-visualizer/domain"
)

// InteractiveReviewer suspends agent reports until manual review
type InteractiveReviewer struct {
	mu            sync.Mutex
	reports       []domain.ReportRequest
	pending       *domain.PendingReport
	decisionChan  chan domain.ReviewDecision
	checkerSvc    *CheckerService
	baseDir       string
}

// NewInteractiveReviewer creates a new blocking Reviewer
func NewInteractiveReviewer() *InteractiveReviewer {
	cwd, _ := os.Getwd()
	return &InteractiveReviewer{
		reports:      make([]domain.ReportRequest, 0),
		decisionChan: make(chan domain.ReviewDecision, 1),
		checkerSvc:   NewCheckerService(cwd),
		baseDir:      cwd,
	}
}

// Review overrides and evaluates the Agent's report.
// Blocks until user interacts in dashboard.
func (r *InteractiveReviewer) Review(req domain.ReportRequest) (*domain.ReportResponse, error) {
	log.Printf("收到 Agent 回调汇报: 打分项=%s, 文件=%s. 正在启动多方 TDD Review 并挂起等待人工...\n", req.SkillName, req.FilePath)
	
	// 1. Run Automated Checker
	checkResults := r.checkerSvc.Evaluate(r.baseDir, req.FilePath)

	r.mu.Lock()
	r.reports = append(r.reports, req)
	
	r.pending = &domain.PendingReport{
		Request:       req,
		CheckerResult: checkResults,
	}
	r.mu.Unlock()

	// 2. Block until humanity decides from UI dashboard
	decision := <-r.decisionChan

	// 3. Clear out pending status
	r.mu.Lock()
	r.pending = nil
	r.mu.Unlock()

	return &domain.ReportResponse{
		Approved: decision.Approved,
		Feedback: decision.Feedback,
	}, nil
}

// GetReports returns the history of Agent reports
func (r *InteractiveReviewer) GetReports() []domain.ReportRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	copyReports := make([]domain.ReportRequest, len(r.reports))
	copy(copyReports, r.reports)
	return copyReports
}

// GetPending returns the currently blocked report
func (r *InteractiveReviewer) GetPending() *domain.PendingReport {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.pending
}

// SubmitDecision resolves the blocked Review request
func (r *InteractiveReviewer) SubmitDecision(decision domain.ReviewDecision) {
	select {
	case r.decisionChan <- decision:
	default:
		// No-op if no one is waiting
	}
}
