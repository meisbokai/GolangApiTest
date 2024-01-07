package v1_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/meisbokai/GolangApiTest/internal/constants"
	V1Domains "github.com/meisbokai/GolangApiTest/internal/http/domains/v1"
	V1Handlers "github.com/meisbokai/GolangApiTest/internal/http/handlers/v1"
	"github.com/meisbokai/GolangApiTest/internal/http/requests"
	V1Usecases "github.com/meisbokai/GolangApiTest/internal/usecases/v1"
	"github.com/meisbokai/GolangApiTest/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/meisbokai/GolangApiTest/internal/mocks"
)

//TODO: figure out how to use Suit
// type userUsecaseSuite struct {
// 	suite.Suite
// 	repository V1Domains.UserRepository
// 	jwtService JWTService
// }

// func (suite *userUsecaseSuite) SetupSuite() {
// 	// This runs once before all the tests are run
// 	jwtServiceMock = mocks.NewJWTService(t)

// }

// func (suite *userUsecaseSuite) TearDownSuite() {
// 	// This runs once after all the tests are run
// }

// func TestUserSuite(t *testing.T) {
// 	suite.Run(t, new(userUsecaseSuite))
// }

var (
	jwtServiceMock        *mocks.JWTService
	userRepoMock          *mocks.UserRepository
	userUsecase           V1Domains.UserUsecase
	usersDataFromDB       []V1Domains.UserDomain
	userDataFromDB        V1Domains.UserDomain
	userAdminDataFromDB   V1Domains.UserDomain
	userUpdatedDataFromDB V1Domains.UserDomain
	userUpdatedResponse   V1Handlers.BaseResponse
	userDeletedResponse   V1Handlers.BaseResponse
)

func setup(t *testing.T) {
	jwtServiceMock = mocks.NewJWTService(t)
	userRepoMock = mocks.NewUserRepository(t)
	userUsecase = V1Usecases.NewUserUsecase(userRepoMock, jwtServiceMock)
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

	userAdminDataFromDB = V1Domains.UserDomain{
		ID:        "ddfcea5c-d919-4a8f-a631-4ace39337s3a",
		Username:  "testuser1",
		Email:     "testuser1@example.com",
		RoleID:    1,
		Password:  "11111",
		CreatedAt: time.Now(),
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

}

func TestCreateUser(t *testing.T) {
	setup(t)
	req := requests.UserCreateRequest{
		Username: "testuser4",
		Email:    "testuser4@example.com",
		Password: "44444",
	}

	badpasswordReq := requests.UserCreateRequest{
		Username: "testuser4",
		Email:    "testuser4@example.com",
		Password: "4",
	}
	t.Run("Test 1 | Success Create New User", func(t *testing.T) {
		pass, _ := util.GenerateHash("44444")

		userRepoMock.Mock.On("CreateUser", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(nil).Once()
		userRepoMock.Mock.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userDataFromDB, nil).Once()
		result, statusCode, err := userUsecase.CreateUser(context.Background(), req.ToV1Domain())

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, statusCode)
		assert.NotEqual(t, "", result.ID)
		assert.Equal(t, 2, result.RoleID)
		assert.Equal(t, true, util.ValidateHash("44444", pass))
		assert.NotNil(t, result.CreatedAt)
	})

	t.Run("Test 2 | Failure When Creating New User", func(t *testing.T) {
		userRepoMock.Mock.On("CreateUser", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(constants.ErrUnexpected).Once()
		result, statusCode, err := userUsecase.CreateUser(context.Background(), req.ToV1Domain())

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.Equal(t, "", result.ID)
	})

	t.Run("Test 3 | Password not valid", func(t *testing.T) {
		_, statusCode, err := userUsecase.CreateUser(context.Background(), badpasswordReq.ToV1Domain())
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("password cannot be less than 4 characters"), err)

		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})

}

func TestDeleteUser(t *testing.T) {
	setup(t)
	t.Run("Test 1 | Success Delete User", func(t *testing.T) {
		userRepoMock.Mock.On("GetUserByID", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userDataFromDB, nil).Once()
		userRepoMock.Mock.On("DeleteUser", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(nil).Once()

		result, statusCode, err := userUsecase.DeleteUser(context.Background(), "testuser3@example.com")

		actualResponse := V1Handlers.BaseResponse{
			Status:  true,
			Message: "user deleted",
			Data:    map[string]interface{}{"user": result.Username},
		}

		assert.Nil(t, err)
		assert.Equal(t, userDeletedResponse, actualResponse)
		assert.Equal(t, http.StatusOK, statusCode)
	})
}

func TestLogin(t *testing.T) {
	setup(t)
	t.Run("Test 1 | Success Login", func(t *testing.T) {
		req := requests.UserLoginRequest{
			Email:    "testuser3@example.com",
			Password: "33333",
		}

		userDataFromDB.Password, _ = util.GenerateHash(userDataFromDB.Password)

		userRepoMock.Mock.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userDataFromDB, nil).Once()

		jwtServiceMock.Mock.On("GenerateToken",
			mock.AnythingOfType("string"), // userId
			mock.AnythingOfType("string"), // username
			mock.AnythingOfType("bool"),   // isAdmin
			mock.AnythingOfType("string"), // email, password
			mock.AnythingOfType("string")).Return("mockJwtToken3", nil).Once()

		result, statusCode, err := userUsecase.Login(context.Background(), req.ToV1Domain())

		assert.NotNil(t, result)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Nil(t, err)
		assert.Contains(t, result.Token, "JwtToken3")
	})

	t.Run("Test 2 | Success Login (Admin)", func(t *testing.T) {
		req := requests.UserLoginRequest{
			Email:    "testuser1@example.com",
			Password: "11111",
		}
		userAdminDataFromDB.Password, _ = util.GenerateHash(userAdminDataFromDB.Password)

		userRepoMock.Mock.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userAdminDataFromDB, nil).Once()

		jwtServiceMock.Mock.On("GenerateToken",
			mock.AnythingOfType("string"), // userId
			mock.AnythingOfType("string"), // username
			mock.AnythingOfType("bool"),   // isAdmin
			mock.AnythingOfType("string"), // email, password
			mock.AnythingOfType("string")).Return("mockJwtToken1", nil).Once()

		result, statusCode, err := userUsecase.Login(context.Background(), req.ToV1Domain())

		assert.NotNil(t, result)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Nil(t, err)
		assert.Contains(t, result.Token, "JwtToken1")

	})

	t.Run("Test 3 | Invalid Credential", func(t *testing.T) {
		req := requests.UserLoginRequest{
			Email:    "testuser3@test.com",
			Password: "22222",
		}
		userDataFromDB.Password, _ = util.GenerateHash(userDataFromDB.Password)

		userRepoMock.Mock.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userDataFromDB, nil).Once()

		result, statusCode, err := userUsecase.Login(context.Background(), req.ToV1Domain())

		assert.Equal(t, V1Domains.UserDomain{}, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Equal(t, errors.New("invalid email or password(hash)"), err)
		assert.Equal(t, "", result.Token)
	})

}

func TestGetUserByEmail(t *testing.T) {
	setup(t)
	t.Run("Test 1 | Success Get User Data By Email", func(t *testing.T) {
		userRepoMock.Mock.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userDataFromDB, nil).Once()

		result, statusCode, err := userUsecase.GetUserByEmail(context.Background(), "testuser1@example.com")

		assert.Nil(t, err)
		assert.Equal(t, userDataFromDB, result)
		assert.Equal(t, http.StatusOK, statusCode)
	})

	t.Run("Test 2 | User doesn't exist", func(t *testing.T) {
		userRepoMock.Mock.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(V1Domains.UserDomain{}, errors.New("email not found")).Once()

		result, statusCode, err := userUsecase.GetUserByEmail(context.Background(), "dontexistuser@example.com")

		assert.Equal(t, V1Domains.UserDomain{}, result)
		assert.Equal(t, errors.New("email not found"), err)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})
}

func TestGetAllUsers(t *testing.T) {
	setup(t)
	t.Run("Test 1 | Success Get All Users", func(t *testing.T) {
		userRepoMock.Mock.On("GetAllUsers", mock.Anything).Return(usersDataFromDB, nil).Once()

		result, statusCode, err := userUsecase.GetAllUsers(context.Background())

		assert.Nil(t, err)
		assert.Equal(t, usersDataFromDB, result)
		assert.Equal(t, http.StatusOK, statusCode)
	})

}

func TestUpdateUserEmail(t *testing.T) {
	setup(t)
	t.Run("Test 1 | Success Update Email", func(t *testing.T) {
		userRepoMock.Mock.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userDataFromDB, nil).Once()
		userRepoMock.Mock.On("UpdateUserEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain"), mock.Anything).Return(nil).Once()
		userRepoMock.Mock.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userUpdatedDataFromDB, nil).Once()

		result, statusCode, err := userUsecase.UpdateUserEmail(context.Background(), userDataFromDB.Email, userUpdatedDataFromDB.Email)

		actualResponse := V1Handlers.BaseResponse{
			Status:  true,
			Message: "Update success",
			Data:    result,
		}

		assert.Nil(t, err)
		assert.Equal(t, userUpdatedResponse, actualResponse)
		assert.Equal(t, http.StatusOK, statusCode)
	})

}

func TestGetUserByID(t *testing.T) {
	setup(t)
	t.Run("Test 1 | Success Get User Data By ID", func(t *testing.T) {
		userRepoMock.Mock.On("GetUserByID", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(userDataFromDB, nil).Once()

		result, statusCode, err := userUsecase.GetUserByID(context.Background(), "ddfcea5c-d919-4a8f-a631-4ace39337s3a")

		assert.Nil(t, err)
		assert.Equal(t, userDataFromDB, result)
		assert.Equal(t, http.StatusOK, statusCode)
	})

	t.Run("Test 2 | User doesn't exist", func(t *testing.T) {
		userRepoMock.Mock.On("GetUserByID", mock.Anything, mock.AnythingOfType("*v1.UserDomain")).Return(V1Domains.UserDomain{}, errors.New("id not found")).Once()

		result, statusCode, err := userUsecase.GetUserByID(context.Background(), "")

		assert.Equal(t, V1Domains.UserDomain{}, result)
		assert.Equal(t, errors.New("id not found"), err)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})
}
