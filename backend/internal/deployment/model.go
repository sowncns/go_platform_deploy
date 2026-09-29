package deployment

import "time"

type DeploymentStatus string

const (
	DeployQueued     DeploymentStatus = "queued"
	DeployInProgress DeploymentStatus = "in_progress"
	DeploySuccess    DeploymentStatus = "success"
	DeployFailed     DeploymentStatus = "failed"
)

type Deployment struct {
	ID         uint             `json:"id"`
	ProjectID  uint             `json:"project_id"`
	CommitSHA  string           `json:"commit_sha"`
	CommitMsg  string           `json:"commit_message"`
	ImageTag   string           `json:"image_tag"` // Tag trên GHCR
	Status     DeploymentStatus `json:"status"`
	BuildLogs  string           `json:"build_logs,omitempty"`
	DeployedAt *time.Time       `json:"deployed_at"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}
