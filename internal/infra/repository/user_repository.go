package repository

import (
	"database/sql"

	"github.com/eduardolima806/my-chat-server/internal/domain"
)

const IdError = int32(-1)

type UserRepository struct {
	Db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		Db: db,
	}
}

func (userRepo *UserRepository) Save(user *domain.User) (int32, error) {
	lastInsertId := 0
	err := userRepo.Db.QueryRow("INSERT INTO app_user (username, displayname, email, password, created) VALUES ($1,$2,$3,$4,$5) RETURNING id",
		user.UserName, user.DisplayName, user.Email, user.Password, user.Created).Scan(&lastInsertId)
	if err != nil {
		return IdError, err
	}

	return int32(lastInsertId), nil
}

func (userRepo *UserRepository) GetUserByUserName(userName string) (*domain.User, error) {
	user := domain.User{}
	err := userRepo.Db.QueryRow("SELECT id, username, displayname, email, password, created FROM app_user WHERE username = $1", userName).Scan(
		&user.ID, &user.UserName, &user.DisplayName, &user.Email, &user.Password, &user.Created)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
