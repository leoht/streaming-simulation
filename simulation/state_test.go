package simulation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialStateIsSignedUp(t *testing.T) {
	userId := "ebb92b43-2113-4947-be5b-69db05928127"
	state := NewUserState(userId)

	if state.FSM.Current() != "signed_up" {
		t.Errorf("Initial state should be signed_up, got: %s", state.FSM.Current())
	}
}

func TestCanTransitionToValidEvents(t *testing.T) {
	userId := "ebb92b43-2113-4947-be5b-69db05928127"
	state := NewUserState(userId)

	_, err := state.Trigger("sign_in")

	assert.Nil(t, err)

	_, err2 := state.Trigger("view_page")

	assert.Nil(t, err2)
}

func TestCannotTransitionToInvalidEvents(t *testing.T) {
	userId := "ebb92b43-2113-4947-be5b-69db05928127"
	state := NewUserState(userId)

	_, err := state.Trigger("add_to_cart")

	assert.Error(t, err)

	_, err2 := state.Trigger("order")

	assert.Error(t, err2)
}
