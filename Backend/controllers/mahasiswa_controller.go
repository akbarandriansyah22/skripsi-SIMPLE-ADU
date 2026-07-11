package controller

import (
	"net/http"

	dto "backend/DTO"
	service "backend/services"

	"github.com/gin-gonic/gin"
)

type MahasiswaController struct {
	service *service.MahasiswaService
}

func NewMahasiswaController() *MahasiswaController {
	return &MahasiswaController{
		service: service.NewMahasiswaService(),
	}
}

func (c *MahasiswaController) Profile(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	result, err := c.service.GetProfile(uint64(userID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile mahasiswa berhasil diambil", "data": result})
}

func (c *MahasiswaController) UpdateProfile(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	var req dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := c.service.UpdateProfile(uint64(userID), req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile mahasiswa berhasil diubah"})
}
