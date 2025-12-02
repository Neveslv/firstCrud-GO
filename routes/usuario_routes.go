package routes

import (
	"CrudGO/handlers"

	"github.com/gin-gonic/gin"
)

func SetupUsuarioRoutes(r *gin.Engine) {
	usuarios := r.Group("/api")
	{
		usuarios.POST("/usuarios", handlers.CriarUsuario)
		usuarios.GET("/usuarios", handlers.ListarUsuarios)
		usuarios.GET("/usuarios/:id", handlers.BuscarUsuarioPorID)
		usuarios.PUT("/usuarios/:id", handlers.AtualizarUsuario)
		usuarios.DELETE("/usuarios/:id", handlers.DeletarUsuario)
		usuarios.POST("usuarios/:id", handlers.DeletarUsuario)
	}
}
