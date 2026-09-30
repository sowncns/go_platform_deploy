package deployment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sowncns/k3s-deploy-platform/internal/deployment/dto"
)


type Handler struct {
	service DeploymentService
}

func NewHandler(service DeploymentService) *Handler {
	return &Handler{service: service}
}


 func (h * Handler) CreateDeployment (c *gin.Context) {
	clusterID := c.Param("clusterId")

	var req dto.DeploymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	deployment , err := h.service.Create(c.Request.Context(),clusterID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, deployment)

 }