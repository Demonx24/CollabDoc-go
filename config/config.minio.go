package config

type Minio struct {
	Endpoint        string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	AccessKeyID     string `mapstructure:"accessKeyID" json:"accessKeyID" yaml:"accessKeyID"`
	SecretAccessKey string `mapstructure:"secretAccessKey" json:"secretAccessKey" yaml:"secretAccessKey"`
	UseSSL          bool   `mapstructure:"useSSL" json:"useSSL" yaml:"useSSL"`
	Bucket          string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
}
