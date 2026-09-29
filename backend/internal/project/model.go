package project

import (
	"time"
)

type ProjectStatus string

const (
	StatusPending  ProjectStatus = "pending"
	StatusBuilding ProjectStatus = "building"
	StatusRunning  ProjectStatus = "running"
	StatusFailed   ProjectStatus = "failed"
)

type Project struct {
	ID           uint              `json:"id"`
	UserID       uint              `json:"user_id"`
	Name         string            `json:"name"`
	RepoFullName string            `json:"repo_full_name"` // ví dụ: "octocat/hello-world"
	GitURL       string            `json:"git_url"`
	Branch       string            `json:"branch"`
	Port         int               `json:"port"`          // Port app lắng nghe trong container
	Domain       string            `json:"domain"`        // Domain map vào Ingress
	Namespace    string            `json:"namespace"`     // K8s namespace (thường đặt theo project slug)
	CurrentImage string            `json:"current_image"` // Image tag đang chạy trên k3s
	Status       ProjectStatus     `json:"status"`
	WebhookID    int64             `json:"webhook_id"` // Để tự động xoá hook trên GitHub khi xoá project
	Env          map[string]string `json:"-"`          // Biến môi trường của project, giấu khỏi JSON response để bảo mật

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
