package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sowncns/k3s-deploy-platform/internal/auth"
	"github.com/sowncns/k3s-deploy-platform/internal/deployment"
	"github.com/sowncns/k3s-deploy-platform/internal/github"
	"github.com/sowncns/k3s-deploy-platform/internal/middleware"
	"github.com/sowncns/k3s-deploy-platform/internal/project"
)

type Handlers struct {
     Project    *project.Handler
     Deployment *deployment.Handler
     GitHub     *github.Handler
     AuthRepo   auth.AuthRepository
}


func Setup(h Handlers) *gin.Engine {
    r := gin.Default()

    // Health check cho K3s liveness/readiness probe
    r.GET("/healthz", func(c *gin.Context) { c.Status(200) })

    api := r.Group("/api/v1")
    {
        // Đăng nhập GitHub không cần token (dùng để lấy token) nên để ngoài middleware
        h.GitHub.RegisterPublicRoutes(api)

        // Các route còn lại yêu cầu header: Authorization: Bearer <github_access_token>
        protected := api.Group("")
        protected.Use(middleware.AuthMiddleware(h.AuthRepo))
        {
            h.Project.RegisterRoutes(protected)
            h.Deployment.RegisterRoutes(protected)
            h.GitHub.RegisterRoutes(protected)
        }
    }

    return r
}