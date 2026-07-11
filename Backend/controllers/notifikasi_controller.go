package controller

import (
	"net/http"

	service "backend/services"

	"github.com/gin-gonic/gin"
)

type NotifikasiController struct {
	service *service.NotifikasiService
}

func NewNotifikasiController() *NotifikasiController {
	return &NotifikasiController{
		service: service.NewNotifikasiService(),
	}
}

func (c *NotifikasiController) Mine(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	result, err := c.service.GetByUserID(uint64(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Notifikasi berhasil diambil", "data": result})
}

func (c *NotifikasiController) Unread(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	result, err := c.service.GetUnreadByUserID(uint64(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Notifikasi belum dibaca berhasil diambil", "data": result})
}

func (c *NotifikasiController) MarkAsRead(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.service.MarkAsRead(uint64(userID), id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Notifikasi berhasil ditandai dibaca"})
}
