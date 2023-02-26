// This file contains the models for the authentication process

package models

import "time"

type Post struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"index"`
	Title     string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:null"`
	UpdatedAt time.Time `gorm:"default:null"`
	DeletedAt time.Time `gorm:"default:null"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
