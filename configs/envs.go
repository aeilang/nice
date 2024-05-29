package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBName                 string
	JWTSecret              string
	JWTExperationInMinites int

	MaxConns                   int
	MinConns                   int
	MaxConnLifeTimeInMinites   int
	MaxConnIdleTimeInMinites   int
	HealthCheckPeriodInMinites int
	ConnectTimeoutInSeconds    int
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "8080"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", "password"),
		DBName:                 getEnv("DB_NAME", "test_db"),
		JWTSecret:              getEnv("JWT_SECRET", "not-so-secret-now-is-it?"),
		JWTExperationInMinites: getEnvAsInt("JWT_EXPIRATION_IN_MIMITES", 3600*24*7),

		MaxConns:                   getEnvAsInt("MAX_CONNS", 4),
		MinConns:                   getEnvAsInt("MIN_CONNS", 1),
		MaxConnLifeTimeInMinites:   getEnvAsInt("MAX_CONN_LIFE_TIME_IN_MINITES", 60),
		MaxConnIdleTimeInMinites:   getEnvAsInt("MAX_CONN_IDLE_TIME_IN_MINITES", 30),
		HealthCheckPeriodInMinites: getEnvAsInt("HEALTH_CHECK_PERIOD_IN_MINITES", 1),
		ConnectTimeoutInSeconds:    getEnvAsInt("CONNECT_TIME_OUT_IN_MINITES", 5),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return int(i)
	}

	return fallback
}

func DBConfig() (*pgxpool.Config, error) {
	env := initConfig()
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.DBUser, env.DBPassword, env.PublicHost, env.Port, env.DBName)

	dbConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	dbConfig.MaxConns = int32(env.MaxConns)
	dbConfig.MinConns = int32(env.MinConns)
	dbConfig.MaxConnLifetime = time.Duration(env.MaxConnLifeTimeInMinites) * time.Minute
	dbConfig.MaxConnIdleTime = time.Duration(env.MaxConnIdleTimeInMinites) * time.Minute
	dbConfig.HealthCheckPeriod = time.Duration(env.HealthCheckPeriodInMinites) * time.Minute
	dbConfig.ConnConfig.ConnectTimeout = time.Duration(env.ConnectTimeoutInSeconds) * time.Second

	return dbConfig, nil
}
