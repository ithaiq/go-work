package mq

type MqCli interface {
	Push(topic string, key string, body []byte) error
	Receive(topic string, key string) ([]byte, error)
}
