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
		api.PATCH("/pengaduan/:id/forward", admin.ForwardToPimpinan)
		api.POST("/pengaduan/:id/reanalyze", admin.ReanalyzeAI)
		api.GET("/unit", admin.GetUnits)
	}
}
