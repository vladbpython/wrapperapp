package tools

import (
	"io/ioutil"
	"log"

	"github.com/goccy/go-yaml"
	"github.com/kelseyhightower/envconfig"
	"github.com/subosito/gotenv"
)

//Парсер Yaml конфига
func LoadYamlConfig(filePath string, config interface{}) {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("[ERROR]: Read yaml config file: %v ", err)
	}

	errParser := yaml.Unmarshal(yamlFile, config)
	if errParser != nil {
		log.Fatalf("[ERROR]: On parse yaml config file: #%v ", err)
	}

}

func ReadEnvConfig(filePath string) {
	err := gotenv.Load(filePath)
	if err != nil {
		log.Fatalf("[ERROR]: Read enviroment config file: %v ", err)
	}
}

func ParseEnvConfig(section_prefix string, config interface{}) {
	err := envconfig.Process(section_prefix, config)
	if err != nil {
		log.Fatalf("[ERROR]: On parse enviroment config file: #%v ", err)
	}
}
