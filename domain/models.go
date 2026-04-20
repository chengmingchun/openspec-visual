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

// AgentSkill defines a capability of the Agent
type AgentSkill struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Trigger     string `yaml:"trigger"`
}

// AgentConfig maps to openspec.yaml
type AgentConfig struct {
	Endpoint           string       `yaml:"endpoint"`
	GlobalInstructions string       `yaml:"global_instructions"`
	Skills             []AgentSkill `yaml:"skills"`
}

// ReviewDecision represents the user's manual response
type ReviewDecision struct {
	Approved bool   `json:"approved"`
	Feedback string `json:"feedback"`
}

// TDDRule represents a single rule logic for TDD validation
type TDDRule struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Regex       string `yaml:"regex"`
}

// TDDConfig represents tdd_rules.yaml
type TDDConfig struct {
	Rules []TDDRule `yaml:"rules"`
}

// CheckerResult represents the result of evaluating a TDDRule
type CheckerResult struct {
	RuleName string `json:"rule_name"`
	Passed   bool   `json:"passed"`
	Message  string `json:"message"`
}

// PendingReport wraps the agent report with automated checker results
type PendingReport struct {
	Request       ReportRequest   `json:"request"`
	CheckerResult []CheckerResult `json:"checker_results"`
}
