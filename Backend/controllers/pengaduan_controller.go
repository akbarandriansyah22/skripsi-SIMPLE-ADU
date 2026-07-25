package controller

import (
	"errors"
	"net/http"
	"strconv"

	dto "backend/DTO"
	service "backend/services"
	"backend/utils"

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

	req, ok := bindCreatePengaduanRequest(ctx)
	if !ok {
		return
	}

	result, err := c.service.Create(userID, req)
	if err != nil {
		_ = utils.DeleteUploadedFile(req.Lampiran)
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	message := "Pengaduan berhasil dibuat"
	if result.AIStatus == "pending" {
		message = "Pengaduan berhasil dibuat, analisis AI menunggu pemrosesan"
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": true, "message": message, "data": result})
}

func bindCreatePengaduanRequest(ctx *gin.Context) (dto.CreatePengaduanRequest, bool) {
	if ctx.ContentType() != "multipart/form-data" && ctx.ContentType() != "application/x-www-form-urlencoded" {
		var req dto.CreatePengaduanRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return req, false
		}
		return req, true
	}

	kategoriID, err := strconv.ParseUint(ctx.PostForm("kategori_id"), 10, 64)
	if err != nil || kategoriID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "kategori_id wajib diisi"})
		return dto.CreatePengaduanRequest{}, false
	}

	lampiran := ctx.PostForm("lampiran")
	var upload utils.UploadedFile
	if _, err := ctx.FormFile("lampiran"); err == nil {
		uploaded, err := utils.UploadFileMetadata(ctx, "lampiran")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "gagal upload lampiran"})
			return dto.CreatePengaduanRequest{}, false
		}
		upload = uploaded
		lampiran = uploaded.Path
	}

	req := dto.CreatePengaduanRequest{
		KategoriID:       uint(kategoriID),
		Judul:            ctx.PostForm("judul"),
		Deskripsi:        ctx.PostForm("deskripsi"),
		Lampiran:         lampiran,
		LampiranNamaAsli: upload.Original,
		LampiranMimeType: upload.MIMEType,
		LampiranUkuran:   upload.Size,
	}
	if req.Judul == "" || req.Deskripsi == "" {
		if upload.Path != "" {
			_ = utils.DeleteUploadedFile(upload.Path)
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "judul dan deskripsi wajib diisi"})
		return req, false
	}

	return req, true
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
