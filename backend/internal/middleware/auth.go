package middleware

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/sowncns/k3s-deploy-platform/internal/auth"

)

// AuthMiddleware kiểm tra token gửi lên từ Postman
func AuthMiddleware(authRepo auth.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy Header Authorization từ request (Postman gửi lên)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "thiếu Authorization header"})
			c.Abort()
			return
		}

		// 2. Tách lấy chuỗi token (bỏ chữ "Bearer ")
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "format token không hợp lệ (Bearer <token>)"})
			c.Abort()
			return
		}
		accessToken := tokenParts[1]

		// 3. Tìm user trong database dựa vào AccessToken đã lưu lúc Callback
		user, err := authRepo.FindByAccessToken(c.Request.Context(), accessToken)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token không hợp lệ hoặc hết hạn"})
			c.Abort()
			return
		}

		// 4. SET "user_id" VÀO CONTEXT CHO HÀM LIST REPOSITORIES DÙNG
		c.Set("user_id", uint(user.ID)) // uint để khớp với c.GetUint trong các handler project

		// Cho phép request đi tiếp vào Handler chính
		c.Next()
	}
}