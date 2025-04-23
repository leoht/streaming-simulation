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

func CreateRandomEvent(availableUserIds []string, availableEventNames []string) Event {

	// Pick random values
	userId := availableUserIds[rand.Intn(len(availableUserIds))]
	eventName := availableEventNames[rand.Intn(len(availableEventNames))]

	eventId, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}

	return Event{eventId.String(), userId, eventName}
}
