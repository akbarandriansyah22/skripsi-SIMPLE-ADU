package routes

import (
	controller "backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	admin := controller.NewAdminController()

	api := r.Group("/api/admin")
	api.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("petugas", "admin", "admin_fakultas", "Admin Fakultas"))
	{
		api.GET("/dashboard", admin.Dashboard)
		api.GET("/pengaduan", admin.GetAllPengaduan)
		api.GET("/pengaduan/:id", admin.GetPengaduanByID)
		api.PATCH("/pengaduan/:id/status", admin.UpdateStatus)
		api.PATCH("/pengaduan/:id/unit", admin.AssignUnit)
		api.PATCH("/pengaduan/:id/forward", admin.ForwardToPimpinan)
		api.GET("/unit", admin.GetUnits)
	}
}
