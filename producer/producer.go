package producer

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var userSimulations []*UserSimulation = []*UserSimulation{}

func GetSimulations() []*UserSimulation {
	return userSimulations
}

// Start random user simulations and record produced events into
// the PostgresSQL database (TODO)
func Start() {

	fmt.Println("Starting Kafka producer...")

	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		// User-specific properties that you must set
		"bootstrap.servers": "<BOOTSTRAP SERVERS>",
		"sasl.username":     "<CLUSTER API KEY>",
		"sasl.password":     "<CLUSTER API SECRET>",

		// Fixed properties
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	// Go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go monitorEventsFromKafka(kafkaProducer)

	userIds := readUserIds()

	// For now let's start just one user simulation.
	userId := userIds[rand.Intn(len(userIds))]

	fmt.Println("Starting user simulation...")

	simulation := NewUserSimulation(userId)
	simulation.Start([]string{"sign_in", "view_page"})
	userSimulations = append(userSimulations, &simulation)

	for {
		select {
		case event := <-simulation.outgoingEvents:
			fmt.Printf("Attempting to send %v\n", event)
		}
	}
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
