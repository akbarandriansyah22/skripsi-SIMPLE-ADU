package routes

import (
	controller "backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func KasubagRoutes(r *gin.Engine) {
	kasubag := controller.NewKasubagController()
	api := r.Group("/api/kasubag")
	api.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("kasubag"))
	{
		api.GET("/dashboard", kasubag.Dashboard)
		api.GET("/pengaduan", kasubag.List)
		api.GET("/pengaduan/:id", kasubag.Detail)
		api.PATCH("/pengaduan/:id/proses", kasubag.StartProcess)
		api.PATCH("/pengaduan/:id/mulai-proses", kasubag.StartProcess)
		api.POST("/pengaduan/:id/respon", kasubag.AddResponse)
		api.PATCH("/pengaduan/:id/selesai", kasubag.Finish)
		api.PATCH("/pengaduan/:id/kembalikan", kasubag.ReturnToAdmin)
	}
}
