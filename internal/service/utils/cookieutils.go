package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWT creates a new JWT for a given user ID.
// It sets an expiration time based on a predefined duration and encodes the user's unique identifier in the claims.
//
// Parameters:
//
//	userID: the user's unique identifier to be embedded in the JWT.
//	jwtKeySecret: the secret key used for signing the JWT.
//
// Returns:
//
//	A string containing the signed JWT or an error if the JWT could not be generated.
func GenerateJWT(userID string, jwtKeySecret string) (string, error) {
	expirationTime := time.Now().Add(config.TokenExpirationInHour * time.Hour)
	claims := &models.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKeySecret))

	return tokenString, err
}

// VerifyJWT checks the validity of a JWT string using the specified secret key.
// It ensures that the token is valid, correctly signed, and not expired.
//
// Parameters:
//
//	tokenString: the JWT string to verify.
//	jwtKeySecret: the secret key used for signing the JWT.
//
// Returns:
//
//	The decoded claims if the JWT is valid or an error if there is a problem with the JWT.
func VerifyJWT(tokenString string, jwtKeySecret string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtKeySecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, config.ErrTokenInvalid
	}

	return claims, nil
}

// SetJWTInCookie sets a JWT in an HTTP response cookie after generating it for the given user ID.
// If the JWT cannot be generated, it sends an HTTP 500 Internal Server Error.
//
// Parameters:
//
//	w: the HTTP response writer to use for setting the cookie.
//	userID: the user's unique identifier for whom the JWT is generated.
func SetJWTInCookie(w http.ResponseWriter, userID string) {
	tokenString, err := GenerateJWT(userID, config.JwtKeySecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: tokenString,
	})
}

// GetUserIDFromCookie retrieves the user ID from a JWT stored in a cookie.
// If the cookie cannot be found or the JWT is invalid, it returns an error.
//
// Parameters:
//
//	r: the HTTP request from which to retrieve the cookie.
//
// Returns:
//
//	The user ID extracted from the JWT or an error if the cookie is missing or the JWT is invalid.
func GetUserIDFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}
	claims, err := VerifyJWT(cookie.Value, config.JwtKeySecret)
	if err != nil {
		return "", err
	}
	logger.Info("Successfully authorized: ", claims.UserID)
	return claims.UserID, nil
}
