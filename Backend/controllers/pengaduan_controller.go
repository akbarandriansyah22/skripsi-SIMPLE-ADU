package controller

import (
	"errors"
	"net/http"
	"strconv"

	dto "backend/DTO"
	service "backend/services"

	"github.com/gin-gonic/gin"
)

type PengaduanController struct {
	service *service.PengaduanService
}

func NewPengaduanController() *PengaduanController {
	return &PengaduanController{
		service: service.NewPengaduanService(),
	}
}

func (c *PengaduanController) Create(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	var req dto.CreatePengaduanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	result, err := c.service.Create(userID, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": true, "message": "Pengaduan berhasil dibuat", "data": result})
}

func (c *PengaduanController) MyPengaduan(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	result, err := c.service.GetByUserID(uint64(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Pengaduan berhasil diambil", "data": result})
}

func (c *PengaduanController) Detail(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	role, ok := getRole(ctx)
	if !ok {
		return
	}

	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}

	result, err := c.service.GetByIDForRole(id, uint64(userID), role)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Pengaduan berhasil diambil", "data": result})
}

func (c *PengaduanController) DetailByTicket(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	role, ok := getRole(ctx)
	if !ok {
		return
	}

	result, err := c.service.GetByKodeTiketForRole(ctx.Param("kode"), uint64(userID), role)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Pengaduan berhasil diambil", "data": result})
}

func (c *PengaduanController) Update(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.UpdatePengaduanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	result, err := c.service.Update(uint64(userID), id, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Pengaduan berhasil diubah", "data": result})
}

func (c *PengaduanController) AddRespon(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	role, ok := getRole(ctx)
	if !ok {
		return
	}

	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}

	var req struct {
		Pesan string `json:"pesan" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := c.service.AddRespon(userID, role, id, req.Pesan); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": true, "message": "Respon berhasil ditambahkan"})
}

func (c *PengaduanController) Finish(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.service.Finish(uint64(userID), id); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Pengaduan selesai"})
}

func getUserID(ctx *gin.Context) (uint, bool) {
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User tidak ditemukan"})
		return 0, false
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User tidak valid"})
		return 0, false
	}

	return userID, true
}

func getRole(ctx *gin.Context) (string, bool) {
	roleValue, exists := ctx.Get("role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Role tidak ditemukan"})
		return "", false
	}

	role, ok := roleValue.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Role tidak valid"})
		return "", false
	}

	return role, true
}

func getIDParam(ctx *gin.Context, name string) (uint64, bool) {
	id, err := strconv.ParseUint(ctx.Param(name), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return 0, false
	}

	return id, true
}
