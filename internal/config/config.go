package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

const defConfigPath = "../configs/config.yaml"

type server struct {
	Port int `yaml:"port"`
}

type logger struct {
	Level string `yaml:"level"`
}

type dataSource struct {
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	DB                 string `yaml:"db"`
	MaxConnection      int    `yaml:"max_connection"`
	MinConnection      int    `yaml:"min_connection"`
	ConnectionLifeTime string `yaml:"connection_life_time"`
	User               string
	Password           string
}

type Config struct {
	Server     server     `yaml:"server"`
	DataSource dataSource `yaml:"data_source"`
	Logger     logger     `yaml:"logger"`
}

func New() *Config {

	cfgFilePath := os.Getenv("ISTOK_AUTHORIZATION_CONFIG_PATH")
	if len(cfgFilePath) < 2 {
		cfgFilePath = defConfigPath
		log.Printf("config file's path not set in env: path=%s", cfgFilePath)
	}

	file, err := os.ReadFile(cfgFilePath)
	if err != nil {
		log.Fatalf("can't read config file: error=%s", err.Error())
	}

	cfg := &Config{}
	cfg.DataSource.User = os.Getenv("DB_USER")
	cfg.DataSource.Password = os.Getenv("DB_PASSWORD")

	err = yaml.Unmarshal(file, cfg)
	if err != nil {
		log.Printf("can't parse yaml file: error=%s", err.Error())
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
