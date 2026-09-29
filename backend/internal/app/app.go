// package app

// import (
// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"github.com/sowncns/k3s-deploy-platform/internal/config"
// 	"github.com/sowncns/k3s-deploy-platform/internal/k3s"
// 	"github.com/sowncns/k3s-deploy-platform/internal/project"
// 	"github.com/sowncns/k3s-deploy-platform/internal/router"
// )

// // BuildHandlers khởi tạo toàn bộ Repo -> Service -> Handler của hệ thống
// func BuildHandlers(pool *pgxpool.Pool, k3sClient *kubernetes.Client, cfg *config.Config) router.Handlers {
// 	// 1. User Module
// 	userRepo := user.NewRepository(pool)
// 	userService := user.NewService(userRepo)
// 	userHandler := user.NewHandler(userService)

// 	// 2. Project Module
// 	projectRepo := project.NewRepository(pool)
// 	projectService := project.NewService(projectRepo, k3sClient)
// 	projectHandler := project.NewHandler(projectService)

// 	// Sau này thêm Deployment, Secret, Webhook chỉ cần new ở đây:
// 	// depRepo := deployment.NewRepository(pool)
// 	// depService := deployment.NewService(depRepo, k3sClient)
// 	// depHandler := deployment.NewHandler(depService)

// 	return router.Handlers{
// 		User:    userHandler,
// 		Project: projectHandler,
// 		// Deployment: depHandler,
// 	}
// }