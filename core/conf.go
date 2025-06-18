package core

import (
	"CollabDoc-go/config"
	"CollabDoc-go/utils"
	"log"

	"gopkg.in/yaml.v3"
)

// InitConf 从 YAML 文件加载配置
func InitConf() *config.Config {
	c := &config.Config{}
	yamlConf, err := utils.LoadYAML()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML configuration: %v", err)
	}
	return c
}
