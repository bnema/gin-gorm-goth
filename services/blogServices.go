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
	if post.Title == "" || post.Content == "" {
		return post, fmt.Errorf("Title and content are required")
	} else if post.UserID != userID {
		return post, fmt.Errorf("You are not the owner of this post")
	} else if post.UserID == "" {
		return post, fmt.Errorf("You are not logged in")
	}
	config.DB.Create(&post)

	return post, nil
}

// UpdatePost updates a post in the database from the given data (title, content and userID)
// Also handle errors if the post doesn't exist
func UpdatePost(id string, title string, content string, userID string) (models.Post, error) {
	var post models.Post
	config.DB.Where("id = ? AND user_id = ?", id, userID).First(&post)

	if post.ID != id {
		return post, fmt.Errorf("Post not found")
	} else if post.UserID != userID {
		return post, fmt.Errorf("You are not the owner of this post")
	}

	post.Title = title
	post.Content = content
	post.UserID = userID
	post.UpdatedAt = time.Now()
	config.DB.Save(&post)
	return post, nil
}

// DeletePost deletes a post in the database from the given ID
func DeletePost(id string, userID string) error {
	var post models.Post
	config.DB.Where("id = ? AND user_id = ?", id, userID).First(&post)

	if post.ID != id {
		return fmt.Errorf("Post not found")
	} else if post.UserID != userID {
		return fmt.Errorf("You are not the owner of this post")
	}

	config.DB.Delete(&post)
	// Return nil if there is no error
	return nil
}
