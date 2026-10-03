package project

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sowncns/k3s-deploy-platform/internal/project/dto"
)


type Handler struct {
	service ProjectService
}


func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	projectGroup := r.Group("/projects")
	
	{
		projectGroup.POST("", h.CreateProject)
		projectGroup.GET("/", h.ListProjectsByUserID)
		projectGroup.GET("/:id", h.FindProjectByID)
		projectGroup.PATCH("/:id", h.UpdateProject)
	}
}

func NewHandler(service ProjectService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) UpdateProject(c *gin.Context) {
    projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid project id",
        })
        return
    }

    userID := c.GetUint("user_id")

    var req dto.UpdateProjectRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    err = h.service.UpdateProject(
        c.Request.Context(),
        uint(projectID),
        userID,
        &req,
    )

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "project updated successfully",
    })
}

func (h *Handler) CreateProject(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := h.service.CreateProject(c.Request.Context(), uint(userID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (h *Handler) FindProjectByID(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	userID := c.GetUint("user_id")

	project, err := h.service.GetProject(c.Request.Context(), uint(projectID), userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *Handler) ListProjectsByUserID(c *gin.Context) {
	userID  := c.GetUint("user_id")
	log.Printf(" user id %d", userID)
	projects, err := h.service.ListProjects(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (h *Handler) DeleteProject(c *gin.Context){

}