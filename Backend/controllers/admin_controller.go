package controller

import (
	"net/http"

	dto "backend/DTO"
	service "backend/services"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	service *service.AdminService
}

func NewAdminController() *AdminController {
	return &AdminController{
		service: service.NewAdminService(),
	}
}

func (c *AdminController) Dashboard(ctx *gin.Context) {
	result, err := c.service.Dashboard()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Dashboard admin berhasil diambil", "data": result})
}

func (c *AdminController) GetAllPengaduan(ctx *gin.Context) {
	result, err := c.service.GetAllPengaduan()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Pengaduan berhasil diambil", "data": result})
}

func (c *AdminController) GetPengaduanByID(ctx *gin.Context) {
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}

	result, err := c.service.GetPengaduanByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Pengaduan berhasil diambil", "data": result})
}

func (c *AdminController) UpdateStatus(ctx *gin.Context) {
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.UpdateStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := c.service.UpdateStatus(id, req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Status pengaduan berhasil diubah"})
}

func (c *AdminController) AssignUnit(ctx *gin.Context) {
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.AssignUnitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := c.service.AssignUnit(id, req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Unit berhasil ditetapkan"})
}

func (c *AdminController) ForwardToPimpinan(ctx *gin.Context) {
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.service.ForwardToPimpinan(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Pengaduan berhasil diteruskan ke pimpinan"})
}

func (c *AdminController) ReanalyzeAI(ctx *gin.Context) {
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}

	result, err := c.service.ReanalyzeAI(id)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"success": false, "message": "Analisis AI gagal diperbarui"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Analisis AI berhasil diperbarui", "data": result})
}

func (c *AdminController) GetUnits(ctx *gin.Context) {
	result, err := c.service.GetUnits()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Unit berhasil diambil", "data": result})
}
