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

func Test_UserName_Is_Or_Not_Is_Valid(t *testing.T) {
	testsCases := []string{
		"ed06",
		"eduardolima806",
		"eduardo@123",
	}

	t.Run("username with less then 5 characters", func(t *testing.T) {
		_, err := NewUser(idUser, testsCases[0], displayName, email, password)
		assert.Error(t, err, "username must has at least 5 alphanumerics characters")
	})

	t.Run("username valid", func(t *testing.T) {
		_, err := NewUser(idUser, testsCases[1], displayName, email, password)
		assert.Nil(t, err, "username valid")
	})

	t.Run("username with special characters", func(t *testing.T) {
		_, err := NewUser(idUser, testsCases[2], displayName, email, password)
		assert.Error(t, err, "username must has at least 5 alphanumerics characters")
	})
}

func Test_Email_Is_Or_Not_Is_Valid(t *testing.T) {

	testsCases := []string{
		"invalidmail.com",
		"validmail@company.com",
	}

	t.Run("email is not valid", func(t *testing.T) {
		_, err := NewUser(idUser, userName, displayName, testsCases[0], password)
		assert.Error(t, err, "email is not valid")
	})

	t.Run("email is valid", func(t *testing.T) {
		_, err := NewUser(idUser, userName, displayName, testsCases[1], password)
		assert.Nil(t, err)
	})
}

func Test_If_Password_Is_Not_Or_Not_Is_Secure(t *testing.T) {

	testsCases := []string{
		"password123",
		"P4$$word",
	}

	t.Run("password is not secure", func(t *testing.T) {
		_, err := NewUser(idUser, userName, displayName, email, testsCases[0])
		assert.Error(t, err, "password is not secure")
	})

	t.Run("password is secure", func(t *testing.T) {
		_, err := NewUser(idUser, userName, displayName, email, testsCases[1])
		assert.Nil(t, err)
	})
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
