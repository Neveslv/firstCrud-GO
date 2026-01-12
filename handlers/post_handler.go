package handlers

import (
	"CrudGO/database"
	"CrudGO/model"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CriarPost(c *gin.Context) {
	var input model.PostInput

	contentType := c.GetHeader("Content-Type")
	if contentType == "application/json" {
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "JSON inválido",
				"detalhes": err.Error(),
			})
			return
		}
	} else {
		input.Titulo = c.PostForm("titulo")
		input.Content = c.PostForm("content")
	}
	if input.Titulo == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "O campo titulo é obrigatório",
		})
		return
	}

	if input.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "O campo content é obrigatório",
		})
		return
	}

	var userID int
	var err error
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	userID = int(userIDInterface.(float64))

	var postID int
	query := "INSERT INTO posts (user_id, titulo, content) VALUES ($1, $2, $3) RETURNING id"
	err = database.DB.QueryRow(query, userID, input.Titulo, input.Content).Scan(&postID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Erro ao criar post",
			"detalhes": err.Error(),
		})
		return
	}
	if contentType == "application/json" {
		c.JSON(http.StatusCreated, gin.H{
			"mensagem": "Post foi criado com sucesso.",
			"id":       postID,
		})
	} else {
		c.Redirect(http.StatusFound, "/posts?msg=Post+publicado!&type=success")
	}
}

func ListarPosts(c *gin.Context) {
	rows, err := database.DB.Query("SELECT p.id, p.user_id, p.titulo,  p.content, p.created_at, u.nome, u.email FROM posts p INNER JOIN usuario u ON p.user_id = u.id ORDER BY p.created_at DESC ")
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

	rows, err := database.DB.Query("SELECT p.id, p.user_id, p.titulo, p.content, p.created_at, u.nome, u.email FROM posts p INNER JOIN usuario u ON p.user_id = u.id WHERE p.user_id=$1 ORDER BY p.created_at DESC ", userID)
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

	contentType := c.GetHeader("Content-Type")

	if contentType == c.GetHeader("application/json") {
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "Dados inválidos",
				"detalhes": err.Error(),
			})
			return
		}
	} else {
		input.Titulo = c.PostForm("titulo")
		input.Content = c.PostForm("content")
	}

	if input.Titulo == "" || input.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Os campos 'titulo' e 'content' são obrigatórios",
		})
		return
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

	if contentType == "application/json" {
		c.JSON(http.StatusOK, gin.H{
			"mensagem": "Post atualizado com sucesso",
		})
	} else {
		c.Redirect(http.StatusFound, "/posts?msg=Post+editado!&type=success")
	}
}

func DeletarPost(c *gin.Context) {
	id := c.Param("id")
	result, err := database.DB.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Erro ao deletar o post",
			"detalhes": err.Error(),
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
	contentType := c.GetHeader("Content-Type")

	if contentType == "application/json" {
		c.JSON(http.StatusOK, gin.H{"mensagem": "Post deletado com sucesso"})
	} else {
		c.Redirect(http.StatusFound, "/posts?msg=Post+removido!&type=success")
	}
}

func ListarPostHTML(c *gin.Context) {
	query := `
		SELECT p.id, p.user_id, p.titulo, p.content, p.created_at, u.nome, u.email 
		FROM posts p 
		LEFT JOIN usuario u ON p.user_id = u.id 
		ORDER BY p.created_at DESC
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "posts.html", gin.H{
			"posts": []model.Post{},
			"error": "Erro ao buscar posts: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	var posts []model.Post

	for rows.Next() {
		var post model.Post
		var nomeTemp sql.NullString
		var emailTemp sql.NullString

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Titulo,
			&post.Content,
			&post.CreatedAt,
			&nomeTemp,
			&emailTemp,
		)

		if err != nil {
			c.HTML(http.StatusInternalServerError, "posts.html", gin.H{
				"posts": []model.Post{},
				"error": "Erro ao processar dados",
			})
			return
		}

		if nomeTemp.Valid {
			post.UserName = nomeTemp.String
		} else {
			post.UserName = "Usuário Desconhecido"
		}

		if emailTemp.Valid {
			post.UserEmail = emailTemp.String
		} else {
			post.UserEmail = ""
		}

		posts = append(posts, post)
	}

	c.HTML(http.StatusOK, "posts.html", gin.H{
		"posts": posts,
	})
}

func ExibirDetalhesPostHTML(c *gin.Context) {
	id := c.Param("id")

	tokenString, _ := c.Cookie("token")
	usuarioLogado := tokenString != ""

	var post model.Post
	queryPost := `
		SELECT p.id, p.user_id, p.titulo, p.content, p.created_at, u.nome, u.email 
		FROM posts p 
		LEFT JOIN usuario u ON p.user_id = u.id 
		WHERE p.id = $1`

	err := database.DB.QueryRow(queryPost, id).Scan(&post.ID, &post.UserID, &post.Titulo, &post.Content, &post.CreatedAt, &post.UserName, &post.UserEmail)
	if err != nil {
		c.HTML(http.StatusNotFound, "post_detalhes.html", gin.H{
			"error": "Post inexistente",
		})
		return
	}

	rows, err := database.DB.Query(`
		SELECT c.id, c.user_id, c.content, c.created_at, u.nome, u.email 
		FROM comentarios c 
		LEFT JOIN usuario u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC `, id)

	var comentarios []model.Comentarios

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var comentario model.Comentarios
			err := rows.Scan(
				&comentario.ID,
				&comentario.UserID,
				&comentario.Content,
				&comentario.CreatedAt,
				&comentario.UserNome,
				&comentario.UserEmail,
			)

			if err == nil {
				comentarios = append(comentarios, comentario)
			}
		}
	}
	c.HTML(http.StatusOK, "post_detalhes.html", gin.H{
		"post":        post,
		"comentarios": comentarios,
		"logado":      usuarioLogado,
	})
}
