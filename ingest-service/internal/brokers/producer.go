package brokers

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type Producer interface {
	SendMessage(topic string, key string, message []byte) error
	Close() error
}

type ProducerConfig struct {
	Brokers          []string
	Version          sarama.KafkaVersion
	RequiredAcks     sarama.RequiredAcks
	Idempotent       bool
	CompressionCodec sarama.CompressionCodec
	FlushFrequency   time.Duration
	FlushBytes       int
	FlushMessages    int
	MaxRetries       int
	RetryBackoff     time.Duration
}

func DefaultProducerConfig(brokers []string) ProducerConfig {
	return ProducerConfig{
		Brokers:          brokers,
		Version:          sarama.V3_6_0_0,
		RequiredAcks:     sarama.WaitForAll,
		Idempotent:       true,
		CompressionCodec: sarama.CompressionSnappy,
		FlushFrequency:   5 * time.Millisecond,
		FlushBytes:       64 * 1024,
		FlushMessages:    100,
		MaxRetries:       10,
		RetryBackoff:     200 * time.Millisecond,
	}
}

type producer struct {
	ap      sarama.AsyncProducer
	wg      *sync.WaitGroup
	closeCh chan struct{}
}

func NewProducer(cfg ProducerConfig) (Producer, error) {
	saramaCfg := sarama.NewConfig()

	saramaCfg.Version = cfg.Version

	saramaCfg.Producer.RequiredAcks = cfg.RequiredAcks
	saramaCfg.Producer.Idempotent = cfg.Idempotent
	saramaCfg.Producer.Retry.Max = cfg.MaxRetries
	saramaCfg.Producer.Retry.Backoff = cfg.RetryBackoff

	if cfg.Idempotent {
		saramaCfg.Net.MaxOpenRequests = 1
	}

	saramaCfg.Producer.Compression = cfg.CompressionCodec
	saramaCfg.Producer.Flush.Frequency = cfg.FlushFrequency
	saramaCfg.Producer.Flush.Bytes = cfg.FlushBytes
	saramaCfg.Producer.Flush.Messages = cfg.FlushMessages

	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Producer.Return.Errors = true

	ap, err := sarama.NewAsyncProducer(cfg.Brokers, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create async producer: %w", err)
	}

	p := &producer{
		ap:      ap,
		closeCh: make(chan struct{}),
	}

	p.wg.Add(1)
	go p.drain()

	return p, nil

}

func (p *producer) drain() {
	defer p.wg.Done()
	for {
		select {
		case msg, ok := <-p.ap.Successes():
			if !ok {
				return
			}
			log.Printf("[kafka] delivered: topic=%s partition=%d offset=%d",
				msg.Topic, msg.Partition, msg.Offset)

		case err, ok := <-p.ap.Errors():
			if !ok {
				return
			}
			log.Printf("[kafka] delivery failed: topic=%s err=%v",
				err.Msg.Topic, err.Err)
		}
	}
}

func (p *producer) SendMessage(topic string, key string, message []byte) error {
	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	if key != "" {
		msg.Key = sarama.StringEncoder(key)
	}
	select {
	case p.ap.Input() <- &msg:
		return nil
	case <-p.closeCh:
		return fmt.Errorf("kafka producer closed")
	}

}

func (p *producer) Close() error {
	close(p.closeCh)
	p.ap.AsyncClose()
	p.wg.Wait()
	return nil
}
