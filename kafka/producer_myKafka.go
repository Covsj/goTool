package kafka

import (
	"strings"
)

type MyProducer struct {
	producers map[string]*KafkaProducer
}

func NewMyProducer(kafkaBrokers string, topics map[string]string) (*MyProducer, error) {
	brokersList := strings.Split(kafkaBrokers, ",")
	myProducer := &MyProducer{}
	myProducer.producers = map[string]*KafkaProducer{}
	for topicIndex, topicName := range topics {
		topic, err := NewAsyncKafkaProducer(topicName, brokersList)
		if err != nil {
			return myProducer, err
		}
		myProducer.producers[topicIndex] = topic
	}
	return myProducer, nil
}

func (client *MyProducer) SendMessage(topicIndex string, key string, msg string) {
	producer, exist := client.producers[topicIndex]
	if !exist {
		return
	}
	producer.AsyncSendMessage(key, msg)
}
