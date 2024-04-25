package user_usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eduardolima806/my-chat-server/internal/domain"
	"github.com/eduardolima806/my-chat-server/internal/infra/repository"
	"github.com/eduardolima806/my-chat-server/internal/util"
	"github.com/stretchr/testify/assert"
)

func Test_If_Get_Error_To_Create_Invalid_User(t *testing.T) {
	userRepository := repository.NewUserRepository(nil)
	passHasherMock := &util.MockPasswordHasher{}
	userInput := UserInput{UserName: "ed12"}
	ucCreate := NewCreateUserUseCase(userRepository, passHasherMock)
	_, err := ucCreate.Execute(userInput)
	expectedError := domain.CreateError(domain.ErrBadRequest.Error(), "username must has at least 5 alphanumerics characters")
	assert.EqualError(t, err, expectedError.Error())
}

func Test_If_Get_Error_When_Create_An_Existing_User(t *testing.T) {
	db, mock, _ := sqlmock.New()
	userRepository := repository.NewUserRepository(db)
	passHasherMock := &util.MockPasswordHasher{}
	userInput := UserInput{UserName: "eduardolima806", Email: "eduardolima.dev.io@gmail.com", Password: "P4$$w0rd"}
	ucCreate := NewCreateUserUseCase(userRepository, passHasherMock)

	t.Run("username already exists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "username", "displayname", "email", "password", "created"}).AddRow(1, "eduardolima806", "Eduardo Lima", "eduardolima.dev.io@gmail.com", "P4$$w0rd", time.Now())
		mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnRows(rows)

		_, err := ucCreate.Execute(userInput)
		expectedError := domain.CreateError(domain.ErrBadRequest.Error(), "username already exists")
		assert.EqualError(t, err, expectedError.Error())
	})

	t.Run("email already exists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "username", "displayname", "email", "password", "created"}).AddRow(1, "eduardolima806", "Eduardo Lima", "eduardolima.dev.io@gmail.com", "P4$$w0rd", time.Now())
		rows2 := sqlmock.NewRows([]string{"id", "username", "displayname", "email", "password", "created"}).AddRow(1, "eduardolima806", "Eduardo Lima", "eduardolima.dev.io@gmail.com", "P4$$w0rd", time.Now())
		mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnRows(rows)
		mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnRows(rows2)

		userInput.UserName = "eduardo123"
		_, err := ucCreate.Execute(userInput)
		expectedError := domain.CreateError(domain.ErrBadRequest.Error(), fmt.Sprintf("already exists an user with this e-email: %s", userInput.Email))
		assert.EqualError(t, err, expectedError.Error())
	})
}

func Test_If_Get_Error_When_Try_Encrypt_Password(t *testing.T) {
	db, mock, _ := sqlmock.New()
	userRepository := repository.NewUserRepository(db)
	passHasherMock := &util.MockPasswordHasher{}
	userInput := UserInput{UserName: "edulima", Email: "eduardolima@gmail.com", Password: "P4$$w0rd"}
	ucCreate := NewCreateUserUseCase(userRepository, passHasherMock)

	rows := sqlmock.NewRows([]string{"id", "username", "displayname", "email", "password", "created"}).AddRow(1, "eduardolima806", "Eduardo Lima", "eduardolima.dev.io@gmail.com", "P4$$w0rd", time.Now())
	rows2 := sqlmock.NewRows([]string{"id", "username", "displayname", "email", "password", "created"}).AddRow(1, "eduardolima806", "Eduardo Lima", "eduardolima.dev.io@gmail.com", "P4$$w0rd", time.Now())
	mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnRows(rows)
	mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnRows(rows2)

	passHasherMock.On("HashPassword", userInput.Password).Return("", errors.New("encryptation error")).Once()

	userCreateOutput, err := ucCreate.Execute(userInput)
	assert.Nil(t, userCreateOutput)
	expecteError := domain.CreateError(domain.ErrInternalServerError.Error(), "encryptation error")
	assert.EqualError(t, err, expecteError.Error())
}

func Test_User_Is_Created_When_User_No_Existing(t *testing.T) {
	db, mock, _ := sqlmock.New()
	userRepository := repository.NewUserRepository(db)
	passHasherMock := &util.MockPasswordHasher{}
	userInput := UserInput{UserName: "eduardolimaNew", Email: "eduardolima.dev.io@gmail.com", Password: "P4$$w0rd"}
	ucCreate := NewCreateUserUseCase(userRepository, passHasherMock)
	passHasherMock.On("HashPassword", userInput.Password).Return("hashedPassword", nil)

	mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnError(sql.ErrNoRows)
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("INSERT INTO app_user").WillReturnRows(rows)

	userOutput, _ := ucCreate.Execute(userInput)
	assert.Equal(t, int32(1), userOutput.CreatedUserId)
	passHasherMock.AssertExpectations(t)
}
