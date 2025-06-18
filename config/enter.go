package config

type Config struct {
	Email  Email  `json:"email" yaml:"email"`
	Jwt    Jwt    `json:"jwt" yaml:"jwt"`
	Mysql  Mysql  `json:"mysql" yaml:"mysql"`
	Redis  Redis  `json:"redis" yaml:"redis"`
	System System `json:"system" yaml:"system"`
	Zap    Zap    `json:"zap" yaml:"zap"`
	Mongo  Mongo  `json:"mongo" yaml:"mongo"`
}
