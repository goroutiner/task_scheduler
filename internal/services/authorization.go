package services

import (
	"crypto/sha256"
	"errors"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

// CheckJWT проверяет токен, полученный из тела запроса, на валидность.
// Проверка происходит с помощью secret key и сравнения checksum паролей.
func CheckJWT(w http.ResponseWriter, r *http.Request) bool {
	var (
		jwtVal   string
		jwtToken *jwt.Token
	)

	envMap, _ := godotenv.Read("../.env")
	password := envMap["TODO_PASSWORD"]

	if len(password) > 0 {
		cookie, err := r.Cookie("token")
		if err != nil {
			log.Println(err.Error())
			return false
		}

		jwtVal = cookie.Value
		jwtToken, _ = jwt.Parse(jwtVal, func(t *jwt.Token) (interface{}, error) {
			secret := []byte("secret_key")
			return secret, nil
		})
	} else {
		log.Println("Environment variable TODO_PASSWORD is not found")
		return false
	}

	if payLoad, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		checksumRow := payLoad["sum"]
		sum := sha256.Sum256([]byte(password))
		if checksum, ok := checksumRow.([]interface{}); ok {
			if len(checksum) == len(sum) {
				for i, v := range sum {
					if elemSum, ok := checksum[i].(byte); ok && elemSum != v {
						log.Println("Checksum is not valid")
						return false
					}
				}
			}
		} else {
			log.Println("Failed to typecast to checksum")
			return false
		}
	} else {
		log.Println("Failed to typecast to jwt.MapClaims")
		return false
	}

	return jwtToken.Valid
}

// GetJWT генерирует JWT токен, используя TODO_PASSWORD из переменных окружения и
// checksum пароля, вложенного в Claims токена.
func GetJWT(takenMap map[string]string) (string, error) {
	var (
		signedToken string
		jwtToken    *jwt.Token
	)

	envMap, err := godotenv.Read("../.env")
	if err != nil {
		return "", err
	}

	password := takenMap["password"]
	if envMap["TODO_PASSWORD"] == password {
		secret := []byte("secret_key")

		jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sum": sha256.Sum256([]byte(password))})

		signedToken, err = jwtToken.SignedString(secret)
		if err != nil {
			log.Println(err.Error())
			return "", err
		}
	} else {
		return "", errors.New("wrong password")
	}

	return signedToken, nil
}
