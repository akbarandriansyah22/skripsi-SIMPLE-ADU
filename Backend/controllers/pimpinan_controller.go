package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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

func (c *PimpinanController) GetRiwayatUrgensiTinggi(ctx *gin.Context) {
	filter, err := pimpinanHistoryFilter(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	result, err := c.service.GetRiwayatUrgensiTinggi(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Riwayat pengaduan urgensi tinggi berhasil diambil", "data": result})
}

func (c *PimpinanController) GetRiwayatUrgensiTinggiDetail(ctx *gin.Context) {
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}
	result, err := c.service.GetRiwayatUrgensiTinggiDetail(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

func pimpinanHistoryFilter(ctx *gin.Context) (service.PimpinanHistoryFilter, error) {
	filter := service.PimpinanHistoryFilter{Search: ctx.Query("q"), Status: ctx.Query("status")}
	for key, target := range map[string]*uint64{"kategori_id": &filter.KategoriID, "unit_id": &filter.UnitID} {
		value := ctx.Query(key)
		if value == "" {
			continue
		}
		parsed, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return filter, fmt.Errorf("%s tidak valid", key)
		}
		*target = parsed
	}
	if value := ctx.Query("from"); value != "" {
		parsed, err := time.Parse("2006-01-02", value)
		if err != nil {
			return filter, fmt.Errorf("format from harus YYYY-MM-DD")
		}
		filter.From = &parsed
	}
	if value := ctx.Query("to"); value != "" {
		parsed, err := time.Parse("2006-01-02", value)
		if err != nil {
			return filter, fmt.Errorf("format to harus YYYY-MM-DD")
		}
		end := parsed.AddDate(0, 0, 1)
		filter.To = &end
	}
	return filter, nil
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

func (c *PimpinanController) GetAssignmentUnits(ctx *gin.Context) {
	result, err := c.service.GetAssignmentUnits()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Unit penanganan berhasil diambil", "data": result})
}
