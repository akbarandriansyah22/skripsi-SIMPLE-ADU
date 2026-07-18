package middleware

import (
	"net/http"
	"strings"

	"backend/config"
	"backend/models"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header wajib diisi",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Format Authorization harus Bearer token",
			})
			return
		}

		userID, role, err := utils.ParseJWT(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		canonicalRole := utils.CanonicalRole(role)
		if canonicalRole == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Role token tidak valid",
			})
			return
		}
		if config.DB == nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": "database belum siap"})
			return
		}
		var user models.User
		if err := config.DB.Select("is_active").First(&user, userID).Error; err != nil || !user.IsActive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "akun tidak aktif atau tidak ditemukan"})
			return
		}

		c.Set("user_id", userID)
		c.Set("role", canonicalRole)
		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(allowedRoles))
	for _, role := range allowedRoles {
		canonical := utils.CanonicalRole(role)
		if canonical != "" {
			allowed[canonical] = true
		}
	}

	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Role tidak ditemukan",
			})
			return
		}

		role, ok := roleValue.(string)
		if !ok || !allowed[utils.CanonicalRole(role)] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak",
			})
			return
		}

		c.Next()
	}
}
