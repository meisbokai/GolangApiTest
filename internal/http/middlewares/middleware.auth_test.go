package middlewares_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	config "github.com/meisbokai/GolangApiTest/internal/configs"
	"github.com/meisbokai/GolangApiTest/internal/http/middlewares"
	"github.com/meisbokai/GolangApiTest/pkg/jwt"
	jwtPkg "github.com/meisbokai/GolangApiTest/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

var (
	jwtService          jwtPkg.JWTService
	ginMock             *gin.Engine
	authBasicMiddleware gin.HandlerFunc
	authAdminMiddleware gin.HandlerFunc
)

const (
	adminEndpoint = "/admin"
	forEveryone   = "/everyone"
)

func authenticatedHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Authentication successful",
	})
}

func setup(t *testing.T) {
	jwtService = jwt.NewJWTService(config.AppConfig.JWTSecret, config.AppConfig.JWTIssuer, config.AppConfig.JWTExpired)
	authBasicMiddleware = middlewares.NewAuthMiddleware(jwtService, false)
	authAdminMiddleware = middlewares.NewAuthMiddleware(jwtService, true)

	ginMock = gin.New()
	ginMock.GET(forEveryone, authBasicMiddleware, authenticatedHandler)
	ginMock.GET(adminEndpoint, authAdminMiddleware, authenticatedHandler)
}

func generateToken(isAdmin bool) (token string, err error) {
	token, err = jwtService.GenerateToken("ddfcea5c-d919-4a8f-a631-4ace39337s3a", "randomUser", isAdmin, "random@example.com", "12345")
	return
}

func getAdminToken() (string, error) {
	return generateToken(true)
}

func getBasicToken() (string, error) {
	return generateToken(false)
}

func TestAuthMiddleware(t *testing.T) {
	setup(t)
	// Define route

	t.Run("Test 1 | Success Get Admin Handler", func(t *testing.T) {
		token, err := getAdminToken()
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, adminEndpoint, nil)

		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", fmt.Sprintf("jwt %s", token))

		// Perform request
		ginMock.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "Authentication successful")
	})
	t.Run("Test 2 | Invalid Token", func(t *testing.T) {
		token := "inva.lidto.ken"

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, forEveryone, nil)

		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", fmt.Sprintf("jwt %s", token))

		// Perform request
		ginMock.ServeHTTP(w, r)

		body := w.Body.String()
		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "invalid token")
	})
	t.Run("Test 3 | Must contain Authorization Header", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, forEveryone, nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		ginMock.ServeHTTP(w, r)

		body := w.Body.String()
		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "missing authorization header")
	})
	t.Run("Test 4 | Invalid Header Format", func(t *testing.T) {
		token, err := getBasicToken()
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, forEveryone, nil)

		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", fmt.Sprintf("invalid header: %s", token))

		// Perform request
		ginMock.ServeHTTP(w, r)

		body := w.Body.String()
		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "invalid header format")
	})
	t.Run("Test 5 | Must contain jwt", func(t *testing.T) {
		token, err := getBasicToken()
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, forEveryone, nil)

		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", fmt.Sprintf("notjwt %s", token))

		// Perform request
		ginMock.ServeHTTP(w, r)

		body := w.Body.String()
		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "token must contain 'jtw'")
	})
	t.Run("Test 6 | Not Authorize", func(t *testing.T) {
		token, err := getBasicToken()
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, adminEndpoint, nil)

		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", fmt.Sprintf("jwt %s", token))

		// Perform request
		ginMock.ServeHTTP(w, r)

		body := w.Body.String()
		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "you don't have access for this action")
	})
}
