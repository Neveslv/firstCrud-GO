package routes

import (
	"CrudGO/handlers"

	"github.com/gin-gonic/gin"
)

func SetupPostRoutes(r *gin.Engine) {
	posts := r.Group("/api")
	{
		posts.POST("/posts", handlers.CriarPost)
		posts.GET("/posts", handlers.ListarPosts)
		posts.GET("/posts/:id", handlers.BuscarPostPorId)
		posts.PUT("/posts/:id", handlers.AtualizarPost)
		posts.DELETE("/posts/:id", handlers.DeletarPost)
		posts.POST("posts/:id", handlers.DeletarPost)

		posts.GET("/posts/usuario/:user_id", handlers.ListarPostPorUsuario)

		posts.POST("/posts/:id/comentarios", handlers.CriarComentario)
		posts.GET("/posts/:id/comentarios", handlers.ListarComentarioPorPost)
	}
}
