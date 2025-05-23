package simulation

// TODO state machine for user simulation

import (
	"context"

	"github.com/looplab/fsm"
)

type UserState struct {
	UserId    string
	LastEvent *Event
	FSM       *fsm.FSM
}

func NewUserState(userId string) *UserState {
	state := &UserState{
		UserId: userId,
	}

	// State machine diagram:
	// (simple for now)
	//
	// signed_up -> signed_in -> viewed_page -> added_to_cart -> signed_in
	state.FSM = fsm.NewFSM(
		"new",
		fsm.Events{
			{Name: "sign_up", Src: []string{"new"}, Dst: "signed_up"},
			{Name: "sign_in", Src: []string{"signed_up"}, Dst: "signed_in"},
			{Name: "view_page", Src: []string{"signed_in", "added_to_cart"}, Dst: "viewed_page"},
			{Name: "add_to_cart", Src: []string{"viewed_page"}, Dst: "added_to_cart"},
			{Name: "order", Src: []string{"added_to_cart"}, Dst: "signed_in"},
		},
		fsm.Callbacks{
			// "enter_state": func(_ context.Context, e *fsm.Event) { state.enterState(e) },
		},
	)

	return state
}

func (state *UserState) Trigger(eventName string) (*UserState, error) {
	err := state.FSM.Event(context.Background(), eventName)

	return state, err
}
