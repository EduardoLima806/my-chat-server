package user_route

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eduardolima806/my-chat-server/internal/domain"
	"github.com/eduardolima806/my-chat-server/internal/infra/repository"
	"github.com/eduardolima806/my-chat-server/internal/usecase/user_usecase"
	"github.com/eduardolima806/my-chat-server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Create_New_User(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("user bind error", func(t *testing.T) {

		db, _, _ := sqlmock.New()
		passHasherMock := &util.MockPasswordHasher{}
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		userRepo := repository.NewUserRepository(db)

		user := map[string]string{
			"displayName": "Eduardo Lima",
		}

		userJson, _ := json.Marshal(user)

		req, err := http.NewRequestWithContext(c, http.MethodPost, "/users/create-user", bytes.NewBuffer(userJson))
		assert.NoError(t, err)

		c.Request = req

		handler := &userRouter{
			useCase: *user_usecase.NewUserBaseUserCase(userRepo, passHasherMock),
		}

		handler.createUser(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, strings.Contains(rec.Body.String(), "Error to bind user data: "))
	})

	t.Run("user name invalid", func(t *testing.T) {

		db, _, _ := sqlmock.New()
		passHasherMock := &util.MockPasswordHasher{}
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		userRepo := repository.NewUserRepository(db)

		user := map[string]string{
			"userName":    "ed",
			"displayName": "Eduardo Lima",
			"email":       "eduardolima.dev.io@gmail.com",
			"password":    "P4$$w0rd001",
		}

		userJson, _ := json.Marshal(user)

		req, err := http.NewRequestWithContext(c, http.MethodPost, "/users/create-user", bytes.NewBuffer(userJson))
		assert.NoError(t, err)

		c.Request = req

		handler := &userRouter{
			useCase: *user_usecase.NewUserBaseUserCase(userRepo, passHasherMock),
		}

		handler.createUser(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		expectedError := domain.CreateError(domain.ErrBadRequest.Error(), "username must has at least 5 alphanumerics characters")
		actualError := errors.New(rec.Body.String())
		assert.Equal(t, domain.ErrorCodeResponse(expectedError), domain.ErrorCodeResponse(actualError))
	})

	t.Run("email invalid", func(t *testing.T) {

		db, _, _ := sqlmock.New()
		passHasherMock := &util.MockPasswordHasher{}
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		userRepo := repository.NewUserRepository(db)

		user := map[string]string{
			"userName":    "eduardolima806",
			"displayName": "Eduardo Lima",
			"email":       "@gmail.com",
			"password":    "P4$$w0rd001",
		}

		userJson, _ := json.Marshal(user)

		req, err := http.NewRequestWithContext(c, http.MethodPost, "/users/create-user", bytes.NewBuffer(userJson))
		assert.NoError(t, err)

		c.Request = req

		handler := &userRouter{
			useCase: *user_usecase.NewUserBaseUserCase(userRepo, passHasherMock),
		}

		handler.createUser(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		expectedError := domain.CreateError(domain.ErrBadRequest.Error(), "email is not valid")
		actualError := errors.New(rec.Body.String())
		assert.Equal(t, domain.ErrorCodeResponse(expectedError), domain.ErrorCodeResponse(actualError))
	})

	t.Run("password is not secure", func(t *testing.T) {

		db, _, _ := sqlmock.New()
		passHasherMock := &util.MockPasswordHasher{}
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		userRepo := repository.NewUserRepository(db)

		user := map[string]string{
			"userName":    "eduardolima806",
			"displayName": "Eduardo Lima",
			"email":       "eduardo@gmail.com",
			"password":    "pass123",
		}

		userJson, _ := json.Marshal(user)

		req, err := http.NewRequestWithContext(c, http.MethodPost, "/users/create-user", bytes.NewBuffer(userJson))
		assert.NoError(t, err)

		c.Request = req

		handler := &userRouter{
			useCase: *user_usecase.NewUserBaseUserCase(userRepo, passHasherMock),
		}

		handler.createUser(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		expectedError := domain.CreateError(domain.ErrBadRequest.Error(), "password is not secure")
		actualError := errors.New(rec.Body.String())
		assert.Equal(t, domain.ErrorCodeResponse(expectedError), domain.ErrorCodeResponse(actualError))
	})

	t.Run("user is created", func(t *testing.T) {

		db, mock, _ := sqlmock.New()
		passHasherMock := &util.MockPasswordHasher{}
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		userRepo := repository.NewUserRepository(db)

		user := map[string]string{
			"userName":    "joaquim2019",
			"displayName": "Joaquim Lima",
			"email":       "joaquim@gmail.com",
			"password":    "P4$$word",
		}

		mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnError(sql.ErrNoRows)
		mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnError(sql.ErrNoRows)

		passHasherMock.On("HashPassword", user["password"]).Return("hashedPassword", nil)

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("INSERT INTO app_user").WillReturnRows(rows)

		userJson, _ := json.Marshal(user)

		req, err := http.NewRequestWithContext(c, http.MethodPost, "/users/create-user", bytes.NewBuffer(userJson))
		assert.NoError(t, err)

		c.Request = req

		handler := &userRouter{
			useCase: *user_usecase.NewUserBaseUserCase(userRepo, passHasherMock),
		}

		handler.createUser(c)
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedBody := "{\"CreatedUserId\":1}"
		assert.Equal(t, expectedBody, rec.Body.String())
	})
}

func Test_Login_User(t *testing.T) {

	t.Run("login bind error", func(t *testing.T) {

		db, mock, _ := sqlmock.New()
		passHasherMock := &util.MockPasswordHasher{}
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		userRepo := repository.NewUserRepository(db)

		login := map[string]string{
			"loginError":   "eduardo01",
			"passwordTypo": "P4$$w0rd001",
		}

		loginJson, _ := json.Marshal(login)

		mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnError(sql.ErrNoRows)

		req, err := http.NewRequestWithContext(c, http.MethodPost, "/users/login", bytes.NewBuffer(loginJson))
		assert.NoError(t, err)
		c.Request = req

		handler := &userRouter{
			useCase: *user_usecase.NewUserBaseUserCase(userRepo, passHasherMock),
		}

		handler.loginUser(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, strings.Contains(rec.Body.String(), "Error to bind login data: "))
	})

	t.Run("user login does not exists", func(t *testing.T) {

		db, mock, _ := sqlmock.New()
		passHasherMock := &util.MockPasswordHasher{}
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		userRepo := repository.NewUserRepository(db)

		login := map[string]string{
			"login":    "eduardo01",
			"password": "P4$$w0rd001",
		}

		loginJson, _ := json.Marshal(login)

		mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnError(sql.ErrNoRows)

		req, err := http.NewRequestWithContext(c, http.MethodPost, "/users/login", bytes.NewBuffer(loginJson))
		assert.NoError(t, err)
		c.Request = req

		handler := &userRouter{
			useCase: *user_usecase.NewUserBaseUserCase(userRepo, passHasherMock),
		}

		handler.loginUser(c)

		expectedError := domain.CreateError(domain.ErrBadRequest.Error(), "user login does not exists")
		actualError := errors.New(rec.Body.String())
		assert.Equal(t, domain.ErrorCodeResponse(expectedError), domain.ErrorCodeResponse(actualError))
	})

	t.Run("email does not exists", func(t *testing.T) {

		db, mock, _ := sqlmock.New()
		passHasherMock := &util.MockPasswordHasher{}
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		userRepo := repository.NewUserRepository(db)

		login := map[string]string{
			"login":    "eduardo01@test.com",
			"password": "P4$$w0rd001",
		}

		loginJson, _ := json.Marshal(login)

		mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnError(sql.ErrNoRows)

		req, err := http.NewRequestWithContext(c, http.MethodPost, "/users/login", bytes.NewBuffer(loginJson))
		assert.NoError(t, err)
		c.Request = req

		handler := &userRouter{
			useCase: *user_usecase.NewUserBaseUserCase(userRepo, passHasherMock),
		}

		handler.loginUser(c)

		expectedError := domain.CreateError(domain.ErrBadRequest.Error(), "email does not exists")
		actualError := errors.New(rec.Body.String())
		assert.Equal(t, domain.ErrorCodeResponse(expectedError), domain.ErrorCodeResponse(actualError))
	})

	t.Run("password does not match", func(t *testing.T) {

		db, mockDb, _ := sqlmock.New()
		passHasherMock := &util.MockPasswordHasher{}
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		userRepo := repository.NewUserRepository(db)

		login := map[string]string{
			"login":    "eduardolima806",
			"password": "P4$$w0rd00122",
		}

		loginJson, _ := json.Marshal(login)

		rows := sqlmock.NewRows([]string{"id", "username", "displayname", "email", "password", "created"}).AddRow(1, "eduardolima806", "Eduardo Lima", "eduardolima.dev.io@gmail.com", "P4$$w0rd", time.Now())
		rows2 := sqlmock.NewRows([]string{"id", "username", "displayname", "email", "password", "created"}).AddRow(1, "eduardolima806", "Eduardo Lima", "eduardolima.dev.io@gmail.com", "P4$$w0rd", time.Now())
		mockDb.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnRows(rows)
		mockDb.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnRows(rows2)

		passHasherMock.On("VerifyPassword", login["password"], mock.Anything).Return(false)

		req, err := http.NewRequestWithContext(c, http.MethodPost, "/users/login", bytes.NewBuffer(loginJson))
		assert.NoError(t, err)
		c.Request = req

		handler := &userRouter{
			useCase: *user_usecase.NewUserBaseUserCase(userRepo, passHasherMock),
		}

		handler.loginUser(c)

		expectedError := domain.CreateError(domain.ErrBadRequest.Error(), "password does not match")
		actualError := errors.New(rec.Body.String())
		assert.Equal(t, domain.ErrorCodeResponse(expectedError), domain.ErrorCodeResponse(actualError))
	})

	t.Run("login succeed", func(t *testing.T) {

		db, mockDb, _ := sqlmock.New()
		passHasherMock := &util.MockPasswordHasher{}
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		userRepo := repository.NewUserRepository(db)

		login := map[string]string{
			"login":    "eduardolima806",
			"password": "P4$$w0rd001",
		}

		loginJson, _ := json.Marshal(login)

		rows := sqlmock.NewRows([]string{"id", "username", "displayname", "email", "password", "created"}).AddRow(1, "eduardolima806", "Eduardo Lima", "eduardolima.dev.io@gmail.com", "P4$$w0rd", time.Now())
		rows2 := sqlmock.NewRows([]string{"id", "username", "displayname", "email", "password", "created"}).AddRow(1, "eduardolima806", "Eduardo Lima", "eduardolima.dev.io@gmail.com", "P4$$w0rd", time.Now())
		mockDb.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnRows(rows)
		mockDb.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnRows(rows2)

		passHasherMock.On("VerifyPassword", login["password"], mock.Anything).Return(true)

		req, err := http.NewRequestWithContext(c, http.MethodPost, "/users/login", bytes.NewBuffer(loginJson))
		assert.NoError(t, err)
		c.Request = req

		handler := &userRouter{
			useCase: *user_usecase.NewUserBaseUserCase(userRepo, passHasherMock),
		}

		handler.loginUser(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "\"login succeed\"", rec.Body.String())
	})
}
