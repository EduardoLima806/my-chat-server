package db

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/eduardolima806/my-chat-server/config"
	_ "github.com/lib/pq"
)

var lock = &sync.Mutex{}

var connection *sql.DB

func ConnectToPostgresDb(dbConfig config.PG) (*sql.DB, error) {
	var err error
	if connection == nil {
		lock.Lock()
		defer lock.Unlock()
		connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
		connection, err = sql.Open("postgres", connectionString)
		if err != nil {
			return nil, err
		}
	}
	return connection, nil
}
