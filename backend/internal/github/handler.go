package github

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/sowncns/k3s-deploy-platform/internal/auth"
	"github.com/sowncns/k3s-deploy-platform/internal/auth/dto"
)

type Handler struct {
	authService auth.AuthService
	authRepo    auth.AuthRepository
	oauthConfig *oauth2.Config
}

func NewHandler(authService auth.AuthService, authRepo auth.AuthRepository) *Handler {
	return &Handler{
		authService: authService,
		authRepo:    authRepo,
		oauthConfig: NewOAuthConfig(),
	}
}

// RegisterPublicRoutes đăng ký route không cần token (dùng để lấy token).
func (h *Handler) RegisterPublicRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("/auth/github")
	{
		authGroup.GET("", h.Login)
		authGroup.GET("/callback", h.Callback)
	}
}

// RegisterRoutes đăng ký route cần token (phải đi qua AuthMiddleware).
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/github/repos", h.ListRepositories)
}

func (h *Handler) Login(c *gin.Context) {
	url := h.oauthConfig.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) Callback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing github code",
		})
		return
	}

	ctx := c.Request.Context()

	token, err := ExchangeCode(ctx, h.oauthConfig, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client := NewClient(token.AccessToken)

	ghUser, err := client.GetUser(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.CreateOrUpdateUser(ctx, dto.CreateAuthUser{
		ID:          ghUser.GetID(),
		Login:       ghUser.GetLogin(),
		Name:        ghUser.GetName(),
		Email:       ghUser.GetEmail(),
		AvatarURL:   ghUser.GetAvatarURL(),
		AccessToken: token.AccessToken,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: tắt trả access_token trực tiếp trong response khi xong giai đoạn test,
	// client lúc đó nên tự lưu token từ nơi khác (session/cookie) thay vì đọc JSON.
	c.JSON(http.StatusOK, gin.H{
		"user":         user,
		"access_token": token.AccessToken,
	})
}

func (h *Handler) ListRepositories(c *gin.Context) {
	userID := c.GetInt64("user_id")

	user, err := h.authRepo.FindByID(c.Request.Context(), userID)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	if user.AccessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "github access token not found"})
		return
	}

	client := NewClient(user.AccessToken)

	repos, err := client.ListRepositories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, repos)
}
