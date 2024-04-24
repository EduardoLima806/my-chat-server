package user_usecase

import (
	"github.com/eduardolima806/my-chat-server/internal/domain"
	"github.com/eduardolima806/my-chat-server/internal/util"
)

type UserBaseUserCase struct {
	CreateUserUseCase CreateUserUseCaseInterface
	LoginUserUseCase  LoginUserUseCaseInterface
}

func NewUserBaseUserCase(userRepository domain.UserRepositoryInterface, passwordHasher util.PasswordHasher) *UserBaseUserCase {
	return &UserBaseUserCase{
		CreateUserUseCase: NewCreateUserUseCase(userRepository, passwordHasher),
		LoginUserUseCase:  NewLoginUserUseCase(userRepository, passwordHasher),
	}
}
