package routes

import (
	controller "backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	admin := controller.NewAdminController()

	api := r.Group("/api/admin")
	api.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin_fakultas"))
	{
		api.GET("/dashboard", admin.Dashboard)
		api.GET("/pengaduan", admin.GetAllPengaduan)
		api.GET("/pengaduan/:id", admin.GetPengaduanByID)
		api.PATCH("/pengaduan/:id/status", admin.UpdateStatus)
		api.PATCH("/pengaduan/:id/unit", admin.AssignUnit)
		api.PATCH("/pengaduan/:id/validasi", admin.ValidatePengaduan)
		api.POST("/pengaduan/:id/validasi", admin.ValidatePengaduan)
		api.PATCH("/pengaduan/:id/forward", admin.ForwardToPimpinan)
		api.POST("/pengaduan/:id/teruskan-unit", admin.AssignUnit)
		api.POST("/pengaduan/:id/teruskan-pimpinan", admin.ForwardToPimpinan)
		api.POST("/pengaduan/:id/reanalyze", admin.ReanalyzeAI)
		api.GET("/unit/penanganan", admin.GetAssignmentUnits)
	}
	// Unit lookup is read-only and remains available to the legacy Pimpinan UI.
	legacyUnits := r.Group("/api/admin")
	legacyUnits.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin_fakultas", "pimpinan_fakultas"))
	legacyUnits.GET("/unit", admin.GetUnits)
}

func AdminFakultasRoutes(r *gin.Engine) {
	admin := controller.NewAdminController()
	api := r.Group("/api/admin-fakultas")
	api.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin_fakultas"))
	{
		api.GET("/dashboard", admin.Dashboard)
		api.GET("/pengaduan", admin.GetAllPengaduan)
		api.GET("/pengaduan/:id", admin.GetPengaduanByID)
		api.POST("/pengaduan/:id/validasi", admin.ValidatePengaduan)
		api.POST("/pengaduan/:id/teruskan-unit", admin.AssignUnit)
		api.POST("/pengaduan/:id/teruskan-pimpinan", admin.ForwardToPimpinan)
		api.GET("/units", admin.GetUnits)
		api.GET("/unit/penanganan", admin.GetAssignmentUnits)
	}
}
