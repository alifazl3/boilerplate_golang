package handler

import (
	"boilerplate/internal/config"
	"boilerplate/internal/model"
	"boilerplate/internal/repository"
	user_repository "boilerplate/internal/repository/user"
	"boilerplate/internal/utilities"
	"boilerplate/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	UserName         string `json:"username"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phone_number"`
	Password         string `json:"password" binding:"required"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Country          string `json:"country"`
	ProfileIMG       string `json:"profile_img"`
	ValidEmail       string `json:"valid_email"`
	ValidPhoneNumber string `json:"valid_phone_number"`
	TokenFacebook    string `json:"token_facebook"`
	TokenGoogle      string `json:"token_google"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password" binding:"required"`
}

type LogoutRequest struct {
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRepo := repository.NewUserRepository(c.MustGet("DB").(*gorm.DB))

	// Check if the user already exists
	if _, err := userRepo.FindByEmail(req.Email); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}
	passwordSalt := uuid.New().String()

	// Combine the salt and the password
	passwordWithSalt := passwordSalt + req.Password

	// Hash the combined password and salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordWithSalt), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := &model.User{
		UID:          "U_" + utilities.GenerateUID(),
		UserName:     req.UserName,
		Email:        &req.Email,
		Password:     string(hashedPassword),
		PasswordSalt: passwordSalt,
	}

	if err := userRepo.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	conf := c.MustGet("Config").(*config.Config)

	jwtString, session, err := auth.NewSession(conf.JWTSecret, user.UID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	err = userRepo.CreateSession(&session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create session"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "jwt": jwtString})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRepo := repository.NewUserRepository(c.MustGet("DB").(*gorm.DB))
	user, err := userRepo.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	passwordWithSalt := user.PasswordSalt + req.Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordWithSalt)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	conf := c.MustGet("Config").(*config.Config)
	jwtString, session, err := auth.NewSession(conf.JWTSecret, user.UID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	err = userRepo.CreateSession(&session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "jwt": jwtString})
}

func ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRepo := user_repository.NewUserRepository(c.MustGet("DB").(*gorm.DB))
	user, err := userRepo.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Validate old password
	passwordWithSalt := user.PasswordSalt + req.OldPassword
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordWithSalt)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid old password"})
		return
	}

	// Generate new password salt
	newPasswordSalt := uuid.New().String()
	passwordWithNewSalt := newPasswordSalt + req.NewPassword

	// Hash the new password with the new salt
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(passwordWithNewSalt), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	// Update user's password and salt
	user.Password = string(hashedNewPassword)
	user.PasswordSalt = newPasswordSalt

	if err := userRepo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// Deactivate all user sessions
	if err := userRepo.DeactivateUserSessions(user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate user sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func Logout(c *gin.Context) {
	sessionInterface, exists := c.Get("session")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Session not found in context"})
		return
	}

	session, ok := sessionInterface.(*model.Session)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid session type"})
		return
	}

	userRepo := repository.NewUserRepository(c.MustGet("DB").(*gorm.DB))
	if err := userRepo.DeactivateSession(session.SessionId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log out"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
