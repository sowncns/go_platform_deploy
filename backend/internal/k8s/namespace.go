package k8s

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

func (c *Client) EnsureNamespace(
	ctx context.Context,
	name string,
) (*corev1.Namespace, error) {
	namespaces := c.Kube.CoreV1().Namespaces()

	// Kiểm tra namespace đã tồn tại chưa
	current, err := namespaces.Get(
		ctx,
		name,
		metav1.GetOptions{},
	)

	if err == nil {
		return current, nil
	}

	// Không tồn tại → tạo mới
	if apierrors.IsNotFound(err) {
		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}

		created, err := namespaces.Create(
			ctx,
			namespace,
			metav1.CreateOptions{},
		)
		if err != nil {
			return nil, fmt.Errorf(
				"create namespace %s: %w",
				name,
				err,
			)
		}

		return created, nil
	}

	return nil, fmt.Errorf(
		"get namespace %s: %w",
		name,
		err,
	)
}