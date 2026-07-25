package controller

import (
	dto "backend/DTO"
	service "backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
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
