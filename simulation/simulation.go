package simulation

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"
)

/**
* Core business logic
* The simulation package runs a 'simulation'
* and manages many 'user simulations', each of
* them managing a goroutine for a given user ID
* which simulates user activity by emitting 'events'
 */

type Simulation struct {
	availableUserIds []string

	// TODO: replace with a map UserId -> Sim for easier access by UserId
	userSimulations []*UserSimulation
	// Buffered channel to forward events to producer
	producerChannel chan Event
}

// For now there is only one running simulation
// TODO: allow many?
var currentSimulation *Simulation

func StartSimulation() *Simulation {
	userIds := readUserIds()

	currentSimulation = &Simulation{
		userIds,
		[]*UserSimulation{},
		make(chan Event, 10),
	}

	return currentSimulation
}

func AllUserSimulations() []*UserSimulation {
	return currentSimulation.userSimulations
}

func Current() *Simulation {
	return currentSimulation
}

func (s *Simulation) ProducerChannel() chan Event {
	return s.producerChannel
}

func (s *Simulation) nextUserId() string {
	return currentSimulation.availableUserIds[rand.Intn(len(s.availableUserIds))]
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
