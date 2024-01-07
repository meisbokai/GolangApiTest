package v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	config "github.com/meisbokai/GolangApiTest/internal/configs"
	"github.com/meisbokai/GolangApiTest/internal/constants"
	V1Domains "github.com/meisbokai/GolangApiTest/internal/http/domains/v1"
	V1Handlers "github.com/meisbokai/GolangApiTest/internal/http/handlers/v1"
	"github.com/meisbokai/GolangApiTest/internal/http/requests"
	V1Usecases "github.com/meisbokai/GolangApiTest/internal/usecases/v1"
	"github.com/meisbokai/GolangApiTest/internal/util"
	jwtPkg "github.com/meisbokai/GolangApiTest/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/meisbokai/GolangApiTest/internal/mocks"
)

var (
	jwtServiceMock        *mocks.JWTService
	userRepoMock          *mocks.UserRepository
	userUsecase           V1Domains.UserUsecase
	userHandler           V1Handlers.UserHandler
	usersDataFromDB       []V1Domains.UserDomain
	userDataFromDB        V1Domains.UserDomain
	userUpdatedDataFromDB V1Domains.UserDomain
	userUpdatedResponse   V1Handlers.BaseResponse
	userDeletedResponse   V1Handlers.BaseResponse
	ginMock               *gin.Engine
)

func setup(t *testing.T) {
	jwtServiceMock = mocks.NewJWTService(t)
	userRepoMock = mocks.NewUserRepository(t)
	userUsecase = V1Usecases.NewUserUsecase(userRepoMock, jwtServiceMock)
	userHandler = V1Handlers.NewUserHandler(userUsecase)

	usersDataFromDB = []V1Domains.UserDomain{
		{
			ID:        "ddfcea5c-d919-4a8f-a631-4ace39337s3a",
			Username:  "testuser1",
			Email:     "testuser1@example.com",
			RoleID:    1,
			Password:  "11111",
			CreatedAt: time.Now(),
		},
		{
			ID:        "wifff3jd-idhd-0sis-8dua-4fiefie37kfj",
			Username:  "testuser2",
			Email:     "testuser2@example.com",
			RoleID:    2,
			Password:  "22222",
			CreatedAt: time.Now(),
		},
	}

	userCreatedTime := time.Now()
	userDataFromDB = V1Domains.UserDomain{
		ID:        "fjskeie8-jfk8-qke0-sksj-ksjf89e8ehfu",
		Username:  "testuser3",
		Email:     "testuser3@example.com",
		Password:  "33333",
		RoleID:    2,
		CreatedAt: userCreatedTime,
	}

	userUpdatedDataFromDB = V1Domains.UserDomain{
		ID:        "fjskeie8-jfk8-qke0-sksj-ksjf89e8ehfu",
		Username:  "testuser3",
		Email:     "testuser3updated@example.com",
		Password:  "33333",
		RoleID:    2,
		CreatedAt: userCreatedTime,
	}

	userUpdatedResponse = V1Handlers.BaseResponse{
		Status:  true,
		Message: "Update success",
		Data:    userUpdatedDataFromDB,
	}

	userDeletedResponse = V1Handlers.BaseResponse{
		Status:  true,
		Message: "user deleted",
		Data:    map[string]interface{}{"user": "testuser3"},
	}

	// Create Gin Engine
	ginMock = gin.Default()
	ginMock.Use(lazyAuth)
}

func lazyAuth(ctx *gin.Context) {
	// hash
	pass, _ := util.GenerateHash(userDataFromDB.Password)
	// prepare claims
	jwtClaims := jwtPkg.JwtCustomClaim{
		UserID:   userDataFromDB.ID,
		Username: userDataFromDB.Username,
		IsAdmin:  false,
		Email:    userDataFromDB.Email,
		Password: pass,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * time.Duration(config.AppConfig.JWTExpired))},
			Issuer:    userDataFromDB.Username,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	ctx.Set(constants.AuthenticatedClaimKey, jwtClaims)
}

func TestCreateUser(t *testing.T) {
	setup(t)

	ginMock.POST(constants.EndpointV1+"/auth/signup", userHandler.CreateUser)
	t.Run("Test 1 | Success Create New User", func(t *testing.T) {
		// Create login request
		req := requests.UserCreateRequest{
			Username: "testuser3",
			Email:    "testuser3@example.com",
			Password: "33333",
		}
		reqBody, _ := json.Marshal(req)

		userRepoMock.Mock.On("CreateUser", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(nil).Once()
		userRepoMock.Mock.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userDataFromDB, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, constants.EndpointV1+"/auth/signup", bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		ginMock.ServeHTTP(w, r)

		// parsing json to raw text
		response := w.Body.String()

		type UserListResponse struct {
			Status  bool
			Message string
			Data    struct {
				User V1Domains.UserDomain
			}
		}
		var responseJson UserListResponse

		json.Unmarshal([]byte(response), &responseJson)

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Equal(t, "registration user success", responseJson.Message)
		assert.Equal(t, req.Username, responseJson.Data.User.Username)
		assert.Equal(t, req.Email, responseJson.Data.User.Email)

	})

	t.Run("Test 2 | Invalid Email", func(t *testing.T) {
		// Create login request
		req := requests.UserCreateRequest{
			Username: "testuser3",
			Email:    "notValidEmail",
			Password: "33333",
		}
		reqBody, _ := json.Marshal(req)

		// userRepoMock.Mock.On("CreateUser", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(nil).Once()
		// userRepoMock.Mock.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userDataFromDB, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, constants.EndpointV1+"/auth/signup", bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		ginMock.ServeHTTP(w, r)

		// parsing json to raw text
		response := w.Body.String()

		var responseJson V1Handlers.BaseResponse

		json.Unmarshal([]byte(response), &responseJson)

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, responseJson.Message, "is not a valid email address")
		// assert.Equal(t, req.Username, responseJson.Data.User.Username)
		// assert.Equal(t, req.Email, responseJson.Data.User.Email)

	})

}

func TestLogin(t *testing.T) {
	setup(t)

	// Define route
	ginMock.POST(constants.EndpointV1+"/auth/login", userHandler.Login)
	t.Run("Test 1 | Success Login", func(t *testing.T) {
		// hash password field
		var err error
		userDataFromDB.Password, err = util.GenerateHash(userDataFromDB.Password)
		if err != nil {
			t.Error(err)
		}

		// Create login request
		req := requests.UserLoginRequest{
			Email:    "testuser3@example.com",
			Password: "33333",
		}
		reqBody, _ := json.Marshal(req)

		userRepoMock.Mock.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userDataFromDB, nil).Once()
		jwtServiceMock.Mock.On("GenerateToken",
			mock.AnythingOfType("string"), // userId
			mock.AnythingOfType("string"), // username
			mock.AnythingOfType("bool"),   // isAdmin
			mock.AnythingOfType("string"), // email, password
			mock.AnythingOfType("string")).Return("mockJwtToken3", nil).Once()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, constants.EndpointV1+"/auth/login", bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		ginMock.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "login success")
		assert.Contains(t, body, "JwtToken3")
		assert.Contains(t, body, req.Email)
	})
	t.Run("Test 2 | User is Not Exists", func(t *testing.T) {
		req := requests.UserLoginRequest{
			Email:    "dontexistuser@example.com",
			Password: "dontexist",
		}
		reqBody, _ := json.Marshal(req)

		userRepoMock.Mock.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(V1Domains.UserDomain{}, constants.ErrUserNotFound).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, constants.EndpointV1+"/auth/login", bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		ginMock.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "invalid email or password")
	})

	t.Run("Test 3 | Empty Password", func(t *testing.T) {
		// hash password field
		var err error
		userDataFromDB.Password, err = util.GenerateHash(userDataFromDB.Password)
		if err != nil {
			t.Error(err)
		}

		// Create login request
		req := requests.UserLoginRequest{
			Email:    "testuser3@example.com",
			Password: "",
		}
		reqBody, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, constants.EndpointV1+"/auth/login", bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		ginMock.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "is a required field")
	})
}

func TestGetUserData(t *testing.T) {
	setup(t)
	// Define route
	ginMock.GET("/users/self", userHandler.GetSelfUser)

	t.Run("Test 1 | Success Fetched User Data", func(t *testing.T) {
		userRepoMock.Mock.On("GetUserByID", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userDataFromDB, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/users/self", nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		ginMock.ServeHTTP(w, r)

		// parsing json to raw text
		response := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, response, "user data fetched successfully")
		// assert.Contains(t, response, userDataFromDB.ID)
		// assert.Contains(t, response, userDataFromDB.Email)
		// assert.Contains(t, response, userDataFromDB.Token)

	})

	t.Run("Test 2 | Failed to fetch User Data", func(t *testing.T) {
		userRepoMock.Mock.On("GetUserByID", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(V1Domains.UserDomain{}, constants.ErrUserNotFound).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/users/self", nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		ginMock.ServeHTTP(w, r)

		// Assertions
		// Assert status code
		assert.NotEqual(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
	})
}

func TestGetUserByEmail(t *testing.T) {
	setup(t)
	// Define route
	ginMock.GET("/admin/users/email", userHandler.GetUserByEmail)

	t.Run("Test 1 | Success Fetched User Data", func(t *testing.T) {
		userRepoMock.Mock.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userDataFromDB, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/admin/users/email", nil)

		r.Header.Set("Content-Type", "application/json")

		q := r.URL.Query()
		q.Add("email", "testuser3@example.com")
		r.URL.RawQuery = q.Encode()

		// Perform request
		ginMock.ServeHTTP(w, r)

		// parsing json to raw text
		response := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, response, "user data fetched successfully")
		assert.Contains(t, response, userDataFromDB.Email)

	})

	t.Run("Test 2 | Failed to fetch User Data", func(t *testing.T) {
		userRepoMock.Mock.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(V1Domains.UserDomain{}, constants.ErrUserNotFound).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/admin/users/email", nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		ginMock.ServeHTTP(w, r)

		// Assertions
		// Assert status code
		assert.NotEqual(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
	})
}

func TestGetAllUsers(t *testing.T) {
	setup(t)

	ginMock.GET("/admin/users/all", userHandler.GetAllUserData)
	t.Run("Test 1 | Success Fetched User Data", func(t *testing.T) {
		userRepoMock.Mock.On("GetAllUsers", mock.Anything).Return(usersDataFromDB, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/admin/users/all", nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		ginMock.ServeHTTP(w, r)

		// parsing json to raw text
		response := w.Body.String()

		type UserListResponse struct {
			Status  bool
			Message string
			Data    struct {
				User []V1Domains.UserDomain
			}
		}
		var responseJson UserListResponse

		json.Unmarshal([]byte(response), &responseJson)

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, responseJson.Message, "user data fetched successfully")
		assert.Equal(t, len(usersDataFromDB), len(responseJson.Data.User))
		assert.Equal(t, reflect.TypeOf([]V1Domains.UserDomain{}), reflect.TypeOf(responseJson.Data.User))
	})

}

func TestDeleteUser(t *testing.T) {
	setup(t)

	ginMock.DELETE("/users/delete", userHandler.DeleteUser)
	t.Run("Test 1 | Success Delete User", func(t *testing.T) {
		userRepoMock.Mock.On("GetUserByID", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userDataFromDB, nil).Once()
		userRepoMock.Mock.On("DeleteUser", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/users/delete", nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		ginMock.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "user deleted")
		assert.Contains(t, body, userDataFromDB.Username)

	})

}

func TestUpdateUserEmail(t *testing.T) {
	setup(t)

	ginMock.PUT("/users/updateEmail", userHandler.UpdateUserEmail)
	t.Run("Test 1 | Success Update User Email", func(t *testing.T) {
		// Create update request
		req := requests.UserUpdateEmailRequest{
			OldEmail: "testuser3@example.com",
			NewEmail: "testuser3updated@example.com",
		}
		reqBody, _ := json.Marshal(req)

		userRepoMock.Mock.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userDataFromDB, nil).Once()
		userRepoMock.Mock.On("UpdateUserEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain"), mock.Anything).Return(nil).Once()
		userRepoMock.Mock.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userUpdatedDataFromDB, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/users/updateEmail", bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		ginMock.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "Update success")
		assert.Contains(t, body, userUpdatedDataFromDB.Username)
		assert.Contains(t, body, userUpdatedDataFromDB.Email)

	})

	t.Run("Test 2 | Fail Update User Email (Same Email)", func(t *testing.T) {
		// Create update request
		req := requests.UserUpdateEmailRequest{
			OldEmail: "testuser3@example.com",
			NewEmail: "testuser3@example.com",
		}
		reqBody, _ := json.Marshal(req)

		userRepoMock.Mock.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userDataFromDB, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/users/updateEmail", bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		ginMock.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "New Email is the same as old email")
	})

	t.Run("Test 3 | Fail Update User Email (Invalid)", func(t *testing.T) {
		// Create update request
		req := requests.UserUpdateEmailRequest{
			OldEmail: "testuser3@example.com",
			NewEmail: "notValidEmail",
		}
		reqBody, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/users/updateEmail", bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		ginMock.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "not a valid email address")
	})

}
