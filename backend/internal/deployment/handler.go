package deployment


type Handler struct {
	service DeploymentService
}

func NewHandler(service DeploymentService) *Handler {
	return &Handler{service: service}
}