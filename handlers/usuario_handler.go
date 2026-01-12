package handlers

import (
	"CrudGO/database"
	"CrudGO/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CriarUsuario(c *gin.Context) {
	var usuario model.Usuario

	contentType := c.GetHeader("Content-Type")

	if contentType == "application/json" {
		if err := c.ShouldBindJSON(&usuario); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			return
		}
	} else {
		usuario.Nome = c.PostForm("nome")
		usuario.Email = c.PostForm("email")
		senhaPura := c.PostForm("password")

		if senhaPura == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "A senha é obrigatória"})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(senhaPura), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar senha"})
			return
		}
		usuario.Password = string(hash)
	}

	if usuario.Nome == "" || usuario.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome e Email são obrigatórios"})
		return
	}

	var id int
	query := "INSERT INTO usuario (nome, email, password) VALUES ($1, $2, $3) RETURNING id"
	err := database.DB.QueryRow(query, usuario.Nome, usuario.Email, usuario.Password).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário. Email já existe?"})
		return
	}

	if contentType == "application/json" {
		c.JSON(http.StatusCreated, gin.H{"mensagem": "Usuário criado", "id": id})
	} else {
		c.Redirect(http.StatusFound, "/usuarios?msg=Usuário+criado+com+sucesso!&type=success")
	}
}

func ListarUsuarios(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, nome, email, created_at FROM usuario ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuários"})
		return
	}
	defer rows.Close()

	var usuarios []model.Usuario
	for rows.Next() {
		var u model.Usuario
		rows.Scan(&u.ID, &u.Nome, &u.Email, &u.CreatedAt)
		usuarios = append(usuarios, u)
	}
	if usuarios == nil {
		usuarios = []model.Usuario{}
	}
	c.JSON(http.StatusOK, usuarios)
}

func BuscarUsuarioPorID(c *gin.Context) {
	id := c.Param("id")
	var u model.Usuario
	query := "SELECT id, nome, email, created_at FROM usuario WHERE id = $1"
	err := database.DB.QueryRow(query, id).Scan(&u.ID, &u.Nome, &u.Email, &u.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	c.JSON(http.StatusOK, u)
}

func AtualizarUsuario(c *gin.Context) {
	id := c.Param("id")
	var usuario model.Usuario

	contentType := c.GetHeader("Content-Type")
	if contentType == "application/json" {
		c.ShouldBindJSON(&usuario)
	} else {
		usuario.Nome = c.PostForm("nome")
		usuario.Email = c.PostForm("email")
	}

	query := "UPDATE usuario SET nome = $1, email = $2 WHERE id = $3"
	result, err := database.DB.Exec(query, usuario.Nome, usuario.Email, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar"})
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	if contentType == "application/json" {
		c.JSON(http.StatusOK, gin.H{"mensagem": "Atualizado"})
	} else {
		c.Redirect(http.StatusFound, "/usuarios?msg=Usuário+atualizado!&type=success")
	}
}

func DeletarUsuario(c *gin.Context) {
	userIDToken, exists := c.Get("user_id")
	idDaURL := c.Param("id")

	if exists && fmt.Sprintf("%v", userIDToken) == idDaURL {
		c.SetCookie("token", "", -1, "/", "", false, true)
		c.Redirect(http.StatusFound, "/login?msg=Sua+conta+foi+encerrada.&type=info")
		return
	}

	if c.GetHeader("Content-Type") == "application/json" {
		c.JSON(http.StatusOK, gin.H{"mensagem": "Deletado"})
	} else {
		c.Redirect(http.StatusFound, "/usuarios?msg=Usuário+removido!&type=success")
	}
}

func ListarUsuariosHTML(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, nome, email, created_at FROM usuario ORDER BY created_at DESC")
	if err != nil {
		c.HTML(http.StatusOK, "usuarios.html", gin.H{"error": "Erro ao buscar"})
		return
	}
	defer rows.Close()

	var usuarios []model.Usuario
	for rows.Next() {
		var u model.Usuario
		rows.Scan(&u.ID, &u.Nome, &u.Email, &u.CreatedAt)
		usuarios = append(usuarios, u)
	}
	c.HTML(http.StatusOK, "usuarios.html", gin.H{"usuarios": usuarios})
}

func ExibirTelaConfiguracoes(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var usuario model.Usuario
	id := int(userID.(float64))
	querySegura := "SELECT id, nome, email, COALESCE(bio, ''), COALESCE(site, '') FROM usuario WHERE id = $1"

	err := database.DB.QueryRow(querySegura, id).Scan(&usuario.ID, &usuario.Nome, &usuario.Email, &usuario.Bio, &usuario.Site)

	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "configuracoes.html", gin.H{"usuario": usuario})
}

func AtualizarConfiguracoes(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		fmt.Println("ERRO: Usuário não autenticado no contexto")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	nome := c.PostForm("nome")
	email := c.PostForm("email")
	bio := c.PostForm("bio")
	site := c.PostForm("site")

	query := "UPDATE usuario SET nome=$1, email=$2, bio=$3, site=$4 WHERE id=$5"

	_, err := database.DB.Exec(query, nome, email, bio, site, userID)

	if err != nil {
		fmt.Println("ERRO SQL:", err)
		c.Redirect(http.StatusFound, "/configuracoes?msg=Erro+ao+salvar&type=error")
		return
	}
	c.Redirect(http.StatusFound, "/configuracoes?msg=Perfil+atualizado!&type=success")
}

func ExibirPerfilPublico(c *gin.Context) {
	id := c.Param("id")

	var usuario model.Usuario

	query := "SELECT id, nome, email, COALESCE(bio, ''), COALESCE(site, ''), created_at FROM usuario WHERE id = $1"
	err := database.DB.QueryRow(query, id).Scan(&usuario.ID, &usuario.Nome, &usuario.Email, &usuario.Bio, &usuario.Site, &usuario.CreatedAt)

	if err != nil {
		c.HTML(http.StatusNotFound, "home.html", gin.H{
			"error": "Usuário não encontrado",
		})
	}

	queryPost := "SELECT id, titulo, content, created_at FROM posts WHERE user_id = $1 ORDER BY created_at DESC"
	rows, err := database.DB.Query(queryPost, id)
	if err != nil {
		rows = nil
	} else {
		defer rows.Close()
	}

	var posts []model.Post
	if rows != nil {
		for rows.Next() {
			var p model.Post
			if err := rows.Scan(&p.ID, &p.Titulo, &p.Titulo, &p.Content, &p.CreatedAt); err == nil {
				posts = append(posts, p)
			}
		}
	}

	tokenString, _ := c.Cookie("token")
	logado := tokenString != ""

	c.HTML(http.StatusOK, "perfil_publico.html", gin.H{
		"usuario": usuario,
		"posts":   posts,
		"logado":  logado,
	})
}
