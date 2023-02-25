package services

// Filename: authServices.go
// This file handle all the authentication process with the database
import (
	"fmt"
	"go-gorm-gauth/config"
	"go-gorm-gauth/models"
	"os"
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
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func verifyJWT(tokenString string) (*jwt.Token, error) {
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		fmt.Println("Error parsing JWT token:", err)
		return nil, err
	}

	return token, nil
}

func GetSessionByToken(tokenString string) (*models.Session, error) {
	// Verify the JWT token
	_, err := verifyJWT(tokenString)
	if err != nil {
		fmt.Println("Error verifying JWT token:", err)
		return nil, err
	}

	// Get the session from the database
	session := &models.Session{}
	err = config.DB.Where("session_token = ?", tokenString).First(session).Error
	if err != nil {
		fmt.Println("Error getting session by token:", err)
		// If no session is found with the specified token, return an error
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Session not found with token %s", tokenString)
		}
		return nil, err
	}

	fmt.Println("Session found:", session.ID)
	return session, nil
}

func DeleteSession(sessionToken string) error {
	// Verify the JWT token
	_, err := verifyJWT(sessionToken)
	if err != nil {
		fmt.Println("Error verifying JWT token:", err)
		return err
	} else {
		// Export the claims from the JWT token
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(sessionToken, claims, func(token *jwt.Token) (interface{}, error) {

			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			fmt.Println("Error parsing JWT token:", err)
			return err
		}

		// Get the user ID from the claims
		userID := claims["id"].(string)
		// Get the session ID from the database
		session := &models.Session{}

		// Get the user from the database
		user := &models.User{}
		err = config.DB.Where("id = ?", userID).First(user).Error
		if err != nil {
			fmt.Println("Error getting user by id:", err)
			// If no user is found with the specified ID, return an error
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("User not found with id %s", userID)
			}
			return err
		}

		// Delete the session from the database where the user ID matches the user ID in the claims and the session token matches the session token in the request
		err = config.DB.Where("user_id = ? AND session_token = ?", userID, sessionToken).Delete(session).Error
		if err != nil {
			fmt.Println("Error deleting session:", err)
			return err
		}

		fmt.Println("Session deleted")
		return nil
	}
}

func GetAccountByProviderAccountID(providerAccountID string) (*models.Account, error) {
	account := &models.Account{}
	err := config.DB.Where("provider_account_id = ?", providerAccountID).First(account).Error
	if err != nil {
		fmt.Println("Error getting account by provider account id:", err)
		// If no account is found with the specified provider account id, return an error
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Account not found with provider account id %s", providerAccountID)
		}
		return nil, err
	}

	fmt.Println("Account found:", account.ID)
	return account, nil
}

func UpdateSession(session *models.Session) (*models.Session, error) {
	err := config.DB.Save(session).Error
	if err != nil {
		fmt.Println("Error updating session:", err)
		return nil, err
	}

	fmt.Println("Session updated:", session.ID)
	return session, nil
}

func UpdateAccount(account *models.Account) (*models.Account, error) {
	err := config.DB.Save(account).Error
	if err != nil {
		fmt.Println("Error updating account:", err)
		return nil, err
	}

	fmt.Println("Account updated:", account.ID)
	return account, nil
}

func UpdateUser(user *models.User) (*models.User, error) {
	err := config.DB.Save(user).Error
	if err != nil {
		fmt.Println("Error updating user:", err)
		return nil, err
	}

	fmt.Println("User updated:", user.ID)
	return user, nil
}

func GetAccountByUserID(userID string) (*models.Account, error) {
	account := &models.Account{}
	err := config.DB.Where("user_id = ?", userID).First(account).Error
	if err != nil {
		fmt.Println("Error getting account by user id:", err)
		// If no account is found with the specified user id, return an error
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Account not found with user id %s", userID)
		}
		return nil, err
	}

	fmt.Println("Account found:", account.ID)
	return account, nil
}

func GetAccountByID(accountID string) (*models.Account, error) {
	account := &models.Account{}
	err := config.DB.Where("id = ?", accountID).First(account).Error
	if err != nil {
		fmt.Println("Error getting account by id:", err)
		// If no account is found with the specified id, return an error
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Account not found with id %s", accountID)
		}
		return nil, err
	}

	fmt.Println("Account found:", account.ID)
	return account, nil
}

func CheckIfAccountExists(providerAccountID string) bool {
	account := &models.Account{}
	err := config.DB.Where("provider_account_id = ?", providerAccountID).First(account).Error
	if err != nil {
		fmt.Println("Error getting account by provider account id:", err)
		// If no account is found with the specified provider account id, return false
		if err == gorm.ErrRecordNotFound {
			return false
		}
		return false
	}

	fmt.Println("Account found:", account.ID)
	return true
}

func CheckIfUserExists(email string) bool {
	user := &models.User{}
	err := config.DB.Where("email = ?", email).First(user).Error
	if err != nil {
		fmt.Println("Error getting user by email:", err)
		// If no user is found with the specified email, return false
		if err == gorm.ErrRecordNotFound {
			return false
		}
		return false
	}

	fmt.Println("User found:", user.ID)
	return true
}
