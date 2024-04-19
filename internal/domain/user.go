package domain

import (
	"errors"
	"regexp"
	"time"
)

type User struct {
	ID          int32
	UserName    string
	DisplayName string
	Email       string
	Password    string
	Created     time.Time
}

func NewUser(id int32, userName string, displayName string, email string, password string) (*User, error) {
	user := &User{
		ID:          id,
		UserName:    userName,
		DisplayName: displayName,
		Email:       email,
		Password:    password,
		Created:     time.Now(),
	}
	err := user.Validate()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Validate() error {

	userNameRegex := regexp.MustCompile(`^[a-zA-Z0-9]{5,}$`)

	if !userNameRegex.MatchString(u.UserName) {
		return errors.New("username must has at least 5 alphanumerics characters")
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	if !emailRegex.MatchString(u.Email) {
		return errors.New("email is not valid")
	}

	passRules := []string{".{7,}", "[a-z]", "[A-Z]", "[0-9]", "[^\\d\\w]"}
	for _, v := range passRules {
		passRegex := regexp.MustCompile(v)
		if !passRegex.MatchString(u.Password) {
			return errors.New("password is not secure")
		}
	}

	return nil
}
