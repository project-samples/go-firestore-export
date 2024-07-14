package app

import "github.com/core-go/log/zap"

type Config struct {
	Log         log.Config `mapstructure:"log"`
	Credentials string     `mapstructure:"credentials"`
}
