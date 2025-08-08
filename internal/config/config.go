package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const defConfigPath = "../configs/config.yaml"

type server struct {
	Port int `env:"ISTOK_AUTH_SERVER_PORT"`
}

type logger struct {
	Level string `env:"ISTOK_AUTH_LOG_LVL"`
}

type dataSource struct {
	Host               string `env:"ISTOK_AUTH_DB_HOST"`
	Port               int    `env:"ISTOK_AUTH_DB_PORT"`
	Name               string `env:"ISTOK_AUTH_DB_NAME"`
	MaxConnection      int    `env:"ISTOK_AUTH_DB_MAX_CONN"`
	MinConnection      int    `env:"ISTOK_AUTH_DB_MIN_CONN"`
	ConnectionLifeTime string `env:"ISTOK_AUTH_DB_CONN_LIFETIME"`
	User               string `env:"ISTOK_AUTH_DB_USER"`
	Password           string `env:"ISTOK_AUTH_DB_PASSWORD"`
}

type Config struct {
	Server     server
	DataSource dataSource
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

func (c Config) ServerPort() string {
	defPort := ":8080"
	if c.Server.Port < 1024 || c.Server.Port > 65535 {
		log.Printf("invalid port value=%v, will use the default value=%s", c.Server.Port, defPort)
		return defPort
	}
	return fmt.Sprintf(":%s", strconv.Itoa(c.Server.Port))
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

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", d.User, d.Password, d.Host, port, d.DB)
}
