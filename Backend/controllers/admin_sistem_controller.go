package controller

import (
	dto "backend/DTO"
	service "backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AdminSistemController struct{ service *service.AdminSistemService }

func NewAdminSistemController() *AdminSistemController {
	return &AdminSistemController{service: service.NewAdminSistemService()}
}
func (c *AdminSistemController) Dashboard(ctx *gin.Context) {
	data, err := c.service.Dashboard()
	c.respond(ctx, data, err, http.StatusOK)
}
func (c *AdminSistemController) Users(ctx *gin.Context) {
	data, err := c.service.Users()
	c.respond(ctx, data, err, http.StatusOK)
}
func (c *AdminSistemController) User(ctx *gin.Context) {
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}
	data, err := c.service.User(id)
	c.respond(ctx, data, err, http.StatusOK)
}
func (c *AdminSistemController) CreateUser(ctx *gin.Context) {
	var req dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := c.service.CreateUser(req)
	c.respond(ctx, data, err, http.StatusCreated)
}
func (c *AdminSistemController) UpdateUser(ctx *gin.Context) {
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}
	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := c.service.UpdateUser(id, req)
	c.respond(ctx, data, err, http.StatusOK)
}
func (c *AdminSistemController) SetUserStatus(ctx *gin.Context) {
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}
	var req dto.UserStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.respond(ctx, nil, c.service.SetUserStatus(id, req.IsActive), http.StatusOK)
}
func (c *AdminSistemController) ResetPassword(ctx *gin.Context) {
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}
	var req dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.respond(ctx, nil, c.service.ResetPassword(id, req.Password), http.StatusOK)
}
func (c *AdminSistemController) Units(ctx *gin.Context) {
	data, err := c.service.Units()
	c.respond(ctx, data, err, 200)
}
func (c *AdminSistemController) CreateUnit(ctx *gin.Context) {
	var req dto.UnitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := c.service.CreateUnit(req)
	c.respond(ctx, data, err, 201)
}
func (c *AdminSistemController) UpdateUnit(ctx *gin.Context) {
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}
	var req dto.UnitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := c.service.UpdateUnit(id, req)
	c.respond(ctx, data, err, 200)
}
func (c *AdminSistemController) SetUnitStatus(ctx *gin.Context) {
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}
	var req dto.ActiveStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.respond(ctx, nil, c.service.SetUnitStatus(id, req.IsActive), http.StatusOK)
}
func (c *AdminSistemController) Categories(ctx *gin.Context) {
	data, err := c.service.Categories()
	c.respond(ctx, data, err, 200)
}
func (c *AdminSistemController) CreateCategory(ctx *gin.Context) {
	var req dto.CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := c.service.CreateCategory(req)
	c.respond(ctx, data, err, 201)
}
func (c *AdminSistemController) UpdateCategory(ctx *gin.Context) {
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}
	var req dto.CategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := c.service.UpdateCategory(id, req)
	c.respond(ctx, data, err, 200)
}
func (c *AdminSistemController) SetCategoryStatus(ctx *gin.Context) {
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}
	var req dto.ActiveStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.respond(ctx, nil, c.service.SetCategoryStatus(id, req.IsActive), http.StatusOK)
}
func (c *AdminSistemController) ImportTemplate(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/csv; charset=utf-8")
	ctx.Header("Content-Disposition", `attachment; filename="template-mahasiswa.csv"`)
	ctx.String(http.StatusOK, "nama_lengkap,nim,email,program_studi,angkatan\nNama Mahasiswa,20260001,mahasiswa@example.com,Teknik Informatika,2026\n")
}
func (c *AdminSistemController) ImportMahasiswa(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}
	header, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "file CSV atau XLSX wajib diunggah"})
		return
	}
	if header.Size <= 0 || header.Size > 10<<20 {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ukuran file maksimal 10 MB"})
		return
	}
	ext := strings.ToLower(header.Filename)
	if !strings.HasSuffix(ext, ".csv") && !strings.HasSuffix(ext, ".xlsx") {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "format file harus CSV atau XLSX"})
		return
	}
	file, err := header.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "file import tidak dapat dibaca"})
		return
	}
	defer file.Close()
	result, err := c.service.ImportMahasiswa(userID, header.Filename, ext, file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Import mahasiswa selesai", "data": result})
}
func (c *AdminSistemController) ImportHistory(ctx *gin.Context) {
	data, err := c.service.ImportHistory()
	c.respond(ctx, data, err, http.StatusOK)
}
func (c *AdminSistemController) ImportHistoryDetail(ctx *gin.Context) {
	id, ok := getIDParam(ctx, "id")
	if !ok {
		return
	}
	data, err := c.service.ImportHistoryDetail(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}
func (c *AdminSistemController) respond(ctx *gin.Context, data any, err error, successStatus int) {
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	payload := gin.H{"success": true}
	if data != nil {
		payload["data"] = data
	}
	ctx.JSON(successStatus, payload)
}
