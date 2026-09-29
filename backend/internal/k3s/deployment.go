package k3s

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type DeploymentConfig struct {
	Name          string
	Namespace     string
	Image         string
	Replicas      int32
	ContainerPort int32
}

func (c *Client) ApplyDeployment(ctx context.Context,cfg DeploymentConfig,
) (*appsv1.Deployment, error) {
	deployments := c.Kube.AppsV1().
		Deployments(cfg.Namespace)

	deployment := buildDeployment(cfg)

	current, err := deployments.Get(
		ctx,
		cfg.Name,
		metav1.GetOptions{},
	)

	if err != nil {
		if errors.IsNotFound(err) {
			created, err := deployments.Create(
				ctx,
				deployment,
				metav1.CreateOptions{},
			)
			if err != nil {
				return nil, fmt.Errorf(
					"create deployment %s/%s: %w",
					cfg.Namespace,
					cfg.Name,
					err,
				)
			}

			return created, nil
		}

		return nil, fmt.Errorf(
			"get deployment %s/%s: %w",
			cfg.Namespace,
			cfg.Name,
			err,
		)
	}

	deployment.ResourceVersion = current.ResourceVersion

	updated, err := deployments.Update(
		ctx,
		deployment,
		metav1.UpdateOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"update deployment %s/%s: %w",
			cfg.Namespace,
			cfg.Name,
			err,
		)
	}

	return updated, nil
}

func buildDeployment(cfg DeploymentConfig) *appsv1.Deployment {
	labels := map[string]string{
		"app": cfg.Name,
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfg.Name,
			Namespace: cfg.Namespace,
		},

		Spec: appsv1.DeploymentSpec{
			Replicas: &cfg.Replicas,

			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},

			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},

				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  cfg.Name,
							Image: cfg.Image,

							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: cfg.ContainerPort,
									Protocol:      corev1.ProtocolTCP,
								},
							},

							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt32(cfg.ContainerPort),
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