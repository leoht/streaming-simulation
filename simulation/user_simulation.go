package simulation

import (
	"context"
	"log"
	"math/rand"
	"time"
)

/**
* A UserSimulation can be started, stopped or resumed
* to emit semi-random events for a given user.
* The possible chain of events is dictated by the finite state machine
* defined in state.go
 */
type UserSimulation struct {
	*UserState
	Running    bool
	UserId     string
	sentEvents int
	lastEvent  *Event

	// A buffered channel to communicate events produced
	outgoingEvents chan Event
	stopChannel    chan bool
}

// Start a simulation for a user who doesn't have a simulation started yet
func StartNewUserSimulation() *UserSimulation {
	log.Println("Starting user simulation...")

	simulation := NewUserSimulation(currentSimulation.nextUserId())
	simulation.Start()
	currentSimulation.userSimulations = append(currentSimulation.userSimulations, &simulation)

	// Forward each outgoing event to kafka producerc channel
	go func() {
		for simulation.Running {
			event := <-simulation.outgoingEvents
			currentSimulation.producerChannel <- event
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
		NewUserState(userId),
		false,
		userId,
		0,
		nil,
		make(chan Event, 10),
		make(chan bool),
	}
}

// Launches a goroutine which will start emitting events from this user,
// simulating some traffic activity which some logic to it
// (e.g signing up must precede signing in, adding an item to cart must precede buying, etc).
func (us *UserSimulation) Start() {
	go us.doStart()
}

func (us *UserSimulation) Stop() {
	us.stopChannel <- true
}

func (us *UserSimulation) Resume() {
	// TODO resume correctly?
	go us.doLoop()
}

func (us *UserSimulation) doStart() {
	us.Running = true

	signupEvent := NewEvent(us.UserId, "sign_up")
	if err := us.validateAndEmitEvent(signupEvent); err != nil {
		return
	}

	us.doLoop()
}

func (us *UserSimulation) doLoop() {
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
			event := us.nextEvent()

			if err := us.validateAndEmitEvent(event); err != nil {
				continue
			}
		}
	}
}

func (us *UserSimulation) nextEvent() Event {
	return CreateRandomEvent(us.UserId, us.UserState.FSM.AvailableTransitions())
}

func (us *UserSimulation) validateAndEmitEvent(event Event) error {
	if err, _ := us.validateNextEventWithFSM(event); err != nil {
		// If errored on transition, return
		return err
	}

	us.outgoingEvents <- event

	// TODO is this all redundant data?
	us.lastEvent = &event
	us.sentEvents += 1

	return nil
}

func (us *UserSimulation) validateNextEventWithFSM(event Event) (error, Event) {
	err := us.UserState.FSM.Event(context.Background(), event.EventName)

	if err != nil {
		log.Printf("Couldn't perform invalid transition with event: %s (current: %s)", event.EventName, us.UserState.FSM.Current())
	}

	return err, event
}
