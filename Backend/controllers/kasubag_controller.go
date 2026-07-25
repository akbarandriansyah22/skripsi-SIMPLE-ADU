package controller

import (
	"errors"
	"net/http"
	"strings"

	service "backend/services"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type KasubagController struct {
	service *service.KasubagService
}

func NewKasubagController() *KasubagController {
	return &KasubagController{service: service.NewKasubagService()}
}

func (c *KasubagController) Dashboard(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}
	data, err := c.service.Dashboard(uint64(userID))
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (c *KasubagController) List(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}
	data, err := c.service.GetComplaints(uint64(userID))
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (c *KasubagController) Detail(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}
	data, err := c.service.GetComplaint(uint64(userID), id)
	if err != nil {
		c.writeServiceError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (c *KasubagController) StartProcess(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.service.StartProcess(uint64(userID), id); err != nil {
		c.writeServiceError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Aduan mulai diproses"})
}

func (c *KasubagController) AddResponse(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}
	message := ""
	var upload utils.UploadedFile
	if strings.HasPrefix(strings.ToLower(ctx.ContentType()), "multipart/form-data") {
		message = ctx.PostForm("pesan")
		if _, err := ctx.FormFile("lampiran"); err == nil {
			upload, err = utils.UploadFileMetadata(ctx, "lampiran")
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
				return
			}
		}
	} else {
		var request struct {
			Pesan string `json:"pesan" binding:"required"`
		}
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}
		message = request.Pesan
	}
	if strings.TrimSpace(message) == "" {
		if upload.Path != "" {
			_ = utils.DeleteUploadedFile(upload.Path)
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "pesan wajib diisi"})
		return
	}
	if err := c.service.AddResponse(uint64(userID), id, message, upload.Path, upload.Original, upload.MIMEType, upload.Size); err != nil {
		if upload.Path != "" {
			_ = utils.DeleteUploadedFile(upload.Path)
		}
		c.writeServiceError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"success": true, "message": "Tindak lanjut disimpan"})
}

func (c *KasubagController) Finish(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.service.Finish(uint64(userID), id); err != nil {
		c.writeServiceError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Aduan diselesaikan"})
}

func (c *KasubagController) ReturnToAdmin(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}
	var request struct {
		Alasan string `json:"alasan" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := c.service.ReturnToAdmin(uint64(userID), id, request.Alasan); err != nil {
		c.writeServiceError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Aduan dikembalikan ke Admin Fakultas"})
}

func (c *KasubagController) writeServiceError(ctx *gin.Context, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, service.ErrForbidden) {
		status = http.StatusForbidden
	}
	if strings.Contains(strings.ToLower(err.Error()), "record not found") {
		status = http.StatusNotFound
	}
	ctx.JSON(status, gin.H{"success": false, "message": err.Error()})
}
