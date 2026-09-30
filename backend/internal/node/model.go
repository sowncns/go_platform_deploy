package node

import "time"

type NodeRole string

const (
	NodeRoleControlPlane NodeRole = "control-plane"
	NodeRoleWorker       NodeRole = "worker"
)

type NodeStatus string

const (
	NodeStatusReady    NodeStatus = "ready"
	NodeStatusNotReady NodeStatus = "not-ready"
	NodeStatusUnknown  NodeStatus = "unknown"
)

type Node struct {
	ID        string
	ClusterID string

	Name   string
	Role   NodeRole
	Status NodeStatus

	
	CPUCapacity    int64
	MemoryCapacity int64
	PodCapacity    int

	CPUUsage    int64
	MemoryUsage int64
	PodUsage    int

	CreatedAt time.Time
	UpdatedAt time.Time
}