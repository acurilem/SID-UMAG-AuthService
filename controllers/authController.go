package controller

import (
	"log"
	"net/http"

	"github.com/acurilem/SID-UMAG-AuthService/forms"
	"github.com/acurilem/SID-UMAG-AuthService/services"
	"github.com/gin-gonic/gin"
)

// LoginFunc godoc
//
//	@Summary		Loggear
//	@Description	Loggeo dentro del sistema que pasa por LDAP
//	@Param			loginData	body	forms.LoginForm	true	"Nombre de usuario y contraseña del usuario"
//	@Tags			auth
//	@Accept			json
//	@Product		json
//	@Success		200	{object}	smaps.LoginModel
//	@Failure		400	{object}	smaps.ErrorRes	"Datos enviados no cumplen con módelo"
//	@Failure		500	{object}	smaps.ErrorRes	"No se pudo generar el token"
//
//	@Failure		401	{object}	smaps.ErrorRes	"Credenciales inválidas LDAP / Sin info en SAYD"
//
//	@Router			/api/v1/auth/login [post]
func LoginFunc(c *gin.Context) {
	var loginValues *forms.LoginForm
	if err := c.ShouldBindJSON(&loginValues); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Datos enviados no cumplen con módelo",
		})
		return
	}

	var rut string
	var err error
	if authService.IsRutOrPassport(loginValues.Username) {
		rut, err = authService.LoginLdapv2WithRut(loginValues)
	} else {
		rut, err = authService.LoginLdapv2(loginValues)
	}

	if err != nil {
		log.Println("No fue posible encontrar al usuario en LDAP")
		c.AbortWithError(http.StatusUnauthorized, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas LDAP"})
		return
	}
	user, err := services.GetUserInfoFromRutService(rut)
	if err != nil {
		log.Println("No fue posible encontrar al usuario en SAYD / Vista Personas")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sin info en SAYD"})
		return
	}
	token, refreshToken, err1, err2 := authService.LoadJWTAuth(user.NombreCompleto, user.MailInstitucional, user.ID)
	if err1 != nil || err2 != nil {
		//if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token":         token,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

// RefreshToken godoc
//
//	@Summary		Refrescar token JWT
//	@Description	Refresco del token JWT
//	@Param			loginData	body	forms.RefreshTokenForm	true	"Token de refresco"
//	@Tags			auth
//	@Accept			json
//	@Product		json
//	@Success		200	{object}	smaps.LoginModel
//	@Failure		503	{object}	smaps.ErrorRes	"Error en el servidor"
//	@Failure		500	{object}	smaps.ErrorRes	"No se pudo regenerar el token"
//
//	@Failure		401	{object}	smaps.ErrorRes	"Token de recuperación inválido"
//
//	@Router			/api/v1/auth/refresh [post]
func RefreshToken(c *gin.Context) {
	var tokenRB *forms.RefreshTokenForm
	if err := c.ShouldBindJSON(&tokenRB); err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	claims, err := services.ValidateRefreshToken(tokenRB.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de recuperación inválido"})
		c.Abort()
		return
	}
	user, err := services.GetUserInfoFromCodPersonaService(uint(claims.RefreshCodPersona))
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusServiceUnavailable,
			gin.H{"error": "Error en el servidor"},
		)
		return
	}

	token, refreshToken, err1, err2 := authService.LoadJWTAuth(user.NombreCompleto, user.MailInstitucional, user.ID)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo regenerar el token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token":        token,
		"refreshToken": refreshToken,
		"user":         user,
	})
}
