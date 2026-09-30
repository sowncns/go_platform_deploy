package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sowncns/k3s-deploy-platform/internal/auth"
	"github.com/sowncns/k3s-deploy-platform/internal/deployment"
	"github.com/sowncns/k3s-deploy-platform/internal/project"
)

type Handlers struct {
     Project    *project.Handler
     Deployment *deployment.Handler
     Auth *auth.Handler
}


func Setup(h Handlers) *gin.Engine {
    r := gin.Default()

    // Health check cho K3s liveness/readiness probe
    r.GET("/healthz", func(c *gin.Context) { c.Status(200) })

    api := r.Group("/api/v1")
    {
        h.Project.RegisterRoutes(api)
        h.Auth.RegisterRoutes(api)
    }

    return r
}