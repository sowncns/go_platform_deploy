package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sowncns/k3s-deploy-platform/internal/k3s"
)

func main() {
	ctx := context.Background()

	client, err := k3s.NewClient()
	if err != nil {
		log.Fatal(err)
	}


	deployment, err := client.ApplyDeployment(
		ctx,
		k3s.DeploymentConfig{
			Name:          "nginx",
			Namespace:     "platform-dev",
			Image:         "nginx:latest",
			Replicas:      2,
			ContainerPort: 80,
		},
	)

	namespace , err := client.EnsureNamespace(ctx,"test-platform")

	ingress, err := client.ApplyIngress(ctx,
	k3s.IngressConfig{
		Name :		"my-api",
		Namespace:	"test-platform",
		Host :		"api.example.com",
		ServiceName:	"test",
		ServicePort:	876,
	},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("deployment created:", deployment.Name)
	fmt.Println("Namesapce Check:", namespace.Name)
	fmt.Println("Ingress Check:", ingress.Name)


	

}