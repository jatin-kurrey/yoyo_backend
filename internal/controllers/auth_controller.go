package controllers

import (
	"yoyo-server/internal/middleware"
	"yoyo-server/internal/services"
	"yoyo-server/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	auth  *services.AuthService
	audit *services.AuditService
}

func NewAuthController(auth *services.AuthService, audit *services.AuditService) *AuthController {
	return &AuthController{auth: auth, audit: audit}
}

func (ctl *AuthController) Login(c *gin.Context) {
	var input services.LoginInput
	if !bindAndValidate(c, &input) {
		return
	}

	result, err := ctl.auth.Login(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	ctl.audit.Log(c.Request.Context(), &result.User.ID, "login", "auth", nil, c.ClientIP())
	utils.OK(c, "Login successful.", result)
}

func (ctl *AuthController) Me(c *gin.Context) {
	user, ok := middleware.CurrentAdmin(c)
	if !ok {
		utils.Unauthorized(c, "Authentication is required.")
		return
	}
	utils.OK(c, "Authenticated admin loaded.", user)
}

func (ctl *AuthController) Logout(c *gin.Context) {
	ctl.audit.Log(c.Request.Context(), currentAdminID(c), "logout", "auth", nil, c.ClientIP())
	utils.OK(c, "Logout successful.", nil)
}
