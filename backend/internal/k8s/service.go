package k8s

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type ServiceConfig struct {
	Name       string
	Namespace  string
	Port       int32
	TargetPort int32
}

func (c *Client) ApplyService(
	ctx context.Context,
	cfg ServiceConfig,
) (*corev1.Service, error) {
	services := c.Kube.CoreV1().
		Services(cfg.Namespace)

	service := buildService(cfg)

	current, err := services.Get(
		ctx,
		cfg.Name,
		metav1.GetOptions{},
	)

	if err != nil {
		if apierrors.IsNotFound(err) {
			created, err := services.Create(
				ctx,
				service,
				metav1.CreateOptions{},
			)
			if err != nil {
				return nil, fmt.Errorf(
					"create service %s/%s: %w",
					cfg.Namespace,
					cfg.Name,
					err,
				)
			}

			return created, nil
		}

		return nil, fmt.Errorf(
			"get service %s/%s: %w",
			cfg.Namespace,
			cfg.Name,
			err,
		)
	}

	service.ResourceVersion = current.ResourceVersion

	updated, err := services.Update(
		ctx,
		service,
		metav1.UpdateOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"update service %s/%s: %w",
			cfg.Namespace,
			cfg.Name,
			err,
		)
	}

	return updated, nil
}

func buildService(cfg ServiceConfig) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfg.Name,
			Namespace: cfg.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": cfg.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       cfg.Port,
					TargetPort: intstr.FromInt32(cfg.TargetPort),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}
}
