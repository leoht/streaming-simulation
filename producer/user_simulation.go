package producer

import (
	"log"
	"math/rand"
	"time"
)

type UserSimulation struct {
	running    bool
	userId     string
	sentEvents int
	lastEvent  *Event

	// A channel for the simulation to communicate events it produced (for logging and saving into the SQL DB)
	outgoingEvents chan Event
	stopChannel    chan bool
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
func (us *UserSimulation) Start(userId string, availableEventNames []string) {
	go us.doStart(userId, availableEventNames)
}

func (us *UserSimulation) doStart(userId string, availableEventNames []string) {
	us.running = true

loop:
	for {
		select {
		case <-us.stopChannel:
			log.Printf("Stopping simulation for %s", us.userId)
			us.running = false
			break loop
		default:
			signupEvent := NewEvent(userId, "sign_up")
			signupEvent.Validate()

			// Send first sign up event for a start
			us.outgoingEvents <- signupEvent

			// Sleep some random time
			seconds := rand.Intn(5)
			time.Sleep(time.Duration(seconds) * time.Second)

			// Send other event now
			event := CreateRandomEvent(userId, availableEventNames)
			us.outgoingEvents <- event
		}
	}
}
