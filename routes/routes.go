package routes

import (
	"CrudGO/handlers"
	"CrudGO/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.GET("/", handlers.ExibirHomeBlog)

	r.GET("/posts/:id/detalhes", handlers.ExibirDetalhesPostHTML)

	r.GET("/perfil/:id", handlers.ExibirPerfilPublico)

	r.GET("/usuarios/novo", func(c *gin.Context) {
		c.HTML(http.StatusOK, "usuario_form.html", nil)
	})

	protected := r.Group("/")
	protected.Use(middleware.AutentMiddleware())
	{
		protected.GET("/dashboard", func(c *gin.Context) {
			userID, userNome, isAdmin, _ := handlers.GetDadosUsuario(c)
			c.HTML(http.StatusOK, "index.html", gin.H{
				"is_admin":  isAdmin,
				"user_id":   userID,
				"user_nome": userNome,
			})
		})

		protected.GET("/posts", handlers.ListarPostHTML)

		protected.GET("/posts/novo", func(c *gin.Context) {
			c.HTML(http.StatusOK, "post_form.html", nil)
		})

		protected.GET("/configuracoes", handlers.ExibirTelaConfiguracoes)
		protected.POST("/configuracoes", handlers.AtualizarConfiguracoes)
	}

	admin := r.Group("/")
	admin.Use(middleware.AdminAuthMiddleware())
	{
		// CRUD USUÁRIOS: Listagem e ações de admin
		admin.GET("/usuarios", handlers.ListarUsuariosHTML)
	}
}
