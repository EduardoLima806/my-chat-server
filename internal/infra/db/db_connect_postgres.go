package db

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

var lock = &sync.Mutex{}

var connection *sql.DB

func ConnectToPostgresDb(dbConfig DBConfig) (*sql.DB, error) {
	var err error
	if connection == nil {
		lock.Lock()
		defer lock.Unlock()
		connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbConfig.UserName, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DatabaseName)
		connection, err = sql.Open("postgres", connectionString)
		if err != nil {
			return nil, err
		}
	}
	return connection, nil
}
