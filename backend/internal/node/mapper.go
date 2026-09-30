package node

import (
    "time"

    corev1 "k8s.io/api/core/v1"
)

func mapK8sNode(clusterID string, k8sNode *corev1.Node) *Node {
    return &Node{
        ID:              string(k8sNode.UID),
        ClusterID:       clusterID,
        Name:            k8sNode.Name,
        Role:            getNodeRole(k8sNode),
        Status:           getNodeStatus(k8sNode),
        CPUCapacity:     k8sNode.Status.Capacity.Cpu().MilliValue(),
        MemoryCapacity:  k8sNode.Status.Capacity.Memory().Value(),
        PodCapacity:     int(k8sNode.Status.Capacity.Pods().Value()),
        CreatedAt:       k8sNode.CreationTimestamp.Time,
        UpdatedAt:       time.Now(),
    }
}

func getNodeRole(node *corev1.Node) NodeRole {
    if _, ok := node.Labels["node-role.kubernetes.io/control-plane"]; ok {
        return NodeRoleControlPlane
    }

    if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {
        return NodeRoleControlPlane
    }

    return NodeRoleWorker
}



func getNodeStatus(node *corev1.Node) NodeStatus {
    for _, condition := range node.Status.Conditions {
        if condition.Type == corev1.NodeReady {
            if condition.Status == corev1.ConditionTrue {
                return NodeStatusReady
            }

            return NodeStatusNotReady
        }
    }

    return NodeStatusUnknown
}