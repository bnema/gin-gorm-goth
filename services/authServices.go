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

func CreateNewUser(authUser goth.User) (*models.User, error) {
	user := &models.User{
		ID:    cuid.New(),
		Name:  authUser.Name,
		Email: authUser.Email,
		Image: authUser.AvatarURL,
	}

	err := config.DB.Create(user).Error
	if err != nil {
		fmt.Println("Error creating new user:", err)
		return nil, err
	}

	fmt.Println("New user created:", user.ID)
	return user, nil
}

func CreateNewAccount(authUser goth.User) (*models.Account, error) {

	// Check the email, retrieve user.id from database and push it to the account as UserID
	userFromDB, err := GetUserByEmail(authUser.Email)
	if err != nil {
		fmt.Println("Error getting user by email:", err)
		return nil, err
	}

	account := &models.Account{
		ID:                cuid.New(),
		UserID:            userFromDB.ID,
		Provider:          authUser.Provider,
		ProviderAccountID: authUser.UserID,
		Type:              "oauth",
		Scope:             "identify email",
		TokenType:         "Bearer",
		RefreshToken:      authUser.RefreshToken,
		AccessToken:       authUser.AccessToken,
		ExpiresAt:         authUser.ExpiresAt,
	}

	// Create the account in the database
	err = config.DB.Create(account).Error
	if err != nil {
		fmt.Println("Error creating new account:", err)
		return nil, err
	}

	fmt.Println("New account created:", account.ID)
	return account, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := config.DB.Where("email = ?", email).First(user).Error
	if err != nil {
		fmt.Println("Error getting user by email:", err)
		// If no user is found with the specified email address, return an error
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("User not found with email %s", email)
		}
		return nil, err
	}

	fmt.Println("User found:", user.ID)
	return user, nil
}
