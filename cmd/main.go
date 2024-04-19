package main

import (
	"fmt"
	"log"

	"github.com/eduardolima806/my-chat-server/internal/infra/db"
	"github.com/eduardolima806/my-chat-server/internal/infra/repository"
	"github.com/eduardolima806/my-chat-server/internal/usecase"
	"github.com/eduardolima806/my-chat-server/internal/util"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello My Chat Server!")

	loadEnvFile()

	dbConfig := db.CreateDbConfigFromEnv(*db.NewEnvDBConfig())
	conn, err := db.ConnectToPostgresDb(*dbConfig)
	if err != nil {
		log.Fatal("Error to connect to database!", err)
	}
	defer conn.Close()

	userRepo := repository.NewUserRepository(conn)
	passwordHasher := &util.DefaultPasswordHasher{}
	createUserUc := usecase.NewCreateUserUseCase(userRepo, passwordHasher)
	input := usecase.UserInput{UserName: "eduardolima806", DisplayName: "Eduardo Lima", Email: "eduardolima.dev.io@gmail.com", Password: "P4$$w0rd001"}
	output, err := createUserUc.Execute(input)
	if err != nil {
		fmt.Printf("Error to create user: %s. %s\n", input.UserName, err)
	} else {
		fmt.Printf("User %s created with the id %d.\n", input.UserName, output.CreatedUserId)
	}
	loginInput := usecase.LoginInput{Login: "eduardolima806", Password: "P4$$w0rd001"}
	loginUserUc := usecase.NewLoginUserUseCase(userRepo, passwordHasher)
	loginOutput, err := loginUserUc.Execute(loginInput)
	if err != nil {
		log.Fatalf("Error to login user: %s. %s", loginInput.Login, err)
	}
	if loginOutput.IsSucceed {
		fmt.Printf("Succeed login for user %s", loginInput.Login)
	} else {
		fmt.Printf("Error to login for user %s. Err: %d", loginInput.Login, loginOutput.ErrorType)
	}
}

func loadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
