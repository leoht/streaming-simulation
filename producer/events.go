package producer

import (
	"log"
	"math/rand"

	"github.com/google/uuid"
)

type Event struct {
	Id        string
	UserId    string
	EventName string
}

func NewEvent(userId string, eventName string) Event {
	eventId, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}

	return Event{
		eventId.String(),
		userId,
		eventName,
	}
}

func CreateRandomEvent(userId string, availableEventNames []string) Event {

	// Pick random values
	eventName := availableEventNames[rand.Intn(len(availableEventNames))]

	eventId, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}

	return Event{eventId.String(), userId, eventName}
}

func (e Event) Validate() {
	// TODO: Validate event or log warning?
}
