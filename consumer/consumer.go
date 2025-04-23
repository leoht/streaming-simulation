package consumer

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-msk-iam-sasl-signer-go/signer"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
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

func NewConsumer(receiveChannel chan string, sigchan chan os.Signal) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVER_URL"),

		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "OAUTHBEARER",
		"group.id":          "analytics",
		"auto.offset.reset": "earliest"})

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	bearerToken := createToken()

	err = c.SetOAuthBearerToken(bearerToken)
	if err != nil {
		panic(err)
	}

	topic := "events"
	err = c.SubscribeTopics([]string{topic}, nil)

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			// fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
			// 	*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))

			// event := producer.Event{}
			// err = json.Unmarshal([]byte(ev.Value), &event)

			if err != nil {
				log.Fatal(err)
			}

			receiveChannel <- string(ev.Value)
		}
	}

	c.Close()

}
