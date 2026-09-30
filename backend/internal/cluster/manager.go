package cluster

import (
	"context"
	"fmt"
	"sync"

	"github.com/sowncns/k3s-deploy-platform/internal/k8s"
)

// Manager quan ly ket noi (k8s.Client) toi nhieu cluster, dong bo voi
// metadata luu trong DB thong qua Repository.
type Manager struct {
	repo ClusterRepository
	mu      sync.RWMutex
	clients map[string]*k8s.Client
}

func NewManager(repo ClusterRepository) *Manager {
	return &Manager{
		repo:    repo,
		clients: make(map[string]*k8s.Client),
	}
}

// Register luu metadata cluster vao DB, dung kubeconfig tai CredentialRef
// de tao client, roi cache client trong bo nho.
func (m *Manager) Register(ctx context.Context, c *Cluster) error {
	client, err := k8s.NewClientFromKubeconfigPath(c.CredentialRef)
	if err != nil {
		return fmt.Errorf("create client for cluster %s: %w", c.ID, err)
	}

	if err := m.repo.Create(ctx, c); err != nil {
		return fmt.Errorf("save cluster %s: %w", c.ID, err)
	}

	m.mu.Lock()
	m.clients[c.ID] = client
	m.mu.Unlock()

	return nil
}

// Get tra ve k8s client dang cache cho cluster. Neu chua co trong bo nho
// (vd sau restart), load lai tu DB + kubeconfig.
func (m *Manager) Get(ctx context.Context, clusterID string) (*k8s.Client, error) {
	m.mu.RLock()
	client, ok := m.clients[clusterID]
	m.mu.RUnlock()
	if ok {
		return client, nil
	}

	c, err := m.repo.FindByID(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	client, err = k8s.NewClientFromKubeconfigPath(c.CredentialRef)
	if err != nil {
		return nil, fmt.Errorf("create client for cluster %s: %w", c.ID, err)
	}

	m.mu.Lock()
	m.clients[c.ID] = client
	m.mu.Unlock()

	return client, nil
}

// List tra ve metadata cua tat ca cluster (khong bao gom k8s client).
func (m *Manager) List(ctx context.Context) ([]*Cluster, error) {
	return m.repo.List(ctx)
}

// UpdateStatus cap nhat trang thai cluster trong DB.
func (m *Manager) UpdateStatus(ctx context.Context, clusterID string, status ClusterStatus) error {
	return m.repo.UpdateStatus(ctx, clusterID, status)
}

// Remove xoa cluster khoi DB va bo cache client.
func (m *Manager) Remove(ctx context.Context, clusterID string) error {
	if err := m.repo.Delete(ctx, clusterID); err != nil {
		return err
	}

	m.mu.Lock()
	delete(m.clients, clusterID)
	m.mu.Unlock()

	return nil
}
