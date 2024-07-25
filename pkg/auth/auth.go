package auth

import (
	"boilerplate/internal/model"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Jwt struct {
	Id              int       `json:"id"`
	Uid             string    `json:"uid"`
	SessionId       string    `json:"session_id"`
	TokenExpiration time.Time `json:"exp"`
	jwt.StandardClaims
}

func (j Jwt) Valid() error {
	if j.TokenExpiration.Before(time.Now()) {
		return fmt.Errorf("token has expired")
	}
	return nil
}

// NewSession generates a new JWT and session.
func NewSession(jwtSecret string, Uid string, userID int) (string, model.Session, error) {
	sessionID, err := uuid.NewUUID()
	if err != nil {
		return "", model.Session{}, fmt.Errorf("failed to generate session ID: %v", err)
	}

	token, err := uuid.NewUUID()
	if err != nil {
		return "", model.Session{}, fmt.Errorf("failed to generate token: %v", err)
	}

	tokenExpiration := time.Now().AddDate(1, 0, 0)

	session := model.Session{
		UserId:          userID,
		UserToken:       token.String(),
		SessionId:       sessionID.String(),
		TokenExpiration: tokenExpiration,
	}

	jwtClaims := Jwt{
		Id:              userID,
		Uid:             Uid,
		SessionId:       sessionID.String(),
		TokenExpiration: tokenExpiration,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiration.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	secret := jwtSecret + token.String()
	tokenString, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", model.Session{}, fmt.Errorf("failed to sign JWT: %v", err)
	}

	return tokenString, session, nil
}
