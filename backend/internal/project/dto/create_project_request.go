package dto

type CreateProjectRequest struct {
	Name   string            `json:"name" binding:"required"`
	GitURL string            `json:"git_url" binding:"required,url"`
	Branch string            `json:"branch" binding:"required"`
	Env    map[string]string `json:"env"`
}
