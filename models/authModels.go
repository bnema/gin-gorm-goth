// File: authModels.go
// This file contains the models for the authentication process

package models

import "time"

type User struct {
	ID            string    `gorm:"primaryKey"`
	Name          string    `gorm:"unique"`
	Email         string    `gorm:"unique"`
	EmailVerified time.Time `gorm:"default:null"`
	Image         string    `gorm:"default:null"`
	IsAdmin       bool      `gorm:"default:false"`
	IsOnline      bool      `gorm:"default:false"`
	Accounts      []Account `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Sessions      []Session `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Account struct {
	ID                string    `gorm:"primaryKey"`
	UserID            string    `gorm:"index"`
	Type              string    `gorm:"not null"`
	Provider          string    `gorm:"not null"`
	ProviderAccountID string    `gorm:"not null"`
	RefreshToken      string    `gorm:"default:null"`
	AccessToken       string    `gorm:"default:null"`
	ExpiresAt         time.Time `gorm:"default:null"`
	TokenType         string    `gorm:"default:null"`
	Scope             string    `gorm:"default:null"`
	IDToken           string    `gorm:"default:null"`
	SessionState      string    `gorm:"default:null"`
	User              User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Session struct {
	ID           string    `gorm:"primaryKey"`
	SessionToken string    `gorm:"unique"`
	UserID       string    `gorm:"index"`
	Expires      time.Time `gorm:"default:null"`
	User         User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	IsOnline     bool      `gorm:"default:false"`
}
