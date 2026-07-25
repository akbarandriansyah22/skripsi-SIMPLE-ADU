package controller

import (
	"net/http"

	dto "backend/DTO"
	service "backend/services"

	"github.com/gin-gonic/gin"
)

type PimpinanController struct {
	service *service.PimpinanService
}

func NewPimpinanController() *PimpinanController {
	return &PimpinanController{
		service: service.NewPimpinanService(),
	}
}

func (c *PimpinanController) Dashboard(ctx *gin.Context) {
	result, err := c.service.Dashboard()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Dashboard pimpinan berhasil diambil", "data": result})
}

func (c *PimpinanController) GetUrgensiTinggi(ctx *gin.Context) {
	result, err := c.service.GetUrgensiTinggi()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Pengaduan urgensi tinggi berhasil diambil", "data": result})
}

func (c *PimpinanController) GetPengaduan(ctx *gin.Context) {
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}
	result, err := c.service.GetPengaduan(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

func (c *PimpinanController) Monitoring(ctx *gin.Context) {
	result, err := c.service.Monitoring()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

func (c *PimpinanController) CreateDisposisi(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	pengaduanID, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.DisposisiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := c.service.CreateDisposisi(userID, pengaduanID, req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": true, "message": "Disposisi berhasil dibuat"})
}

func (c *PimpinanController) MyDisposisi(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	result, err := c.service.GetDisposisiByPimpinan(uint64(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Disposisi berhasil diambil", "data": result})
}
