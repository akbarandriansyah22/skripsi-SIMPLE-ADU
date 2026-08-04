package controller

import (
	"net/http"

	dto "backend/DTO"
	service "backend/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *service.AuthService
}

// ==========================
// PROFILE
// ==========================

func (c *AuthController) Profile(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak ditemukan",
		})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User tidak valid",
		})
		return
	}

	result, err := c.service.Profile(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile berhasil diambil",
		"data":    result,
	})
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

func (c *AuthController) ChangePassword(ctx *gin.Context) {
	value, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User tidak ditemukan"})
		return
	}
	userID, ok := value.(uint)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User tidak valid"})
		return
	}
	var req dto.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := c.service.ChangePassword(userID, req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Password berhasil diubah"})
}
