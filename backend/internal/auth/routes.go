package auth

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	
	{
		authGroup.GET("/github",h.GitHubLogin)
		authGroup.GET("/github/callback",h.GitHubCallback)
	}
}
