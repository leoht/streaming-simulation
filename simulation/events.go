package simulation

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

func CreateRandomEvent(userId string, possibleEvents []string) Event {
	eventName := possibleEvents[rand.Intn(len(possibleEvents))]

	eventId, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}

	return Event{eventId.String(), userId, eventName}
}
