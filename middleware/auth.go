package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/citiaps/SID-UMAG-AuthService/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_KEY"))

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := services.GetToken(c)
		token, err := jwt.ParseWithClaims(tokenString, &services.Sid3Claims{}, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		}, jwt.WithLeeway(5*time.Second))

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Error de autenticación"})
			return
		}

		if claims, ok := token.Claims.(*services.Sid3Claims); ok && token.Valid {
			if claims.CodPersona > 0 {
				c.Next()
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token no valido"})
				c.Abort()
				return
			}

		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
	}
}
