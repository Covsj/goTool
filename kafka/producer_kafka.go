package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

type KafkaProducer struct {
	asyncProducer sarama.AsyncProducer
	topic         string
}

func NewAsyncKafkaProducer(topic string, kafkaBrokers []string) (*KafkaProducer, error) {
	kafkaProducer := &KafkaProducer{}
	// producer config
	config := sarama.NewConfig()
	// 本地commit成功返回
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond
	config.Producer.MaxMessageBytes = 20971520

	producer, err := sarama.NewAsyncProducer(kafkaBrokers, config)
	if err != nil {
		return nil, err
	}

	kafkaProducer.asyncProducer = producer
	kafkaProducer.topic = topic

	go func() {
		for err := range kafkaProducer.asyncProducer.Errors() {
			fmt.Println("[NewAsyncKafkaProducer] kafka producer error:", err)
		}
	}()

	return kafkaProducer, nil
}

func (p *KafkaProducer) AsyncSendMessage(key string, value string) {

	p.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}
}
