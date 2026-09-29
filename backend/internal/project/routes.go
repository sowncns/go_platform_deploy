package project

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	projectGroup := r.Group("/projects")
	
	{
		projectGroup.POST("", h.CreateProject)
		projectGroup.GET("", h.ListProjectsByUserID)
		projectGroup.GET("/:id", h.FindProjectByID)
		projectGroup.PATCH("/:id", h.UpdateProject)
	}
}
