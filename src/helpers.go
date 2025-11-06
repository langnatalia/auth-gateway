package authgateway

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	secretKey = "your-secret-key"
)

type AuthenticationError struct {
	Code    string
	Message string
}

func NewAuthenticationError(code, message string) *AuthenticationError {
	return &AuthenticationError{
		Code:    code,
		Message: message,
	}
}

func GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 72).Unix(),
			"iss": "auth-gateway",
			"user_id": userID,
		},
	)
	return token.SignedString([]byte(secretKey))
}

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return "", NewAuthenticationError("invalid_token", "invalid token")
	}
	return claims["user_id"].(string), nil
}

func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func SendError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(message))
}

func SplitQueryString(s string) map[string]string {
	pairs := strings.Split(s, "&")
	result := make(map[string]string)
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}
	return result
}

func LogError(err error) {
	log.Printf("error: %v\n", err)
}