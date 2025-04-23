package producer

import (
	"slices"
	"testing"
)

func TestCreateEventWithAvailableUserIdsAndEventNames(t *testing.T) {
	availableIds := []string{"ebb92b43-2113-4947-be5b-69db05928127", "c1986785-1e08-4cbe-878b-b31b61a06ae5", "b7e989bc-7755-4b0a-8647-29cf684e3150", "c5c167cc-9e76-4169-9984-b455130d932e"}
	availableEventNames := []string{"sign_in", "sign_up"}

	for i := 0; i < 100; i++ {
		event := CreateRandomEvent(availableIds, availableEventNames)
		if !slices.Contains(availableIds, event.UserId) {
			t.Errorf(`Event contains invalid user ID %s`, event.UserId)
		}
		if !slices.Contains(availableEventNames, event.EventName) {
			t.Errorf(`Event contains invalid event name %s`, event.EventName)
		}
	}
}
