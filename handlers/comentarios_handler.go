package handlers

import (
	"CrudGO/database"
	"CrudGO/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CriarComentario(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID do post inválido",
		})
		return
	}
	var input model.ComentarioInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Dados inválidos",
			"detalhes": err.Error(),
		})
		return
	}

	userIDstr := c.Query("user_id")
	if userIDstr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id é obrigatório",
		})
		return
	}

	userID, _ := strconv.Atoi(userIDstr)

	var comentarioID int
	query := "INSERT INTO comentarios (post_id, user_id, content) VALUES ($1, $2, $3) RETURNING id"
	err = database.DB.QueryRow(query, postID, userID, input.Content).Scan(&comentarioID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Erro ao criar um comentario",
			"detalhes": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensagem": "O comentario foi criado com sucesso",
		"id":       comentarioID,
	})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Erro ao buscar os comentarios",
			"detalhes": err.Error(),
		})
		return
	}
	defer rows.Close()

	var comentarios []model.Comentarios

	for rows.Next() {
		var comentario model.Comentarios
		err := rows.Scan(
			&comentario.ID, &comentario.PostID, &comentario.UserID,
			&comentario.Content, &comentario.CreatedAt,
			&comentario.UserNome, &comentario.UserEmail,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":    "Erro ao processar dados",
				"detalhes": err.Error(),
			})
			return
		}
		comentarios = append(comentarios, comentario)
	}

	if comentarios == nil {
		comentarios = []model.Comentarios{}
	}

	c.JSON(http.StatusOK, comentarios)
}

func BuscarComentarioPorId(c *gin.Context) {
	ID := c.Param("id")
	var comentario model.Comentarios
	query := `
        SELECT 
            co.id, co.post_id, co.user_id, co.content, co.created_at,
            u.nome, u.email
        FROM comentarios co
        INNER JOIN usuario u ON co.user_id = u.id
        WHERE co.id = $1
    `
	err := database.DB.QueryRow(query, ID).Scan(&comentario.ID, &comentario.PostID, &comentario.UserID, &comentario.Content, &comentario.CreatedAt, &comentario.UserNome, &comentario.UserEmail)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Comentario não encontrado",
		})
		return
	}
	c.JSON(http.StatusOK, comentario)
}

func AtualizarComentario(c *gin.Context) {
	id := c.Param("id")
	var input model.ComentarioInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Dados inválidos",
			"detalhes": err.Error(),
		})
		return
	}

	result, err := database.DB.Exec(`
        UPDATE comentarios 
        SET content = $1
        WHERE id = $2
    `, input.Content, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao atualizar o comentário",
		})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Comentario não encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensagem": "Comentario atualizado com sucesso",
	})
}

func DeletarComentario(c *gin.Context) {
	id := c.Param("id")
	result, err := database.DB.Exec("DELETE FROM comentarios WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Comentario não encontrado",
		})
		return
	}

	rowsaffected, _ := result.RowsAffected()
	if rowsaffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Comentario não encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensagem": "Comentario deletado com sucesso",
	})
}
