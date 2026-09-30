package node

import (
	"context"
	"fmt"
	"github.com/sowncns/k3s-deploy-platform/internal/cluster"
)

type Service struct {
	clusterManager *cluster.Manager
	repo           NodeRepository
}

func NewService(
	clusterManager *cluster.Manager,
	repo NodeRepository,
) *Service {
	return &Service{
		clusterManager: clusterManager,
		repo:           repo,
	}
}

func (s *Service) SyncClusterNodes(
	ctx context.Context,
	clusterID string,
) error {
	client, err := s.clusterManager.Get(ctx, clusterID)
	if err != nil {
		return fmt.Errorf(
			"get cluster %s: %w",
			clusterID,
			err,
		)
	}

	nodes, err := client.ListNodes(ctx)
	if err != nil {
		return fmt.Errorf(
			"list nodes for cluster %s: %w",
			clusterID,
			err,
		)
	}
	
	for i := range nodes {
		node := mapK8sNode(
			clusterID,
			&nodes[i],
		)
		if err := s.repo.CreateOrUpdate(ctx, node); err != nil {
			return fmt.Errorf(
				"sync node %s: %w",
				node.Name,
				err,
			)
		}
	}

	return nil
}