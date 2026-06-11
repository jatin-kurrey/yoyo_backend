package middleware

import (
	"strings"

	"yoyo-server/internal/config"
	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
)

func AdminAuth(cfg *config.Config, users *repositories.AdminUserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			utils.Unauthorized(c, "Authentication token is required.")
			c.Abort()
			return
		}

		claims, err := utils.ParseAdminToken(strings.TrimPrefix(header, "Bearer "), cfg.JWTSecret)
		if err != nil {
			utils.Unauthorized(c, "Invalid or expired token.")
			c.Abort()
			return
		}

		user, err := users.FindByID(c.Request.Context(), claims.UserID)
		if err != nil || !user.IsActive {
			utils.Unauthorized(c, "Admin account is not active.")
			c.Abort()
			return
		}

		c.Set("adminUser", user)
		c.Set("adminUserID", user.ID)
		c.Set("adminRole", string(user.Role))
		c.Next()
	}
}

func RequireRoles(roles ...models.AdminRole) gin.HandlerFunc {
	allowed := map[string]bool{}
	for _, role := range roles {
		allowed[string(role)] = true
	}

	return func(c *gin.Context) {
		role, _ := c.Get("adminRole")
		if !allowed[role.(string)] {
			utils.Forbidden(c, "You do not have permission to perform this action.")
			c.Abort()
			return
		}
		c.Next()
	}
}

func CurrentAdmin(c *gin.Context) (*models.AdminUser, bool) {
	value, ok := c.Get("adminUser")
	if !ok {
		return nil, false
	}
	user, ok := value.(*models.AdminUser)
	return user, ok
}
