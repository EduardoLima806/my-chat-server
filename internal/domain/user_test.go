package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	idUser      = 1
	userName    = "eduardolima806"
	displayName = "Eduardo Lima"
	email       = "eduardolima.dev.io@gmail.com"
	password    = "Pa$$w0rd"
)

func Test_If_Get_Error_If_UserName_Is_Not_Valid(t *testing.T) {
	_, err := NewUser(idUser, "ed06", displayName, email, password)
	assert.Error(t, err, "username must has at least 5 characters")
}

func Test_If_Get_Error_If_Email_Is_Not_Valid(t *testing.T) {
	_, err := NewUser(idUser, userName, displayName, "invalidmail.com", password)
	assert.Error(t, err, "email is not valid")
}

func Test_If_Get_Error_If_Password_Is_Not_Secure(t *testing.T) {
	_, err := NewUser(idUser, userName, displayName, email, "password123")
	assert.Error(t, err, "password is not secure")
}

func Test_If_Get_No_Error_For_Valid_User(t *testing.T) {
	user, err := NewUser(idUser, userName, displayName, email, password)
	assert.Nil(t, err)
	expectedUser := User{
		ID:          idUser,
		UserName:    userName,
		DisplayName: displayName,
		Email:       email,
		Password:    password,
	}
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.UserName, user.UserName)
	assert.Equal(t, expectedUser.DisplayName, user.DisplayName)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Password, user.Password)
}
