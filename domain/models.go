package domain

// SpecTask represents a single task in OpenSpec tasks.md
type SpecTask struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// LLMConfig represents the AI configuration
type LLMConfig struct {
	APIKey  string `json:"apiKey"`
	BaseURL string `json:"baseUrl"`
	Model   string `json:"model"`
}

// FileNode represents a node in the file tree
type FileNode struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	IsDir    bool        `json:"isDir"`
	Children []*FileNode `json:"children"`
}

// ReportRequest is the payload sent by the Agent
type ReportRequest struct {
	SkillName string `json:"skill_name"`
	Status    string `json:"status"`
	FilePath  string `json:"file_path"`
}

// ReportResponse is what the system responds to the Agent
type ReportResponse struct {
	Approved bool   `json:"approved"`
	Feedback string `json:"feedback,omitempty"`
}

// Reviewer defines the interface for reviewing Agent tasks
type Reviewer interface {
	Review(req ReportRequest) (*ReportResponse, error)
}
