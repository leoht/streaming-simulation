package producer

import (
	"slices"
	"testing"
	"time"
)

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
	simulation.Start(userId, []string{"sign_up", "sign_in"})

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
