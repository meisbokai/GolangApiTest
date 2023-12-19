package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
)

func CreateJWTTokenString(user *User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"username":  user.Username,
	}

	secret := os.Getenv("JWTTokenSecret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))

}

func validateJWTTokenString(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWTTokenSecret")

	// Example from https://pkg.go.dev/github.com/golang-jwt/jwt#example-Parse-Hmac
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

}

// Decorate the HandlerFunc with JWT Authentication
func WithJWTAuth(handlerFunction http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("JWT middleware")

		// Get token string from the request header
		tokenString := r.Header.Get("jwt-token")

		// Validate token string
		token, err := validateJWTTokenString(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}

		// Get user details from database
		userId, err := getID(r)
		if err != nil {
			permissionDenied(w)
			return
		}
		user, err := s.GetUserByID(userId)
		if err != nil {
			permissionDenied(w)
			return
		}

		// Compare user against claims
		claims := token.Claims.(jwt.MapClaims)
		if user.Username != claims["username"] {
			permissionDenied(w)
			return
		}

		handlerFunction(w, r)
	}
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "Access denied"})
}

// Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzQXQiOjE1MDAwLCJ1c2VybmFtZSI6InVzZXJXaXRoVG9rZW4ifQ.ygdLMUekKwOStdEnCJKXJkQZSiAA9nPIJMmayeXzY-A
