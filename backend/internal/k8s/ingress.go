package k3s

import (
	"context"
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/util/intstr"
)

type IngressConfig struct {
	Name        string
	Namespace   string
	Host        string
	ServiceName string
	ServicePort int32
}

func (c *Client) ApplyIngress(ctx context.Context, cfg IngressConfig,
) (*networkingv1.Ingress, error) {
	ingresses := c.Kube.NetworkingV1().
		Ingresses(cfg.Namespace)

	ingress := buildIngress(cfg)

	current, err := ingresses.Get(
		ctx,
		cfg.Name,
		metav1.GetOptions{},
	)

	if err != nil {
		if apierrors.IsNotFound(err) {
			created, err := ingresses.Create(
				ctx,
				ingress,
				metav1.CreateOptions{},
			)
			if err != nil {
				return nil, fmt.Errorf(
					"create ingress %s/%s: %w",
					cfg.Namespace,
					cfg.Name,
					err,
				)
			}

			return created, nil
		}

		return nil, fmt.Errorf(
			"get ingress %s/%s: %w",
			cfg.Namespace,
			cfg.Name,
			err,
		)
	}

	ingress.ResourceVersion = current.ResourceVersion

	updated, err := ingresses.Update(
		ctx,
		ingress,
		metav1.UpdateOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"update ingress %s/%s: %w",
			cfg.Namespace,
			cfg.Name,
			err,
		)
	}

	return updated, nil
}

func buildIngress(cfg IngressConfig) *networkingv1.Ingress {
	pathType := networkingv1.PathTypePrefix

	return &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfg.Name,
			Namespace: cfg.Namespace,
		},

		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: cfg.Host,

					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,

									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: cfg.ServiceName,

											Port: networkingv1.ServiceBackendPort{
												Number: cfg.ServicePort,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
