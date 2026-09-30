package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service AuthService
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GitHubLogin(c *gin.Context) {
	url:= h.service.GetGitHubLoginURL()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) GitHubCallback(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing github code",
		})
		return
	}

	user, err := h.service.LoginWithGitHub(
		c.Request.Context(),
		code,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}