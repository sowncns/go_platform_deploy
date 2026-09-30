package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) ListNodes(
	ctx context.Context,
) ([]corev1.Node, error) {
	nodes, err := c.Kube.CoreV1().
		Nodes().
		List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("list kubernetes nodes: %w", err)
	}

	return nodes.Items, nil
}