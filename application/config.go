package application

import (
	"os"
	"strconv"
)

type Config struct {
	// Listening port
	ServerPort uint16
	//Config for redis database (non-relational)
	RedisAddress string
	//Config for MySQL database (relational)
	User              string
	Password          string
	MySqlAddress      string
	MySqlDatabaseName string
}

func LoadConfig() Config {
	cfg := Config{
		ServerPort:        3000,
		RedisAddress:      "localhost:6379",
		User:              "user",
		Password:          "password",
		MySqlAddress:      "localhost:3306",
		MySqlDatabaseName: "SystemFiles",
	}

	if serverPort, exists := os.LookupEnv("SERVER_PORT"); exists {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
		}
	}

	if redisAddr, exists := os.LookupEnv("REDIS_ADDR"); exists {
		cfg.RedisAddress = redisAddr
	}

	if user, exists := os.LookupEnv("SQL_USER"); exists {
		cfg.User = user
	}

	if password, exists := os.LookupEnv("SQL_PASSWORD"); exists {
		cfg.Password = password
	}

	if MySqlAddress, exists := os.LookupEnv("MYSQL_ADDRESS"); exists {
		cfg.MySqlAddress = MySqlAddress
	}

	if MySqlDatabaseName, exists := os.LookupEnv("DATABASE_NAME"); exists {
		cfg.MySqlDatabaseName = MySqlDatabaseName
	}

	return cfg
}
