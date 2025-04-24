package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"leohetsch.com/simulation/simulation"

	"github.com/aws/aws-msk-iam-sasl-signer-go/signer"
	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func createToken() kafka.OAuthBearerToken {
	token, tokenExpirationTime, err := signer.GenerateAuthToken(context.TODO(), os.Getenv("AWS_REGION"))
	if err != nil {
		panic(err)
	}
	seconds := tokenExpirationTime / 1000
	nanoseconds := (tokenExpirationTime % 1000) * 1000000
	bearerToken := kafka.OAuthBearerToken{
		TokenValue: token,
		Expiration: time.Unix(seconds, nanoseconds),
	}

	return bearerToken
}

var userIds []string

// Start random user simulations and record produced events into
// the PostgresSQL database (TODO)
func Start(producerInChannel chan simulation.Event) {

	fmt.Println("Starting Kafka producer...")

	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVER_URL"),

		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "OAUTHBEARER",
		"client.id":         "simulation-producer",
		"acks":              "all",
	})

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	bearerToken := createToken()

	err = kafkaProducer.SetOAuthBearerToken(bearerToken)
	if err != nil {
		panic(err)
	}

	// Go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go monitorEventsFromKafka(kafkaProducer)

	for {
		select {
		case event := <-producerInChannel:
			topic := os.Getenv("TOPIC_NAME")
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

func monitorEventsFromKafka(p *kafka.Producer) {
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
