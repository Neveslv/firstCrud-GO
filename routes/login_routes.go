package routes

import (
	"CrudGO/handlers"

	"github.com/gin-gonic/gin"
)

func SetupLoginRoutes(r *gin.Engine) {
	r.GET("/login", handlers.ExibirTelaLogin)
	r.POST("/login", handlers.Login)
	r.GET("/logout", handlers.Logout)
}
