package usecase

import (
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
	// TODO: Check if the username not exits before create
	idUserCreated, err := cUser.UserRepository.Save(user)
	if err != nil {
		return nil, err
	}
	return &UserOutput{
		CreatedUserId: idUserCreated,
	}, nil
}
