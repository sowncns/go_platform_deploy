package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/sowncns/k3s-deploy-platform/internal/app"
	"github.com/sowncns/k3s-deploy-platform/internal/cluster"
	"github.com/sowncns/k3s-deploy-platform/internal/config"
	"github.com/sowncns/k3s-deploy-platform/internal/database"
	"github.com/sowncns/k3s-deploy-platform/internal/router"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()
	db, err := database.NewPostgresPool(ctx, cfg)
	
	if err != nil {
		log.Fatal(err)
	}
	
	defer db.Close()

	// Create cluster manager
	clusterManager := cluster.NewManager(cluster.NewRepository(db))

	// Register cluster (kubeconfig mac dinh ~/.kube/config)
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// if err := clusterManager.Register(ctx, &cluster.Cluster{
	// 	ID:            "cluster-003",
	// 	Name:          "cluster-003",
	// 	Provider:      cluster.ClusterProviderK3s,
	// 	Status:        cluster.ClusterStatusActive,
	// 	CredentialRef: kubeconfig,
	// }); err != nil {
	// 	log.Fatal(err)
	// }

	handlers := app.BuildHandlers(db, clusterManager, cfg)

	r := router.Setup(handlers)
	if err := r.Run(":"+ config.Load().AppPort); err != nil {
		log.Fatal(err)
	}
}