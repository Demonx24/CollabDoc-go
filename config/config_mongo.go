package config

type Mongo struct {
	Host       string `mapstructure:"host" json:"host" yaml:"host"`
	Port       int    `mapstructure:"port" json:"port" yaml:"port"`
	Username   string `mapstructure:"username" json:"username" yaml:"username"`
	Password   string `mapstructure:"password" json:"password" yaml:"password"`
	Database   string `mapstructure:"database" json:"database" yaml:"database"`
	AuthSource string `mapstructure:"authSource" json:"authSource" yaml:"authSource"`
	ReplicaSet string `mapstructure:"replicaSet" json:"replicaSet" yaml:"replicaSet"`
	SSL        bool   `mapstructure:"ssl" json:"ssl" yaml:"ssl"`
	Enabled    bool   `mapstructure:"enabled" json:"enabled" yaml:"enabled"` // 建议加上开关字段
}
