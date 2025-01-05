package services

import (
	"crypto/sha256"
	"errors"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
)

// CheckJWT проверяет токен, полученный из тела запроса, на валидность.
// Проверка происходит с помощью secret key и сравнения checksum паролей.
func CheckJWT(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			jwtVal   string
			jwtToken *jwt.Token
		)

		password := EnvMap["TODO_PASSWORD"]

		if len(password) == 0 {
			log.Println("Environment variable TODO_PASSWORD is not found")
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
		jwtToken, _ = jwt.Parse(jwtVal, func(t *jwt.Token) (interface{}, error) {
			secret := []byte("secret_key")
			return secret, nil
		})

		payLoad, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Failed to typecast to jwt.MapClaims")
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return

		}

		checksumRow := payLoad["sum"]
		sum := sha256.Sum256([]byte(password))
		checksum, ok := checksumRow.([]interface{})
		if !ok {
			log.Println("Failed to typecast to checksum")
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		if len(checksum) == len(sum) {
			for i, v := range sum {
				if elemSum, ok := checksum[i].(byte); ok && elemSum != v {
					log.Println("Checksum is not valid")
					http.Error(w, "Authentication required", http.StatusUnauthorized)
					return
				}
			}
		}

		if !jwtToken.Valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}

// GetJWT генерирует JWT токен, используя TODO_PASSWORD из переменных окружения и
// checksum пароля, вложенного в Claims токена.
func GetJWT(takenMap map[string]string) (string, error) {
	var (
		signedToken string
		jwtToken    *jwt.Token
		err         error
	)

	password := takenMap["password"]
	if EnvMap["TODO_PASSWORD"] != password {
		return "", errors.New("wrong password")
	}
	secret := []byte("secret_key")

	jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sum": sha256.Sum256([]byte(password))})

	signedToken, err = jwtToken.SignedString(secret)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return signedToken, nil
}
