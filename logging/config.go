package logging

type ConfigFileLogger struct {
	DirPath   string `yaml:"directory_path" envconfig:"directory_path" default:"log"`
	MaxSize   int    `yaml:"max_size" envconfig:"max_size"`
	MaxRotate int    `yaml:"max_rotate" envconfig:"max_rotate"`
	Gzip      bool   `yaml:"gzip" envconfig:"gzip"`
}

type Config struct {
	Debug      bool             `yaml:"debug" envconfig:"debug"`
	FileMode   bool             `yaml:"file_mode" envconfig:"file_mode"`
	StdMode    bool             `yaml:"std_mode" envconfig:"std_mode"`
	FileConfig ConfigFileLogger `yaml:"file_mode_settings" envconfig:"file_mode_settings"`
}
