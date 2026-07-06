package controller

import (
	"net/http"

	dto "backend/DTO"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *service.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		service: service.NewAuthService(),
	}
}

// ==========================
// REGISTER
// ==========================

func (c *AuthController) Register(ctx *gin.Context) {

	var req dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})

		return
	}

	result, err := c.service.Register(req)

	if err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registrasi berhasil",
		"data":    result,
	})

}

// ==========================
// LOGIN
// ==========================

func (c *AuthController) Login(ctx *gin.Context) {

	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})

		return
	}

	result, err := c.service.Login(req)

	if err != nil {

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"data":    result,
	})

}