package services

import (
	"crypto/sha256"
	"errors"
	"log"
	"net/http"
	"task_scheduler/internal/config"

	"github.com/golang-jwt/jwt"
)

type AuthService struct{}

func GetAuthService() *AuthService {
	return &AuthService{}
}

// CheckJWTMiddleware проверяет токен, полученный из тела запроса, на валидность.
// Проверка происходит с помощью secret key и сравнения checksum паролей.
func CheckJWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			jwtVal   string
			jwtToken *jwt.Token
		)

		if len(config.Password) == 0 {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		jwtVal = cookie.Value
		jwtToken, err = jwt.Parse(jwtVal, func(t *jwt.Token) (interface{}, error) {
			secret := []byte("secret_key")
			return secret, nil
		})
		if err != nil {
			log.Println("Failed to parse jwt from cookie")
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		payLoad, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Failed to typecast to jwt.MapClaims")
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		checkSum, ok := payLoad["sum"].([]interface{})
		if !ok {
			log.Println("Failed to typecast to checksum")
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		sum := sha256.Sum256([]byte(config.Password))

		if len(checkSum) != len(sum) {
			log.Println("Checksum is not valid")
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		for i, v := range sum {
			if elemSum, ok := checkSum[i].(byte); ok && elemSum != v {
				log.Println("Checksum is not valid")
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}

		if !jwtToken.Valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// GetJWT генерирует JWT токен, используя PASSWORD из переменных окружения и
// checksum пароля, вложенного в Claims токена.
func (a *AuthService) GetJWT(password string) (string, error) {
	var (
		signedToken string
		jwtToken    *jwt.Token
		err         error
	)

	if config.Password != password {
		return "", errors.New("wrong password")
	}

	secret := []byte("secret_key")

	jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sum": sha256.Sum256([]byte(password))})

	signedToken, err = jwtToken.SignedString(secret)

	return signedToken, err
}
