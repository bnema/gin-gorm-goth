package services

// Filename: blogServices.go
// This file handle all the services for the blog (creating, updating and deleting posts)
import (
	"go-gorm-gauth/config"
	"go-gorm-gauth/models"
	"strings"
	"time"

	"github.com/lucsky/cuid"
)

// GetAllPosts returns all the posts in the database
func GetAllPosts() []models.Post {
	var posts []models.Post
	config.DB.Find(&posts)
	return posts
}

// GetPostByID returns a single post by ID
func GetPostByTitle(title string) models.Post {
	// Transform _ to space
	title = strings.ReplaceAll(title, "-", " ")
	var post models.Post
	config.DB.Where("title = ?", title).First(&post)
	return post
}

// CreatePost creates a new post in the database from the given data (title, content and userID)
func CreatePost(title string, content string, userID string) models.Post {
	post := models.Post{
		ID:        cuid.New(),
		Title:     title,
		Content:   content,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	config.DB.Create(&post)
	return post
}

// UpdatePost updates a post in the database from the given data (title, content and userID)
func UpdatePost(id string, title string, content string, userID string) models.Post {
	var post models.Post
	config.DB.Where("id = ?", id).First(&post)
	post.Title = title
	post.Content = content
	post.UserID = userID
	post.UpdatedAt = time.Now()
	config.DB.Save(&post)
	return post
}

// DeletePost deletes a post in the database from the given ID
func DeletePost(id string, userID string) {
	var post models.Post
	config.DB.Where("id = ? AND user_id = ?", id, userID).First(&post)
	config.DB.Delete(&post)
}
