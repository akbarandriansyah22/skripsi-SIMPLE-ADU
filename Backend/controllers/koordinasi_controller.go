package controller

import (
	"errors"
	"net/http"
	"strconv"

	"backend/config"
	"backend/models"
	service "backend/services"
	"backend/utils"
	"github.com/gin-gonic/gin"
)

type KoordinasiController struct{ service *service.KoordinasiService }

func NewKoordinasiController() *KoordinasiController {
	return &KoordinasiController{service: service.NewKoordinasiService()}
}

func (c *KoordinasiController) List(ctx *gin.Context) {
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
	result, err := c.service.List(userID, role, uint(id))
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrForbidden) {
			status = http.StatusForbidden
		}
		ctx.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}
	for index := range result {
		if result[index].Lampiran != "" {
			result[index].LampiranURL = "/api/pengaduan/" + strconv.FormatUint(id, 10) + "/koordinasi/" + strconv.FormatUint(uint64(result[index].ID), 10) + "/lampiran"
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}
func (c *KoordinasiController) Create(ctx *gin.Context) {
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
	message, parent := ctx.PostForm("pesan"), (*uint)(nil)
	if message == "" && ctx.ContentType() != "multipart/form-data" {
		var body struct {
			Pesan    string `json:"pesan"`
			ParentID *uint  `json:"parent_id"`
		}
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}
		message, parent = body.Pesan, body.ParentID
	}
	var attachment utils.UploadedFile
	if _, err := ctx.FormFile("lampiran"); err == nil {
		attachment, err = utils.UploadFileMetadata(ctx, "lampiran")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}
	}
	if value := ctx.PostForm("parent_id"); value != "" {
		parsed, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "parent_id tidak valid"})
			return
		}
		parentValue := uint(parsed)
		parent = &parentValue
	}
	if err := c.service.Create(userID, role, uint(id), parent, message, attachment.Path, attachment.Original, attachment.MIMEType, attachment.Size); err != nil {
		if attachment.Path != "" {
			_ = utils.DeleteUploadedFile(attachment.Path)
		}
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrForbidden) {
			status = http.StatusForbidden
		}
		ctx.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"success": true, "message": "Koordinasi berhasil dikirim"})
}
func (c *KoordinasiController) Attachment(ctx *gin.Context) {
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
	messageID, ok := getIDParam(ctx, "messageId")
	if !ok {
		return
	}
	var row models.KoordinasiInternal
	if err := config.DB.Preload("Pengaduan.HasilAI").First(&row, messageID).Error; err != nil || row.PengaduanID != uint(id) {
		ctx.Status(http.StatusNotFound)
		return
	}
	if _, err := c.service.List(userID, role, uint(id)); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			ctx.Status(http.StatusForbidden)
		} else {
			ctx.Status(http.StatusBadRequest)
		}
		return
	}
	attachment, err := utils.ResolveUploadedFile(row.Lampiran)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.File(attachment)
}
