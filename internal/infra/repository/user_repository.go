package repository

import (
	"database/sql"

	"github.com/eduardolima806/my-chat-server/internal/domain"
)

const IdError = -1

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

func (userRepo *UserRepository) GetUserByID(id int32) (*domain.User, error) {
	// TODO: Implements this method
	return nil, nil
}
