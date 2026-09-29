package dto

type UpdateProjectRequest struct {
	Name   *string           `json:"name"`
	GitURL *string           `json:"git_url" binding:"omitempty,url"`
	Branch *string           `json:"branch"`
	Env    map[string]string `json:"env"`
}
