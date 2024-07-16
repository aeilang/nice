package configs

import (
	"os"
	"strconv"

	"github.com/aeilang/nice/utils"
	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string `validate:"required"`
	Port       string `validate:"required"`
	DBUser     string `validate:"required"`
	DBPassword string `validate:"required"`
	DBName     string `validate:"required"`

	JWTAccessSecret              string `validate:"required"`
	JWTAccessExperationInMinites int    `validate:"gte=0"`

	JWTRefreshSecret            string `validate:"required"`
	JWTRefreshExperationInHours int    `validate:"gte=0"`

	RedisAddr string `validate:"required"`

	MailHost     string `validate:"required"`
	MailPort     int    `validate:"gte=0"`
	MailUsername string `validate:"required,email"`
	MailPassword string `validate:"required"`
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	config := Config{
		PublicHost: getEnv("PUBLIC_HOST"),
		Port:       getEnv("PORT"),
		DBUser:     getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		DBName:     getEnv("DB_NAME"),

		JWTAccessSecret:              getEnv("JWT_ACCESS_SECRET"),
		JWTAccessExperationInMinites: getEnvAsInt("JWT_ACCESS_EXPIRATION_IN_MINITES"),

		JWTRefreshSecret:            getEnv("JWT_REFRESH_SECRET"),
		JWTRefreshExperationInHours: getEnvAsInt("JWT_REFRESH_EXPIRATION_IN_HOURS"),

		RedisAddr: getEnv("REDIS_ADDR"),

		MailHost:     getEnv("MAIL_HOST"),
		MailPort:     getEnvAsInt("MAIL_PORT"),
		MailUsername: getEnv("MAIL_USERNAME"),
		MailPassword: getEnv("MAIL_PASSWORD"),
	}

	if err := utils.Validate.Struct(config); err != nil {
		panic(err)
	}

	return config
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return ""
}

func getEnvAsInt(key string) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return -1
		}
		return int(i)
	}

	return -1
}
