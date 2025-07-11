package config

type Config struct {
	Email       Email       `json:"email" yaml:"email"`
	Jwt         Jwt         `json:"jwt" yaml:"jwt"`
	Mysql       Mysql       `json:"mysql" yaml:"mysql"`
	Redis       Redis       `json:"redis" yaml:"redis"`
	System      System      `json:"system" yaml:"system"`
	Zap         Zap         `json:"zap" yaml:"zap"`
	Mongo       Mongo       `json:"mongodb" yaml:"mongodb"`
	Upload      Upload      `json:"upload" yaml:"upload"`
	Captcha     Captcha     `json:"captcha" yaml:"captcha"`
	EmailGoogle EmailGoogle `json:"email_google" yaml:"email_google"`
	Website     Website     `json:"website" yaml:"website"`
	Minio       Minio       `json:"minio" yaml:"MinIO"`
	Kafka       Kafka       `json:"kafka" yaml:"kafka"`
}
