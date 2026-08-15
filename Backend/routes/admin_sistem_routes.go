package routes

import (
	controller "backend/controllers"
	"backend/middleware"
	"github.com/gin-gonic/gin"
)

func AdminSistemRoutes(r *gin.Engine) {
	c := controller.NewAdminSistemController()
	api := r.Group("/api/admin-sistem")
	api.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin_sistem"))
	{
		api.GET("/dashboard", c.Dashboard)
		api.GET("/users", c.Users)
		api.GET("/users/:id", c.User)
		api.POST("/users", c.CreateUser)
		api.PATCH("/users/:id", c.UpdateUser)
		api.PATCH("/users/:id/status", c.SetUserStatus)
		api.PATCH("/users/:id/reset-password", c.ResetPassword)
		api.GET("/units", c.Units)
		api.POST("/units", c.CreateUnit)
		api.PATCH("/units/:id", c.UpdateUnit)
		api.PATCH("/units/:id/status", c.SetUnitStatus)
		api.GET("/categories", c.Categories)
		api.POST("/categories", c.CreateCategory)
		api.PATCH("/categories/:id", c.UpdateCategory)
		api.PATCH("/categories/:id/status", c.SetCategoryStatus)
		api.GET("/mahasiswa/template", c.ImportTemplate)
		api.POST("/mahasiswa/import", c.ImportMahasiswa)
		api.GET("/mahasiswa/import-history", c.ImportHistory)
		api.GET("/mahasiswa/import-history/:id", c.ImportHistoryDetail)
	}
	// Frontend compatibility: older Admin Sistem screens used /api/admin/users.
	legacy := r.Group("/api/admin")
	legacy.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin_sistem"))
	legacy.GET("/users", c.Users)
}
