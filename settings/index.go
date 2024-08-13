package settings

import (
	"fmt"
	"os"
)

type settings struct {
	JWT_KEY     string
	LDAP_SERVER string
	LDAP_USER   string
	LDAP_PASS   string
	DB_HOST     string
	DB_PORT     string
	DB_INSTANCE string
	DB_USER     string
	DB_PASS     string
	DB_DB       string
}

func (settings *settings) validate() {
	var errors []string

	if settings.JWT_KEY == "" {
		errors = append(errors, "JWT_KEY Missing")
	}
	if settings.LDAP_SERVER == "" {
		errors = append(errors, "LDAP_SERVER Missing")
	}
	if settings.LDAP_USER == "" {
		errors = append(errors, "LDAP_ADMIN_USERNAME Missing")
	}
	if settings.LDAP_PASS == "" {
		errors = append(errors, "LDAP_ADMIN_PASS Missing")
	}
	if settings.DB_HOST == "" {
		errors = append(errors, "DB_HOST Missing")
	}
	if settings.DB_PORT == "" {
		errors = append(errors, "DB_PORT Missing")
	}
	if settings.DB_INSTANCE == "" {
		errors = append(errors, "DB_INSTANCE Missing")
	}
	if settings.DB_USER == "" {
		errors = append(errors, "DB_USER Missing")
	}
	if settings.DB_PASS == "" {
		errors = append(errors, "DB_PASS Missing")
	}
	if settings.DB_DB == "" {
		errors = append(errors, "DB_DB Missing")
	}

	if errors != nil {
		fmt.Printf("errors: %v\n", errors)
		panic("venv validation failed")
	}
}

var settingsData *settings

func NewSettings() *settings {
	if settingsData == nil {
		// Load and validate
		settingsData = &settings{
			JWT_KEY:     os.Getenv("JWT_KEY"),
			LDAP_SERVER: os.Getenv("LDAP_SERVER"),
			LDAP_USER:   os.Getenv("LDAP_ADMIN_USERNAME"),
			LDAP_PASS:   os.Getenv("LDAP_ADMIN_PASS"),
			DB_HOST:     os.Getenv("DB_HOST"),
			DB_PORT:     os.Getenv("DB_PORT"),
			DB_INSTANCE: os.Getenv("DB_INSTANCE"),
			DB_USER:     os.Getenv("DB_USER"),
			DB_PASS:     os.Getenv("DB_PASS"),
			DB_DB:       os.Getenv("DB_DB"),
		}

		settingsData.validate()
	}

	return settingsData
}
