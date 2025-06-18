package config

type Mongo struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Database   string `yaml:"database"`
	AuthSource string `yaml:"authSource"`
	ReplicaSet string `yaml:"replicaSet"`
	SSL        bool   `yaml:"ssl"`
}
