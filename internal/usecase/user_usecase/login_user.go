package user_usecase

import (
	"database/sql"
	"regexp"

	"github.com/eduardolima806/my-chat-server/internal/domain"
	"github.com/eduardolima806/my-chat-server/internal/util"
)

type LoginInput struct {
	Login    string
	Password string
}

type LoginErrorType struct {
	Code        uint8
	Description string
}

var (
	UserLoginNotExists   = LoginErrorType{0, "user login does not exists"}
	EmailNotExists       = LoginErrorType{1, "email does not exists"}
	PasswordDoesNotMatch = LoginErrorType{2, "password does not match"}
)

type LoginOuput struct {
	IsSucceed bool
	ErrorType LoginErrorType
}

type LoginUserUseCase struct {
	UserRepository domain.UserRepositoryInterface
	PasswordHasher util.PasswordHasher
}

type LoginUserUseCaseInterface interface {
	Execute(input LoginInput) (*LoginOuput, error)
}

func NewLoginUserUseCase(userRepo domain.UserRepositoryInterface, passwordHasher util.PasswordHasher) *LoginUserUseCase {
	return &LoginUserUseCase{
		UserRepository: userRepo,
		PasswordHasher: passwordHasher,
	}
}

func (uc *LoginUserUseCase) Execute(loginInput LoginInput) (*LoginOuput, error) {

	userToCheck, err := uc.UserRepository.GetUserByUserNameOrEmail(loginInput.Login)

	if err != nil {
		if err == sql.ErrNoRows {
			errType := UserLoginNotExists
			if checkIsEmail(loginInput.Login) {
				errType = EmailNotExists
			}
			return &LoginOuput{
				IsSucceed: false,
				ErrorType: errType,
			}, nil
		} else {
			return nil, domain.CreateError(domain.ErrInternalServerError.Error(), "could not possible to fetch user")
		}
	}

	if userToCheck != nil &&
		!uc.PasswordHasher.VerifyPassword(loginInput.Password, userToCheck.Password) {
		return &LoginOuput{
			IsSucceed: false,
			ErrorType: PasswordDoesNotMatch,
		}, nil
	}

	return &LoginOuput{
		IsSucceed: true,
	}, nil
}

func checkIsEmail(login string) bool {
	emailRegex := regexp.MustCompile(domain.EmailRegex)
	return emailRegex.Match([]byte(login))
}
