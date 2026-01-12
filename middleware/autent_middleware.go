package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AutentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil

		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.Redirect(http.StatusFound, "/login?msg=Sessão+expirada")
				c.Abort()
				return
			}
			isAdmin, ok := claims["is_admin"].(bool)
			if !ok || !isAdmin {
				c.Redirect(http.StatusFound, "/?msg+Acesso+restrito+a+admin&type=error")
				c.Abort()
				return
			}

			c.Set("user_id", claims["sub"])
			c.Set("is_admin", isAdmin)
			c.Next()
		} else {
			c.Redirect(http.StatusFound, "/login?msg=Token+inválido")
			c.Abort()
		}
	}
}
