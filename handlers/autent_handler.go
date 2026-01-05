package handlers

import (
	"CrudGO/database"
	"CrudGO/model"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	var usuario model.Usuario
	query := "SELECT id, nome, email, password FROM usuario WHERE email=$1"
	err := database.DB.QueryRow(query, email).Scan(&usuario.ID, &usuario.Nome, &usuario.Email, &usuario.Password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Email ou senha inválidos",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(password))
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Email ou senha inválidos",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usuario.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"error": "Erro ao gerar token",
		})
		return
	}

	c.SetCookie("token", tokenString, 3600*24, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func ExibirTelaLogin(c *gin.Context) {
	_, err := c.Cookie("token")
	if err == nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "login.html", nil)
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

func RegistrarUsuario(c *gin.Context) {
	var usuario model.Usuario
	usuario.Nome = c.PostForm("nome")
	usuario.Email = c.PostForm("email")
	normPassword := c.PostForm("password")

	if usuario.Nome == "" || usuario.Email == "" || normPassword == "" {
		c.HTML(http.StatusBadRequest, "registrar.html", gin.H{
			"error": "Preencha todos os campos",
		})
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(normPassword), bcrypt.DefaultCost)
	usuario.Password = string(hashPassword)

	query := "INSERT INTO usuario (nome, email, password) VALUES ($1, $2, $3)"
	_, err := database.DB.Exec(query, usuario.Nome, usuario.Email, usuario.Password)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "registrar.html", gin.H{
			"error": "Erro ao criar conta.",
		})
		return

	}

	c.Redirect(http.StatusFound, "/login?msg=Conta+criada!+Faça+login+para+continuar.&type=sucess")
}

func ExibirTelaRegistro(c *gin.Context) {
	c.HTML(http.StatusOK, "registrar.html", nil)
}
