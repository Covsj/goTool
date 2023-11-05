package gotool_kafka

import (
	"strings"
)

type MyProducer struct {
	Producers map[string]*producer
}

func NewMyProducer(kafkaBrokers string, topics map[string]string) (*MyProducer, error) {
	brokersList := strings.Split(kafkaBrokers, ",")
	myProducer := &MyProducer{}
	myProducer.Producers = map[string]*producer{}
	for topicIndex, topicName := range topics {
		topic, err := newAsyncKafkaProducer(topicName, brokersList)
		if err != nil {
			return myProducer, err
		}
		myProducer.Producers[topicIndex] = topic
	}
	return myProducer, nil
}

func (client *MyProducer) SendMessage(topicIndex string, key string, msg string) {
	producer, exist := client.Producers[topicIndex]
	if !exist {
		return
	}
	producer.SyncSendMessage(key, msg)
}
