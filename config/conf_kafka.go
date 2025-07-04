package config

type Kafka struct {
	Enabled       bool     `yaml:"enabled"`        // 是否启用 Kafka
	Brokers       []string `yaml:"brokers"`        // Kafka 集群地址列表
	TopicPrefix   string   `yaml:"topic_prefix"`   // Topic 前缀
	DiffTopic     string   `yaml:"diff_topic"`     // 用于文档差异推送的 Topic
	ConsumerGroup string   `yaml:"consumer_group"` // 消费者组
	Username      string   `yaml:"username"`       // Kafka 用户名（如启用 SASL）
	Password      string   `yaml:"password"`       // Kafka 密码（如启用 SASL）
}
