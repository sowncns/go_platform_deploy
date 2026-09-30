package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sowncns/k3s-deploy-platform/internal/auth"
	"github.com/sowncns/k3s-deploy-platform/internal/cluster"
	"github.com/sowncns/k3s-deploy-platform/internal/config"
	"github.com/sowncns/k3s-deploy-platform/internal/deployment"
	"github.com/sowncns/k3s-deploy-platform/internal/github"
	"github.com/sowncns/k3s-deploy-platform/internal/router"
)

// BuildHandlers khởi tạo toàn bộ Repo -> Service -> Handler của hệ thống
func BuildHandlers(pool *pgxpool.Pool, clusterManager *cluster.Manager, cfg *config.Config) router.Handlers {

	depRepo := deployment.NewRepository(pool)
	depService := deployment.NewService(clusterManager, depRepo)
	depHandler := deployment.NewHandler(depService)


	

	githubOAuth := &github.OAuthConfig{
		ClientID:     cfg.GitHubClientID,
		ClientSecret: cfg.GitHubClientSecret,
		RedirectURL:  cfg.GitHubRedirectURL,
	}
	authRepo := auth.NewRepository(pool)
	authService := auth.NewService(githubOAuth,authRepo)
	authHandler := auth.NewHandler(authService)

	return router.Handlers{
		Deployment: depHandler,
		Auth : authHandler,
		
	}
}