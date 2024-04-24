package app

import (
	"database/sql"
	"fmt"

	"github.com/eduardolima806/my-chat-server/config"
	v1 "github.com/eduardolima806/my-chat-server/internal/controller/http/v1"
	"github.com/eduardolima806/my-chat-server/internal/infra/db"
	"github.com/eduardolima806/my-chat-server/internal/infra/repository"
	"github.com/eduardolima806/my-chat-server/internal/usecase/user_usecase"
	"github.com/eduardolima806/my-chat-server/internal/util"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {

	fmt.Printf("Running %s %s\n", cfg.App.Name, cfg.App.Version)

	handler := gin.Default()
	conn, err := db.ConnectToPostgresDb(cfg.PG)

	if err != nil {
		fmt.Errorf("failed to connect to database %w", err)
	}

	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			panic("ERROR CLOSING POSTGRES CONNECTION")
		}
	}(conn)

	userRepo := repository.NewUserRepository(conn)
	passwordHasher := &util.DefaultPasswordHasher{}
	userUseCase := user_usecase.NewUserBaseUserCase(userRepo, passwordHasher)
	v1.NewRouter(handler, *userUseCase)
	// TODO: Should implements in pkg/httpserver ?
	handler.Run()
}
