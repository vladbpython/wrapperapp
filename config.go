package wrapperapp

import "github.com/vladbpython/wrapperapp/logging"

type ConfigSystem struct {
	AppName string         `yaml:"name" envconfig:"name" default:"WrapperApplication"`
	Debug   uint8          `yaml:"debug" envconfig:"debug"`
	Logger  logging.Config `yaml:"logging"`
}

type ConfigWrapper struct {
	System ConfigSystem `yaml:"system"`
}

type WrapperStructConfig struct {
	MaxRetries uint `yaml:"max_retries" envconfig:"max_retries"`
	MaxWait    uint `yaml:"max_wait" envconfig:"max_wait"`
}
