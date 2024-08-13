package services

import (
	"time"

	"github.com/citiaps/SID-UMAG-AuthService/settings"
)

// Settings
var settingsData = settings.NewSettings()

// JWT Config
var secretKey = []byte(settingsData.JWT_KEY)
var tokenDuration = 1 * time.Hour
var refreshTokenDuration = 12 * time.Hour
