package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prodanov17/znk/internal/config"
	"github.com/prodanov17/znk/internal/utils"
)

type contextKey string

const UserKey contextKey = "userID"

var secretKey = []byte(config.Env.JWTSecret)

func WithJWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetTokenFromRequest(r)

		token, err := VerifyToken(tokenString)

		if err != nil {
			permissionDenied(w, r)
			return
		}

		userID, err := GetUserIDFromToken(token)
		if err != nil {
			log.Printf("Failed to convert userID to int: %v", err)
			permissionDenied(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, userID)
		r = r.WithContext(ctx)

		next(w, r)
	})
}

func CreateToken(userId int) (string, error) {
	expiration := time.Second * time.Duration(config.Env.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userID": userId,
			"exp":    time.Now().Add(expiration).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	return token, nil
}

func GetUserIDFromToken(token *jwt.Token) (int, error) {
	claims := token.Claims.(jwt.MapClaims)

	userIDFloat, ok := claims["userID"].(float64)
	if !ok {
		return 0, fmt.Errorf("failed to assert userID as float64")
	}

	userID := int(userIDFloat)

	if userID == 0 {
		return 0, fmt.Errorf("no userID found in token")
	}

	return userID, nil
}

func permissionDenied(w http.ResponseWriter, r *http.Request) {
	utils.WriteError(w, r, http.StatusForbidden, errors.New("Unauthorized"))
}
