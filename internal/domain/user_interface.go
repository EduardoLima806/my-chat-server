package domain

type UserRepositoryInterface interface {
	Save(user *User) (int32, error)
	GetUserByUserNameOrEmail(userNameOrEmail string) (*User, error)
}
