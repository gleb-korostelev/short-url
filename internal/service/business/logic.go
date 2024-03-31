package business

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateShortPath() string {
	b := make([]byte, config.Length)
	for i := range b {
		b[i] = config.Letters[rand.Intn(len(config.Letters))]
	}
	return string(b)
}

func SaveURLs(save models.URLData) error {
	file, err := os.OpenFile(config.BaseFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	data, err := json.Marshal(save)
	if err != nil {
		return err
	}
	_, err = writer.WriteString(string(data) + "\n")
	if err != nil {
		return err
	}
	return writer.Flush()
}

func LoadURLs(path string, shortURL string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var urlData models.URLData
		if err := json.Unmarshal([]byte(scanner.Text()), &urlData); err != nil {
			return "", err
		}
		if urlData.ShortURL == shortURL && !urlData.DeletedFlag {
			return urlData.OriginalURL, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", config.ErrNotFound
}

func LoadUserURLs(path string, userID string) ([]models.AllUserURL, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var urls []models.AllUserURL
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var urlData models.URLData
		var data models.AllUserURL
		if err := json.Unmarshal([]byte(scanner.Text()), &urlData); err != nil {
			return nil, err
		}
		if urlData.UUID.String() == userID && !urlData.DeletedFlag {
			data.OriginalURL = urlData.OriginalURL
			data.ShortURL = config.BaseURL + "/" + urlData.ShortURL
			urls = append(urls, data)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

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

func SetJWTInCookie(w http.ResponseWriter, userID string) {
	tokenString, err := GenerateJWT(userID, config.JwtKeySecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(config.TokenExpirationInHour * time.Hour),
	})
}

func GetUserIDFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", err
		}
		return "", err
	}
	claims, err := VerifyJWT(cookie.Value, config.JwtKeySecret)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

func MarkURLsAsDeletedInFile(path, userID string, shortURLs []string) error {
	file, err := os.OpenFile(config.BaseFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(file)

	for scanner.Scan() {
		var urlData models.URLData
		if err := json.Unmarshal([]byte(scanner.Text()), &urlData); err != nil {
			return err
		}
		if urlData.UUID.String() == userID && CheckURL(urlData.ShortURL, shortURLs) {
			urlData.DeletedFlag = true
		}
		data, err := json.Marshal(urlData)
		if err != nil {
			logger.Errorf("error marshalling json: %w", err)
			return err
		}
		_, err = writer.WriteString(string(data) + "\n")
		if err != nil {
			logger.Errorf("error writing file: %w", err)
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return writer.Flush()
}

func CheckURL(check string, findlist []string) bool {
	for _, find := range findlist {
		if find == check {
			return true
		}
	}
	return false
}
