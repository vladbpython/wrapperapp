package wrapperapp

import (
	"github.com/vladbpython/wrapperapp/monitoring"
)

type ConfigLogger struct {
	DirPath   string `yaml:"directory_path" envconfig:"directory_path" default:"/var/log/WrapperApplication"`
	MaxSize   int    `yaml:"max_size" envconfig:"max_size"`
	MaxRotate int    `yaml:"max_rotate" envconfig:"max_rotate"`
	Gzip      bool   `yaml:"gzip" envconfig:"gzip"`
	StdMode   bool   `yaml:"std_mode" envconfig:"std_mode"`
}
type ConfigSystem struct {
	AppName    string                      `yaml:"name" envconfig:"name" default:"WrapperApplication"`
	Debug      uint8                       `yaml:"debug" envconfig:"debug"`
	Logger     ConfigLogger                `yaml:"logging"`
	Monitoring monitoring.ConfigMinotiring `yaml:"monitoring"`
}

type ConfigWrapper struct {
	System ConfigSystem `yaml:"system"`
}
