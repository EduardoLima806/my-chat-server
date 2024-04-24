package user_usecase

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eduardolima806/my-chat-server/internal/infra/repository"
	"github.com/eduardolima806/my-chat-server/internal/util"
	"github.com/stretchr/testify/assert"
)

func Test_If_UserName_Login_Success(t *testing.T) {
	loginInput := LoginInput{Login: "eduardolima806", Password: "P4$$w0rd001"}
	db, mock, _ := sqlmock.New()
	userRepository := repository.NewUserRepository(db)
	ucLogin := NewLoginUserUseCase(userRepository, &util.DefaultPasswordHasher{})

	rows := sqlmock.NewRows([]string{"id", "username", "displayname", "email", "password", "created"}).AddRow(1, "eduardolima806", "Eduardo Lima", "eduardolima.dev.io@gmail.com", "$2a$14$dI1.i3EBN4Zl0FuVj.gDdOA4QzN6Bg9DVrXlsQJCFemwrnTj8OPh.", time.Now())
	mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnRows(rows)

	loginOutput, _ := ucLogin.Execute(loginInput)
	assert.NotNil(t, loginOutput)
	assert.True(t, loginOutput.IsSucceed)
}

func Test_If_UserName_Login_Is_Empty_Not_Error(t *testing.T) {
	loginInput := LoginInput{Login: "", Password: "P4$$w0rd001"}
	db, mock, _ := sqlmock.New()
	userRepository := repository.NewUserRepository(db)
	ucLogin := NewLoginUserUseCase(userRepository, &util.DefaultPasswordHasher{})

	mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnError(sql.ErrNoRows)

	_, err := ucLogin.Execute(loginInput)
	assert.Nil(t, err)
}

func Test_Error_When_UserName_Login_Not_Exists(t *testing.T) {
	loginInput := LoginInput{Login: "eduardolima", Password: "P4$$w0rd001"}
	db, mock, _ := sqlmock.New()
	userRepository := repository.NewUserRepository(db)
	ucLogin := NewLoginUserUseCase(userRepository, &util.DefaultPasswordHasher{})

	mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnError(sql.ErrNoRows)

	loginOutput, _ := ucLogin.Execute(loginInput)
	assert.False(t, loginOutput.IsSucceed)
	assert.Equal(t, UserLoginNotExists, loginOutput.ErrorType)
}

func Test_If_Get_Error_When_Try_Fetch_User(t *testing.T) {
	loginInput := LoginInput{Login: "eduardolima", Password: "P4$$w0rd001"}
	db, mock, _ := sqlmock.New()
	userRepository := repository.NewUserRepository(db)
	ucLogin := NewLoginUserUseCase(userRepository, &util.DefaultPasswordHasher{})

	mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnError(errors.New("an internal error"))

	loginOutput, err := ucLogin.Execute(loginInput)
	assert.Nil(t, loginOutput)
	assert.NotNil(t, err)
}

func Test_Error_When_Email_Login_Not_Exists(t *testing.T) {
	loginInput := LoginInput{Login: "eduardolima.dev.io@gmail.com", Password: "P4$$w0rd001"}
	db, mock, _ := sqlmock.New()
	userRepository := repository.NewUserRepository(db)
	ucLogin := NewLoginUserUseCase(userRepository, &util.DefaultPasswordHasher{})

	mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnError(sql.ErrNoRows)

	loginOutput, _ := ucLogin.Execute(loginInput)
	assert.False(t, loginOutput.IsSucceed)
	assert.Equal(t, EmailNotExists, loginOutput.ErrorType)
}

func Test_Error_When_Password_Does_Not_Match(t *testing.T) {
	loginInput := LoginInput{Login: "eduardolima.dev.io@gmail.com", Password: "P4$$w0rd001Not"}
	db, mock, _ := sqlmock.New()
	userRepository := repository.NewUserRepository(db)
	ucLogin := NewLoginUserUseCase(userRepository, &util.DefaultPasswordHasher{})

	rows := sqlmock.NewRows([]string{"id", "username", "displayname", "email", "password", "created"}).AddRow(1, "eduardolima806", "Eduardo Lima", "eduardolima.dev.io@gmail.com", "$2a$14$dI1.i3EBN4Zl0FuVj.gDdOA4QzN6Bg9DVrXlsQJCFemwrnTj8OPh.", time.Now())
	mock.ExpectQuery("SELECT id, username, displayname, email, password, created FROM app_user").WillReturnRows(rows)

	loginOutput, _ := ucLogin.Execute(loginInput)
	assert.False(t, loginOutput.IsSucceed)
	assert.Equal(t, PasswordDoesNotMatch, loginOutput.ErrorType)
}
