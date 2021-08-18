package tools

import (
	"io/ioutil"
	"log"

	"github.com/goccy/go-yaml"
)

//Парсер Yaml конфига
func LoadYamlConfig(filePath string, config interface{}) {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("[ERROR]: Config file: %v ", err)
	}

	errParser := yaml.Unmarshal(yamlFile, config)
	if errParser != nil {
		log.Fatalf("[ERROR]: On parse config file: #%v ", err)
	}

}
