package config

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type token struct {
	Secret    string `env:"ISTOK_AUTH_TOKEN_SECRET" env-required`
	RTokenTTL string `env:"ISTOK_AUTH_TOKEN_REFRESH_TTL" env-default:"5m"`
	ATokenTTL string `env:"ISTOK_AUTH_TOKEN_ACCESS_TTL" env-default:"1M"`
}

type server struct {
	Port int `env:"ISTOK_AUTH_SERVER_PORT" env-required`
}

type logger struct {
	Level string `env:"ISTOK_AUTH_LOG_LVL" env-default:"INFO"`
}

type dataSource struct {
	Host               string `env:"ISTOK_AUTH_DB_HOST" env-required`
	Port               int    `env:"ISTOK_AUTH_DB_PORT" env-required`
	Name               string `env:"ISTOK_AUTH_DB_NAME" env-required`
	MaxConnection      int    `env:"ISTOK_AUTH_DB_MAX_CONN" env-default:"20"`
	MinConnection      int    `env:"ISTOK_AUTH_DB_MIN_CONN" env-default:"5"`
	ConnectionLifeTime string `env:"ISTOK_AUTH_DB_CONN_LIFETIME" env-default:"1h"`
	User               string `env:"ISTOK_AUTH_DB_USER" env-required`
	Password           string `env:"ISTOK_AUTH_DB_PASSWORD" env-required`
}

type Config struct {
	Server     server
	DataSource dataSource
	Token      token
	Logger     logger
}

func New() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("error reading configuration from virtual environment (env): %s", err.Error())
	}
	return cfg
}

func (s server) ServerPort() string {
	defPort := ":8080"
	if s.Port < 1024 || s.Port > 65535 {
		log.Printf("invalid port value=%v, will use the default value=%s", s.Port, defPort)
		return defPort
	}
	return fmt.Sprintf(":%s", strconv.Itoa(s.Port))
}

func (d dataSource) LifeTime() time.Duration {
	t, err := time.ParseDuration(d.ConnectionLifeTime)
	if err != nil {
		log.Printf("invalid connection_life_time value=%s, will use the default value=%s", d.ConnectionLifeTime, "1h")
		return time.Hour
	}
	return t
}

func (d dataSource) PostgresConnString() string {
	port := strconv.Itoa(d.Port)
	if d.Port < 1024 || d.Port > 65535 {
		port = "5432"
		log.Printf("invalid postgres port value=%v, will use the default value=%s", d.Port, port)
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", d.User, d.Password, d.Host, port, d.Name)
}
