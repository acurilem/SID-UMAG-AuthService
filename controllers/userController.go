package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// User godoc
//
//	@Summary		Obtener usuario
//	@Description	Obtener usuario por medio del token JWT en la cabecera Authorization
//	@Tags			auth
//	@Accept			json
//	@Product		json
//	@Success		200	{object}	models.User
//	@Failure		401	{object}	smaps.ErrorRes	"Sin info en SAYD"
//
//	@Router			/api/v1/auth/user [get]
func User(c *gin.Context) {
	user, err := authService.GetUser(c)
	if err != nil {
		log.Println("No fue posible encontrar al usuario en SAYD / Vista Personas")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sin info en SAYD"})
		return
	}
	c.JSON(http.StatusOK, user)
}
