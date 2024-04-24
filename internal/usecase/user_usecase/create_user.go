package user_usecase

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/eduardolima806/my-chat-server/internal/domain"
	"github.com/eduardolima806/my-chat-server/internal/util"
)

type UserInput struct {
	UserName    string
	DisplayName string
	Email       string
	Password    string
}

type UserOutput struct {
	CreatedUserId int32
}

type CreateUserUseCaseInterface interface {
	Execute(input UserInput) (*UserOutput, error)
}

type CreateUserUseCase struct {
	UserRepository domain.UserRepositoryInterface
	PasswordHasher util.PasswordHasher
}

const IdDummy = 0

func NewCreateUserUseCase(userRepository domain.UserRepositoryInterface, passwordHasher util.PasswordHasher) *CreateUserUseCase {
	return &CreateUserUseCase{
		UserRepository: userRepository,
		PasswordHasher: passwordHasher,
	}
}

func (cUser *CreateUserUseCase) Execute(userInput UserInput) (*UserOutput, error) {
	user, err := domain.NewUser(IdDummy, userInput.UserName, userInput.DisplayName, userInput.Email, userInput.Password)
	if err != nil {
		return nil, err
	}

	err = checkIfUserExists(cUser, userInput)
	if err != nil {
		return nil, err
	}

	user.Password, err = cUser.PasswordHasher.HashPassword(user.Password)

	if err != nil {
		return nil, errors.New("could not be possible encrypt password")
	}

	idUserCreated, err := cUser.UserRepository.Save(user)

	if err != nil {
		return nil, err
	}

	return &UserOutput{
		CreatedUserId: idUserCreated,
	}, nil
}

func checkIfUserExists(cUser *CreateUserUseCase, userInput UserInput) error {

	userToCheck, err := cUser.UserRepository.GetUserByUserNameOrEmail(userInput.UserName)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if userToCheck != nil && userInput.UserName == userToCheck.UserName {
		return errors.New("username already exists")
	}

	userToCheck, err = cUser.UserRepository.GetUserByUserNameOrEmail(userInput.Email)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if userToCheck != nil && userToCheck.Email == userInput.Email {
		return fmt.Errorf("already exists an user with this e-email: %s", userInput.Email)
	}

	return nil
}
