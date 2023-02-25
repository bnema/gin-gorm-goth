package services

// Filename: authServices.go
// This file handle all the authentication process with the database
import (
	"fmt"
	"go-gorm-gauth/config"
	"go-gorm-gauth/models"

	"github.com/lucsky/cuid"
	"github.com/markbates/goth"
	"gorm.io/gorm"
)

// This file handle all the authentication process
// CreateUser, CreateAccount, CreateSession (for now)

// Get the informations from goth.User and create a new user
// Create a new user with the given auth user data
func CreateNewUser(authUser goth.User) (*models.User, error) {
	// check if we received a valid user
	fmt.Println(authUser)
	dbConn := config.DB
	newID := cuid.New()
	newUser := models.User{
		ID:    newID,
		Name:  authUser.Name,
		Email: authUser.Email,
		Image: authUser.AvatarURL,
	}
	if err := dbConn.Create(newUser).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	return &newUser, nil
}
