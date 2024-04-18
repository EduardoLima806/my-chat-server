package db

import "os"

type EnvDBConfig struct {
	host         string
	port         string
	userName     string
	password     string
	databaseName string
}

const (
	DB_HOST      = "DB_HOST"
	DB_PORT      = "DB_PORT"
	DB_USER_NAME = "DB_USERNAME"
	DB_PASSWORD  = "DB_PASSWORD"
	DB_DATABASE  = "DB_DATABASE"
)

func NewEnvDBConfig() *EnvDBConfig {
	return &EnvDBConfig{
		host:         os.Getenv(DB_HOST),
		port:         os.Getenv(DB_PORT),
		userName:     os.Getenv(DB_USER_NAME),
		password:     os.Getenv(DB_PASSWORD),
		databaseName: os.Getenv(DB_DATABASE),
	}
}

func (dbEnv *EnvDBConfig) GetHost() string {
	return dbEnv.host
}

func (dbEnv *EnvDBConfig) GetPort() string {
	return dbEnv.port
}

func (dbEnv *EnvDBConfig) GetUserName() string {
	return dbEnv.userName
}

func (dbEnv *EnvDBConfig) GetPassword() string {
	return dbEnv.password
}

func (dbEnv *EnvDBConfig) GetDatabase() string {
	return dbEnv.databaseName
}
