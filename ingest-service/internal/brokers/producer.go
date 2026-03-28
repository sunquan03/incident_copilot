package brokers

type Producer interface {
	SendMessage(topic string, message []byte) error
	Close() error
}
