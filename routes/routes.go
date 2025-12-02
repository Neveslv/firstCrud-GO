package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
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
	r.GET("/posts/novo", func(c *gin.Context) {
		c.HTML(http.StatusOK, "post_form.html", nil)
	})
	r.GET("/posts/:id/detalhes", func(c *gin.Context) {
		c.HTML(http.StatusOK, "post_detalhes.html", nil)
	})

}
