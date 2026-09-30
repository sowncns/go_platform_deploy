package deployment

import (
	"context"

	"github.com/sowncns/k3s-deploy-platform/internal/cluster"
	"github.com/sowncns/k3s-deploy-platform/internal/deployment/dto"
	"github.com/sowncns/k3s-deploy-platform/internal/k8s"
	appsv1 "k8s.io/api/apps/v1"
)

type DeploymentService interface {
	Create(ctx context.Context, clusterID string, req dto.DeploymentRequest) (*appsv1.Deployment, error)
}

type Service struct {
	clusterManager *cluster.Manager
	repo           DeploymentRepository
}

func NewService(clusterManager *cluster.Manager, repo DeploymentRepository) *Service {
	return &Service{
		clusterManager: clusterManager,
		repo:           repo,
	}
}

func (s *Service) Create(ctx context.Context, clusterID string, req dto.DeploymentRequest) (*appsv1.Deployment, error) {
	dbDeployment := &Deployment{
		ProjectID: req.ProjectID,
		CommitSHA: req.CommitSHA,
		ImageTag:  req.ImageTag,
		Status:    DeployQueued,
	}

	if err := s.repo.Create(ctx, dbDeployment); err != nil {
		return nil, err
	}

	client, err := s.clusterManager.Get(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	k8sDeployment, err := client.ApplyDeployment(
		ctx,
		k8s.DeploymentConfig{
			Name:          req.Name,
			Namespace:     req.Namespace,
			Image:         req.Image,
			Replicas:      req.Replicas,
			ContainerPort: req.ContainerPort,
		},
	)
	if err != nil {
		return nil, err
	}

	if _, err := client.ApplyService(ctx,
		k8s.ServiceConfig{
			Name:       req.Name,
			Namespace:  req.Namespace,
			Port:       req.Port,
			TargetPort: req.TargetPort,
		},
	); err != nil {
		return nil, err
	}

	if _, err := client.ApplyIngress(ctx,
		k8s.IngressConfig{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Host:        req.Host,
			ServiceName: req.ServiceName,
			ServicePort: req.ServicePort,
		},
	); err != nil {
		return nil, err
	}

	return k8sDeployment, nil
}
