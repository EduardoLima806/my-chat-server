package domain

type UserRepositoryInterface interface {
	Save(user *User) (int32, error)
	GetUserByID(userID int32) (*User, error)
}
