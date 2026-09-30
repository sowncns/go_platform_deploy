package deployment

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	deployments := rg.Group("/clusters")
	{
		deployments.POST("/:clusterId/deployments", h.CreateDeployment)
	}
}
