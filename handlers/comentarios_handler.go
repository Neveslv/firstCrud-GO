package handlers

import (
	"CrudGO/database"
	"CrudGO/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CriarComentario(c *gin.Context) {
	postID := c.Param("id")

	var comentario model.Comentarios

	contentType := c.GetHeader("Content-Type")

	if contentType == "application/json" {
		if err := c.ShouldBindJSON(&comentario); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido", "detalhes": err.Error()})
			return
		}
	} else {
		comentario.Content = c.PostForm("content")
		comentario.UserID, _ = strconv.Atoi(c.PostForm("user_id"))
	}

	if comentario.Content == "" || comentario.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Conteúdo e User ID são obrigatórios"})
		return
	}

	var newID int
	query := "INSERT INTO comentarios (post_id, user_id, content) VALUES ($1, $2, $3) RETURNING id"

	err := database.DB.QueryRow(query, postID, comentario.UserID, comentario.Content).Scan(&newID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar comentário", "detalhes": err.Error()})
		return
	}

	if contentType == "application/json" {
		c.JSON(http.StatusCreated, gin.H{
			"mensagem": "Comentário criado com sucesso",
			"id":       newID,
		})
	} else {
		c.Redirect(http.StatusFound, "/posts/"+postID+"/detalhes?msg=Comentário+enviado!&type=success")
	}
}

func ListarComentarioPorPost(c *gin.Context) {
	postId := c.Param("id")
	rows, err := database.DB.Query(`
        SELECT 
            co.id, co.post_id, co.user_id, co.content, co.created_at,
            u.nome, u.email
        FROM comentarios co
        INNER JOIN usuario u ON co.user_id = u.id
        WHERE co.post_id = $1
        ORDER BY co.created_at ASC
    `, postId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar comentários", "detalhes": err.Error()})
		return
	}
	defer rows.Close()

	var comentarios []model.Comentarios

	for rows.Next() {
		var com model.Comentarios
		err := rows.Scan(
			&com.ID, &com.PostID, &com.UserID, &com.Content, &com.CreatedAt,
			&com.UserNome, &com.UserEmail,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar dados"})
			return
		}
		comentarios = append(comentarios, com)
	}

	if comentarios == nil {
		comentarios = []model.Comentarios{}
	}

	c.JSON(http.StatusOK, comentarios)
}

func BuscarComentarioPorId(c *gin.Context) {
	id := c.Param("id")
	var com model.Comentarios

	query := `
        SELECT co.id, co.post_id, co.user_id, co.content, co.created_at, u.nome, u.email
        FROM comentarios co
        INNER JOIN usuario u ON co.user_id = u.id
        WHERE co.id = $1
    `
	err := database.DB.QueryRow(query, id).Scan(
		&com.ID, &com.PostID, &com.UserID, &com.Content, &com.CreatedAt,
		&com.UserNome, &com.UserEmail,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comentário não encontrado"})
		return
	}
	c.JSON(http.StatusOK, com)
}

func AtualizarComentario(c *gin.Context) {
	id := c.Param("id")
	var comentario model.Comentarios

	contentType := c.GetHeader("Content-Type")

	if contentType == "application/json" {
		if err := c.ShouldBindJSON(&comentario); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
			return
		}
	} else {
		comentario.Content = c.PostForm("content")
	}

	if comentario.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Conteúdo é obrigatório"})
		return
	}

	var postID int
	database.DB.QueryRow("SELECT post_id FROM comentarios WHERE id = $1", id).Scan(&postID)

	query := "UPDATE comentarios SET content = $1 WHERE id = $2"
	result, err := database.DB.Exec(query, comentario.Content, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comentário não encontrado"})
		return
	}

	if contentType == "application/json" {
		c.JSON(http.StatusOK, gin.H{"mensagem": "Comentário atualizado"})
	} else {
		c.Redirect(http.StatusFound, "/posts/"+strconv.Itoa(postID)+"/detalhes?msg=Comentário+editado!&type=success")
	}
}

func DeletarComentario(c *gin.Context) {
	id := c.Param("id")

	var postID int
	err := database.DB.QueryRow("SELECT post_id FROM comentarios WHERE id = $1", id).Scan(&postID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comentário não encontrado"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM comentarios WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comentário não encontrado"})
		return
	}

	contentType := c.GetHeader("Content-Type")

	if contentType == "application/json" {
		c.JSON(http.StatusOK, gin.H{"mensagem": "Comentário deletado"})
	} else {

		c.Redirect(http.StatusFound, "/posts/"+strconv.Itoa(postID)+"/detalhes?msg=Comentário+apagado!&type=success")
	}
}
