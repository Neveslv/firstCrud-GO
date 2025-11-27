package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"mensagem": "API CrudGO está rodando!",
			"versao":   "1.0",
			"rotas": gin.H{
				"usuarios":    "/api/usuarios",
				"posts":       "/api/posts",
				"comentarios": "/api/comentarios",
			},
		})
	})

	r.GET("/usuarios", func(c *gin.Context) {
		c.HTML(http.StatusOK, "usuarios.html", nil)
	})
	r.GET("/usuarios/novo", func(c *gin.Context) {
		c.HTML(http.StatusOK, "usuario_form.html", nil)
	})

	r.GET("/posts", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts.html", nil)
	})
	r.GET("/post/novo", func(c *gin.Context) {
		c.HTML(http.StatusOK, "post_form.html", nil)
	})
	r.GET("/post/:id/detalhes", func(c *gin.Context) {
		c.HTML(http.StatusOK, "post_detalhes.html", nil)
	})

}
