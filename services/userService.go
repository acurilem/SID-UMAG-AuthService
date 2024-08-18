package services

import (
	"github.com/acurilem/SID-UMAG-AuthService/config"
	"github.com/acurilem/SID-UMAG-AuthService/models"
)

// Función que va a buscar la info de un usuario desde vista_persona
func GetUserInfoFromCodPersonaService(codPersona uint) (models.User, error) {
	db := config.Database
	var user models.User
	var categories []models.Role

	//Este codigo funciona con SQL Server 2012 en adelante
	//db.Model(&user).Where("Rut = ?", rut).First(&user)

	db.Raw("SELECT TOP 1 * FROM vista_personas WHERE cod_persona = ? ", codPersona).Preload("roles").Scan(&user)

	if user.Rut != "" {
		db.Raw(`
			SELECT c.* 
			FROM detalle_categoria dc
			LEFT JOIN categorias c ON dc.cod_categoria = c.cod_categoria
			WHERE dc.cod_persona = ?
			`, user.ID).Scan(&categories)
		user.Roles = categories
	}

	/*
		db.Where("Rut = ? OR Rut = ? OR Rut = ?", rut, rut[1:], "0"+rut).Order("cod_persona").First(&user)
		db.Model(&user).Association("Roles").Find(&user.Roles)*/

	return user, nil
}

func GetUserInfoFromRutService(rut string) (models.User, error) {
	db := config.Database
	var user models.User
	var categories []models.Role

	db.Where("Rut = ?", rut).Or("Rut = ?", rut[1:]).Or("Rut = ?", "0"+rut).Find(&user)
	//fmt.Println(user)
	//db.Model(&user).Where("Rut = ?", rut).First(&user)

	//db.Raw("SELECT TOP 1 * FROM vista_personas WHERE Rut = ?  or Rut = ? or Rut = ? ORDER BY cod_persona", rut, rut[1:], "0"+rut).Preload("roles").Scan(&user)

	if user.Rut != "" {
		db.Raw(`
			SELECT c.* 
			FROM detalle_categoria dc
			LEFT JOIN categorias c ON dc.cod_categoria = c.cod_categoria
			WHERE dc.cod_persona = ?
			`, user.ID).Scan(&categories)
		user.Roles = categories
	}

	/*
		db.Where("Rut = ? OR Rut = ? OR Rut = ?", rut, rut[1:], "0"+rut).Order("cod_persona").First(&user)
		db.Model(&user).Association("Roles").Find(&user.Roles)*/

	return user, nil
}
