package producer

import (
	"encoding/json"
	"fmt"
	"log"

	"leohetsch.com/simulation/simulation"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Kafka (MSK) implementation

type KafkaClient struct {
	topicName string
	producer  KafkaProducer
}

type KafkaProducer interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
	SetOAuthBearerToken(token kafka.OAuthBearerToken) error
	Events() chan kafka.Event
}

func NewKafkaClient(topicName string, kafkaProducer KafkaProducer) KafkaClient {
	bearerToken := createToken()
	err := kafkaProducer.SetOAuthBearerToken(bearerToken)
	if err != nil {
		panic(err)
	}

	return KafkaClient{topicName, kafkaProducer}
}

func (c KafkaClient) Start(producerInChannel chan simulation.Event) {

	// Go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go monitorEventsFromKafka(c.producer)

	waitForEvents(c.topicName, c.producer, producerInChannel)
}

func waitForEvents(topic string, kafkaProducer KafkaProducer, producerInChannel <-chan simulation.Event) {
	for {
		select {
		case event := <-producerInChannel:
			log.Printf("Attempting to send %s to topic for user %s...\n", event.EventName, event.UserId)

			key := event.Id
			data, _ := json.Marshal(event)
			kafkaProducer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Key:            []byte(key),
				Value:          data,
			}, nil)
		}

	}
}

func monitorEventsFromKafka(p KafkaProducer) {
	for e := range p.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
			} else {
				fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
					*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
			}
		}
	}
}
