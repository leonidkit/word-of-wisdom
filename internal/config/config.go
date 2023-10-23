package config

import "time"

type ServerConfig struct {
	Env          string        `env:"ENVIRONMENT,notEmpty"`
	Addr         string        `env:"SERVER_ADDR,notEmpty"`
	LogLevel     string        `env:"LOGGER_LEVEL,notEmpty"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT"`
	KeepAlive    time.Duration `env:"SERVER_KEEP_ALIVE"`
}

func (sc *ServerConfig) IsProduction() bool {
	return sc.Env == "prod"
}
