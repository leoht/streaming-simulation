package simulation

import (
	"slices"
	"testing"
	"time"
)

func receiveOnOutgoingEventsChannel(simulation UserSimulation) {
	for {
		<-simulation.outgoingEvents
	}
}

func TestCreateEventWithAvailableUserIdsAndEventNames(t *testing.T) {
	userId := "ebb92b43-2113-4947-be5b-69db05928127"
	availableEventNames := []string{"sign_in", "sign_up"}

	for i := 0; i < 100; i++ {
		event := CreateRandomEvent(userId, availableEventNames)

		if !slices.Contains(availableEventNames, event.EventName) {
			t.Errorf(`Event contains invalid event name %s`, event.EventName)
		}
	}
}

func TestCreateUserSimulationSendsSignupEvent(t *testing.T) {
	userId := "ebb92b43-2113-4947-be5b-69db05928127"
	simulation := NewUserSimulation(userId)
	simulation.Start()

	// TODO: better way to do this?
	select {
	case event := <-simulation.outgoingEvents:
		if event.EventName != "sign_up" {
			t.Errorf("First event received from simulation was not sign_up")
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Did not receive sign up event from user simulation after 1 second")
	}
}

func TestCreateUserSimulationSendsSignupThenOtherEvent(t *testing.T) {
	userId := "ebb92b43-2113-4947-be5b-69db05928127"
	simulation := NewUserSimulation(userId)
	simulation.Start()

	var gotSignup = false
	var gotOther = false

	select {
	case event := <-simulation.outgoingEvents:
		if event.EventName == "sign_up" {
			gotSignup = true
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Did not receive signup from user simulation after 5 second")
	}

	select {
	case event := <-simulation.outgoingEvents:
		// fmt.Println(event)
		if event.EventName == "sign_in" || event.EventName == "view_page" {
			gotOther = true

		}
	case <-time.After(6 * time.Second):
		t.Errorf("Did not receive other event from user simulation after 5 second")
	}

	if !(gotSignup && gotOther) {
		t.Errorf("Did not receive signup or other event: gotSignup = %v, gotOther = %v", gotSignup, gotOther)
	}
}

func TestStopUserSimulation(t *testing.T) {
	userId := "ebb92b43-2113-4947-be5b-69db05928127"
	simulation := NewUserSimulation(userId)
	simulation.Start()

	// Needed otherwise sending outgoing messages channel is blocking
	// TODO: improve this?
	go receiveOnOutgoingEventsChannel(simulation)

	// Stop() is blocking because of sending
	// to channel - improve this?
	go func() {
		simulation.Stop()
	}()

	time.Sleep(time.Duration(1) * time.Second)

	if simulation.Running {
		t.Errorf("Did not stop simulation")
	}
}

func TestStopAndResumeUserSimulation(t *testing.T) {
	userId := "ebb92b43-2113-4947-be5b-69db05928127"
	simulation := NewUserSimulation(userId)
	simulation.Start()

	// Needed otherwise sending outgoing messages channel is blocking
	// TODO: improve this?
	go receiveOnOutgoingEventsChannel(simulation)

	// Stop() is blocking because of sending
	// to channel - improve this?
	go func() {
		simulation.Stop()
	}()

	time.Sleep(time.Duration(100) * time.Millisecond)

	if simulation.Running {
		t.Errorf("Did not stop simulation")
	}

	go func() {
		simulation.Resume()
	}()

	time.Sleep(time.Duration(100) * time.Millisecond)

	if !simulation.Running {
		t.Errorf("Did not resume simulation")
	}
}
