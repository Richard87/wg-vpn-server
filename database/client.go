package database

import (
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Name       string `json:"name"`
	AllowedIp4 string `json:"allowedIp4"`
	PublicKey  string `json:"publicKey"`
}
