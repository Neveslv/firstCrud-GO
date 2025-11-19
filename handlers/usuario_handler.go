package handlers

import (
	"CrudGO/database"
	"CrudGO/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CriarUsuario(c *gin.Context) {
	var usuario model.Usuario

	if err := c.ShouldBindJSON(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	query := "INSERT INTO usuario (nome, email) VALUES ($1, $2) RETURNING id, created_at"
	err := database.DB.QueryRow(query, usuario.Nome, usuario.Email).Scan(&usuario.ID, &usuario.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao criar usuário",
		})
		return
	}

	c.JSON(http.StatusCreated, usuario)

}

func ListarUsuarios(c *gin.Context) {
	rows, err := database.DB.Query("SELECT * FROM usuario ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	var usuarios []model.Usuario

	for rows.Next() {
		var usuario model.Usuario
		if err := rows.Scan(&usuario.ID, &usuario.Nome, &usuario.Email, &usuario.CreatedAt); err != nil {
			continue
		}
		usuarios = append(usuarios, usuario)
	}

	c.JSON(http.StatusOK, usuarios)
}

func BuscarPorId(c *gin.Context) {
	id := c.Param("id")
	var usuario model.Usuario

	query := "SELECT id, nome, email, created_at FROM users WHERE id = $1"
	err := database.DB.QueryRow(query, id).Scan(&usuario.ID, &usuario.Nome, &usuario.Email, &usuario.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuario não encontrado",
		})
		return
	}
	c.JSON(http.StatusOK, usuario)
}
