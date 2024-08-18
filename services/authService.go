package services

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/acurilem/SID-UMAG-AuthService/forms"
	"github.com/acurilem/SID-UMAG-AuthService/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	ldap2 "gopkg.in/ldap.v2"
)

type authService struct{}

type Sid3Claims struct {
	CodPersona uint   `json:"cod_persona"`
	Email      string `json:"email"`
	jwt.RegisteredClaims
}

type Sid3RefreshClaims struct {
	RefreshCodPersona uint `json:"refresh_cod_persona"`
	jwt.RegisteredClaims
}

func (*authService) IsRutOrPassport(text string) bool {
	regex := regexp.MustCompile(`[0-9]+`) // Si tiene uno o mas número o es Rut o pasaporte
	return regex.MatchString(text)
}

func (authService *authService) LoginLdapv2WithRut(login *forms.LoginForm) (string, error) {
	// Connect to LDAP
	ldapServer := settingsData.LDAP_SERVER
	l, err := ldap2.Dial("tcp", ldapServer)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	bindUsername := "cn=" + settingsData.LDAP_USER + ",dc=umag,dc=cl"
	err = l.Bind(bindUsername, settingsData.LDAP_PASS)
	if err != nil {
		log.Println("No fue posible encontrar al usuario en LDAP")
		return "", err
	}

	userDN := "ou=usuarios,dc=umag,dc=cl"
	attributes := []string{"uid", "gecos", "employeetype", "employeenumber"}
	searchRequest := ldap2.NewSearchRequest(
		userDN,
		ldap2.ScopeWholeSubtree, ldap2.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(employeenumber=%s)", login.Username), // Filtro basado en 'employeenumber'
		attributes,
		nil,
	)
	sr, err := l.Search(searchRequest)
	var rut string
	if err != nil {
		log.Println("No fue posible encontrar al usuario en LDAP (3)")
		return "", err
	}
	if len(sr.Entries) == 0 {
		fmt.Println("No se encontraron resultados.")
	} else {
		entry := sr.Entries[0]
		username := entry.GetAttributeValue("uid")
		login.Username = username
		//fmt.Println(login.Username)
		rut, err = authService.LoginLdapv2(login)
	}
	return rut, err
}

func (*authService) LoginLdapv2(login *forms.LoginForm) (string, error) {
	ldapServer := settingsData.LDAP_SERVER
	l, err := ldap2.Dial("tcp", ldapServer)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	bindUsername := "uid=" + login.Username + ",ou=usuarios,dc=umag,dc=cl"
	err = l.Bind(bindUsername, login.Password)
	if err != nil {
		log.Println("No fue posible encontrar al usuario (1)")
		return "", err
		//log.Fatal(err)
	}
	userDN := "uid=" + login.Username + ",ou=usuarios,dc=umag,dc=cl"
	attributes := []string{"gecos", "employeetype", "employeenumber"}
	searchRequest := ldap2.NewSearchRequest(
		userDN,
		ldap2.ScopeBaseObject,
		ldap2.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)",
		attributes,
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Println("No fue posible encontrar al usuario (2)")
		return "", err
	}
	entry := sr.Entries[0]
	rut := entry.GetAttributeValue("employeeNumber")
	return rut, err
}

func (*authService) LoadJWTAuth(name string, email string, codPersona uint) (string, string, error, error) {

	claims := Sid3Claims{
		codPersona,
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    name,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err1 := token.SignedString(secretKey)

	refreshClaims := Sid3RefreshClaims{
		codPersona,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenDuration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err2 := refreshToken.SignedString(secretKey)
	return tokenString, refreshTokenString, err1, err2
}

func GetToken(c *gin.Context) string {
	tokenString := c.GetHeader("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sin autenticación"})
		c.Abort()
		return ""
	}
	return tokenString
}

func secretKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(secretKey), nil
}

func ValidateJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, secretKeyFunc)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Autentificación no valida")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("Error al obtener los detalles de la auntenticación")
	}
}

func ValidateRefreshToken(tokenString string) (Sid3RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Sid3RefreshClaims{}, secretKeyFunc)
	if err != nil {
		fmt.Println(err)
		return Sid3RefreshClaims{}, err
	}

	if !token.Valid {
		return Sid3RefreshClaims{}, fmt.Errorf("Autenticación no valida (2)")
	}

	if claims, ok := token.Claims.(*Sid3RefreshClaims); ok {
		return *claims, nil
	} else {
		return Sid3RefreshClaims{}, fmt.Errorf("Error al obtener los detalles de la auntenticación(2)")
	}
}

func (*authService) GetUser(c *gin.Context) (models.User, error) {
	token := GetToken(c)
	claims, err := ValidateJWTToken(token)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{"error": "No se encontró información del usuario"})
		return models.User{}, err
	}
	user, err := GetUserInfoFromCodPersonaService(uint(claims["cod_persona"].(float64)))
	return user, err
}

func NewAuthService() *authService {
	return &authService{}
}
