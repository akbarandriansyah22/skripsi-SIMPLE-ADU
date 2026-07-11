package routes

import (
	controller "backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func NotifikasiRoutes(r *gin.Engine) {
	notifikasi := controller.NewNotifikasiController()

	api := r.Group("/api/notifikasi")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("", notifikasi.Mine)
		api.GET("/unread", notifikasi.Unread)
		api.PATCH("/:id/read", notifikasi.MarkAsRead)
	}
}
