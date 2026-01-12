package routes

import (
	"CrudGO/handlers"
	"CrudGO/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.GET("/usuarios/novo", func(c *gin.Context) {
		c.HTML(http.StatusOK, "usuario_form.html", nil)
	})

	r.GET("/perfil/:id", handlers.ExibirPerfilPublico)

	admin := r.Group("/")
	admin.Use(middleware.AutentMiddleware())
	{
		admin.GET("/usuarios", handlers.ListarUsuariosHTML)
	}

	protected := r.Group("/")
	protected.Use(middleware.AutentMiddleware())
	{
		protected.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})
		protected.GET("/posts", handlers.ListarPostHTML)
		protected.GET("/posts/novo", func(c *gin.Context) {
			c.HTML(http.StatusOK, "post_form.html", nil)
		})
		protected.GET("/posts/:id/detalhes", handlers.ExibirDetalhesPostHTML)

		protected.GET("/configuracoes", handlers.ExibirTelaConfiguracoes)
		protected.POST("/configuracoes", handlers.AtualizarConfiguracoes)
	}
}
