package routes

import (
	"CrudGO/handlers"

	"github.com/gin-gonic/gin"
)

func SetupComentarioRoutes(r *gin.Engine) {
	comentarios := r.Group("/api")
	{
		comentarios.GET("/comentarios/:id", handlers.BuscarComentarioPorId)
		comentarios.PUT("/comentarios/:id", handlers.AtualizarComentario)
		comentarios.DELETE("/comentarios/:id", handlers.DeletarComentario)
	}
}
