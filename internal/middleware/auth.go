package middleware

import (
	"boilerplate/internal/config"
	"boilerplate/internal/repository"
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware(conf *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse the token without validation to extract session_id
		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		sessionID, ok := claims["session_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: session_id not found"})
			c.Abort()
			return
		}

		// Initialize the AuthRepository
		authRepo := repository.NewAuthRepository(db)

		// Get the session using the AuthRepository
		session, err := authRepo.GetUserTokenBySessionId(sessionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session or session expired"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			c.Abort()
			return
		}

		// Construct the secret using conf.JWTSecret and session.UserToken
		secret := conf.JWTSecret + session.UserToken

		// Validate the token with the reconstructed secret
		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract userID from claims
		idFloat, ok := claims["id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: userID not found"})
			c.Abort()
			return
		}

		userID := int(idFloat)

		// Initialize the UserRepository
		userRepo := repository.NewUserRepository(db)

		// Get the user data using the UserRepository
		user, err := userRepo.FindByID(userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			c.Abort()
			return
		}

		// Set user and session structs in context
		c.Set("user", user)
		c.Set("session", session)

		c.Next()
	}
}
