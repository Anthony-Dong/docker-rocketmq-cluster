package config

import "time"

type RocketMqConfig struct {
	Host       []string `json:"host"`
	AccessKey  string   `json:"access_key"`
	SecretKey  string   `json:"secret_key"`
	RetryTimes int      `json:"retry_times"`
	GroupName  string   `json:"group_name"` // 默认值 DEFAULT_CONSUMER
}

type RocketMqConsumer struct {
	RocketMqConfig
	Topic string `json:"topic"`
}
type RocketMqProducer struct {
	RocketMqConfig
	Topic          string        `json:"topic"`
	SendMsgTimeout time.Duration `json:"send_msg_timeout"`
}
