package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sowncns/k3s-deploy-platform/internal/cluster"
	"github.com/sowncns/k3s-deploy-platform/internal/config"
	"github.com/sowncns/k3s-deploy-platform/internal/database"
	"github.com/sowncns/k3s-deploy-platform/internal/k8s"
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

	if err := clusterManager.Register(ctx, &cluster.Cluster{
		ID:            "cluster-001",
		Name:          "cluster-001",
		Provider:      cluster.ClusterProviderK3s,
		Status:        cluster.ClusterStatusActive,
		CredentialRef: kubeconfig,
	}); err != nil {
		log.Fatal(err)
	}

	// Get cluster client
	client, err := clusterManager.Get(ctx, "cluster-001")
	if err != nil {
		log.Fatal(err)
	}

	// Test deployment
	deployment, err := client.ApplyDeployment(
		ctx,
		k8s.DeploymentConfig{
			Name:          "nginx",
			Namespace:     "platform-dev",
			Image:         "nginx:latest",
			Replicas:      2,
			ContainerPort: 80,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deployment created:", deployment.Name)
	fmt.Println("Namespace:", deployment.Namespace)
}