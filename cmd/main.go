package main

import (
	"fmt"
	"log"

	"github.com/eduardolima806/my-chat-server/internal/infra/db"
	"github.com/eduardolima806/my-chat-server/internal/infra/repository"
	"github.com/eduardolima806/my-chat-server/internal/usecase"
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
	createUserUc := usecase.NewCreateUserUseCase(userRepo)
	input := usecase.UserInput{UserName: "eduardolima806", DisplayName: "Eduardo Lima", Email: "eduardolima.dev.io@gmail.com", Password: "P4$$w0rd001"}
	output, err := createUserUc.Execute(input)
	if err != nil {
		log.Fatalf("Error to create user: %s. %s", input.UserName, err)
	}
	fmt.Printf("User %s created with the id %d.", input.UserName, output.CreatedUserId)
}

func loadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
