package wrapperapp

import (
	"github.com/vladbpython/wrapperapp/monitoring"
)

type ConfigLogger struct {
	DirPath   string `yaml:"directory_path"`
	MaxSize   int    `yaml:"max_size"`
	MaxRotate int    `yaml:"max_rotate"`
	Gzip      bool   `yaml:"gzip"`
	StdMode   bool   `yaml:"std_mode"`
}
type ConfigSystem struct {
	AppName    string                      `yaml:"name"`
	Debug      uint8                       `yaml:"debug"`
	Logger     ConfigLogger                `yaml:"logging"`
	Monitoring monitoring.ConfigMinotiring `yaml:"monitoring"`
}

type ConfigWrapper struct {
	System ConfigSystem `yaml:"system"`
}
