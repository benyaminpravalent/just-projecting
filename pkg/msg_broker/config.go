package msgbroker

type KafkaConfig struct {
	ProducerReturnSuccess bool   `json:"producerReturnSuccess"`
	WriteTimeout          int    `json:"writeTimeout"`
	MaxRetry              int    `json:"maxRetry"`
	Version               string `json:"version"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	UrlKafkaList          string `json:"urlKafkaList"`
	ConsumerGroup         string `json:"consumerGroup"`
	ExampleTopic          string `json:"exampleTopic"`
}
