package node

import "context"

type NodeRepository interface {
	Create(ctx context.Context, node *Node) error
	Update(ctx context.Context, node *Node) error

	CreateOrUpdate(ctx context.Context, node *Node) error

	FindByID(ctx context.Context, id string) (*Node, error)

	ListByClusterID(
		ctx context.Context,
		clusterID string,
	) ([]*Node, error)

	DeleteByClusterID(
		ctx context.Context,
		clusterID string,
	) error
}


