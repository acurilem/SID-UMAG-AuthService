package config

import (
	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var Database = func() (db *gorm.DB) {
	dbHost := settingsData.DB_HOST
	dbPort := settingsData.DB_PORT
	dbInstance := settingsData.DB_INSTANCE
	dbUser := settingsData.DB_USER
	dbPassword := settingsData.DB_PASS
	dbDatabase := settingsData.DB_DB

	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?instance=%s&database=%s&encrypt=disable&connection+timeout=30",
		dbUser, dbPassword, dbHost, dbPort, dbInstance, dbDatabase,
	)
	if db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{}); err != nil {
		fmt.Println("Error en la Conexion", err)
		panic(err)
	} else {
		fmt.Println("Conexion sqlserver exitosa")
		return db
	}
}()
