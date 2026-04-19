package service

import (
	"log"
	"sync"

	"openspec-visualizer/domain"
)

// MockReviewer is a mock implementation of the Reviewer interface
type MockReviewer struct {
	mu      sync.Mutex
	reports []domain.ReportRequest
}

// NewMockReviewer creates a new MockReviewer
func NewMockReviewer() *MockReviewer {
	return &MockReviewer{
		reports: make([]domain.ReportRequest, 0),
	}
}

// Review evaluates the Agent's report.
// It acts as a State Machine interceptor.
// Expected to be replaced by a real scoring engine later.
func (r *MockReviewer) Review(req domain.ReportRequest) (*domain.ReportResponse, error) {
	log.Printf("收到 Agent 回调汇报: 打分项=%s, 状态=%s, 文件=%s\n", req.SkillName, req.Status, req.FilePath)
	
	r.mu.Lock()
	r.reports = append(r.reports, req)
	r.mu.Unlock()

	// Mock implementation: always approve
	return &domain.ReportResponse{
		Approved: true,
		Feedback: "通过代码规范检查，允许继续！",
	}, nil
}

// GetReports returns the history of Agent reports
func (r *MockReviewer) GetReports() []domain.ReportRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Return a copy to prevent race conditions
	copyReports := make([]domain.ReportRequest, len(r.reports))
	copy(copyReports, r.reports)
	return copyReports
}
