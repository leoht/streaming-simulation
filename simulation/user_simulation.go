package simulation

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type UserSimulation struct {
	Running    bool
	UserId     string
	sentEvents int
	lastEvent  *Event

	// A channel for the simulation to communicate events it produced (for logging and saving into the SQL DB)
	outgoingEvents chan Event
	stopChannel    chan bool
}

// Start a simulation for a user who doesn't have a simulation started yet
func StartNewSimulation(userIds []string, kafkaProducerInChannel chan Event) *UserSimulation {
	userId := userIds[rand.Intn(len(userIds))]

	fmt.Println("Starting user simulation...")

	simulation := NewUserSimulation(userId)
	simulation.Start([]string{"sign_in", "view_page"})
	currentSimulation.userSimulations = append(currentSimulation.userSimulations, &simulation)

	go func() {
		for simulation.Running {
			select {
			case event := <-simulation.outgoingEvents:
				kafkaProducerInChannel <- event
			}
		}
	}()

	return &simulation
}

func StopSimulationForUser(userId string) *UserSimulation {
	for _, simulation := range currentSimulation.userSimulations {
		if simulation.UserId == userId && simulation.Running {
			simulation.Stop()
			return simulation
		}
	}

	return nil
}

func ResumeSimulationForUser(userId string) *UserSimulation {
	for _, simulation := range currentSimulation.userSimulations {
		if simulation.UserId == userId && !simulation.Running {
			simulation.Resume()
			return simulation
		}
	}

	return nil
}

func NewUserSimulation(userId string) UserSimulation {
	return UserSimulation{
		false,
		userId,
		0,
		nil,
		make(chan Event),
		make(chan bool),
	}
}

// Launches a goroutine which will start emitting events from this user,
// simulating some traffic activity which some logic to it
// (e.g signing up must precede signing in, adding an item to cart must precede buying, etc).
func (us *UserSimulation) Start(availableEventNames []string) {
	go us.doStart(availableEventNames)
}

func (us *UserSimulation) Stop() {
	us.stopChannel <- true
}

func (us *UserSimulation) Resume() {
	// TODO resume correctly
	go us.doStart([]string{"sign_in", "view_page"})
}

func (us *UserSimulation) doStart(availableEventNames []string) {
	us.Running = true

	signupEvent := NewEvent(us.UserId, "sign_up")
	signupEvent.Validate()

	// Send first sign up event for a start
	us.outgoingEvents <- signupEvent

	// Then loop

	for us.Running {
		select {
		case <-us.stopChannel:
			log.Printf("Stopping simulation for %s", us.UserId)
			us.Running = false
		default:
			// Sleep some random time

			// TODO improve/parameterise these values
			min := 2
			max := 6
			seconds := rand.Intn(max-min) + min
			time.Sleep(time.Duration(seconds) * time.Second)

			// Send other event now
			event := CreateRandomEvent(us.UserId, availableEventNames)
			us.outgoingEvents <- event
		}
	}
}
