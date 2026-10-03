package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sowncns/k3s-deploy-platform/internal/auth"
	"github.com/sowncns/k3s-deploy-platform/internal/cluster"
	"github.com/sowncns/k3s-deploy-platform/internal/config"
	"github.com/sowncns/k3s-deploy-platform/internal/deployment"
	"github.com/sowncns/k3s-deploy-platform/internal/github"
	"github.com/sowncns/k3s-deploy-platform/internal/router"
	"github.com/sowncns/k3s-deploy-platform/internal/project"
)

// BuildHandlers khởi tạo toàn bộ Repo -> Service -> Handler của hệ thống
func BuildHandlers(pool *pgxpool.Pool, clusterManager *cluster.Manager, cfg *config.Config) router.Handlers {

	depRepo := deployment.NewRepository(pool)
	depService := deployment.NewService(clusterManager, depRepo)
	depHandler := deployment.NewHandler(depService)


	
	authRepo := auth.NewRepository(pool)
	authService := auth.NewService(authRepo)
	githubHandler := github.NewHandler(authService, authRepo)



	projectRepo := project.NewRepository(pool)
	projectService := project.NewService(projectRepo)
	projectHandler := project.NewHandler(projectService)
	return router.Handlers{
		Deployment: depHandler,
		GitHub:     githubHandler,
		Project:    projectHandler,
		AuthRepo:   authRepo,
	}
}