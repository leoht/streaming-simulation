package producer

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

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

var kafkaProducerInChannel chan Event = make(chan Event)
var userSimulations []*UserSimulation = []*UserSimulation{}
var userIds []string

func GetSimulations() []*UserSimulation {
	return userSimulations
}

// Start random user simulations and record produced events into
// the PostgresSQL database (TODO)
func Start() {

	userIds = readUserIds()

	fmt.Println("Starting Kafka producer...")

	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		// User-specific properties that you must set
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVER_URL"),

		// Fixed properties
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
		case event := <-kafkaProducerInChannel:
			topic := os.Getenv("TOPIC_NAME")
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

// Start a simulation for a user who doesn't have a simulation started yet
func StartNewSimulation() *UserSimulation {
	userId := userIds[rand.Intn(len(userIds))]

	fmt.Println("Starting user simulation...")

	simulation := NewUserSimulation(userId)
	simulation.Start([]string{"sign_in", "view_page"})
	userSimulations = append(userSimulations, &simulation)

	go func() {
		for {
			select {
			case event := <-simulation.outgoingEvents:
				kafkaProducerInChannel <- event
			}
		}
	}()

	return &simulation
}

func StopSimulationForUser(userId string) *UserSimulation {
	for _, simulation := range userSimulations {
		if simulation.UserId == userId && simulation.Running {
			simulation.Stop()
			return simulation
		}
	}

	return nil
}

func ResumeSimulationForUser(userId string) *UserSimulation {
	for _, simulation := range userSimulations {
		if simulation.UserId == userId && !simulation.Running {
			simulation.Resume()
			return simulation
		}
	}

	return nil
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

func readUserIds() []string {
	contents, err := os.ReadFile("users.txt")
	if err != nil {
		log.Fatal("could not read users.txt ")
	}

	return splitLines(string(contents))
}

func splitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}
