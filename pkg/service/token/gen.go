package token

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// func GenerateToken(id int, exp time.Duration, role string) (string, error) {
func GenerateAccessToken(userID int) (string, error) {
	// Create a new JWT claims instance
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    "user",
		"exp":     time.Now().Add(time.Hour * 15).Unix(), // Token expiration time (24 hours)
		"iat":     time.Now().Unix(),                     // Token issuance time
	}

	// Create the JWT token with the claims and sign it with the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("access-token-src"))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
