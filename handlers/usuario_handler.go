package handlers

import (
	"CrudGO/database"
	"CrudGO/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CREATE
func CriarUsuario(c *gin.Context) {
	var input model.UsuarioInput
	id := c.Param("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Dados inválidos",
			"detalhes": err.Error(),
		})
		return
	}

	query := "INSERT INTO usuario (nome, email) VALUES ($1, $2) RETURNING id"
	err := database.DB.QueryRow(query, input.Nome, input.Email).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Erro ao criar usuário",
			"detalhes": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensagem": "Usuário criado com sucesso",
		"id":       id,
	})
}

func ListarUsuarios(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, nome, email, created_at
		FROM usuario
		ORDER BY created_at DESC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Erro ao buscar usuários",
			"detalhes": err.Error(),
		})
		return
	}
	defer rows.Close()

	var usuarios []model.Usuario

	for rows.Next() {
		var usuario model.Usuario
		err := rows.Scan(
			&usuario.ID, &usuario.Nome,
			&usuario.Email, &usuario.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":    "Erro ao processar dados",
				"detalhes": err.Error(),
			})
			return
		}
		usuarios = append(usuarios, usuario)
	}

	if usuarios == nil {
		usuarios = []model.Usuario{}
	}

	c.JSON(http.StatusOK, usuarios)
}

func BuscarUsuarioPorID(c *gin.Context) {
	id := c.Param("id")

	var usuario model.Usuario
	query := `
		SELECT id, nome, email, created_at
		FROM usuario
		WHERE id = $1
	`
	err := database.DB.QueryRow(query, id).Scan(
		&usuario.ID, &usuario.Nome,
		&usuario.Email, &usuario.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuário não encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, usuario)
}

func AtualizarUsuario(c *gin.Context) {
	id := c.Param("id")
	var input model.UsuarioInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Dados inválidos",
			"detalhes": err.Error(),
		})
		return
	}

	result, err := database.DB.Exec(`
		UPDATE usuario 
		SET nome = $1, email = $2
		WHERE id = $3
	`, input.Nome, input.Email, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Erro ao atualizar usuário",
			"detalhes": err.Error(),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuário não encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensagem": "Usuário atualizado com sucesso",
	})
}

func DeletarUsuario(c *gin.Context) {
	id := c.Param("id")

	result, err := database.DB.Exec("DELETE FROM usuario WHERE id = $1", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Erro ao deletar usuário",
			"detalhes": err.Error(),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuário não encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensagem": "Usuário deletado com sucesso",
	})
}
