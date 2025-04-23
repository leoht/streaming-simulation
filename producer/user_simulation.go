package producer

import (
	"math/rand"
	"time"
)

type UserSimulation struct {
	userId     string
	sentEvents int
	lastEvent  *Event

	// A channel for the simulation to communicate events it produced (for logging and saving into the SQL DB)
	outgoingEvents chan Event
}

func NewUserSimulation(userId string) UserSimulation {
	return UserSimulation{
		userId,
		0,
		nil,
		make(chan Event),
	}
}

// Launches a goroutine which will start emitting events from this user,
// simulating some traffic activity which some logic to it
// (e.g signing up must precede signing in, adding an item to cart must precede buying, etc).
func (us *UserSimulation) Start(userId string, availableEventNames []string) {
	go us.doStart(userId, availableEventNames)
}

func (us *UserSimulation) doStart(userId string, availableEventNames []string) {
	for {
		event := NewEvent(userId, "sign_up")
		event.Validate()

		// Send first sign up event for a start
		us.outgoingEvents <- event

		// Sleep some random time
		seconds := rand.Intn(5)
		time.Sleep(time.Duration(seconds) * time.Second)
	}
}
