package services

// Filename: authServices.go
// This file handle all the authentication process with the database
import (
	"fmt"
	"go-gorm-gauth/config"
	"go-gorm-gauth/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
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

func CreateNewSession(authUser goth.User) (*models.Session, error) {
	// Check the email, retrieve user.id from database and push it to the account as UserID
	userFromDB, err := GetUserByEmail(authUser.Email)
	if err != nil {
		fmt.Println("Error getting user by email:", err)
		return nil, err
	}
	// Call the generateJWT function to create a new JWT token for the new session
	token, err := generateJWT(userFromDB)
	if err != nil {
		fmt.Println("Error generating JWT token:", err)
		return nil, err
	}

	session := &models.Session{
		ID:           cuid.New(),
		UserID:       userFromDB.ID,
		SessionToken: token,
		Expires:      time.Now().Add(time.Hour * 24 * 7),
		IsOnline:     true,
	}

	err = config.DB.Create(session).Error
	if err != nil {
		fmt.Println("Error creating new session:", err)
		return nil, err
	}

	fmt.Println("New session created:", session.ID)
	return session, nil
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

func generateJWT(user *models.User) (string, error) {
	// Create the JWT claims, which includes the user ID, email, and expiration time
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		// 7days expiration
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	// Create the JWT token using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	return token.SignedString([]byte("secret"))
}
