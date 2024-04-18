package usecase

import (
	"database/sql"
	"errors"

	"github.com/eduardolima806/my-chat-server/internal/domain"
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

type CreateUserUseCase struct {
	UserRepository domain.UserRepositoryInterface
}

const IdDummy = 0

func NewCreateUserUseCase(userRepository domain.UserRepositoryInterface) *CreateUserUseCase {
	return &CreateUserUseCase{
		UserRepository: userRepository,
	}
}

func (cUser *CreateUserUseCase) Execute(userInput UserInput) (*UserOutput, error) {
	user, err := domain.NewUser(IdDummy, userInput.UserName, userInput.DisplayName, userInput.Email, userInput.Password)
	if err != nil {
		return nil, err
	}

	// TODO: Check if exists e-mail as well ?
	userToCheck, err := cUser.UserRepository.GetUserByUserName(userInput.UserName)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if userToCheck != nil {
		return nil, errors.New("username already exists")
	}

	idUserCreated, err := cUser.UserRepository.Save(user)

	if err != nil {
		return nil, err
	}

	return &UserOutput{
		CreatedUserId: idUserCreated,
	}, nil
}
