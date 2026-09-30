package cluster

import "time"

type ClusterStatus string

const (
	ClusterStatusProvisioning ClusterStatus = "provisioning"
	ClusterStatusActive       ClusterStatus = "active"
	ClusterStatusInactive     ClusterStatus = "inactive"
	ClusterStatusDraining     ClusterStatus = "draining"
	ClusterStatusError        ClusterStatus = "error"
)

type ClusterProvider string

const (
	ClusterProviderK3s        ClusterProvider = "k3s"
	ClusterProviderKubernetes  ClusterProvider = "kubernetes"
	ClusterProviderEKS        ClusterProvider = "eks"
	ClusterProviderGKE        ClusterProvider = "gke"
	ClusterProviderAKS        ClusterProvider = "aks"
)

type Cluster struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Provider ClusterProvider `json:"provider"`
	Status   ClusterStatus   `json:"status"`

	// Kubernetes API endpoint.
	Endpoint string `json:"endpoint"`

	// Geographic / logical location.
	Region string `json:"region"`
	Zone   string `json:"zone,omitempty"`

	// Reference to secret storage.
	// Never store kubeconfig/token directly here.
	CredentialRef string `json:"-"`

	// Cluster capacity.
	Capacity ClusterCapacity `json:"capacity"`

	// Current resource usage.
	Usage ClusterUsage `json:"usage"`

	// Number of Kubernetes nodes.
	NodeCount int `json:"node_count"`

	// Health information.
	Health ClusterHealth `json:"health"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ClusterCapacity struct {
	CPU    int64 `json:"cpu_millicores"`
	Memory int64 `json:"memory_bytes"`

	Pods int `json:"pods"`
}

type ClusterUsage struct {
	CPU    int64 `json:"cpu_millicores"`
	Memory int64 `json:"memory_bytes"`

	Pods int `json:"pods"`
}

type ClusterHealth struct {
	Healthy       bool       `json:"healthy"`
	LastCheckedAt *time.Time `json:"last_checked_at,omitempty"`
	Message       string     `json:"message,omitempty"`
}