package repository

import (
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eduardolima806/my-chat-server/internal/domain"
	"github.com/stretchr/testify/assert"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func Test_If_The_User_Is_Saved(t *testing.T) {
	const insertQuery = "INSERT INTO app_user"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userRepo := NewUserRepository(db)
	userDomain, _ := domain.NewUser(1, "eduardolima806", "Eduardo Lima", "eduardolima.dev.io@gmail.com", "P4$$w0rd")

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(insertQuery).WithArgs(userDomain.UserName, userDomain.DisplayName, userDomain.Email, userDomain.Password, AnyTime{}).WillReturnRows(rows)
	var createdId int32
	if createdId, err = userRepo.Save(userDomain); err != nil {
		t.Errorf("error was not expected while insert user: %s", err)
	}

	assert.Equal(t, int32(1), createdId)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_If_Get_Error_The_User_Is_Not_Saved(t *testing.T) {
	const insertQuery = "INSERT INTO app_user"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userRepo := NewUserRepository(db)
	userDomain, _ := domain.NewUser(1, "eduardolima806", "Eduardo Lima", "eduardolima.dev.io@gmail.com", "P4$$w0rd")

	mock.ExpectQuery(insertQuery).WithArgs(userDomain.UserName, userDomain.DisplayName, userDomain.Email, userDomain.Password, AnyTime{}).WillReturnError(errors.New("error to insert user"))
	var createdId int32
	if createdId, err = userRepo.Save(userDomain); err != nil {
		assert.EqualError(t, err, "error to insert user")
		assert.Equal(t, IdError, createdId)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_If_The_User_Fetched_When_Search_By_UserName(t *testing.T) {
	const selectQuery = "SELECT id, username, displayname, email, password, created FROM app_user"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	timestamp := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "username", "displayname", "email", "password", "created"}).AddRow(1, "eduardolima806", "Eduardo Lima", "eduardolima.dev.io@gmail.com", "P4$$w0rd", timestamp)
	mock.ExpectQuery(selectQuery).WithArgs("eduardolima806").WillReturnRows(rows)
	userRepo := NewUserRepository(db)
	fetchedUser, err := userRepo.GetUserByUserNameOrEmail("eduardolima806")

	if err != nil {
		t.Errorf("error was not expected while fetching user: %s", err)
	}

	expectedUser := &domain.User{ID: 1,
		UserName:    "eduardolima806",
		DisplayName: "Eduardo Lima",
		Email:       "eduardolima.dev.io@gmail.com",
		Password:    "P4$$w0rd",
		Created:     timestamp}

	assert.EqualValues(t, expectedUser, fetchedUser)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_If_Get_Error_When_Search_By_UserName(t *testing.T) {
	const selectQuery = "SELECT id, username, displayname, email, password, created FROM app_user"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery(selectQuery).WithArgs("eduardolima806").WillReturnError(errors.New("error to get user"))
	userRepo := NewUserRepository(db)

	_, err = userRepo.GetUserByUserNameOrEmail("eduardolima806")

	assert.EqualError(t, err, "error to get user")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
