// This file contains the models for the authentication process

package models

import "time"

type Post struct {
	ID string `gorm:"primaryKey"`
	// Foreign key (UserID -> users.id)
	UserID    string    `gorm:"index" binding:"required"`
	Title     string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:null"`
	UpdatedAt time.Time `gorm:"default:null"`
	DeletedAt time.Time `gorm:"default:null"`
}
