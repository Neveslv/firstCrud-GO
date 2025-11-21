package handlers

import (
	"CrudGO/database"
	"CrudGO/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CriarPost(c *gin.Context) {
	var input model.PostInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    err.Error(),
			"detalhes": err.Error(),
		})
		return
	}

	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id é obrigatório",
		})
		return
	}
	userID, _ := strconv.Atoi(userIDStr)

	var postID int
	query := "INSERT INTO posts (user_id, titulo, content) VALUES ($1, $2, $3) RETURNING id"
	err := database.DB.QueryRow(query, userID, input.Titulo, input.Content).Scan(&postID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Erro ao criar post",
			"detalhes": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensagem": "Post foi criado com sucesso.",
		"id":       postID,
	})
}

func ListarPosts(c *gin.Context) {
	rows, err := database.DB.Query("SELECT p.id, p.user_id, p.content, p.created_at, u.nome, u.email FROM posts p INNER JOIN usuario u ON p.user_id = u.id ORDER BY p.created_at DESC ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Erro ao buscar posts",
			"detalhes": err.Error(),
		})
		return
	}
	defer rows.Close()

	var posts []model.Post

	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Titulo, &post.Content, &post.CreatedAt, &post.UserName, &post.UserEmail); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":    "Erro ao processar dados",
				"detalhes": err.Error(),
			})
			return
		}
		posts = append(posts, post)
	}

	if posts == nil {
		posts = []model.Post{}
	}

	c.JSON(http.StatusOK, posts)
}

func BuscarPostPorId(c *gin.Context) {
	id := c.Param("id")
	var post model.Post

	query := "SELECT p.id, p.user_id, p.titulo, p.content, p.created_at, u.nome, u.email FROM posts p INNER JOIN usuario u ON p.user_id = u.id WHERE p.id = $1"
	err := database.DB.QueryRow(query, id).Scan(&post.ID, &post.UserID, &post.Titulo, &post.Content, &post.CreatedAt, &post.UserName, &post.UserEmail)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post não encontrado",
		})
		return
	}
	c.JSON(http.StatusOK, post)
}

func ListarPostPorUsuario(c *gin.Context) {
	userID := c.Param("user_id")

	rows, err := database.DB.Query("SELECT p.id, p.user_id, p.content, p.created_at, u.nome, u.email FROM posts p INNER JOIN usuario u ON p.user_id = u.id WHERE p.user_id=$1 ORDER BY p.created_at DESC ", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Erro ao buscar posts",
			"detalhes": err.Error(),
		})
		return
	}
	defer rows.Close()

	var posts []model.Post

	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Titulo, &post.Content, &post.CreatedAt, &post.UserName, &post.UserEmail); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":    "Erro ao processar dados",
				"detalhes": err.Error(),
			})
			return
		}
		posts = append(posts, post)
	}
	if posts == nil {
		posts = []model.Post{}

		c.JSON(http.StatusOK, posts)
	}
}

func AtualizarPost(c *gin.Context) {
	id := c.Param("id")

	var input model.PostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Dados inválidos",
			"detalhes": err.Error(),
		})
	}

	result, err := database.DB.Exec("UPDATE posts SET titulo = $1, content = $2 WHERE id = $3", input.Titulo, input.Content, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao atualizar o post",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post não encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensagem": "Post atualizado com sucesso",
	})
}

func DeletarPost(c *gin.Context) {
	id := c.Param("id")
	result, err := database.DB.Exec("DELETE post WHERE id = $1`", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao deletar o post",
		})
		return
	}

	rowsaffected, _ := result.RowsAffected()
	if rowsaffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post não encontrado",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"mensagem": "Post deletado com sucesso",
	})
}
