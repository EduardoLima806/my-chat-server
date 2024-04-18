package domain

type UserRepositoryInterface interface {
	Save(user *User) (int32, error)
	GetUserByUserName(userName string) (*User, error)
}
