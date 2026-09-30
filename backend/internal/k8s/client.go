package k8s

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)
// tao interface cua k8s
type Client struct {
	Kube kubernetes.Interface
}
// maau trar ve con tro client
func NewClient() (*Client, error) {
	config, err := buildConfig()
	if err != nil {
		return nil, fmt.Errorf("build kubernetes config: %w", err)
	}

	return newClientFromConfig(config)
}


func NewClientFromKubeconfigPath(path string) (*Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		return nil, fmt.Errorf("build kubernetes config from %s: %w", path, err)
	}

	return newClientFromConfig(config)
}

func newClientFromConfig(config *rest.Config) (*Client, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("create kubernetes client: %w", err)
	}

	return &Client{
		Kube: clientset,
	}, nil
}

func buildConfig() (*rest.Config, error) {
	//lay cfg cua kube /var/lib/rancher.....
	kubeconfig := os.Getenv("KUBECONFIG")

	if kubeconfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		kubeconfig = filepath.Join(
			home,
			".kube",
			"config",
		)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}