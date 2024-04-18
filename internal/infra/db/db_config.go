package db

type DBConfig struct {
	Host         string
	Port         string
	UserName     string
	Password     string
	DatabaseName string
}

func CreateDbConfigFromEnv(dbEnv EnvDBConfig) *DBConfig {
	return &DBConfig{
		Host:         dbEnv.GetHost(),
		Port:         dbEnv.GetPort(),
		UserName:     dbEnv.GetUserName(),
		Password:     dbEnv.GetPassword(),
		DatabaseName: dbEnv.GetDatabase(),
	}
}
