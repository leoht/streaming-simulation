package simulation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialStateIsNew(t *testing.T) {
	userId := "ebb92b43-2113-4947-be5b-69db05928127"
	state := NewUserState(userId)

	if state.FSM.Current() != "new" {
		t.Errorf("Initial state should be new, got: %s", state.FSM.Current())
	}
}

func TestCanTransitionToValidEvents(t *testing.T) {
	userId := "ebb92b43-2113-4947-be5b-69db05928127"
	state := NewUserState(userId)

	_, err := state.Trigger("sign_up")

	assert.Nil(t, err)

	_, err2 := state.Trigger("sign_in")

	assert.Nil(t, err2)
	assert.Equal(t, "signed_in", state.FSM.Current(), "State should be signed_in")
}

func TestCannotTransitionToInvalidEvents(t *testing.T) {
	userId := "ebb92b43-2113-4947-be5b-69db05928127"
	state := NewUserState(userId)

	_, err := state.Trigger("add_to_cart")

	assert.Error(t, err)

	_, err2 := state.Trigger("order")

	assert.Error(t, err2)
}
