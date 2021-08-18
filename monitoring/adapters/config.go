package adapters

type ConfigAdapter struct {
	Adapter   string        `yaml:"adapter"`
	Host      string        `yaml:"host"`
	Token     string        `yaml:"token"`
	Merchants []interface{} `yaml:"merchants"`
}
