package models

import "time"

type Account struct {
	ID                string `gorm:"primaryKey"`
	UserID            string
	Type              string
	Provider          string
	ProviderAccountID string
	RefreshToken      string
	AccessToken       string
	ExpiresAt         time.Time
	TokenType         string
	Scope             string
	IDToken           string
	SessionState      string
	User              User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Session struct {
	ID           string `gorm:"primaryKey"`
	SessionToken string `gorm:"unique"`
	UserID       string
	Expires      time.Time
	User         User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type User struct {
	ID            string `gorm:"primaryKey"`
	Name          string `gorm:"unique"`
	Email         string `gorm:"unique"`
	EmailVerified time.Time
	Image         string
	Accounts      []Account `gorm:"foreignKey:UserID"`
	Sessions      []Session `gorm:"foreignKey:UserID"`
	IsAdmin       bool      `gorm:"default:false"`
	IsOnline      bool      `gorm:"default:false"`
}
