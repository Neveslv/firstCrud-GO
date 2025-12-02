package handlers

import (
	"CrudGO/database"
	"CrudGO/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CriarUsuario(c *gin.Context) {
	var usuario model.Usuario

	contentType := c.GetHeader("Content-Type")

	if contentType == "application/json" {
		if err := c.ShouldBindJSON(&usuario); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "JSON inválido",
				"detalhes": err.Error(),
			})
			return
		}
	} else {
		usuario.Nome = c.PostForm("nome")
		usuario.Email = c.PostForm("email")
	}
	if usuario.Nome == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "O campo 'nome' é obrigatório",
		})
		return
	}

	if usuario.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "O campo 'email' é obrigatório",
		})
		return
	}

	var id int
	query := "INSERT INTO usuario (nome, email) VALUES ($1, $2) RETURNING id"
	err := database.DB.QueryRow(query, usuario.Nome, usuario.Email).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Erro ao criar usuário",
			"detalhes": err.Error(),
		})
		return
	}

	if contentType == "application/json" {
		c.JSON(http.StatusCreated, gin.H{
			"mensagem": "Usuário criado com sucesso",
			"id":       id,
			"nome":     usuario.Nome,
			"email":    usuario.Email,
		})
	} else {
		c.Redirect(http.StatusFound, "/usuarios")
	}
}

func ListarUsuarios(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, nome, email, created_at FROM usuario ORDER BY created_at DESC")
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
		err := rows.Scan(&usuario.ID, &usuario.Nome, &usuario.Email, &usuario.CreatedAt)
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

	query := "SELECT id, nome, email, created_at FROM usuario WHERE id = $1"
	err := database.DB.QueryRow(query, id).Scan(&usuario.ID, &usuario.Nome, &usuario.Email, &usuario.CreatedAt)

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
	var usuario model.Usuario

	contentType := c.GetHeader("Content-Type")

	if contentType == "application/json" {
		if err := c.ShouldBindJSON(&usuario); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":    "JSON inválido",
				"detalhes": err.Error(),
			})
			return
		}
	} else {
		usuario.Nome = c.PostForm("nome")
		usuario.Email = c.PostForm("email")
	}

	if usuario.Nome == "" || usuario.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Os campos 'nome' e 'email' são obrigatórios",
		})
		return
	}

	query := "UPDATE usuario SET nome = $1, email = $2 WHERE id = $3"
	result, err := database.DB.Exec(query, usuario.Nome, usuario.Email, id)

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

	if contentType == "application/json" {
		c.JSON(http.StatusOK, gin.H{
			"mensagem": "Usuário atualizado com sucesso",
		})
	} else {
		c.Redirect(http.StatusFound, "/usuarios")
	}
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

func ListarUsuariosHTML(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, nome, email, created_at FROM usuario ORDER BY created_at DESC ")
	if err != nil {
		c.HTML(http.StatusOK, "usuario.html", gin.H{
			"usuarios": []model.Usuario{},
			"error":    "Erro ao buscar o usuario",
		})
		return
	}
	defer rows.Close()

	var usuarios []model.Usuario

	for rows.Next() {
		var usuario model.Usuario
		err := rows.Scan(&usuario.ID, &usuario.Nome, &usuario.Email, &usuario.CreatedAt)
		if err != nil {
			c.HTML(http.StatusOK, "usuarios.html", gin.H{
				"usuarios": []model.Usuario{},
				"error":    "Erro ao processar dados",
			})
			return
		}

		usuarios = append(usuarios, usuario)

	}

	c.HTML(http.StatusOK, "usuarios.html", gin.H{
		"usuarios": usuarios,
	})
}
