package main

import (
	"log"

	"github.com/eduardolima806/my-chat-server/config"
	"github.com/eduardolima806/my-chat-server/internal/app"
)

func main() {
	cgf, err := config.NewConfig()

	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cgf)
}
