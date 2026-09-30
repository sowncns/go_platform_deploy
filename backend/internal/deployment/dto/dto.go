package dto

type DeploymentRequest struct {
	Name          string
	Namespace     string
	Image         string
	Replicas      int32
	ContainerPort int32
	Host          string
	ServiceName   string
	ServicePort   int32
	Port          int32
	TargetPort    int32
	ProjectID     uint
	CommitSHA     string
	ImageTag      string
}
