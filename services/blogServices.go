package services

// Filename: blogServices.go
// This file handle all the services for the blog (creating, updating and deleting posts)
import (
	"fmt"
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
func CreatePost(title string, content string, userID string) (models.Post, error) {
	post := models.Post{
		ID:        cuid.New(),
		Title:     title,
		Content:   content,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// Handle errors gracefully
	if post.Title == "" {
		return post, fmt.Errorf("Title is required")
	} else if post.Content == "" {
		return post, fmt.Errorf("Content is required")
	} else if post.UserID != userID || post.UserID == "" {
		return post, fmt.Errorf("Error creating post, verification failed")
	}

	config.DB.Create(&post)

	return post, nil
}

// UpdatePost updates a post in the database from the given data (title, content and userID)
func UpdatePost(id string, title string, content string, userID string) (models.Post, error) {
	var post models.Post
	config.DB.Where("id = ? AND user_id = ?", id, userID).First(&post)

	if post.UserID != userID || post.UserID == "" {
		return post, fmt.Errorf("Error updating post, verification failed")
	}

	post.Title = title
	post.Content = content
	post.UserID = userID
	post.UpdatedAt = time.Now()
	config.DB.Save(&post)
	return post, nil
}

// DeletePost deletes a post in the database from the given postID and userID
func DeletePost(id string, userID string) error {
	var post models.Post
	config.DB.Where("id = ? AND user_id = ?", id, userID).First(&post)

	if post.UserID != userID || post.UserID == "" {
		return fmt.Errorf("Error deleting post, verification failed")
	}

	config.DB.Delete(&post)
	// Return nil if there is no error
	return nil
}
