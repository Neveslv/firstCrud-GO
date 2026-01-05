package routes

import (
	"CrudGO/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRegistroRoutes(r *gin.Engine) {
	r.GET("/registrar", handlers.ExibirTelaRegistro)
	r.POST("/registrar", handlers.RegistrarUsuario)

}
