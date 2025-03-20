package services_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"task_scheduler/internal/config"
	"task_scheduler/internal/services"
	"testing"

	"github.com/stretchr/testify/require"
)

var signedToken string

// TestMain выполняется перед запуском всех тестов.
func TestMain(m *testing.M) {
	authService := services.GetAuthService()

	var err error
	config.Password = "valid_password"
	signedToken, err = authService.GetJWT(config.Password)
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	os.Exit(code)
}

// TestGetJWT тестирует метод GetJWT сервиса аутентификации.
func TestGetJWT(t *testing.T) {
	authService := services.GetAuthService()

	t.Run("valid password", func(t *testing.T) {
		testValidPassword := "valid_password"

		testSignedToken, err := authService.GetJWT(testValidPassword)

		require.NoError(t, err)
		require.NotEmpty(t, testSignedToken)
	})

	t.Run("invalid password", func(t *testing.T) {
		testInValidPassword := "invalid_password"

		signedToken, err := authService.GetJWT(testInValidPassword)

		require.Error(t, err)
		require.Empty(t, signedToken)
	})
}

// TestCheckJWTMiddleware тестирует middleware.
func TestCheckJWTMiddleware(t *testing.T) {
	// Создание обработчика для тестирования
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	handler := services.CheckJWTMiddleware(nextHandler)

	t.Run("valid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: signedToken})
		respRec := httptest.NewRecorder()

		handler.ServeHTTP(respRec, req)

		expectedResponse := "OK"
		actualResponse := respRec.Body.String()

		require.Equal(t, http.StatusOK, respRec.Code)
		require.Contains(t, actualResponse, expectedResponse)
	})

	t.Run("invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: "invalid_token"})
		respRec := httptest.NewRecorder()

		handler.ServeHTTP(respRec, req)

		expectedResponse := "Authentication required"
		actualResponse := respRec.Body.String()

		require.Equal(t, http.StatusUnauthorized, respRec.Code)
		require.Contains(t, actualResponse, expectedResponse)
	})
}
